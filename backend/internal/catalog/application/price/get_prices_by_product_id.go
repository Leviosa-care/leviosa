package price

import (
	"context"
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
		return nil, fmt.Errorf("failed to check product existence: %w", err)
	}

	prices, err := s.repo.GetPricesByProductID(ctx, productIDStr, true) // List only active prices by default for external facing
	if err != nil {
		return nil, fmt.Errorf("failed to list prices: %w", err)
	}

	return prices, nil
}
