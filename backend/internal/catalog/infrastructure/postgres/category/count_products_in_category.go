package categoryRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *CategoryRepository) CountProductsInCategory(ctx context.Context, categoryID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM catalog.products WHERE category_id = $1;` // Adjust table/column names if different

	var count int
	err := r.pool.QueryRow(ctx, query, categoryID).Scan(&count)
	if err != nil {
		return 0, errs.NewDBQueryErr(fmt.Errorf("failed to count products in category %s: %w", categoryID, err))
	}
	return count, nil
}
