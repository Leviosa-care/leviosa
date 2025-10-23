package price

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// GetPriceByStripeID retrieves a price from the internal database using its Stripe Price ID.
// This is crucial for handling webhooks or other Stripe-initiated events.
func (s *PriceService) GetPriceByStripeID(ctx context.Context, stripePriceID string) (*domain.Price, error) {
	if stripePriceID == "" {
		return nil, errs.NewInvalidValueErr("Stripe Price ID cannot be empty")
	}

	p, err := s.repo.GetPriceByStripeID(ctx, stripePriceID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "price by Stripe ID")
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("failed to get price by Stripe ID %s from database: %w", stripePriceID, err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled error getting price by Stripe ID %s: %w", stripePriceID, err))
		}
	}

	return p, nil
}
