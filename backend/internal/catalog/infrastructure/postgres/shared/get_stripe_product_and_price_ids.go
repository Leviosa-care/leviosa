package sharedRepository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *SharedRepository) GetStripeProductAndPriceIDs(ctx context.Context, productID uuid.UUID) (string, []string, error) {
	query := `
	SELECT
		p.stripe_product_id,
		pr.stripe_price_id
	FROM catalog.products p
	LEFT JOIN catalog.prices pr ON p.id = pr.product_id
	WHERE p.id = $1;`

	var stripeProductID sql.NullString
	var stripePriceIDs []string

	rows, err := r.pool.Query(ctx, query, productID)
	if err != nil {
		return "", nil, errs.ClassifyPgError(fmt.Sprintf("get Stripe IDs for product %s", productID), err)
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		found = true
		var priceID sql.NullString
		if err := rows.Scan(&stripeProductID, &priceID); err != nil {
			return "", nil, errs.NewDBQueryErr(fmt.Errorf("failed to scan Stripe IDs for product %s: %w", productID, err))
		}
		if priceID.Valid {
			stripePriceIDs = append(stripePriceIDs, priceID.String)
		}
	}

	if err := rows.Err(); err != nil {
		return "", nil, errs.NewDBQueryErr(fmt.Errorf("error iterating over Stripe IDs for product %s: %w", productID, err))
	}

	if !found {
		// No rows were returned, which means the product was not found.
		return "", nil, errs.NewRepositoryNotFoundErr(nil, fmt.Sprintf("Stripe IDs for product %q", productID))
	}

	// Check if stripe_product_id is valid (could be NULL if the product has been created without it).
	if !stripeProductID.Valid {
		return "", nil, errs.NewRepositoryNotFoundErr(nil, fmt.Sprintf("stripe product ID for product %q", productID))
	}

	return stripeProductID.String, stripePriceIDs, nil
}
