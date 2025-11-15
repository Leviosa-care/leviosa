package availabilityRepository

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	pool   *pgxpool.Pool
	schema string
}

func New(ctx context.Context, pool *pgxpool.Pool) ports.AvailabilityRepository {
	return &Repository{pool: pool, schema: "booking"}
}
