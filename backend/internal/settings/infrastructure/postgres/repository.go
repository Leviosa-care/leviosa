package postgres

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/settings/ports"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type repository struct {
	pool   *pgxpool.Pool
	schema string
}

func New(ctx context.Context, pool *pgxpool.Pool) ports.SettingsRepository {
	return &repository{pool: pool, schema: "settings"}
}
