package pricePayment

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stripe/stripe-go/v82"
)

func (s *service) GetPricesByProductID(ctx context.Context, productID string, opts *domain.PriceListOptions) ([]*domain.PaymentPrice, error) {
	params := &stripe.PriceListParams{
		Product: stripe.String(productID),
	}

	if opts != nil {
		if opts.Active != nil {
			params.Active = stripe.Bool(*opts.Active)
		}
		if opts.Limit > 0 {
			params.Limit = stripe.Int64(int64(opts.Limit))
		}
	}

	iter := s.Client.V1Prices.List(ctx, params)
	var prices []*domain.PaymentPrice

	for stripePrice, err := range iter {
		if err != nil {
			return nil, err
		}
		prices = append(prices, mapStripePriceToPaymentPrice(stripePrice))
	}

	return prices, nil
}