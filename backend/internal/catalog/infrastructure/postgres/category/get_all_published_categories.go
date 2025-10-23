package categoryRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *CategoryRepository) GetAllPublishedCategories(ctx context.Context) ([]*domain.Category, error) {
	query := `
	SELECT id, name, description, status, metadata, created_at, updated_at
	FROM catalog.categories
	WHERE status = $1
	ORDER BY name ASC`

	rows, err := r.pool.Query(ctx, query, domain.Published)
	if err != nil {
		return nil, errs.ClassifyPgError("get all categories", err)
	}
	defer rows.Close()

	var categories []*domain.Category

	for rows.Next() {
		var (
			cat          domain.Category
			metadataJSON []byte
		)

		if err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Description,
			&cat.Status,
			&metadataJSON,
			&cat.CreatedAt,
			&cat.UpdatedAt,
		); err != nil {
			return nil, errs.NewDBQueryErr(fmt.Errorf("failed to scan row: %w", err))
		}

		if metadataJSON != nil {
			if err := json.Unmarshal(metadataJSON, &cat.Metadata); err != nil {
				return nil, errs.NewInvalidInputErr(fmt.Errorf("failed to unmarshal category metadata: %w", err))
			}
		}

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
