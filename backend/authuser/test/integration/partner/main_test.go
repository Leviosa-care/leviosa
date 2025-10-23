package partner_test

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

	partnerHandler "github.com/Leviosa-care/authuser/internal/adapters/http/partner"
	partnerRepository "github.com/Leviosa-care/authuser/internal/adapters/postgres/partner"
	userRepository "github.com/Leviosa-care/authuser/internal/adapters/postgres/user"
	authRabbitMQ "github.com/Leviosa-care/authuser/internal/adapters/rabbitmq"
	sessionRepository "github.com/Leviosa-care/authuser/internal/adapters/redis/session"
	authPayment "github.com/Leviosa-care/authuser/internal/adapters/stripe"
	"github.com/Leviosa-care/authuser/internal/application/catalog"
	"github.com/Leviosa-care/authuser/internal/application/partner"
	"github.com/Leviosa-care/authuser/internal/application/user"
	"github.com/Leviosa-care/authuser/internal/ports"

	"github.com/Leviosa-care/core/auth/session"
	mq "github.com/Leviosa-care/core/contracts/rabbitmq"
	"github.com/Leviosa-care/core/contracts/services"
	"github.com/Leviosa-care/core/ctxutil"
	"github.com/Leviosa-care/core/envmode"
	"github.com/Leviosa-care/core/logger"
	"github.com/Leviosa-care/core/messaging/rabbitmq"
	"github.com/Leviosa-care/core/middleware"
	"github.com/Leviosa-care/core/middleware/auth"
	"github.com/Leviosa-care/core/migrations"
	tu "github.com/Leviosa-care/core/testutils"
	"github.com/hengadev/encx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

var (
	pgContainer     *tu.PostgresContainer
	testPool        *pgxpool.Pool
	redisContainer  *tu.RedisContainer
	redisClient     *redis.Client
	crypto          encx.CryptoService
	userRepo        ports.UserRepository
	partnerRepo     ports.PartnerRepository
	catalogCache    ports.CatalogCache
	vaultSetup      *tu.ServiceVaultSetup // Enhanced Vault setup with per-service keys
	authCtx         *tu.AuthTestContext   // Authentication context for user/session tests
	authSessionRepo session.SessionRepository
	sessionRepo     ports.SessionRepository
	userSvc         ports.UserService
	partnerSvc      ports.PartnerService
	catalogSvc      ports.CatalogService
	handler         partnerHandler.Handler
	testServerURL   string           // Global variable to hold the URL of the running test server
	testServer      *http.Server     // To allow graceful shutdown
	testMQConn      *amqp.Connection // RabbitMQ connection for test verification
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

	ctx = context.WithValue(ctx, ctxutil.LoggerKey, testLogger)

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

	// Setup RabbitMQ
	rabbit, err := tu.SetupRabbitMQ(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to setup RabbitMQ: %v", err)
	}
	defer tu.TeardownRabbitMQ(ctx, nil, rabbit)

	ch, conn, err := rabbit.NewChannel()
	if err != nil {
		log.Fatalf("Failed to create channel: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	// Store connection for test verification
	testMQConn = conn

	// Setup RabbitMQ exchanges and queues needed for catalog
	if err := setupCatalogQueues(ch); err != nil {
		log.Fatalf("Failed to setup catalog queues: %v", err)
	}

	// stripe container
	stripeContainer, err := tu.SetupStripeMock(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to setup stripe container: %v", err)
	}
	defer tu.TeardownStripeMock(ctx, nil, stripeContainer)

	// Setup enhanced Vault testcontainer with per-service encryption (GDPR compliant)
	log.Println("Setting up enhanced Vault testcontainer with per-service keys...")
	serviceNames := []string{services.AuthUser} // Only authuser service for this test
	vaultSetup, err = tu.SetupServiceVault(ctx, nil, serviceNames)
	if err != nil {
		log.Fatalf("Failed to setup service Vault container: %v", err)
	}
	defer tu.TeardownVault(ctx, nil, vaultSetup.VaultContainer)

	// Set environment variables for Vault (for encx library)
	os.Setenv("VAULT_ADDR", vaultSetup.VaultContainer.HTTPSEndpoint)
	os.Setenv("VAULT_TOKEN", vaultSetup.VaultContainer.RootToken)

	// Get service-specific crypto service (GDPR compliant)
	var exists bool
	crypto, exists = vaultSetup.GetServiceCrypto(services.AuthUser)
	if !exists {
		log.Fatal("AuthUser service crypto not found in vault setup")
	}
	if crypto == nil {
		log.Fatal("AuthUser crypto service is nil")
	}
	log.Printf("✓ AuthUser service crypto initialized with per-service encryption key")

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

	log.Printf("✓ AuthTestContext initialized for user authentication testing")

	payment := authPayment.NewService("sk_test_123456789012345678901234", stripeContainer.URL)

	// Initialize application layers
	userRepo = userRepository.New(ctx, testPool)
	userSvc = user.New(userRepo, crypto, payment)

	partnerRepo = partnerRepository.New(ctx, testPool)

	// Initialize catalog service with cache
	catalogCache = catalog.NewCatalogCache()
	catalogSvc, err = catalog.New(ctx, catalogCache, testMQConn)
	if err != nil {
		log.Fatalf("Failed to create catalog service: %v", err)
	}

	// Initialize partner service (includes catalog consumer)
	partnerSvc, err = partner.New(ctx, partnerRepo, userRepo, catalogSvc, testMQConn, crypto, payment)
	if err != nil {
		log.Fatalf("Failed to create partner service: %v", err)
	}

	sessionRepo = sessionRepository.New(redisClient)

	authSessionRepo = session.NewRedisSessionRepository(redisClient)
	authmw := auth.NewSessionAuthMiddleware(authSessionRepo, crypto, nil)

	handler = partnerHandler.New(partnerSvc, authmw)

	// Set required environment variables for logger middleware
	os.Setenv("CLIENT_IP_HEADER", "X-Forwarded-For")
	os.Setenv("LOGGING_SALT", "test_logging_salt_12345")

	// HTTP server setup with logger middleware
	router := http.NewServeMux()
	handler.RegisterRoutes(router)

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

// setupCatalogQueues sets up the RabbitMQ exchanges and queues needed for catalog integration
func setupCatalogQueues(ch *amqp.Channel) error {
	log.Println("Setting up catalog queues...")

	// Declare catalog exchange
	if err := rabbitmq.DeclareExchange(ch, mq.CatalogExchangeName, "topic"); err != nil {
		return fmt.Errorf("declare catalog exchange: %w", err)
	}

	// Declare catalog queue for partner service
	if err := rabbitmq.DeclareQueue(ch, mq.CatalogQueueName); err != nil {
		return fmt.Errorf("declare catalog queue: %w", err)
	}

	// Bind catalog queue to exchange with all routing keys
	routingKeys := []string{
		mq.CategoryCreatedRoutingKey,
		mq.CategoryUpdatedRoutingKey,
		mq.CategoryDeletedRoutingKey,
		mq.ProductCreatedRoutingKey,
		mq.ProductUpdatedRoutingKey,
		mq.ProductDeletedRoutingKey,
	}

	for _, key := range routingKeys {
		if err := rabbitmq.BindQueue(ch, mq.CatalogQueueName, key, mq.CatalogExchangeName); err != nil {
			return fmt.Errorf("bind catalog queue with key %s: %w", key, err)
		}
	}

	log.Println("Catalog queues setup completed")
	return nil
}
