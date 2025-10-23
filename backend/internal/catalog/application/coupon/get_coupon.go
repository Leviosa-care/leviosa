package coupon

import (
	"context"
	"errors"
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
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "coupon not found")
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("failed to get coupon: %w", err))
		}
	}

	return coupon.ToResponse(), nil
}

func (s *CouponService) GetCouponByStripeID(ctx context.Context, stripeCouponID string) (*domain.CouponResponse, error) {
	if stripeCouponID == "" {
		return nil, errs.NewInvalidValueErr("stripe coupon ID cannot be empty")
	}

	coupon, err := s.repo.GetCouponByStripeID(ctx, stripeCouponID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "coupon not found")
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("failed to get coupon by stripe ID: %w", err))
		}
	}

	return coupon.ToResponse(), nil
}

func (s *CouponService) GetAllCoupons(ctx context.Context) ([]*domain.CouponResponse, error) {
	coupons, err := s.repo.GetAllCoupons(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed for get all coupons: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error for get all coupons: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error during get all coupons: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during get all coupons: %w", err))
		}
	}

	// Convert slice of coupons to responses
	responses := make([]*domain.CouponResponse, len(coupons))
	for i, coupon := range coupons {
		responses[i] = coupon.ToResponse()
	}
	return responses, nil
}

func (s *CouponService) GetValidCoupons(ctx context.Context) ([]*domain.CouponResponse, error) {
	coupons, err := s.repo.GetValidCoupons(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed for get valid coupons: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error for get valid coupons: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error during get valid coupons: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during get valid coupons: %w", err))
		}
	}

	// Convert slice of coupons to responses
	responses := make([]*domain.CouponResponse, len(coupons))
	for i, coupon := range coupons {
		responses[i] = coupon.ToResponse()
	}
	return responses, nil
}

