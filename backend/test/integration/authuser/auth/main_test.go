package auth_test

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

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/application/aggregator"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/application/catalog"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/application/otp"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/application/partner"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/application/session"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/application/user"
	partnerRepository "github.com/Leviosa-care/leviosa/backend/internal/authuser/infrastructure/postgres/partner"
	userRepository "github.com/Leviosa-care/leviosa/backend/internal/authuser/infrastructure/postgres/user"
	authRabbitMQ "github.com/Leviosa-care/leviosa/backend/internal/authuser/infrastructure/rabbitmq"
	otpRepository "github.com/Leviosa-care/leviosa/backend/internal/authuser/infrastructure/redis/otp"
	sessionRepository "github.com/Leviosa-care/leviosa/backend/internal/authuser/infrastructure/redis/session"
	authPayment "github.com/Leviosa-care/leviosa/backend/internal/authuser/infrastructure/stripe"
	aggregatorHandler "github.com/Leviosa-care/leviosa/backend/internal/authuser/interface/auth"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"

	// Catalog services for partner validation
	catalogApp "github.com/Leviosa-care/leviosa/backend/internal/catalog/application/category"
	catalogProductApp "github.com/Leviosa-care/leviosa/backend/internal/catalog/application/product"
	categoryRepository "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/category"
	productRepository "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/product"
	sharedRepository "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/shared"
	pricePayment "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/stripe/price"
	productPayment "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/stripe/product"
	catalogPorts "github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"

	authsession "github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	mq "github.com/Leviosa-care/leviosa/backend/internal/common/contracts/rabbitmq"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/services"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/envmode"
	"github.com/Leviosa-care/leviosa/backend/internal/common/logger"
	"github.com/Leviosa-care/leviosa/backend/internal/common/messaging/rabbitmq"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
	"github.com/Leviosa-care/leviosa/backend/internal/common/migrations"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"
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
	catalogCache    *catalog.CatalogCache
	otpRepo         ports.OTPRepository
	sessionRepo     ports.SessionRepository
	vaultSetup      *tu.ServiceVaultSetup // Enhanced Vault setup with per-service keys
	authCtx         *tu.AuthTestContext   // Authentication context for user/session tests
	authSessionRepo authsession.SessionRepository
	service         ports.AuthAggregatorService
	authHandler     aggregatorHandler.Handler
	testServerURL   string                                // Global variable to hold the URL of the running test server
	testServer      *http.Server                          // To allow graceful shutdown
	testMQConn      *amqp.Connection                      // RabbitMQ connection for test verification
	testNotifier    *td.MockNotificationService            // Mock notification service for OTP verification

	// Catalog services for partner validation
	categoryService catalogPorts.PublicCategoryService
	productService  catalogPorts.PublicProductService
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

	// Setup
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

	// stripe container
	stripeContainer, err := tu.SetupStripeMock(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to setup stripe container: %v", err)
	}
	defer tu.TeardownStripeMock(ctx, nil, stripeContainer)

	// Setup enhanced Vault testcontainer with per-service encryption (GDPR compliant)
	log.Println("Setting up enhanced Vault testcontainer with per-service keys...")
	serviceNames := []string{services.AuthUser} // Only settings service for this test
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

	payment := authPayment.NewService("sk_test_123456789012345678901234", stripeContainer.URL, "")

	// Setup RabbitMQ exchanges and queues needed for settings and OTP notifications
	if err := setupAllRabbitMQQueues(ctx, ch); err != nil {
		log.Fatalf("Failed to setup RabbitMQ queues: %v", err)
	}

	// Initialize application layers
	userRepo = userRepository.New(ctx, testPool)
	userService := user.New(userRepo, crypto, payment)

	catalogCache = catalog.NewCatalogCache()
	// Note: Old catalog service creation removed as it's no longer needed by partner service

	// Initialize catalog repositories and services for partner validation
	sharedRepo := sharedRepository.New(ctx, testPool)
	categoryRepo := categoryRepository.New(ctx, testPool)
	productRepo := productRepository.New(ctx, testPool)

	// Create Stripe services for catalog (needed by ProductService)
	catalogStripeService := productPayment.NewProduct("sk_test_123456789012345678901234", stripeContainer.URL)
	catalogPriceStripeService := pricePayment.NewPrice("sk_test_123456789012345678901234", stripeContainer.URL)

	// Initialize catalog services
	categoryService = catalogApp.New(categoryRepo, sharedRepo)
	productService = catalogProductApp.New(productRepo, sharedRepo, catalogStripeService, catalogPriceStripeService, nil)

	otpRepo = otpRepository.New(redisClient)
	testNotifier = td.NewMockNotificationService()
	otpService, err := otp.New(ctx, otpRepo, crypto, testNotifier)
	if err != nil {
		log.Fatalf("Failed to create OTP service: %v", err)
	}

	sessionRepo = sessionRepository.New(redisClient)
	sessionService := session.New(ctx, sessionRepo, crypto)

	// Create partner repository and service with new catalog services
	partnerRepo := partnerRepository.New(ctx, testPool)
	partnerService, err := partner.New(ctx, partnerRepo, userRepo, productService, categoryService, crypto, payment)
	if err != nil {
		log.Fatalf("Failed to create partner service: %v", err)
	}

	service = aggregator.New(otpService, userService, sessionService, partnerService, nil)

	authSessionRepo = authsession.NewRedisSessionRepository(redisClient)
	authmw := auth.NewSessionAuthMiddleware(authSessionRepo, crypto, nil)

	authHandler = aggregatorHandler.New(service, authmw)

	// Set required environment variables for logger middleware
	os.Setenv("CLIENT_IP_HEADER", "X-Forwarded-For")
	os.Setenv("LOGGING_SALT", "test_logging_salt_12345")

	// HTTP server setup with logger middleware
	router := http.NewServeMux()
	authHandler.RegisterRoutes(router)

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

// setupAllRabbitMQQueues sets up all RabbitMQ exchanges and queues needed for the auth service
func setupAllRabbitMQQueues(ctx context.Context, ch *amqp.Channel) error {
	// Setup authuser queues (catalog consumer)
	log.Println("Setting up authuser queues...")
	if err := authRabbitMQ.Setup(ctx, ch); err != nil {
		return fmt.Errorf("setup authuser queues: %w", err)
	}

	// Setup settings queues (for consuming settings updates)
	log.Println("Setting up settings queues...")
	if err := setupSettingsQueues(ch); err != nil {
		return fmt.Errorf("setup settings queues: %w", err)
	}

	log.Println("All RabbitMQ queues and exchanges setup completed")
	return nil
}

// setupSettingsQueues sets up the RabbitMQ exchanges and queues needed for settings consumption
func setupSettingsQueues(ch *amqp.Channel) error {
	// Declare settings exchange
	if err := rabbitmq.DeclareExchange(ch, mq.SettingsExchangeName, "direct"); err != nil {
		return fmt.Errorf("declare settings exchange: %w", err)
	}

	// Declare OTP settings queue (only the one we need for the OTP consumer)
	if err := rabbitmq.DeclareQueue(ch, mq.OTPSettingsQueueName); err != nil {
		return fmt.Errorf("declare OTP settings queue: %w", err)
	}

	// Bind OTP settings queue to exchange
	if err := rabbitmq.BindQueue(ch, mq.OTPSettingsQueueName, mq.SettingsRoutingKey, mq.SettingsExchangeName); err != nil {
		return fmt.Errorf("bind OTP settings queue: %w", err)
	}

	return nil
}
