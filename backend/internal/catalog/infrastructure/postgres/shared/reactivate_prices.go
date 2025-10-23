package sharedRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// ReactivatePrices updates a price's `is_active` status to TRUE in the database.
// This reactivates previously deactivated prices by their Stripe price IDs.
func (r *SharedRepository) ReactivatePrices(ctx context.Context, stripePriceIDs []string) error {
	if len(stripePriceIDs) == 0 {
		return nil // No IDs to reactivate, so no operation needed.
	}

	query := fmt.Sprintf(`
		UPDATE %s.prices
		SET is_active = TRUE
		WHERE stripe_price_id = ANY($1);
	`, r.schema)
	commandTag, err := r.pool.Exec(ctx, query, stripePriceIDs)
	if err != nil {
		return errs.ClassifyPgError("reactivate prices", err) // Classify the PostgreSQL error
	}

	if commandTag.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(nil, "prices to reactivate")
	}

	return nil
}