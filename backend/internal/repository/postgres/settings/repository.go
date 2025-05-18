package settingsRepository

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type repository struct {
	DB *sql.DB
}

func (r *repository) GetDB() *sql.DB {
	return r.DB
}

func New(ctx context.Context, db *sql.DB) (*repository, error) {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("pgx"); err != nil {
		return nil, fmt.Errorf("setting dialect for settings repository: %w", err)
	}
	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return nil, fmt.Errorf("running all migrations for settings repository: %w", err)
	}
	return &repository{db}, nil
}
