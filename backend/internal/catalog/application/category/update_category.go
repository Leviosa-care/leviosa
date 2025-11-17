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

	if err := s.repo.UpdateCategory(ctx, category); err != nil {
		switch {
		case errors.Is(err, errs.ErrUniqueViolation):
			// Handle unique constraint violations (e.g., if category name or slug must be unique)
			return errs.NewConflictErr(err)
		default:
			return fmt.Errorf("update category: %w", err)
		}
	}

	return nil
}
