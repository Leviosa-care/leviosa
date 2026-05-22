package product_test

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

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/application/aggregator"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/application/image"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/application/price"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/application/product"
	imageRepository "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/image"
	priceRepository "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/price"
	productRepository "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/product"
	sharedRepository "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/shared"
	imageMedia "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/s3/image"
	pricePayment "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/stripe/price"
	productPayment "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/stripe/product"
	productHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/product"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
	authsession "github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/services"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/envmode"
	"github.com/Leviosa-care/leviosa/backend/internal/common/logger"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
	"github.com/Leviosa-care/leviosa/backend/internal/common/migrations"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hengadev/encx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
)

var (
	pgContainer           *tu.PostgresContainer
	testPool              *pgxpool.Pool
	redisClient           *redis.Client
	crypto                encx.CryptoService
	s3Client              *s3.Client
	productPaymentGateway ports.ProductPaymentGateway
	pricePaymentGateway   ports.PricePaymentGateway
	repo                  ports.ProductRepository
	sharedRepo            ports.SharedRepository
	authCtx               *tu.AuthTestContext // Authentication context for user/session tests
	handler               productHandler.Handler
	testServerURL         string       // Global variable to hold the URL of the running test server
	testServer            *http.Server // To allow graceful shutdown
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error

	// Create and configure logger for tests
	loggerHandler, err := logger.SetHandler("debug", "dev")
	if err != nil {
		log.Fatalf("Failed to create logger handler: %v", err)
	}
	testLogger := slog.New(loggerHandler)
	slog.SetDefault(testLogger) // Set as default for the application

	ctx = context.WithValue(ctx, ctxutil.LoggerKey, testLogger)

	// postgres container
	pgContainer, err = tu.SetupPostgres(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to setup postgres container: %v", err)
	}
	defer tu.TeardownPostgres(ctx, nil, pgContainer)

	// DB
	log.Println("Creating pgxpool...")
	// Use a context with timeout for pool creation
	poolCtx, poolCancel := context.WithTimeout(ctx, 10*time.Second)
	defer poolCancel()
	// ParseConfig is useful for setting pool options from connection string
	pgCfg, err := pgxpool.ParseConfig(pgContainer.ConnectionString)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse pgxpool config: %v", err))
	}
	// Optional: pgCfgure pool settings for tests
	pgCfg.MaxConns = 5
	pgCfg.MinConns = 1
	pgCfg.MaxConnLifetime = 30 * time.Minute
	pgCfg.MaxConnIdleTime = 5 * time.Minute

	testPool, err = pgxpool.NewWithConfig(poolCtx, pgCfg) // Use NewWithConfig
	if err != nil {
		tu.TeardownPostgres(ctx, nil, pgContainer)
		panic(fmt.Sprintf("Failed to open test database pool: %v", err))
	}
	log.Println("pgxpool created.")

	// Ping the database to ensure connections are established
	if err = testPool.Ping(poolCtx); err != nil {
		panic(fmt.Sprintf("Failed to ping database pool: %v", err))
	}
	log.Println("Database pool ping successful.")

	// migrations for schema and table
	log.Println("Applying database migrations...")
	goose.SetBaseFS(migrations.FS)
	if err = goose.SetDialect("pgx"); err != nil {
		log.Fatalf("Setting dialect for migrations: %s\n", err)
	}

	gooseDB, err := sql.Open("pgx", testPool.Config().ConnString())
	if err != nil {
		panic(fmt.Sprintf("Failed to open temp *sql.DB for goose migrations: %v", err))
	}
	defer gooseDB.Close() // Close the temporary DB connection

	if err = goose.UpContext(ctx, gooseDB, "."); err != nil { // Use gooseDB for migrations
		panic(fmt.Sprintf("running all migrations: %s\n", err))
	}
	log.Println("Migrations applied.")

	log.Println("Setting S3 testcontainer...")
	// s3 container
	localstackContainer, err := tu.SetupLocalstack(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to setup S3 container: %v", err)
	}
	defer tu.TeardownLocalstack(ctx, nil, localstackContainer)
	log.Println("S3 testcontainer et.")

	// s3 config
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			"test", // Access Key ID
			"test", // Secret Access Key
			"",     // Session Token (empty for Localstack)
		)),
	)
	if err != nil {
		log.Fatalf("Load default S3 configuration: %s\n", err)
	}
	// s3 client
	s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(localstackContainer.S3Endpoint)
		o.UsePathStyle = true // Required for Localstack
		o.Region = "us-east-1"
	})
	// Create a test bucket in Localstack S3
	_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(td.BUCKETNAME), // Use your test bucket name
	})
	if err != nil {
		log.Fatalf("Failed to create test S3 bucket: %v", err)
	}
	log.Println("Test S3 bucket created.")

	// Redis container
	redisContainer, err := tu.SetupRedis(ctx, nil)
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
	serviceNames := []string{services.Catalog} // Only catalog service for this test
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
	crypto, exists = vaultSetup.GetServiceCrypto(services.Catalog)
	if !exists {
		log.Fatal("Catalog service crypto not found in vault setup")
	}
	if crypto == nil {
		log.Fatal("Catalog crypto service is nil")
	}
	log.Printf("✓ Catalog service crypto initialized with per-service encryption key")

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

	// stripe container
	stripeContainer, err := tu.SetupStripeMock(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to setup stripe container: %v", err)
	}
	defer tu.TeardownStripeMock(ctx, nil, stripeContainer)

	pricePaymentGateway = pricePayment.NewPrice("sk_test_123456789012345678901234", stripeContainer.URL)
	productPaymentGateway = productPayment.NewProduct("sk_test_123456789012345678901234", stripeContainer.URL)

	repo = productRepository.New(ctx, testPool)
	sharedRepo = sharedRepository.New(ctx, testPool)
	imageRepo := imageRepository.New(ctx, testPool)
	priceRepo := priceRepository.New(ctx, testPool)
	mediaRepo := imageMedia.New(ctx, s3Client, td.BUCKETNAME)

	productService := product.New(repo, sharedRepo, productPaymentGateway, pricePaymentGateway, nil)
	imageService := image.New(imageRepo, mediaRepo, sharedRepo)
	priceService := price.New(priceRepo, sharedRepo, pricePaymentGateway)

	productAggregator := aggregator.NewProductPricesAggregatorService(productService, priceService, imageService)

	authSessionRepo := authsession.NewRedisSessionRepository(redisClient)
	authmw := auth.NewSessionAuthMiddleware(authSessionRepo, crypto, nil)

	handler = productHandler.New(productService, productAggregator, authmw)

	router := http.NewServeMux()
	handler.RegisterRoutes(router)

	// Use the enhanced AttachLogger middleware from core package
	loggerMiddleware := middleware.AttachLogger(envmode.Dev, testLogger)

	listener, err := net.Listen("tcp", ":0") // ":0" tells OS to pick a random available port
	if err != nil {
		log.Fatalf("Failed to listen for test server: %v", err)
	}
	testServerURL = "http://" + listener.Addr().String()
	testServer = &http.Server{Handler: loggerMiddleware(router)}

	// Run the server in a goroutine
	go func() {
		if err := testServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Test server failed to serve: %v", err)
		}
	}()
	log.Printf("Test HTTP server started at %s", testServerURL)

	// Give the server a moment to start up fully
	time.Sleep(100 * time.Millisecond)

	// Run tests
	code := m.Run()

	log.Println("Shutting down test HTTP server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := testServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Test server shutdown failed: %v", err)
	}
	log.Println("Test HTTP server shut down.")

	// Exit with the test result code
	os.Exit(code) // Commented out to allow cleanup before exiting in some environments
}
