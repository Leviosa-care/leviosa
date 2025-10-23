package priceRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// CreatePrice inserts a new Price into the database.
func (r *PriceRepository) CreatePrice(ctx context.Context, price *domain.Price) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.prices (id, product_id, stripe_price_id, amount, currency, interval, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW());
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		price.ID,
		price.ProductID,
		price.StripePriceID,
		price.Amount,
		price.Currency,
		price.Interval,
		price.IsActive,
	)
	if err != nil {
		return errs.ClassifyPgError("create price", err)
	}
	return nil
}
