package pricePayment

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

func (s *service) GetPrice(ctx context.Context, priceID string) (*domain.PaymentPrice, error) {
	stripePrice, err := s.Client.V1Prices.Retrieve(ctx, priceID, nil)
	if err != nil {
		return nil, err
	}

	return mapStripePriceToPaymentPrice(stripePrice), nil
}