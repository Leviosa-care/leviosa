package productRepository

import (
	"context"
	"errors"
	"fmt"

	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres"
)

func (p *repository) RemoveProduct(ctx context.Context, productID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1;", pg.QualifiedTable(p.schema, "products"))
	result, err := p.DB.ExecContext(ctx, query, productID)
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
		return rp.NewNotDeletedErr(err, "product")
	}
	return nil
}
