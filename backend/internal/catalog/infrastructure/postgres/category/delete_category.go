package categoryRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *CategoryRepository) DeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	query := `
	DELETE FROM catalog.categories
	WHERE id = $1
	`

	commandTag, err := r.pool.Exec(ctx, query, categoryID)
	if err != nil {
		return errs.ClassifyPgError("delete category ", err)
	}

	if commandTag.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(fmt.Errorf("failed to delete category with ID %s", categoryID), "category")
	}
	return nil
}
