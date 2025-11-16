package price

import (
	"context"
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
		return nil, fmt.Errorf("get price by Stripe ID %s: %w", stripePriceID, err)
	}

	return p, nil
}
