package promotionCodePayment

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stripe/stripe-go/v82/promotioncode"
)

func (s *service) GetPromotionCode(ctx context.Context, stripePromotionID string) (*domain.PromotionCode, error) {
	stripePromotionCode, err := promotioncode.Get(stripePromotionID, nil)
	if err != nil {
		return nil, err
	}

	return mapStripePromotionCodeToDomainPromotionCode(stripePromotionCode), nil
}