package category

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

func (s *CategoryService) GetAllPublishedCategories(ctx context.Context) ([]*domain.Category, error) {
	categories, err := s.repo.GetAllPublishedCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all published categories: %w", err)
	}

	return categories, nil
}
