package product

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *ProductService) GetProductByID(ctx context.Context, productIDStr string) (*domain.ProductRes, error) {
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return nil, errs.NewInvalidValueErr("product ID is not a valid UUID")
	}
	product, err := s.sharedRepo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("get product with ID %s: %w", productID, err)
	}

	category, err := s.sharedRepo.GetCategoryByID(ctx, product.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("get category with ID %s: %w", productID, err)
	}

	return domain.ToProductRes(product, category), nil
}
