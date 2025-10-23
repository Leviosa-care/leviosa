package pricePayment

import (
	"context"

	"github.com/stripe/stripe-go/v82"
)

func (s *service) ReactivatePrices(ctx context.Context, stripePriceIDs []string) error {
	for _, priceID := range stripePriceIDs {
		params := &stripe.PriceUpdateParams{
			Active: stripe.Bool(true),
		}

		_, err := s.Client.V1Prices.Update(ctx, priceID, params)
		if err != nil {
			return err
		}
	}

	return nil
}

