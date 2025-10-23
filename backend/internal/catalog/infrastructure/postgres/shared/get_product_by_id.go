package sharedRepository

// TODO:
// GetProductByID(ctx context.Context, productID uuid.UUID) (*domain.ProductRes, error)

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *SharedRepository) GetProductByID(ctx context.Context, productID uuid.UUID) (*domain.Product, error) {
	query := `
	SELECT
		id,
		name,
		description,
		category_id,
		duration,
		created_at,
		updated_at,
		status,
		availability,
		buffer_time,
		cancellation_hours,
		stripe_product_id,
		metadata
	FROM catalog.products
	WHERE id = $1
	`

	var (
		pr       domain.Product
		metaProd []byte
	)

	err := r.pool.QueryRow(ctx, query, productID).Scan(
		&pr.ID,
		&pr.Name,
		&pr.Description,
		&pr.CategoryID,
		&pr.Duration,
		&pr.CreatedAt,
		&pr.UpdatedAt,
		&pr.Status,
		&pr.Availability,
		&pr.BufferTime,
		&pr.CancellationHours,
		&pr.StripeProductID,
		&metaProd,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NewRepositoryNotFoundErr(err, "product")
		}
		return nil, errs.ClassifyPgError("get products by ID", err)
	}

	if metaProd != nil {
		if err := json.Unmarshal(metaProd, &pr.Metadata); err != nil {
			return nil, errs.NewInvalidInputErr(fmt.Errorf("failed to unmarshal product metadata: %w", err))
		}
	}

	return &pr, nil
}
