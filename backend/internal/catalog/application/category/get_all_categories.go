package category

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

func (s *CategoryService) GetAllCategories(ctx context.Context) ([]*domain.Category, error) {
	categories, err := s.repo.GetAllCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all categories: %w", err)
	}

	return categories, nil
}
