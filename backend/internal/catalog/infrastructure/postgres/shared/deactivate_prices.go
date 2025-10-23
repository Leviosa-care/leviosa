package sharedRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// DeactivatePrices updates a price's `is_active` status to FALSE in the database.
// This is an alternative to a generic `UpdatePrice` if deactivation is very common.
// For consistency, I'd suggest calling `UpdatePrice` with `Active: stripe.Bool(false)`.
// However, if you explicitly want a `DeactivatePrices` func in repo:
func (r *SharedRepository) DeactivatePrices(ctx context.Context, priceIDs []string) error {
	if len(priceIDs) == 0 {
		return nil // No IDs to deactivate, so no operation needed.
	}

	query := fmt.Sprintf(`
		UPDATE %s.prices
		SET is_active = FALSE
		WHERE stripe_price_id = ANY($1);
	`, r.schema)
	commandTag, err := r.pool.Exec(ctx, query, priceIDs)
	if err != nil {
		return errs.ClassifyPgError("deactivate prices", err) // Classify the PostgreSQL error
	}

	if commandTag.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(nil, "prices to deactivate")
	}

	return nil
}
