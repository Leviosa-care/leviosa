package coupon

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *CouponService) DeactivateCoupon(ctx context.Context, couponID string) error {
	id, err := uuid.Parse(couponID)
	if err != nil {
		return errs.NewInvalidValueErr("invalid coupon ID format")
	}

	// Check if coupon exists
	_, err = s.repo.GetCouponByID(ctx, id)
	if err != nil {
		return fmt.Errorf("validate coupon: %w", err)
	}

	if err := s.repo.DeactivateCoupon(ctx, id); err != nil {
		return fmt.Errorf("deactivate coupon: %w", err)
	}

	return nil
}
