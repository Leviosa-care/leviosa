package aggregator

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// CreateProductWithPrice creates a new product and an associated price.
// If the price creation fails, it attempts to remove the newly created product.
// The function returns a composite error if the rollback also fails.
func (s *ProductAggregatorService) CreateProductWithPrice(ctx context.Context, request *domain.CreateProductWithPriceRequest) (string, string, error) {
	productID, err := s.productService.CreateProduct(ctx, &request.Product)
	if err != nil {
		// domain error are alredy handled
		return "", "", err
	}
	priceID, err := s.priceService.CreatePrice(ctx, productID, &request.Price)
	if err != nil {
		// rollback if error
		// if err := s.productService.RemoveProduct(ctx, productID); err != nil {
		if rollbackErr := s.productService.RemoveProduct(ctx, productID); !errors.Is(rollbackErr, errs.ErrDomainNotFound) {
			return "", "", fmt.Errorf("failed to create price (%w), AND product rollback also failed: %w", err, rollbackErr)
		}
		return "", "", err
	}
	return productID, priceID, nil
}