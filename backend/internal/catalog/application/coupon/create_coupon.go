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
		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			return "", errs.NewInvalidValueErr(fmt.Sprintf("coupon data: %v", err))
		case errors.Is(err, errs.ErrUniqueViolation):
			return "", errs.NewAlreadyExistsError(err, "coupon with this name")
		case errors.Is(err, errs.ErrNotNullViolation):
			return "", errs.NewInvalidValueErr(fmt.Sprintf("missing required data for coupon: %v", err))
		case errors.Is(err, errs.ErrCheckViolation):
			return "", errs.NewInvalidValueErr(fmt.Sprintf("coupon data failed check constraint: %v", err))
		case errors.Is(err, errs.ErrDBQuery):
			return "", errs.NewQueryFailedErr(fmt.Errorf("repository query failed for coupon: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return "", errs.NewUnexpectedError(fmt.Errorf("database connection error for coupon: %w", err))
		case errors.Is(err, errs.ErrContext):
			return "", errs.NewUnexpectedError(fmt.Errorf("context error during coupon creation: %w", err))
		default:
			return "", errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during coupon creation: %w", err))
		}
	}

	return couponID, nil
}
