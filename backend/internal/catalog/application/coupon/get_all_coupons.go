package coupon

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

func (s *CouponService) GetAllCoupons(ctx context.Context) ([]*domain.CouponResponse, error) {
	coupons, err := s.repo.GetAllCoupons(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all coupons: %w", err)
	}

	// Convert slice of coupons to responses
	responses := make([]*domain.CouponResponse, len(coupons))
	for i, coupon := range coupons {
		responses[i] = coupon.ToResponse()
	}
	return responses, nil
}
