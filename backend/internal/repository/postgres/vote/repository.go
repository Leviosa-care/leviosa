package voteRepository

import (
	"context"
	"database/sql"

	"github.com/hengadev/leviosa/internal/domain/vote"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type repository struct {
	DB     *sql.DB
	schema string
}

func (v *repository) GetDB() *sql.DB {
	return v.DB
}

func New(ctx context.Context, db *sql.DB) (vote.ReadWriter, error) {
	return &repository{DB: db, schema: "votes"}, nil
}
