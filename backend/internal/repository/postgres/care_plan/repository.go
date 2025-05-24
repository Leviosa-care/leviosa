package carePlanRepository

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

type repository struct {
	DB *sql.DB
}

func (r *repository) GetDB() *sql.DB {
	return r.DB
}

func New(ctx context.Context, db *sql.DB) (*repository, error) {
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("pgx"); err != nil {
		return nil, fmt.Errorf("setting dialect for care plan repository: %w", err)
	}
	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return nil, fmt.Errorf("running all migrations for care plan repository: %w", err)
	}

	return &repository{db}, nil
}
