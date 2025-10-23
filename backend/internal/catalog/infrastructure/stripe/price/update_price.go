package pricePayment

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stripe/stripe-go/v82"
)

func (s *service) UpdatePrice(ctx context.Context, stripePriceID string, req domain.UpdateStripePriceRequest) (*domain.PaymentPrice, error) {
	params := &stripe.PriceUpdateParams{}

	if req.Active != nil {
		params.Active = stripe.Bool(*req.Active)
	}

	if req.Nickname != nil {
		params.Nickname = stripe.String(*req.Nickname)
	}

	if req.Metadata != nil {
		params.Metadata = req.Metadata
	}

	stripePrice, err := s.Client.V1Prices.Update(ctx, stripePriceID, params)
	if err != nil {
		return nil, err
	}

	return mapStripePriceToPaymentPrice(stripePrice), nil
}