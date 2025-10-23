package promotionCode

import (
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
)

// check if *PromotionCodeService implements PromotionCodeService interface
var _ ports.PromotionCodeService = (*PromotionCodeService)(nil)

type PromotionCodeService struct {
	repo       ports.PromotionCodeRepository
	couponRepo ports.CouponRepository
}

func New(repo ports.PromotionCodeRepository, couponRepo ports.CouponRepository) ports.PromotionCodeService {
	return &PromotionCodeService{
		repo:       repo,
		couponRepo: couponRepo,
	}
}
