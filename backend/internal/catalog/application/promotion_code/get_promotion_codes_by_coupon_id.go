package promotionCode

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *PromotionCodeService) GetPromotionCodesByCouponID(ctx context.Context, couponID string) ([]*domain.PromotionCode, error) {
	id, err := uuid.Parse(couponID)
	if err != nil {
		return nil, errs.NewInvalidValueErr("invalid coupon ID format")
	}

	promotionCodes, err := s.repo.GetPromotionCodesByCouponID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get promotion codes by coupon ID: %w", err)
	}

	return promotionCodes, nil
}
