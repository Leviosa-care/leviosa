package metricsRepository

import (
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository implements ports.MetricsRepository
type Repository struct {
	pool *pgxpool.Pool
}

// New creates a new metrics repository
func New(pool *pgxpool.Pool) ports.MetricsRepository {
	return &Repository{
		pool: pool,
	}
}
