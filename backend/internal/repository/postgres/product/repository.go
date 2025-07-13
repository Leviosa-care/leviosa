package productRepository

import (
	"context"
	"database/sql"

	"github.com/hengadev/leviosa/internal/domain/product"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type repository struct {
	DB     *sql.DB
	schema string
}

func (u *repository) GetDB() *sql.DB {
	return u.DB
}

func New(ctx context.Context, db *sql.DB) (productService.ReadWriter, error) {
	return &repository{DB: db, schema: "products"}, nil
}
