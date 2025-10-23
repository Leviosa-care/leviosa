package category

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *CategoryService) UpdateCategory(ctx context.Context, category *domain.UpdateCategoryRequest) error {
	if err := category.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr("category")
	}

	err := s.repo.UpdateCategory(ctx, category)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			// Category with the given ID was not found by the repository.
			return errs.NewNotFoundErr(err, "category")
		case errors.Is(err, errs.ErrUniqueViolation):
			// Handle unique constraint violations (e.g., if category name or slug must be unique)
			return errs.NewConflictErr(err)
		case errors.Is(err, errs.ErrCheckViolation):
			// Handle check constraint violations
			return errs.NewInvalidValueErr(err.Error()) // Or a more specific errs error
		case errors.Is(err, errs.ErrDBQuery), errors.Is(err, errs.ErrDatabase):
			// General database query or connection issue.
			return errs.NewQueryFailedErr(fmt.Errorf("repository query failed to update category: %w", err))
		case errors.Is(err, errs.ErrContext):
			// Context timeout or cancellation during DB operation.
			return errs.NewUnexpectedError(fmt.Errorf("context error during category update: %w", err))
		case errors.Is(err, errs.ErrInvalidInput):
			// This would come from generateUpdateQuery if no fields are provided
			return errs.NewInvalidValueErr(err.Error()) // Or errs.NewConflictErr if you consider "no fields" a conflict
		default:
			// Catch any unhandled repository errors.
			return errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during category update: %w", err))
		}
	}

	return nil
}
