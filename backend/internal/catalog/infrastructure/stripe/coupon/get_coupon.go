package couponPayment

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stripe/stripe-go/v82/coupon"
)

func (s *service) GetCoupon(ctx context.Context, stripeCouponID string) (*domain.Coupon, error) {
	stripeCoupon, err := coupon.Get(stripeCouponID, nil)
	if err != nil {
		return nil, err
	}

	return mapStripeCouponToDomainCoupon(stripeCoupon), nil
}