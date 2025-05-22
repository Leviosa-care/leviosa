package testdatabase

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/pressly/goose/v3"
)

func NewSQLite(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory")
	if err != nil {
		return nil, fmt.Errorf("create new in memory database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping in memory: %w", err)
	}
	return db, nil
}

func MigrateSQLite(ctx context.Context, db *sql.DB, version int64) error {
	// setup goose
	goose.SetBaseFS(nil)
	// Set the dialect to SQLite3
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("Failed to set dialect: %w", err)
	}
	if err := goose.UpToContext(ctx, db, os.Getenv("TEST_MIGRATION_PATH"), version); err != nil {
		return fmt.Errorf("test migration up failed with version %d: %w", version, err)
	}
	return nil
}
