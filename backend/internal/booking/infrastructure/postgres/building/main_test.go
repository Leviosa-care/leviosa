package buildingRepository_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	buildingRepository "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/postgres/building"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"

	"github.com/Leviosa-care/leviosa/backend/internal/common/migrations"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

var (
	pgContainer *tu.PostgresContainer
	testPool    *pgxpool.Pool
	repo        ports.BuildingRepository
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error

	// Postgres container
	log.Println("Creating postgres container.")
	pgContainer, err = tu.SetupPostgres(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to setup postgres container: %v", err))
	}
	defer tu.TeardownPostgres(ctx, nil, pgContainer)
	log.Println("Postgres container successfully created.")

	// DB
	log.Println("Creating pgxpool...")
	// Use a context with timeout for pool creation
	poolCtx, poolCancel := context.WithTimeout(ctx, 10*time.Second)
	defer poolCancel()
	// ParseConfig is useful for setting pool options from connection string
	config, err := pgxpool.ParseConfig(pgContainer.ConnectionString)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse pgxpool config: %v", err))
	}
	// Optional: Configure pool settings for tests
	config.MaxConns = 5
	config.MinConns = 1
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	testPool, err = pgxpool.NewWithConfig(poolCtx, config) // Use NewWithConfig
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

	repo = buildingRepository.New(ctx, testPool)
	if err != nil {
		panic(fmt.Sprintf("Failed to create repository: %v", err))
	}

	// Run tests
	code := m.Run()

	log.Println("Test(s) executed")
	// Exit with the test result code
	os.Exit(code) // Commented out to allow cleanup before exiting in some environments
	// The `m.Run()` call handles exiting with the correct code in `go test`
}
