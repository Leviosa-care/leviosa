package coupon

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *CouponService) UpdateCoupon(ctx context.Context, couponID string, req *domain.UpdateCouponRequest) error {
	if err := req.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	id, err := uuid.Parse(couponID)
	if err != nil {
		return errs.NewInvalidValueErr("invalid coupon ID format")
	}

	// Check if coupon exists
	_, err = s.repo.GetCouponByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "coupon not found")
		default:
			return errs.NewUnexpectedError(fmt.Errorf("failed to validate coupon: %w", err))
		}
	}

	err = s.repo.UpdateCoupon(ctx, id, req)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "coupon not found")
		case errors.Is(err, errs.ErrInvalidInput):
			return errs.NewInvalidValueErr(fmt.Sprintf("coupon update data: %v", err))
		case errors.Is(err, errs.ErrDBQuery):
			return errs.NewQueryFailedErr(fmt.Errorf("repository query failed for coupon update: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewUnexpectedError(fmt.Errorf("database connection error for coupon update: %w", err))
		case errors.Is(err, errs.ErrContext):
			return errs.NewUnexpectedError(fmt.Errorf("context error during coupon update: %w", err))
		default:
			return errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during coupon update: %w", err))
		}
	}

	return nil
}

func (s *CouponService) DeactivateCoupon(ctx context.Context, couponID string) error {
	id, err := uuid.Parse(couponID)
	if err != nil {
		return errs.NewInvalidValueErr("invalid coupon ID format")
	}

	// Check if coupon exists
	_, err = s.repo.GetCouponByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "coupon not found")
		default:
			return errs.NewUnexpectedError(fmt.Errorf("failed to validate coupon: %w", err))
		}
	}

	err = s.repo.DeactivateCoupon(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "coupon not found")
		case errors.Is(err, errs.ErrDBQuery):
			return errs.NewQueryFailedErr(fmt.Errorf("repository query failed for coupon deactivation: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewUnexpectedError(fmt.Errorf("database connection error for coupon deactivation: %w", err))
		case errors.Is(err, errs.ErrContext):
			return errs.NewUnexpectedError(fmt.Errorf("context error during coupon deactivation: %w", err))
		default:
			return errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during coupon deactivation: %w", err))
		}
	}

	return nil
}

func (s *CouponService) DeleteCoupon(ctx context.Context, couponID string) error {
	id, err := uuid.Parse(couponID)
	if err != nil {
		return errs.NewInvalidValueErr("invalid coupon ID format")
	}

	// Check if coupon exists
	_, err = s.repo.GetCouponByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "coupon not found")
		default:
			return errs.NewUnexpectedError(fmt.Errorf("failed to validate coupon: %w", err))
		}
	}

	err = s.repo.DeleteCoupon(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "coupon not found")
		case errors.Is(err, errs.ErrDBQuery):
			return errs.NewQueryFailedErr(fmt.Errorf("repository query failed for coupon deletion: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewUnexpectedError(fmt.Errorf("database connection error for coupon deletion: %w", err))
		case errors.Is(err, errs.ErrContext):
			return errs.NewUnexpectedError(fmt.Errorf("context error during coupon deletion: %w", err))
		default:
			return errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during coupon deletion: %w", err))
		}
	}

	return nil
}
