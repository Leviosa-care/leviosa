package userRepository

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx"
)

type repository struct {
	DB *sql.DB
}

func (u *repository) GetDB() *sql.DB {
	return u.DB
}

func New(ctx context.Context, db *sql.DB) (*repository, error) {
	return &repository{db}, nil
}
