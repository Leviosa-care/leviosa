package categoryRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *CategoryRepository) GetAllCategories(ctx context.Context) ([]*domain.Category, error) {
	query := `
		SELECT id, name, description, status, created_at, updated_at
		FROM catalog.categories
		ORDER BY name ASC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errs.ClassifyPgError("get all categories", err)
	}
	defer rows.Close()

	var categories []*domain.Category

	for rows.Next() {
		var cat domain.Category
		var statusStr string

		if err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Description,
			&statusStr,
			&cat.CreatedAt,
			&cat.UpdatedAt,
		); err != nil {
			return nil, errs.NewDBQueryErr(fmt.Errorf("failed to scan row: %w", err))
		}

		cat.Status = domain.PublishedStatus(statusStr)

		categories = append(categories, &cat)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.NewDBQueryErr(fmt.Errorf("error during rows iteration: %w", err))
	}

	if len(categories) == 0 {
		return []*domain.Category{}, nil
	}

	return categories, nil
}
