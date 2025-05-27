package productRepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain/product"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres"
)

func (p *repository) AddProduct(ctx context.Context, product *productService.Product) error {
	query := fmt.Sprintf(`
            INSERT INTO %s (
                id,
                name,
                description
            ) VALUES ($1, $2, $3)`, pg.QualifiedTable(p.schema, "products"))

	result, err := p.DB.ExecContext(ctx, query,
		product.ID,
		product.Name,
		product.Description,
	)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return rp.NewContextErr(err)
		default:
			return rp.NewDatabaseErr(err)
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return rp.NewDatabaseErr(err)
	}
	if rowsAffected == 0 {
		return rp.NewNotCreatedErr(err, "product")
	}
	return nil
}
