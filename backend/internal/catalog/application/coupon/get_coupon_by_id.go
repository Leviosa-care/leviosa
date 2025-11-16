package coupon

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *CouponService) GetCouponByID(ctx context.Context, couponID string) (*domain.CouponResponse, error) {
	id, err := uuid.Parse(couponID)
	if err != nil {
		return nil, errs.NewInvalidValueErr("invalid coupon ID format")
	}

	coupon, err := s.repo.GetCouponByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get coupon by ID: %w", err)
	}

	return coupon.ToResponse(), nil
}
