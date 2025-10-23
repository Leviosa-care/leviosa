package promotionCode

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *PromotionCodeService) UpdatePromotionCode(ctx context.Context, promotionCodeID string, request *domain.UpdatePromotionCodeRequest) error {
	if err := request.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	id, err := uuid.Parse(promotionCodeID)
	if err != nil {
		return errs.NewInvalidValueErr("invalid promotion code ID format")
	}

	// Check if promotion code exists
	_, err = s.repo.GetPromotionCodeByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "promotion code not found")
		default:
			return errs.NewUnexpectedError(fmt.Errorf("failed to validate promotion code: %w", err))
		}
	}

	err = s.repo.UpdatePromotionCode(ctx, id, request)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "promotion code not found")
		case errors.Is(err, errs.ErrInvalidInput):
			return errs.NewInvalidValueErr(fmt.Sprintf("promotion code update data: %v", err))
		case errors.Is(err, errs.ErrDBQuery):
			return errs.NewQueryFailedErr(fmt.Errorf("repository query failed for promotion code update: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewUnexpectedError(fmt.Errorf("database connection error for promotion code update: %w", err))
		case errors.Is(err, errs.ErrContext):
			return errs.NewUnexpectedError(fmt.Errorf("context error during promotion code update: %w", err))
		default:
			return errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during promotion code update: %w", err))
		}
	}

	return nil
}

func (s *PromotionCodeService) DeactivatePromotionCode(ctx context.Context, promotionCodeID string) error {
	id, err := uuid.Parse(promotionCodeID)
	if err != nil {
		return errs.NewInvalidValueErr("invalid promotion code ID format")
	}

	// Check if promotion code exists
	_, err = s.repo.GetPromotionCodeByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "promotion code not found")
		default:
			return errs.NewUnexpectedError(fmt.Errorf("failed to validate promotion code: %w", err))
		}
	}

	err = s.repo.DeactivatePromotionCode(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "promotion code not found")
		case errors.Is(err, errs.ErrDBQuery):
			return errs.NewQueryFailedErr(fmt.Errorf("repository query failed for promotion code deactivation: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewUnexpectedError(fmt.Errorf("database connection error for promotion code deactivation: %w", err))
		case errors.Is(err, errs.ErrContext):
			return errs.NewUnexpectedError(fmt.Errorf("context error during promotion code deactivation: %w", err))
		default:
			return errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during promotion code deactivation: %w", err))
		}
	}

	return nil
}

func (s *PromotionCodeService) DeletePromotionCode(ctx context.Context, promotionCodeID string) error {
	id, err := uuid.Parse(promotionCodeID)
	if err != nil {
		return errs.NewInvalidValueErr("invalid promotion code ID format")
	}

	// Check if promotion code exists
	_, err = s.repo.GetPromotionCodeByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "promotion code not found")
		default:
			return errs.NewUnexpectedError(fmt.Errorf("failed to validate promotion code: %w", err))
		}
	}

	err = s.repo.DeletePromotionCode(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "promotion code not found")
		case errors.Is(err, errs.ErrDBQuery):
			return errs.NewQueryFailedErr(fmt.Errorf("repository query failed for promotion code deletion: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewUnexpectedError(fmt.Errorf("database connection error for promotion code deletion: %w", err))
		case errors.Is(err, errs.ErrContext):
			return errs.NewUnexpectedError(fmt.Errorf("context error during promotion code deletion: %w", err))
		default:
			return errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during promotion code deletion: %w", err))
		}
	}

	return nil
}

func (s *PromotionCodeService) DeletePromotionCodesByCouponID(ctx context.Context, couponID string) error {
	id, err := uuid.Parse(couponID)
	if err != nil {
		return errs.NewInvalidValueErr("invalid coupon ID format")
	}

	err = s.repo.DeletePromotionCodesByCouponID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrDBQuery):
			return errs.NewQueryFailedErr(fmt.Errorf("repository query failed for promotion codes deletion: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewUnexpectedError(fmt.Errorf("database connection error for promotion codes deletion: %w", err))
		case errors.Is(err, errs.ErrContext):
			return errs.NewUnexpectedError(fmt.Errorf("context error during promotion codes deletion: %w", err))
		default:
			return errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during promotion codes deletion: %w", err))
		}
	}

	return nil
}

