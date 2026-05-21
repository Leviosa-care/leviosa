package admin

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	bookingService "github.com/Leviosa-care/leviosa/backend/internal/booking/application/booking"
	availabilityPostgres "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/postgres/availability"
	bookingPostgres "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/postgres/booking"
	bookingHandler "github.com/Leviosa-care/leviosa/backend/internal/booking/interface/booking"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"

	productService "github.com/Leviosa-care/leviosa/backend/internal/catalog/application/product"
	priceService "github.com/Leviosa-care/leviosa/backend/internal/catalog/application/price"
	pricePostgres "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/price"
	productPostgres "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/product"
	sharedPostgres "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/shared"
	pricePayment "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/stripe/price"
	productPayment "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/stripe/product"

	authsession "github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/services"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/envmode"
	"github.com/Leviosa-care/leviosa/backend/internal/common/logger"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
	"github.com/Leviosa-care/leviosa/backend/internal/common/migrations"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"

	"github.com/hengadev/encx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
)

var (
	pgContainer      *tu.PostgresContainer
	testPool         *pgxpool.Pool
	redisContainer   *tu.RedisContainer
	redisClient      *redis.Client
	crypto           encx.CryptoService
	bookingRepo      ports.BookingRepository
	availabilityRepo ports.AvailabilityRepository
	paymentService   ports.PaymentService
	service          ports.BookingService
	authCtx          *tu.AuthTestContext
	handler          bookingHandler.Handler
	testServerURL    string
	testServer       *http.Server
	authSessionRepo  authsession.SessionRepository
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	loggerHandler, err := logger.SetHandler("debug", "dev")
	if err != nil {
		log.Fatalf("Failed to create logger handler: %v", err)
	}
	testLogger := slog.New(loggerHandler)
	slog.SetDefault(testLogger)

	pgContainer, err = tu.SetupPostgres(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to setup postgres container: %v", err)
	}
	defer tu.TeardownPostgres(ctx, nil, pgContainer)

	poolCtx, poolCancel := context.WithTimeout(ctx, 10*time.Second)
	defer poolCancel()

	pgCfg, err := pgxpool.ParseConfig(pgContainer.ConnectionString)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse pgxpool config: %v", err))
	}

	pgCfg.MaxConns = 5
	pgCfg.MinConns = 1
	pgCfg.MaxConnLifetime = 30 * time.Minute
	pgCfg.MaxConnIdleTime = 5 * time.Minute

	testPool, err = pgxpool.NewWithConfig(poolCtx, pgCfg)
	if err != nil {
		tu.TeardownPostgres(ctx, nil, pgContainer)
		panic(fmt.Sprintf("Failed to open test database pool: %v", err))
	}

	if err = testPool.Ping(poolCtx); err != nil {
		panic(fmt.Sprintf("Failed to ping database pool: %v", err))
	}

	goose.SetBaseFS(migrations.FS)
	if err = goose.SetDialect("pgx"); err != nil {
		log.Fatalf("Setting dialect for migrations: %s\n", err)
	}

	gooseDB, err := sql.Open("pgx", testPool.Config().ConnString())
	if err != nil {
		panic(fmt.Sprintf("Failed to open temp *sql.DB for goose migrations: %v", err))
	}
	defer gooseDB.Close()

	if err = goose.UpContext(ctx, gooseDB, "."); err != nil {
		panic(fmt.Sprintf("running all migrations: %s\n", err))
	}

	redisContainer, err = tu.SetupRedis(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to setup redis container: %v", err))
	}
	defer tu.TeardownRedis(ctx, nil, redisContainer)

	redisClient = redisContainer.NewClient()

	if err = redisClient.Ping(ctx).Err(); err != nil {
		tu.TeardownRedis(ctx, nil, redisContainer)
		panic(fmt.Sprintf("Failed to ping Redis: %v", err))
	}

	serviceNames := []string{services.Booking}
	vaultSetup, err := tu.SetupServiceVault(ctx, nil, serviceNames)
	if err != nil {
		log.Fatalf("Failed to setup service Vault container: %v", err)
	}
	defer tu.TeardownVault(ctx, nil, vaultSetup.VaultContainer)

	os.Setenv("VAULT_ADDR", vaultSetup.VaultContainer.HTTPSEndpoint)
	os.Setenv("VAULT_TOKEN", vaultSetup.VaultContainer.RootToken)

	var exists bool
	crypto, exists = vaultSetup.GetServiceCrypto(services.Booking)
	if !exists {
		log.Fatal("Booking service crypto not found in vault setup")
	}
	if crypto == nil {
		log.Fatal("Booking crypto service is nil")
	}

	authCtx = &tu.AuthTestContext{
		Pool:   testPool,
		Redis:  redisClient,
		Crypto: crypto,
	}

	ctx = context.WithValue(ctx, ctxutil.LoggerKey, testLogger)

	bookingRepo = bookingPostgres.New(ctx, testPool)
	availabilityRepo = availabilityPostgres.New(ctx, testPool)
	paymentService = NewMockPaymentService()

	productRepo := productPostgres.New(ctx, testPool)
	sharedRepo := sharedPostgres.New(ctx, testPool)

	stripeTestKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeTestKey == "" {
		stripeTestKey = "sk_test_dummy_key_for_testing"
	}
	stripeBaseURL := os.Getenv("STRIPE_API_BASE_URL")

	productStripe := productPayment.NewProduct(stripeTestKey, stripeBaseURL)
	priceStripe := pricePayment.NewPrice(stripeTestKey, stripeBaseURL)

	catalogProductService := productService.New(productRepo, sharedRepo, productStripe, priceStripe)

	priceRepo := pricePostgres.New(ctx, testPool)
	catalogPriceService := priceService.New(priceRepo, sharedRepo, priceStripe)

	service = bookingService.New(bookingRepo, availabilityRepo, paymentService, catalogProductService, catalogPriceService, nil, crypto)

	authSessionRepo = authsession.NewRedisSessionRepository(redisClient)
	authmw := auth.NewSessionAuthMiddleware(authSessionRepo, crypto, nil)

	handler = bookingHandler.New(service, paymentService, authmw)

	os.Setenv("CLIENT_IP_HEADER", "X-Forwarded-For")
	os.Setenv("LOGGING_SALT", "test_logging_salt_12345")

	router := http.NewServeMux()
	handler.RegisterRoutes(router)

	loggerMiddleware := middleware.AttachLogger(envmode.Dev, testLogger)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Failed to listen for test server: %v", err)
	}
	testServerURL = "http://" + listener.Addr().String()
	testServer = &http.Server{Handler: loggerMiddleware(router)}

	go func() {
		if err := testServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Test server failed to serve: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	code := m.Run()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := testServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Test server shutdown failed: %v", err)
	}

	os.Exit(code)
}
