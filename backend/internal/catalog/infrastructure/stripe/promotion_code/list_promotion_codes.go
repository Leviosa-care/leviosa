package promotionCodePayment

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/promotioncode"
)

func (s *service) ListPromotionCodes(ctx context.Context, stripeCouponID string) ([]*domain.PromotionCode, error) {
	params := &stripe.PromotionCodeListParams{
		Coupon: stripe.String(stripeCouponID),
	}

	iter := promotioncode.List(params)
	var promotionCodes []*domain.PromotionCode

	for iter.Next() {
		stripePromotionCode := iter.PromotionCode()
		promotionCodes = append(promotionCodes, mapStripePromotionCodeToDomainPromotionCode(stripePromotionCode))
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return promotionCodes, nil
}