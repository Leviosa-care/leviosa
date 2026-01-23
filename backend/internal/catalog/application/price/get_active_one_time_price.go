package price

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

// GetActiveOneTimePriceByProductID retrieves the active one-time price for a product
// in the specified currency. This is primarily used by the booking service to
// determine the price for a booking.
//
// Parameters:
//   - productID: The product to get the price for
//   - currency: The currency code (e.g., "EUR", "USD") - case insensitive
//
// Returns:
//   - The active one-time price if found
//   - ErrRepositoryNotFound if no matching price exists
//   - ErrInvalidValue if the product ID or currency is invalid
func (s *PriceService) GetActiveOneTimePriceByProductID(ctx context.Context, productIDStr string, currency string) (*domain.Price, error) {
	// Validate product ID
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("product ID is invalid: %s", err.Error()))
	}

	// Validate currency format (3-letter uppercase ISO code)
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if len(currency) != 3 {
		return nil, errs.NewInvalidValueErr("currency must be a 3-letter ISO code")
	}

	// Verify product exists
	_, err = s.sharedRepo.GetProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, fmt.Errorf("product not found: %w", errs.ErrRepositoryNotFound)
		}
		return nil, fmt.Errorf("failed to check product existence: %w", err)
	}

	// Get all active prices for this product
	prices, err := s.repo.GetPricesByProductID(ctx, productIDStr, true)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, fmt.Errorf("no prices found for product: %w", errs.ErrRepositoryNotFound)
		}
		return nil, fmt.Errorf("failed to get prices: %w", err)
	}

	// Find the active one-time price with matching currency
	for _, price := range prices {
		if price.IsActive && price.Interval == domain.OneTime && price.Currency == currency {
			return price, nil
		}
	}

	return nil, fmt.Errorf("no active one-time price found for product in %s: %w", currency, errs.ErrRepositoryNotFound)
}
