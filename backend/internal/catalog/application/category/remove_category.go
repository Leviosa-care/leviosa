package category

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *CategoryService) RemoveCategory(ctx context.Context, categoryIDStr string) error {
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		return errs.NewInvalidValueErr("category ID is required")
	}

	productCount, err := s.repo.CountProductsInCategory(ctx, categoryID)
	if err != nil {
		return errs.NewUnexpectedError(fmt.Errorf("fail to get products count for category with ID %s: %w", categoryID, err))
	}

	if productCount > 0 {
		// Return a specific errs error indicating the conflict
		return errs.NewConflictErr(errs.ErrCategoryHasProducts) // Or just errs.ErrCategoryHasProducts directly
	}

	if err := s.repo.DeleteCategory(ctx, categoryID); err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return errs.NewNotFoundErr(err, fmt.Sprintf("category with ID %s for deletion", categoryID))
		}
		return errs.NewUnexpectedError(fmt.Errorf("failed to delete category with ID '%s': %w", categoryID, err))
	}
	return nil
}
