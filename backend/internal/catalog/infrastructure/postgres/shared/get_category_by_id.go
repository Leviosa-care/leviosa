package sharedRepository

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

func (r *SharedRepository) GetCategoryByID(ctx context.Context, categoryID uuid.UUID) (*domain.Category, error) {
	query := `
	SELECT id, name, description, status, metadata, created_at, updated_at
	FROM catalog.categories
	WHERE id = $1`

	var (
		category     domain.Category
		metadataJSON []byte
	)

	err := r.pool.QueryRow(ctx, query, categoryID).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.Status,
		&metadataJSON,
		&category.CreatedAt,
		&category.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NewRepositoryNotFoundErr(err, "category")
		}
		return nil, errs.ClassifyPgError("get category by ID", err)
	}

	// Decode JSONB metadata
	if metadataJSON != nil {
		if err := json.Unmarshal(metadataJSON, &category.Metadata); err != nil {
			return nil, errs.NewInvalidInputErr(fmt.Errorf("failed to unmarshal category metadata: %w", err))
		}
	} else {
		category.Metadata = make(map[string]any)
	}

	return &category, nil
}
