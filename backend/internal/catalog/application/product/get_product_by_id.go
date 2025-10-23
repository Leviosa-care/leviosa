package product

import (
	"context"
	"errors"
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
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.NewNotFoundErr(err, fmt.Sprintf("product with ID %s", productID))
		}
		return nil, errs.NewUnexpectedError(fmt.Errorf("failed to retrieve product with ID %s: %w", productID, err))
	}
	category, err := s.sharedRepo.GetCategoryByID(ctx, product.CategoryID)
	if err != nil {

		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.NewNotFoundErr(err, fmt.Sprintf("category with ID %s", productID))
		}
		return nil, errs.NewUnexpectedError(fmt.Errorf("failed to retrieve category with ID %s: %w", productID, err))
	}

	return domain.ToProductRes(product, category), nil
}
