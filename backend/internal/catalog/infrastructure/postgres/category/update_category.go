package categoryRepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *CategoryRepository) UpdateCategory(ctx context.Context, category *domain.UpdateCategoryRequest) error {
	query, args, err := generateUpdateQuery(category.ID, category)
	if err != nil {
		// If generateUpdateQuery returns "no fields provided", handle it or let it propagate.
		// If it's a marshalling error, it's an invalid input error from the repo's perspective.
		if errors.Is(err, errs.ErrNoFieldsForUpdate) { // Check for the specific error string (or make it a custom error)
			return errs.NewInvalidInputErr(err)
		}
		return errs.NewInternalErr(fmt.Errorf("failed to generate update query: %w", err))
	}

	commandTag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return errs.ClassifyPgError("update category", err)
	}

	if commandTag.RowsAffected() == 0 {
		// If 0 rows were affected, it means the category with the given ID was not found.
		return errs.NewRepositoryNotFoundErr(errors.New(""), "category")
	}

	return nil
}
