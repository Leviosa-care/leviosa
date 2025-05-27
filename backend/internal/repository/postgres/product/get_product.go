package productRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain/product"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres"
)

func (p *repository) GetProduct(ctx context.Context, productID string) (*productService.Product, error) {
	var product productService.Product
	query := fmt.Sprintf(`
        SELECT 
            name,
            description
        FROM %s
        WHERE id = $1;`, pg.QualifiedTable(p.schema, "products"))
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
