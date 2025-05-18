package productRepository

import (
	"context"
	"errors"

	rp "github.com/hengadev/leviosa/internal/repository"
)

// AddPriceID updates the price_id field for a specific product in the database.
// It takes a context for cancellation/timeout, a productID to identify the product,
// and a priceID to be set. Returns an error if the update fails, the product
// doesn't exist, or if there are any database connectivity issues.
func (p *repository) AddPriceID(ctx context.Context, productID, priceID string) error {
	query := `
        UPDATE products
        SET price_id = ?
        WHERE id = ?;
    `
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
