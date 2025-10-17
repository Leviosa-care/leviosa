package partner_test

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"log/slog"
// 	"net"
// 	"net/http"
// 	"os"
// 	"testing"
// 	"time"
//
// 	partnerHandler "github.com/Leviosa-care/authuser/internal/adapters/http/partner"
// 	specializationHandler "github.com/Leviosa-care/authuser/internal/adapters/http/specialization"
// 	partnerRepository "github.com/Leviosa-care/authuser/internal/adapters/postgres/partner"
// 	specializationRepository "github.com/Leviosa-care/authuser/internal/adapters/postgres/specialization"
// 	userRepository "github.com/Leviosa-care/authuser/internal/adapters/postgres/user"
// 	sessionRepository "github.com/Leviosa-care/authuser/internal/adapters/redis/session"
// 	authPayment "github.com/Leviosa-care/authuser/internal/adapters/stripe"
// 	partnerService "github.com/Leviosa-care/authuser/internal/application/partner"
// 	specializationService "github.com/Leviosa-care/authuser/internal/application/specialization"
// 	userService "github.com/Leviosa-care/authuser/internal/application/user"
// 	"github.com/Leviosa-care/authuser/internal/ports"
//
// 	"github.com/Leviosa-care/core/auth/session"
// 	"github.com/Leviosa-care/core/envmode"
// 	"github.com/Leviosa-care/core/logger"
// 	"github.com/Leviosa-care/core/middleware"
// 	"github.com/Leviosa-care/core/middleware/auth"
// 	"github.com/Leviosa-care/core/migrations"
// 	tu "github.com/Leviosa-care/core/testutils"
// 	"github.com/hengadev/encx"
// 	"github.com/hengadev/encx/providers/hashicorpvault"
// 	"github.com/jackc/pgx/v5/pgxpool"
// 	"github.com/pressly/goose/v3"
// 	"github.com/redis/go-redis/v9"
// )
//
// var (
// 	pgContainer           *tu.PostgresContainer
// 	testPool              *pgxpool.Pool
// 	redisContainer        *tu.RedisContainer
// 	testClient            *redis.Client
// 	crypto                encx.CryptoService
// 	userRepo              ports.UserRepository
// 	partnerRepo           ports.PartnerRepository
// 	specializationRepo    ports.SpecializationRepository
// 	authSessionRepo       session.SessionRepository
// 	sessionRepo           ports.SessionRepository
// 	userSvc               ports.UserService
// 	partnerSvc            ports.PartnerService
// 	specializationSvc     ports.SpecializationService
// 	partnerHandler        partnerHandler.Handler
// 	specializationHandler specializationHandler.Handler
// 	testServerURL         string       // Global variable to hold the URL of the running test server
// 	testServer            *http.Server // To allow graceful shutdown
// )
//
// func TestMain(m *testing.M) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
//
// 	// Create and configure logger for tests
// 	loggerHandler, err := logger.SetHandler("debug", "dev")
// 	if err != nil {
// 		log.Fatalf("Failed to create logger handler: %v", err)
// 	}
// 	testLogger := slog.New(loggerHandler)
// 	slog.SetDefault(testLogger) // Set as default for the application
//
// 	// Postgres container
// 	pgContainer, err = tu.SetupPostgres(ctx, nil)
// 	if err != nil {
// 		log.Fatalf("Failed to setup postgres container: %v", err)
// 	}
// 	defer tu.TeardownPostgres(ctx, nil, pgContainer)
//
// 	// DB Pool
// 	log.Println("Creating pgxpool...")
// 	poolCtx, poolCancel := context.WithTimeout(ctx, 10*time.Second)
// 	defer poolCancel()
//
// 	pgCfg, err := pgxpool.ParseConfig(pgContainer.ConnectionString)
// 	if err != nil {
// 		panic(fmt.Sprintf("Failed to parse pgxpool config: %v", err))
// 	}
//
// 	// Configure pool settings for tests
// 	pgCfg.MaxConns = 5
// 	pgCfg.MinConns = 1
// 	pgCfg.MaxConnLifetime = 30 * time.Minute
// 	pgCfg.MaxConnIdleTime = 5 * time.Minute
//
// 	testPool, err = pgxpool.NewWithConfig(poolCtx, pgCfg)
// 	if err != nil {
// 		tu.TeardownPostgres(ctx, nil, pgContainer)
// 		panic(fmt.Sprintf("Failed to open test database pool: %v", err))
// 	}
// 	log.Println("pgxpool created.")
//
// 	// Ping database
// 	if err = testPool.Ping(poolCtx); err != nil {
// 		panic(fmt.Sprintf("Failed to ping database pool: %v", err))
// 	}
// 	log.Println("Database pool ping successful.")
//
// 	// Database migrations
// 	log.Println("Applying database migrations...")
// 	goose.SetBaseFS(migrations.FS)
// 	if err = goose.SetDialect("pgx"); err != nil {
// 		log.Fatalf("Setting dialect for migrations: %s\n", err)
// 	}
//
// 	gooseDB, err := sql.Open("pgx", testPool.Config().ConnString())
// 	if err != nil {
// 		panic(fmt.Sprintf("Failed to open temp *sql.DB for goose migrations: %v", err))
// 	}
// 	defer gooseDB.Close()
//
// 	if err = goose.UpContext(ctx, gooseDB, "."); err != nil {
// 		panic(fmt.Sprintf("running all migrations: %s\n", err))
// 	}
// 	log.Println("Migrations applied.")
//
// 	// stripe container
// 	stripeContainer, err := tu.SetupStripeMock(ctx, nil)
// 	if err != nil {
// 		log.Fatalf("Failed to setup stripe container: %v", err)
// 	}
// 	defer tu.TeardownStripeMock(ctx, nil, stripeContainer)
//
// 	// Setup Vault testcontainer
// 	log.Println("Setting up Vault testcontainer...")
// 	vaultContainer, err := tu.SetupVault(ctx, nil)
// 	if err != nil {
// 		log.Fatalf("Failed to setup Vault container: %v", err)
// 	}
// 	defer tu.TeardownVault(ctx, nil, vaultContainer)
//
// 	// Set environment variables for Vault
// 	os.Setenv("VAULT_ADDR", vaultContainer.HTTPSEndpoint)
// 	os.Setenv("VAULT_TOKEN", vaultContainer.RootToken)
//
// 	// Crypto service
// 	log.Println("Creating crypto service...")
// 	kms, err := hashicorpvault.New()
// 	if err != nil {
// 		log.Fatalf("Failed to create vault provider: %v", err)
// 	}
//
// 	crypto, err = encx.New(
// 		ctx,
// 		kms,
// 		tu.EncryptionKey,
// 		"secret/data/pepper",
// 	)
// 	if err != nil {
// 		log.Fatalf("Failed to create crypto service: %v", err)
// 	}
// 	if crypto == nil {
// 		log.Fatal("Crypto service is nil after creation")
// 	}
// 	log.Println("Crypto service created successfully")
//
// 	// Redis container
// 	redisContainer, err = tu.SetupRedis(ctx, nil)
// 	if err != nil {
// 		panic(fmt.Sprintf("Failed to setup redis container: %v", err))
// 	}
// 	defer tu.TeardownRedis(ctx, nil, redisContainer)
//
// 	// Redis client
// 	log.Println("Creating Redis client...")
// 	testClient = redisContainer.NewClient()
//
// 	// Test Redis connection
// 	if err = testClient.Ping(ctx).Err(); err != nil {
// 		tu.TeardownRedis(ctx, nil, redisContainer)
// 		panic(fmt.Sprintf("Failed to ping Redis: %v", err))
// 	}
// 	log.Println("Redis client connected successfully.")
//
// 	payment := authPayment.NewService("sk_test_123456789012345678901234", stripeContainer.URL)
//
// 	// Initialize repositories
// 	userRepo = userRepository.New(ctx, testPool)
// 	partnerRepo = partnerRepository.New(ctx, testPool)
// 	specializationRepo = specializationRepository.New(ctx, testPool)
//
// 	// Initialize services
// 	userSvc = userService.New(userRepo, crypto, payment)
// 	partnerSvc = partnerService.New(partnerRepo, userRepo, specializationRepo, crypto, payment)
// 	specializationSvc = specializationService.New(specializationRepo, crypto)
//
// 	sessionRepo = sessionRepository.New(testClient)
//
// 	authSessionRepo = session.NewRedisSessionRepository(testClient)
// 	authmw := auth.NewSessionAuthMiddleware(authSessionRepo, crypto, nil)
//
// 	// Initialize handlers
// 	partnerHandler = partnerHandler.New(partnerSvc, authmw)
// 	specializationHandler = specializationHandler.New(specializationSvc, authmw)
//
// 	// Set required environment variables for logger middleware
// 	os.Setenv("CLIENT_IP_HEADER", "X-Forwarded-For")
// 	os.Setenv("LOGGING_SALT", "test_logging_salt_12345")
//
// 	// HTTP server setup with logger middleware
// 	router := http.NewServeMux()
// 	partnerHandler.RegisterRoutes(router)
// 	specializationHandler.RegisterRoutes(router)
//
// 	// Use the enhanced AttachLogger middleware from core package
// 	loggerMiddleware := middleware.AttachLogger(envmode.Dev, testLogger)
//
// 	listener, err := net.Listen("tcp", ":0") // Random available port
// 	if err != nil {
// 		log.Fatalf("Failed to listen for test server: %v", err)
// 	}
// 	testServerURL = "http://" + listener.Addr().String()
// 	testServer = &http.Server{Handler: loggerMiddleware(router)}
//
// 	// Start server in goroutine
// 	go func() {
// 		if err := testServer.Serve(listener); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("Test server failed to serve: %v", err)
// 		}
// 	}()
// 	log.Printf("Test HTTP server started at %s", testServerURL)
//
// 	// Give server time to start
// 	time.Sleep(100 * time.Millisecond)
//
// 	// Run tests
// 	code := m.Run()
//
// 	// Graceful shutdown
// 	log.Println("Shutting down test HTTP server...")
// 	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer shutdownCancel()
// 	if err := testServer.Shutdown(shutdownCtx); err != nil {
// 		log.Fatalf("Test server shutdown failed: %v", err)
// 	}
// 	log.Println("Test HTTP server shut down.")
//
// 	// Exit with test result code
// 	os.Exit(code)
// }
//

