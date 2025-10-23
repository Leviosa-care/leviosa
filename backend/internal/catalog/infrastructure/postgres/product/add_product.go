package productRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *ProductRepository) AddProduct(ctx context.Context, p *domain.Product) (string, error) {
	metadataJSON, err := json.Marshal(p.Metadata)
	if err != nil {
		return "", errs.NewInvalidInputErr(fmt.Errorf("failed to encode product metadata: %w", err))
	}

	query := fmt.Sprintf(`
	INSERT INTO %s.products (
		id,
		name,
		description,
		category_id,
		duration,
		created_at,
		status,
		availability,
		buffer_time,
		cancellation_hours,
		stripe_product_id,
		metadata
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7,
		$8, $9, $10, $11, $12
	) RETURNING id`, r.schema)

	var newID string
	err = r.pool.QueryRow(ctx, query,
		p.ID,
		p.Name,
		p.Description,
		p.CategoryID,
		p.Duration,
		p.CreatedAt,
		p.Status,
		p.Availability,
		p.BufferTime,
		p.CancellationHours,
		p.StripeProductID,
		metadataJSON,
	).Scan(&newID)

	if err != nil {
		return "", errs.ClassifyPgError("insert product", err)
	}

	return newID, nil
}
