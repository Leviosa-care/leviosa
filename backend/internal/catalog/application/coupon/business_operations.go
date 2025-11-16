package coupon

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *CouponService) ValidateCoupon(ctx context.Context, stripeCouponID string) (*domain.CouponResponse, error) {
	if stripeCouponID == "" {
		return nil, errs.NewInvalidValueErr("stripe coupon ID cannot be empty")
	}

	coupon, err := s.repo.GetCouponByStripeID(ctx, stripeCouponID)
	if err != nil {
		return nil, fmt.Errorf("get coupon: %w", err)
	}

	// Validate coupon rules
	if !coupon.IsValid {
		return nil, errs.NewInvalidValueErr("coupon is not valid")
	}

	// Check if coupon is expired
	if coupon.RedeemBy != nil && time.Now().UTC().After(*coupon.RedeemBy) {
		return nil, errs.NewInvalidValueErr("coupon has expired")
	}

	// Check redemption limits
	if coupon.MaxRedemptions != nil && coupon.TimesRedeemed >= *coupon.MaxRedemptions {
		return nil, errs.NewInvalidValueErr("coupon has reached its redemption limit")
	}

	return coupon.ToResponse(), nil
}

func (s *CouponService) IncrementRedemptionCount(ctx context.Context, couponID string) error {
	id, err := uuid.Parse(couponID)
	if err != nil {
		return errs.NewInvalidValueErr("invalid coupon ID format")
	}

	if err := s.repo.IncrementRedemptionCount(ctx, id); err != nil {
		return fmt.Errorf("increment redemption count: %w", err)
	}

	return nil
}

func (s *CouponService) CheckRedemptionLimit(ctx context.Context, couponID string) (bool, error) {
	id, err := uuid.Parse(couponID)
	if err != nil {
		return false, errs.NewInvalidValueErr("invalid coupon ID format")
	}

	coupon, err := s.repo.GetCouponByID(ctx, id)
	if err != nil {
		return false, fmt.Errorf("get coupon: %w", err)
	}

	// Check if there's a redemption limit
	if coupon.MaxRedemptions == nil {
		return true, nil // No limit means it can be redeemed
	}

	// Check if current redemptions exceed or equal the limit
	return coupon.TimesRedeemed < *coupon.MaxRedemptions, nil
}
