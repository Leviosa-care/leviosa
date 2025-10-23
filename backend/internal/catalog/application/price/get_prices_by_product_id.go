package price

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetPricesByProductID retrieves a list of prices for a given internal product ID.
func (s *PriceService) GetPricesByProductID(ctx context.Context, productIDStr string) ([]*domain.Price, error) {
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("product ID is invalid: %s", err.Error()))
	}

	_, err = s.sharedRepo.GetProductByID(ctx, productID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewConflictErr(fmt.Errorf("product with ID '%s' does not exist in database", productID))
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("failed to check product existence: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled error checking product existence: %w", err))
		}
	}

	prices, err := s.repo.GetPricesByProductID(ctx, productIDStr, true) // List only active prices by default for external facing
	if err != nil {
		return nil, errs.NewQueryFailedErr(fmt.Errorf("failed to list prices from database: %w", err))
	}

	return prices, nil
}
