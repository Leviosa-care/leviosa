package pricePayment

import (
	"context"

	"github.com/stripe/stripe-go/v82"
)

func (s *service) DeactivatePrices(ctx context.Context, priceIDs []string) error {
	for _, priceID := range priceIDs {
		params := &stripe.PriceUpdateParams{
			Active: stripe.Bool(false),
		}

		_, err := s.Client.V1Prices.Update(ctx, priceID, params)
		if err != nil {
			return err
		}
	}

	return nil
}