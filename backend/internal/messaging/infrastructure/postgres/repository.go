package postgres

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/messaging/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool *pgxpool.Pool
}

// New creates a new PostgreSQL message repository
func New(ctx context.Context, pool *pgxpool.Pool) ports.MessageRepository {
	return &repository{pool: pool}
}
