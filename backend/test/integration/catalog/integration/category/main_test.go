package category_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	categoryHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/category"
	categoryRepository "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/category"
	imageRepository "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/image"
	sharedRepository "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/shared"
	imageMedia "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/s3/image"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/application/aggregator"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/application/category"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/application/image"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"
	"github.com/Leviosa-care/leviosa/backend/internal/common/migrations"

	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

var (
	pgContainer   *tu.PostgresContainer
	testPool      *pgxpool.Pool
	s3Client      *s3.Client
	repo          ports.CategoryRepository
	sharedRepo    ports.SharedRepository
	imageRepo     ports.ImageRepository
	mediaRepo     ports.ImageMedia
	handler       categoryHandler.Handler
	testServerURL string       // Global variable to hold the URL of the running test server
	testServer    *http.Server // To allow graceful shutdown
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error

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

	repo = categoryRepository.New(ctx, testPool)
	sharedRepo = sharedRepository.New(ctx, testPool)
	mediaRepo = imageMedia.New(ctx, s3Client, td.BUCKETNAME)
	imageRepo = imageRepository.New(ctx, testPool)

	categoryService := category.New(repo, sharedRepo)
	imageService := image.New(imageRepo, mediaRepo, sharedRepo)
	categoryImagesService := aggregator.NewCategoryAggregatorService(categoryService, imageService)
	handler = categoryHandler.New(categoryService, imageService, categoryImagesService)

	router := http.NewServeMux()
	handler.RegisterRoutes(router)

	listener, err := net.Listen("tcp", ":0") // ":0" tells OS to pick a random available port
	if err != nil {
		log.Fatalf("Failed to listen for test server: %v", err)
	}
	testServerURL = "http://" + listener.Addr().String()
	testServer = &http.Server{Handler: router} // Store the server for graceful shutdown

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
