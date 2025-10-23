package category

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *CategoryService) GetAllPublishedCategories(ctx context.Context) ([]*domain.Category, error) {
	categories, err := s.repo.GetAllPublishedCategories(ctx)
	if err != nil {
		return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed to get all published categories: %w", err))
	}

	return categories, nil
}
