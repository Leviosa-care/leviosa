package category

import (
	"context"
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
		return nil, fmt.Errorf("get category with ID %s", categoryIDStr)
	}

	return category, nil
}
