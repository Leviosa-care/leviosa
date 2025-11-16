package promotionCode

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *PromotionCodeService) DeletePromotionCodesByCouponID(ctx context.Context, couponID string) error {
	id, err := uuid.Parse(couponID)
	if err != nil {
		return errs.NewInvalidValueErr("invalid coupon ID format")
	}

	if err := s.repo.DeletePromotionCodesByCouponID(ctx, id); err != nil {
		return fmt.Errorf("delete promotion codes by coupon ID: %w", err)
	}

	return nil
}
