package categoryRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *CategoryRepository) AddCategory(ctx context.Context, category *domain.Category) (string, error) {
	metadataJSON, err := json.Marshal(category.Metadata)
	if err != nil {
		return "", errs.NewInvalidInputErr(fmt.Errorf("failed to encode category metadata: %w", err))
	}

	query := `
	INSERT INTO catalog.categories (
		id,
		name,
		description,
		status,
		created_at,
		metadata
	) VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`
	var newID string
	err = r.pool.QueryRow(ctx, query,
		category.ID,
		category.Name,
		category.Description,
		category.Status,
		category.CreatedAt,
		metadataJSON,
	).Scan(&newID)

	if err != nil {
		return "", errs.ClassifyPgError("insert category", err)
	}

	return newID, nil
}
