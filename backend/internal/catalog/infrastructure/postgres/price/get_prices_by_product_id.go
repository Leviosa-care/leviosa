package priceRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// GetPricesByProductID retrieves a list of Prices for a given internal product ID.
func (r *PriceRepository) GetPricesByProductID(ctx context.Context, productID string, activeOnly bool) ([]*domain.Price, error) {
	baseQuery := `
		SELECT id, product_id, stripe_price_id, amount, currency, interval, is_active, created_at, updated_at
		FROM %s.prices
		WHERE product_id = $1
	`
	query := fmt.Sprintf(baseQuery, r.schema)
	args := []any{productID}

	if activeOnly {
		query += " AND is_active = TRUE"
	}
	query += " ORDER BY created_at DESC;" // Order by creation date, newest first

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, errs.ClassifyPgError("get all prices for product ID", err)
	}
	defer rows.Close()

	var prices []*domain.Price
	for rows.Next() {
		var p domain.Price
		if err := rows.Scan(
			&p.ID,
			&p.ProductID,
			&p.StripePriceID,
			&p.Amount,
			&p.Currency,
			&p.Interval,
			&p.IsActive,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, errs.NewDBQueryErr(fmt.Errorf("failed to scan price row for product %s: %w", productID, err))
		}
		prices = append(prices, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.NewDBQueryErr(fmt.Errorf("error iterating over price rows for product %s: %w", productID, err))
	}

	if len(prices) == 0 {
		return []*domain.Price{}, nil
	}

	return prices, nil
}
