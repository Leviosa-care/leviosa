package coupon

import (
	"context"
	"errors"
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
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "coupon not found")
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("failed to get coupon: %w", err))
		}
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

	err = s.repo.IncrementRedemptionCount(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "coupon not found")
		case errors.Is(err, errs.ErrDBQuery):
			return errs.NewQueryFailedErr(fmt.Errorf("repository query failed for redemption count increment: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewUnexpectedError(fmt.Errorf("database connection error for redemption count increment: %w", err))
		case errors.Is(err, errs.ErrContext):
			return errs.NewUnexpectedError(fmt.Errorf("context error during redemption count increment: %w", err))
		default:
			return errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during redemption count increment: %w", err))
		}
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
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return false, errs.NewNotFoundErr(err, "coupon not found")
		default:
			return false, errs.NewUnexpectedError(fmt.Errorf("failed to get coupon: %w", err))
		}
	}

	// Check if there's a redemption limit
	if coupon.MaxRedemptions == nil {
		return true, nil // No limit means it can be redeemed
	}

	// Check if current redemptions exceed or equal the limit
	return coupon.TimesRedeemed < *coupon.MaxRedemptions, nil
}
