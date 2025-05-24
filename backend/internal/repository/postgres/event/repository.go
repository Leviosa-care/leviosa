package eventRepository

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

func (e *repository) GetDB() *sql.DB {
	return e.DB
}

func New(ctx context.Context, db *sql.DB) (*repository, error) {

	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("pgx"); err != nil {
		return nil, fmt.Errorf("setting dialect for event repository: %w", err)
	}
	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return nil, fmt.Errorf("running all migrations for event repository: %w", err)
	}
	return &repository{db}, nil
}
