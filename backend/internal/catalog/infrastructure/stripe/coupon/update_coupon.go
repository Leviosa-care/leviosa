package couponPayment

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/coupon"
)

func (s *service) UpdateCoupon(ctx context.Context, stripeCouponID string, req *domain.UpdateCouponRequest) (*domain.Coupon, error) {
	params := &stripe.CouponParams{}

	if req.Name != nil {
		params.Name = stripe.String(*req.Name)
	}

	if req.Metadata != nil {
		params.Metadata = req.Metadata
	}

	stripeCoupon, err := coupon.Update(stripeCouponID, params)
	if err != nil {
		return nil, err
	}

	return mapStripeCouponToDomainCoupon(stripeCoupon), nil
}