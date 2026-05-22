package categoryRepository

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *CategoryRepository) AddCategory(ctx context.Context, category *domain.Category) (string, error) {
	query := `
	INSERT INTO catalog.categories (
		id,
		name,
		description,
		status,
		created_at
	) VALUES ($1, $2, $3, $4, $5)
	RETURNING id`
	var newID string
	err := r.pool.QueryRow(ctx, query,
		category.ID,
		category.Name,
		category.Description,
		category.Status,
		category.CreatedAt,
	).Scan(&newID)

	if err != nil {
		return "", errs.ClassifyPgError("insert category", err)
	}

	return newID, nil
}
