package promotionCode_test

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

	promotionCode "github.com/Leviosa-care/leviosa/backend/internal/catalog/application/promotion_code"
	couponRepository "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/coupon"
	promotionCodeRepository "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/promotion_code"
	promotionCodeHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/promotion_code"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"

	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/envmode"
	"github.com/Leviosa-care/leviosa/backend/internal/common/logger"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
	"github.com/Leviosa-care/leviosa/backend/internal/common/migrations"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

var (
	pgContainer   *tu.PostgresContainer
	testPool      *pgxpool.Pool
	promotionRepo ports.PromotionCodeRepository
	couponRepo    ports.CouponRepository
	handler       promotionCodeHandler.Handler
	testServerURL string       // Global variable to hold the URL of the running test server
	testServer    *http.Server // To allow graceful shutdown
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
	// Optional: Configure pool settings for tests
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

	// Initialize repositories
	promotionRepo = promotionCodeRepository.New(ctx, testPool)
	couponRepo = couponRepository.New(ctx, testPool)

	// Initialize services
	promotionCodeService := promotionCode.New(promotionRepo, couponRepo)

	// Initialize handlers
	handler = promotionCodeHandler.New(promotionCodeService)

	// Set up HTTP router
	router := http.NewServeMux()
	handler.RegisterRoutes(router)

	// Use the enhanced AttachLogger middleware from core package
	loggerMiddleware := middleware.AttachLogger(envmode.Dev, testLogger)

	// Start test server
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
	os.Exit(code)
}
