package productRepository

import (
	"context"
	"errors"
	"fmt"

	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres"
)

// AddPriceID updates the price_id field for a specific product in the database.
// It takes a context for cancellation/timeout, a productID to identify the product,
// and a priceID to be set. Returns an error if the update fails, the product
// doesn't exist, or if there are any database connectivity issues.
func (p *repository) AddPriceID(ctx context.Context, productID, priceID string) error {
	query := fmt.Sprintf(`
        UPDATE %s
        SET price_id = $1
        WHERE id = $2;
    `, pg.QualifiedTable(p.schema, "products"))
	result, err := p.DB.ExecContext(ctx, query,
		priceID,
		productID,
	)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return rp.NewContextErr(err)
		default:
			return rp.NewDatabaseErr(err)
		}
	}

	// Check if the insert was successful
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return rp.NewDatabaseErr(err)
	}
	if rowsAffected == 0 {
		return rp.NewNotUpdatedErr(err, "product")
	}
	return nil
}
