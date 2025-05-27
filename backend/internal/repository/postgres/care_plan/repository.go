package carePlanRepository

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type repository struct {
	DB     *sql.DB
	schema string
}

func (r *repository) GetDB() *sql.DB {
	return r.DB
}

func New(ctx context.Context, db *sql.DB) (*repository, error) {
	return &repository{DB: db, schema: "care_plan"}, nil
}
