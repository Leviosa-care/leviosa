package stripeService

import (
	"context"
	"fmt"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/price"
)

func (s *service) RemovePrice(ctx context.Context, priceID string) error {
	price_params := &stripe.PriceParams{
		Active: stripe.Bool(false),
	}
	_, err := price.Update(priceID, price_params)
	if err != nil {
		return fmt.Errorf("Failed to udpate active status to false for price with ID %s: %w", priceID, err)
	}
	return nil
}
