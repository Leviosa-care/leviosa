package coupon

import (
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
)

type CouponService struct {
	repo ports.CouponRepository
}

func NewCouponService(repo ports.CouponRepository) *CouponService {
	return &CouponService{
		repo: repo,
	}
}

// Compile-time check to ensure CouponService implements ports.CouponService
var _ ports.CouponService = (*CouponService)(nil)
