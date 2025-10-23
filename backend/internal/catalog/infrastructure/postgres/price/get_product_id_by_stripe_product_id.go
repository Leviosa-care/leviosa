package priceRepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *PriceRepository) GetProductIDByStripeProductID(ctx context.Context, stripeProductID string) (string, error) {
	var productID string
	query := `SELECT id FROM catalog.products WHERE stripe_product_id = $1;`

	err := r.pool.QueryRow(ctx, query, stripeProductID).Scan(&productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.NewRepositoryNotFoundErr(nil, "product by Stripe ID") // Product not found by Stripe ID
		}
		return "", errs.ClassifyPgError("get product ID by Stripe product ID", err)
	}

	return productID, nil
}
