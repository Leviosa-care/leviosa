package product

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

func (s *ProductService) GetAllProducts(ctx context.Context) ([]*domain.ProductRes, error) {
	products, err := s.repo.GetAllProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all products: %w", err)
	}
	return products, nil
}
