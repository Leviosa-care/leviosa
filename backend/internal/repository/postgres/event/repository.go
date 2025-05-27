package eventRepository

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type repository struct {
	DB *sql.DB
}

func (e *repository) GetDB() *sql.DB {
	return e.DB
}

func New(ctx context.Context, db *sql.DB) (*repository, error) {

	return &repository{db}, nil
}
