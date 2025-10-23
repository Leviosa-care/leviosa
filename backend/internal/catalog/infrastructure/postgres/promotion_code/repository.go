package promotionCodeRepository

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PromotionCodeRepository struct {
	pool   *pgxpool.Pool
	schema string
}

func New(ctx context.Context, pool *pgxpool.Pool) ports.PromotionCodeRepository {
	return &PromotionCodeRepository{pool: pool, schema: "catalog"}
}