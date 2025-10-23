package promotionCode

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *PromotionCodeService) GetPromotionCodeByID(ctx context.Context, promotionCodeID string) (*domain.PromotionCode, error) {
	id, err := uuid.Parse(promotionCodeID)
	if err != nil {
		return nil, errs.NewInvalidValueErr("invalid promotion code ID format")
	}

	promotionCode, err := s.repo.GetPromotionCodeByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "promotion code not found")
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error: %w", err))
		}
	}

	return promotionCode, nil
}

func (s *PromotionCodeService) GetPromotionCodeByCode(ctx context.Context, code string) (*domain.PromotionCode, error) {
	if code == "" {
		return nil, errs.NewInvalidValueErr("promotion code cannot be empty")
	}

	code = strings.ToUpper(strings.TrimSpace(code))

	promotionCode, err := s.repo.GetPromotionCodeByCode(ctx, code)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "promotion code not found")
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error: %w", err))
		}
	}

	return promotionCode, nil
}

func (s *PromotionCodeService) GetPromotionCodeByStripeID(ctx context.Context, stripePromotionID string) (*domain.PromotionCode, error) {
	if stripePromotionID == "" {
		return nil, errs.NewInvalidValueErr("stripe promotion ID cannot be empty")
	}

	promotionCode, err := s.repo.GetPromotionCodeByStripeID(ctx, stripePromotionID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "promotion code not found")
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error: %w", err))
		}
	}

	return promotionCode, nil
}

func (s *PromotionCodeService) GetPromotionCodesByCouponID(ctx context.Context, couponID string) ([]*domain.PromotionCode, error) {
	id, err := uuid.Parse(couponID)
	if err != nil {
		return nil, errs.NewInvalidValueErr("invalid coupon ID format")
	}

	promotionCodes, err := s.repo.GetPromotionCodesByCouponID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error: %w", err))
		}
	}

	return promotionCodes, nil
}

func (s *PromotionCodeService) GetAllPromotionCodes(ctx context.Context) ([]*domain.PromotionCode, error) {
	promotionCodes, err := s.repo.GetAllPromotionCodes(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error: %w", err))
		}
	}

	return promotionCodes, nil
}

func (s *PromotionCodeService) GetActivePromotionCodes(ctx context.Context) ([]*domain.PromotionCode, error) {
	promotionCodes, err := s.repo.GetActivePromotionCodes(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrDBQuery):
			return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewUnexpectedError(fmt.Errorf("database connection error: %w", err))
		case errors.Is(err, errs.ErrContext):
			return nil, errs.NewUnexpectedError(fmt.Errorf("context error: %w", err))
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error: %w", err))
		}
	}

	return promotionCodes, nil
}

