package testdatabase

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"time"

	_ "github.com/jackc/pgx"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Postgres struct {
	container testcontainers.Container
	DB        *sql.DB
	connStr   string
}

func NewPostgres(ctx context.Context) (*Postgres, error) {
	// Create PostgreSQL container
	pgContainer, err := postgres.Run(ctx,
		"postgres:17.5-alpine3.21",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}
	// Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Verify connection
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Postgres{
		container: pgContainer,
		DB:        db,
		connStr:   connStr,
	}, nil
}

// Cleanup closes the database connection and terminates the container
func (tdb *Postgres) CleanupPostgres(ctx context.Context) error {
	if tdb.DB != nil {
		tdb.DB.Close()
	}
	if tdb.container != nil {
		return tdb.container.Terminate(ctx)
	}
	return nil
}

func (tdb *Postgres) PostgresUp(ctx context.Context, migrations embed.FS, version int64) error {
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("pgx"); err != nil {
		return fmt.Errorf("setting dialect for postgres testcontainer: %s\n", err)
	}
	if version == 0 {
		if err := goose.UpContext(ctx, tdb.DB, "."); err != nil {
			return fmt.Errorf("running all migrations for postgres testcontainer: %s\n", err)
		}
	}
	if err := goose.UpToContext(ctx, tdb.DB, ".", version); err != nil {
		return fmt.Errorf("running migrations to version %d for postgres testcontainer: %s\n", version, err)
	}
	return nil
}

func (tdb *Postgres) PostgresDown(ctx context.Context, migrations embed.FS, version int64) error {
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("pgx"); err != nil {
		return fmt.Errorf("setting dialect for postgres testcontainer: %s\n", err)
	}
	if version == 0 {
		if err := goose.DownContext(ctx, tdb.DB, "."); err != nil {
			return fmt.Errorf("running down all migrations for postgres testcontainer: %s\n", err)
		}
	}
	if err := goose.DownToContext(ctx, tdb.DB, ".", version); err != nil {
		return fmt.Errorf("running down migrations to version %d for postgres testcontainer: %s\n", version, err)
	}
	return nil
}
