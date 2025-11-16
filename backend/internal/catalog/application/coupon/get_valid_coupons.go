package coupon

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

func (s *CouponService) GetValidCoupons(ctx context.Context) ([]*domain.CouponResponse, error) {
	coupons, err := s.repo.GetValidCoupons(ctx)
	if err != nil {
		return nil, fmt.Errorf("get valid coupons: %w", err)
	}

	// Convert slice of coupons to responses
	responses := make([]*domain.CouponResponse, len(coupons))
	for i, coupon := range coupons {
		responses[i] = coupon.ToResponse()
	}
	return responses, nil
}
