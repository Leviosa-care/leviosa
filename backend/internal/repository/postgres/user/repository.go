package userRepository

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

type repository struct {
	DB *sql.DB
}

func (u *repository) GetDB() *sql.DB {
	return u.DB
}

func New(ctx context.Context, db *sql.DB) (*repository, error) {
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("pgx"); err != nil {
		return nil, fmt.Errorf("setting dialect for user repository: %w", err)
	}
	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return nil, fmt.Errorf("running all migrations for user repository: %w", err)
	}
	return &repository{db}, nil
}
