package product

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *ProductService) GetAllPublishedProducts(ctx context.Context) ([]*domain.ProductRes, error) {
	products, err := s.repo.GetAllPublishedProducts(ctx)
	if err != nil {
		return nil, errs.NewUnexpectedError(fmt.Errorf("failed to retrieve published products: %w", err))
	}

	return products, nil
}
