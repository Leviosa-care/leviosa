package coupon

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *CouponService) CreateCoupon(ctx context.Context, req *domain.CreateCouponRequest) (string, error) {
	if err := req.Valid(ctx); err != nil {
		return "", errs.NewInvalidValueErr(err.Error())
	}

	now := time.Now().UTC().Truncate(time.Microsecond)

	coupon := &domain.Coupon{
		ID:               uuid.New(),
		StripeCouponID:   fmt.Sprintf("coupon_%s", uuid.New().String()[:12]), // Generate temp Stripe ID
		Name:             req.Name,
		PercentOff:       req.PercentOff,
		AmountOff:        req.AmountOff,
		Currency:         req.Currency,
		Duration:         domain.CouponDuration(req.Duration),
		DurationInMonths: req.DurationInMonths,
		MaxRedemptions:   req.MaxRedemptions,
		TimesRedeemed:    0,
		IsValid:          true,
		RedeemBy:         req.RedeemBy,
		CreatedAt:        now,
		UpdatedAt:        now,
		Metadata:         req.Metadata,
	}

	couponID, err := s.repo.CreateCoupon(ctx, coupon)
	if err != nil {
		return "", fmt.Errorf("create coupon: %w", err)
	}

	return couponID, nil
}
