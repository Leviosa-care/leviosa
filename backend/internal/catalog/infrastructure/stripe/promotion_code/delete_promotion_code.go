package promotionCodePayment

import (
	"context"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/promotioncode"
)

func (s *service) DeletePromotionCode(ctx context.Context, stripePromotionID string) error {
	// Note: Stripe doesn't have a delete operation for promotion codes
	// Instead, we deactivate the promotion code
	params := &stripe.PromotionCodeParams{
		Active: stripe.Bool(false),
	}

	_, err := promotioncode.Update(stripePromotionID, params)
	return err
}