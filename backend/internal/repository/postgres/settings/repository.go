package settingsRepository

import (
	"context"
	"database/sql"

	"github.com/hengadev/leviosa/internal/domain/settings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type repository struct {
	DB     *sql.DB
	schema string
}

func (r *repository) GetDB() *sql.DB {
	return r.DB
}

func New(ctx context.Context, db *sql.DB) (settings.ReadWriter, error) {
	return &repository{DB: db, schema: "settings"}, nil
}
