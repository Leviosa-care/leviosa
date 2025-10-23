package product

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *ProductService) GetAllProducts(ctx context.Context) ([]*domain.ProductRes, error) {
	products, err := s.repo.GetAllProducts(ctx)
	if err != nil {
		return nil, errs.NewQueryFailedErr(fmt.Errorf("failed to retrieve products: %w", err))
	}
	return products, nil
}
