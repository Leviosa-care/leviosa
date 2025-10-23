package price

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetPrice retrieves price details by its internal ID.
func (s *PriceService) GetPrice(ctx context.Context, priceID string) (*domain.Price, error) {
	if err := uuid.Validate(priceID); err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("price ID is invalid: %s", err.Error()))
	}

	p, err := s.repo.GetPrice(ctx, priceID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "price")
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("failed to get price from database: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled error getting price: %w", err))
		}
	}

	// Optionally, fetch from Stripe to confirm latest status if needed
	// paymentPrice, err := s.paymentGateway.GetPrice(ctx, p.StripePriceID)
	// if err != nil {
	//    log.Printf("Warning: Could not get latest Stripe status for price %s: %v", p.StripePriceID, err)
	//    // Decide if this is a critical error or just a warning
	// }

	return p, nil
}
