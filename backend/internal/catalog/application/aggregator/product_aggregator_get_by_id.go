package aggregator

import (
	"context"
	"errors"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *ProductAggregatorService) GetProductByID(ctx context.Context, productIDStr string) (*domain.ProductAggregator, error) {
	image, err := s.imageService.GetActiveImage(ctx, productIDStr, string(domain.ProductType))
	if err != nil {
		if !errors.Is(err, errs.ErrDomainNotFound) {
			return nil, err
		}
		image = nil
	}

	product, err := s.productService.GetProductByID(ctx, productIDStr)
	if err != nil {
		return nil, err
	}
	prices, err := s.priceService.GetPricesByProductID(ctx, productIDStr)
	if err != nil {
		return nil, err
	}

	return &domain.ProductAggregator{
		Product: product,
		Image:   image,
		Prices:  prices,
	}, nil
}