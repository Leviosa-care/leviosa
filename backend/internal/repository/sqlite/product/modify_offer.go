package productRepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain/product"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/pkg/sqliteutil"
)

func (p *Repository) ModifyOffer(
	ctx context.Context,
	productType *productService.Offer,
	whereMap map[string]any,
) error {
	query, values, err := sqliteutil.WriteUpdateQuery(*productType, whereMap)
	if err != nil {
		return rp.NewInternalErr(err)
	}
	result, err := p.DB.ExecContext(ctx, query, values...)
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
		return rp.NewNotUpdatedErr(err, fmt.Sprintf("product named %s", productType.Name))
	}
	return nil
}
