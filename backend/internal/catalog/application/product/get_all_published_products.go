package product

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

func (s *ProductService) GetAllPublishedProducts(ctx context.Context) ([]*domain.ProductRes, error) {
	products, err := s.repo.GetAllPublishedProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all published products: %w", err)
	}

	return products, nil
}
