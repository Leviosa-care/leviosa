package sharedRepository

import (
	"context"
	"errors"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *SharedRepository) GetCategoryByID(ctx context.Context, categoryID uuid.UUID) (*domain.Category, error) {
	query := `
	SELECT id, name, description, status, created_at, updated_at
	FROM catalog.categories
	WHERE id = $1`

	var category domain.Category

	err := r.pool.QueryRow(ctx, query, categoryID).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.Status,
		&category.CreatedAt,
		&category.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NewRepositoryNotFoundErr(err, "category")
		}
		return nil, errs.ClassifyPgError("get category by ID", err)
	}

	return &category, nil
}
