package category

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *CategoryService) GetCategoryByID(ctx context.Context, categoryIDStr string) (*domain.Category, error) {
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	category, err := s.sharedRepo.GetCategoryByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.NewNotFoundErr(err, fmt.Sprintf("category with ID %s", categoryIDStr))
		}
		return nil, errs.NewUnexpectedError(fmt.Errorf("failed to retrieve category with ID %s: %w", categoryIDStr, err))
	}

	return category, nil
}
