package helpers

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

	settingsHandler "github.com/Leviosa-care/settings/internal/adapters/http"
	"github.com/Leviosa-care/settings/internal/adapters/postgres"
	"github.com/Leviosa-care/settings/internal/adapters/rabbitmq"
	"github.com/Leviosa-care/settings/internal/adapters/s3"
	settings "github.com/Leviosa-care/settings/internal/application"
	"github.com/Leviosa-care/settings/internal/ports"
	th "github.com/Leviosa-care/settings/test/helpers"

	"github.com/Leviosa-care/core/migrations"
	tu "github.com/Leviosa-care/core/testutils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hengadev/encx"
	"github.com/hengadev/encx/providers/hashicorpvault"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	pgContainer   *tu.PostgresContainer
	testPool      *pgxpool.Pool
	s3Client      *s3.Client
	repo          ports.SettingsRepository
	mediaRepo     ports.SettingsMedia
	handler       settingsHandler.Handler
	testServerURL string           // Global variable to hold the URL of the running test server
	testServer    *http.Server     // To allow graceful shutdown
	testMQConn    *amqp.Connection // RabbitMQ connection for test verification
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
		Bucket: aws.String(th.BUCKETNAME), // Use your test bucket name
	})
	if err != nil {
		log.Fatalf("Failed to create test S3 bucket: %v", err)
	}
	log.Println("Test S3 bucket created.")

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

	// Setup Vault testcontainer
	log.Println("Setting up Vault testcontainer...")
	vaultContainer, err := tu.SetupVault(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to setup Vault container: %v", err)
	}
	defer tu.TeardownVault(ctx, nil, vaultContainer)

	// Set environment variables for Vault
	os.Setenv("VAULT_ADDR", vaultContainer.HTTPSEndpoint)
	os.Setenv("VAULT_TOKEN", vaultContainer.RootToken)

	// crypto
	kms, err := hashicorpvault.New()
	if err != nil {
		fmt.Println("creating vault:", err)
		return
	}
	crypto, err := encx.New(
		ctx,
		kms,
		tu.EncryptionKey,
		"secret/data/pepper",
	)
	if err != nil {
		log.Printf("Crypto service creation error details: %+v", err)
		log.Fatalf("Failed to create crypto service: %v", err)
	}
	if crypto == nil {
		log.Fatal("Crypto service is nil after creation")
	}
	log.Println("Crypto service created successfully")

	rabbitmq.Setup(ctx, ch)

	repo = postgres.New(ctx, testPool)
	mediaRepo = media.New(ctx, s3Client, th.BUCKETNAME)

	service := settings.New(repo, mediaRepo, crypto, conn)
	handler = settingsHandler.New(service)

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
