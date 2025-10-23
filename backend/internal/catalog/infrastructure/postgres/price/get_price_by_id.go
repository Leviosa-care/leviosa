package priceRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// GetPrice retrieves a single Price by its internal ID.
func (r *PriceRepository) GetPrice(ctx context.Context, priceID string) (*domain.Price, error) {
	query := fmt.Sprintf(`
		SELECT id, product_id, stripe_price_id, amount, currency, interval, is_active, created_at, updated_at
		FROM %s.prices
		WHERE id = $1;
	`, r.schema)
	var p domain.Price
	err := r.pool.QueryRow(ctx, query, priceID).Scan(
		&p.ID,
		&p.ProductID,
		&p.StripePriceID,
		&p.Amount,
		&p.Currency,
		&p.Interval,
		&p.IsActive,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewRepositoryNotFoundErr(nil, "price")
		}
		return nil, errs.ClassifyPgError("get price by ID", err)
	}
	return &p, nil
}
