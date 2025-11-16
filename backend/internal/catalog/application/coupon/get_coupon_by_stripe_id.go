package coupon

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *CouponService) GetCouponByStripeID(ctx context.Context, stripeCouponID string) (*domain.CouponResponse, error) {
	if stripeCouponID == "" {
		return nil, errs.NewInvalidValueErr("stripe coupon ID cannot be empty")
	}

	coupon, err := s.repo.GetCouponByStripeID(ctx, stripeCouponID)
	if err != nil {
		return nil, fmt.Errorf("get coupon by Stripe ID %s: %w", stripeCouponID, err)
	}

	return coupon.ToResponse(), nil
}
