package categoryRepository

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type CategoryRepository struct {
	pool   *pgxpool.Pool
	schema string
}

func New(ctx context.Context, pool *pgxpool.Pool) ports.CategoryRepository {
	return &CategoryRepository{pool: pool, schema: "catalog"}
}
