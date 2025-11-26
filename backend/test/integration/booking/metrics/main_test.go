package metrics_test

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

	metricsService "github.com/Leviosa-care/leviosa/backend/internal/booking/application/metrics"
	metricsPostgres "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/postgres/metrics"
	metricsHandler "github.com/Leviosa-care/leviosa/backend/internal/booking/interface/metrics"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"

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
	pgContainer     *tu.PostgresContainer
	testPool        *pgxpool.Pool
	redisContainer  *tu.RedisContainer
	redisClient     *redis.Client
	crypto          encx.CryptoService
	metricsRepo     ports.MetricsRepository
	service         ports.MetricsService
	authCtx         *tu.AuthTestContext
	handler         *metricsHandler.Handler
	testServerURL   string
	testServer      *http.Server
	authSessionRepo authsession.SessionRepository
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create and configure logger for tests
	loggerHandler, err := logger.SetHandler("debug", "dev")
	if err != nil {
		log.Fatalf("Failed to create logger handler: %v", err)
	}
	testLogger := slog.New(loggerHandler)
	slog.SetDefault(testLogger)

	// Postgres container
	pgContainer, err = tu.SetupPostgres(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to setup postgres container: %v", err)
	}
	defer tu.TeardownPostgres(ctx, nil, pgContainer)

	// DB Pool
	log.Println("Creating pgxpool...")
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
	log.Println("pgxpool created.")

	// Ping database
	if err = testPool.Ping(poolCtx); err != nil {
		panic(fmt.Sprintf("Failed to ping database pool: %v", err))
	}
	log.Println("Database pool ping successful.")

	// Database migrations
	log.Println("Applying database migrations...")
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
	log.Println("Migrations applied.")

	// Redis container
	redisContainer, err = tu.SetupRedis(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to setup redis container: %v", err))
	}
	defer tu.TeardownRedis(ctx, nil, redisContainer)

	// Redis client
	log.Println("Creating Redis client...")
	redisClient = redisContainer.NewClient()

	if err = redisClient.Ping(ctx).Err(); err != nil {
		tu.TeardownRedis(ctx, nil, redisContainer)
		panic(fmt.Sprintf("Failed to ping Redis: %v", err))
	}
	log.Println("Redis client connected successfully.")

	// Setup Vault testcontainer with per-service encryption
	log.Println("Setting up Vault testcontainer...")
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
	log.Printf("✓ Booking service crypto initialized")

	// Initialize AuthTestContext
	authCtx = &tu.AuthTestContext{
		Pool:   testPool,
		Redis:  redisClient,
		Crypto: crypto,
	}

	ctx = context.WithValue(ctx, ctxutil.LoggerKey, testLogger)

	// Initialize application layers
	metricsRepo = metricsPostgres.New(testPool)
	service = metricsService.New(metricsRepo, crypto)

	authSessionRepo = authsession.NewRedisSessionRepository(redisClient)
	authmw := auth.NewSessionAuthMiddleware(authSessionRepo, crypto, nil)

	handler = metricsHandler.New(service, authmw)

	// Set required environment variables
	os.Setenv("CLIENT_IP_HEADER", "X-Forwarded-For")
	os.Setenv("LOGGING_SALT", "test_logging_salt_12345")

	// HTTP server setup
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
	log.Printf("Test HTTP server started at %s", testServerURL)

	time.Sleep(100 * time.Millisecond)

	// Run tests
	code := m.Run()

	// Graceful shutdown
	log.Println("Shutting down test HTTP server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := testServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Test server shutdown failed: %v", err)
	}
	log.Println("Test HTTP server shut down.")

	os.Exit(code)
}
