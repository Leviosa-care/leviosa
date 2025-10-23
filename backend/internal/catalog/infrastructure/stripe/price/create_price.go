package pricePayment

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stripe/stripe-go/v82"
)

func (s *service) CreatePrice(ctx context.Context, input domain.CreateStripePriceRequest) (*domain.PaymentPrice, error) {
	params := &stripe.PriceCreateParams{
		Product:    stripe.String(input.ProductID),
		UnitAmount: stripe.Int64(int64(input.Amount)),
		Currency:   stripe.String(input.Currency),
		Recurring: &stripe.PriceCreateRecurringParams{
			Interval: stripe.String(input.Interval),
		},
		Active: stripe.Bool(input.Active),
	}

	if input.Nickname != "" {
		params.Nickname = stripe.String(input.Nickname)
	}

	if len(input.Metadata) > 0 {
		params.Metadata = input.Metadata
	}

	stripePrice, err := s.Client.V1Prices.Create(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapStripePriceToPaymentPrice(stripePrice), nil
}

func mapStripePriceToPaymentPrice(stripePrice *stripe.Price) *domain.PaymentPrice {
	var interval string
	if stripePrice.Recurring != nil {
		interval = string(stripePrice.Recurring.Interval)
	}

	return &domain.PaymentPrice{
		ID:       stripePrice.ID,
		Product:  stripePrice.Product.ID,
		Amount:   int64(stripePrice.UnitAmount),
		Currency: string(stripePrice.Currency),
		Interval: interval,
		Active:   stripePrice.Active,
		Nickname: stripePrice.Nickname,
		Metadata: stripePrice.Metadata,
	}
}