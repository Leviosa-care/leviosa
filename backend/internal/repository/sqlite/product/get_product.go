package productRepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hengadev/leviosa/internal/domain/product"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (p *repository) GetProduct(ctx context.Context, productID string) (*productService.Product, error) {
	var product productService.Product
	query := `
        SELECT 
            name,
            description
        FROM products
        WHERE id = ?;`
	err := p.DB.QueryRowContext(ctx, query, productID).Scan(
		&product.Name,
		&product.Description,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, rp.NewNotFoundErr(err, "product")
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return nil, rp.NewContextErr(err)
		default:
			return nil, rp.NewDatabaseErr(err)
		}
	}
	product.ID = productID
	return &product, nil
}
