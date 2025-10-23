package couponRepository

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type CouponRepository struct {
	pool   *pgxpool.Pool
	schema string
}

func New(ctx context.Context, pool *pgxpool.Pool) ports.CouponRepository {
	return &CouponRepository{pool: pool, schema: "catalog"}
}