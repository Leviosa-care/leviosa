package price_test

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

	priceHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/price"
	priceRepository "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/price"
	sharedRepository "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/shared"
	pricePayment "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/stripe/price"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/application/price"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/migrations"

	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

var (
	pgContainer         *tu.PostgresContainer
	testPool            *pgxpool.Pool
	pricePaymentGateway ports.PricePaymentGateway
	repo                ports.PriceRepository
	sharedRepo          ports.SharedRepository
	handler             priceHandler.Handler
	testServerURL       string       // Global variable to hold the URL of the running test server
	testServer          *http.Server // To allow graceful shutdown
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
	// Optional: configure pool settings for tests
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

	// stripe container
	stripeContainer, err := tu.SetupStripeMock(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to setup stripe container: %v", err)
	}
	defer tu.TeardownStripeMock(ctx, nil, stripeContainer)

	pricePaymentGateway = pricePayment.NewPrice("sk_test_123456789012345678901234", stripeContainer.URL)

	repo = priceRepository.New(ctx, testPool)
	sharedRepo = sharedRepository.New(ctx, testPool)

	priceService := price.New(repo, sharedRepo, pricePaymentGateway)

	handler = priceHandler.New(priceService)

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
	os.Exit(code)
}

