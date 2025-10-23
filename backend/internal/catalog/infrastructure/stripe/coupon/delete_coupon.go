package couponPayment

import (
	"context"

	"github.com/stripe/stripe-go/v82/coupon"
)

func (s *service) DeleteCoupon(ctx context.Context, stripeCouponID string) error {
	_, err := coupon.Del(stripeCouponID, nil)
	return err
}