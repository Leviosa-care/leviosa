package room_test

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

	buildingService "github.com/Leviosa-care/leviosa/backend/internal/booking/application/building"
	roomService "github.com/Leviosa-care/leviosa/backend/internal/booking/application/room"
	buildingPostgres "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/postgres/building"
	roomPostgres "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/postgres/room"
	buildingHandler "github.com/Leviosa-care/leviosa/backend/internal/booking/interface/building"
	roomHandler "github.com/Leviosa-care/leviosa/backend/internal/booking/interface/room"
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
	buildingRepo    ports.BuildingRepository
	roomRepo        ports.RoomRepository
	buildingSvc     ports.BuildingService
	roomSvc         ports.RoomService
	authCtx         *tu.AuthTestContext // Authentication context for user/session tests
	buildingHdl     buildingHandler.Handler
	roomHdl         roomHandler.Handler
	testServerURL   string       // Global variable to hold the URL of the running test server
	testServer      *http.Server // To allow graceful shutdown
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
	slog.SetDefault(testLogger) // Set as default for the application

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

	// Configure pool settings for tests
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

	// Test Redis connection
	if err = redisClient.Ping(ctx).Err(); err != nil {
		tu.TeardownRedis(ctx, nil, redisContainer)
		panic(fmt.Sprintf("Failed to ping Redis: %v", err))
	}
	log.Println("Redis client connected successfully.")

	// Setup enhanced Vault testcontainer with per-service encryption (GDPR compliant)
	log.Println("Setting up enhanced Vault testcontainer with per-service keys...")
	serviceNames := []string{services.Booking} // Only booking service for this test
	vaultSetup, err := tu.SetupServiceVault(ctx, nil, serviceNames)
	if err != nil {
		log.Fatalf("Failed to setup service Vault container: %v", err)
	}
	defer tu.TeardownVault(ctx, nil, vaultSetup.VaultContainer)

	// Set environment variables for Vault (for encx library)
	os.Setenv("VAULT_ADDR", vaultSetup.VaultContainer.HTTPSEndpoint)
	os.Setenv("VAULT_TOKEN", vaultSetup.VaultContainer.RootToken)

	// Get service-specific crypto service (GDPR compliant)
	var exists bool
	crypto, exists = vaultSetup.GetServiceCrypto(services.Booking)
	if !exists {
		log.Fatal("Booking service crypto not found in vault setup")
	}
	if crypto == nil {
		log.Fatal("Booking crypto service is nil")
	}
	log.Printf("✓ Booking service crypto initialized with per-service encryption key")

	// Log vault setup details for debugging
	log.Printf("✓ Vault setup complete:")
	log.Printf("  - Services: %d", len(vaultSetup.CryptoServices))
	log.Printf("  - API Keys: %d", len(vaultSetup.ServiceKeys))
	log.Printf("  - GDPR compliant per-service encryption: enabled")

	// Initialize AuthTestContext for user/session testing
	authCtx = &tu.AuthTestContext{
		Pool:   testPool,
		Redis:  redisClient,
		Crypto: crypto,
	}

	ctx = context.WithValue(ctx, ctxutil.LoggerKey, testLogger)

	// Initialize application layers
	buildingRepo = buildingPostgres.New(ctx, testPool)
	roomRepo = roomPostgres.New(ctx, testPool)
	buildingSvc = buildingService.New(buildingRepo, crypto)
	roomSvc = roomService.New(roomRepo, buildingRepo, crypto)

	authSessionRepo = authsession.NewRedisSessionRepository(redisClient)
	authmw := auth.NewSessionAuthMiddleware(authSessionRepo, crypto, nil)

	buildingHdl = buildingHandler.New(buildingSvc, authmw)
	roomHdl = roomHandler.New(roomSvc, authmw)

	// Set required environment variables for logger middleware
	os.Setenv("CLIENT_IP_HEADER", "X-Forwarded-For")
	os.Setenv("LOGGING_SALT", "test_logging_salt_12345")

	// HTTP server setup with logger middleware
	router := http.NewServeMux()
	buildingHdl.RegisterRoutes(router)
	roomHdl.RegisterRoutes(router)

	// Use the enhanced AttachLogger middleware from core package
	loggerMiddleware := middleware.AttachLogger(envmode.Dev, testLogger)

	listener, err := net.Listen("tcp", ":0") // Random available port
	if err != nil {
		log.Fatalf("Failed to listen for test server: %v", err)
	}
	testServerURL = "http://" + listener.Addr().String()
	testServer = &http.Server{Handler: loggerMiddleware(router)}

	// Start server in goroutine
	go func() {
		if err := testServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Test server failed to serve: %v", err)
		}
	}()
	log.Printf("Test HTTP server started at %s", testServerURL)

	// Give server time to start
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

	// Exit with test result code
	os.Exit(code)
}
