package priceRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// GetAllPrices retrieves ALL prices from the database.
func (r *PriceRepository) GetAllPrices(ctx context.Context) ([]*domain.Price, error) {
	query := fmt.Sprintf(`
		SELECT 
			id,
			product_id,
			stripe_price_id,
			amount,
			currency,
			interval,
			is_active,
			created_at,
			updated_at
		FROM 
			%s.prices
		ORDER BY 
			created_at ASC
	`, r.schema)

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errs.ClassifyPgError("get all prices", err)
	}
	defer rows.Close()

	var prices []*domain.Price

	for rows.Next() {
		var p domain.Price
		var intervalStr string

		if err := rows.Scan(
			&p.ID,
			&p.ProductID,
			&p.StripePriceID,
			&p.Amount,
			&p.Currency,
			&intervalStr,
			&p.IsActive,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, errs.NewDBQueryErr(fmt.Errorf("failed to scan row: %w", err))
		}

		p.Interval = domain.Interval(intervalStr)

		prices = append(prices, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.NewDBQueryErr(fmt.Errorf("error during rows iteration: %w", err))
	}

	if len(prices) == 0 {
		return []*domain.Price{}, nil
	}

	return prices, nil
}
