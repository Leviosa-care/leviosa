package promotionCode

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *PromotionCodeService) CreatePromotionCode(ctx context.Context, request *domain.CreatePromotionCodeRequest) (string, error) {
	if err := request.Valid(ctx); err != nil {
		return "", errs.NewInvalidValueErr(err.Error())
	}

	// Validate that the coupon exists
	couponID, err := uuid.Parse(request.CouponID)
	if err != nil {
		return "", errs.NewInvalidValueErr("invalid coupon ID format")
	}

	_, err = s.couponRepo.GetCouponByID(ctx, couponID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return "", errs.NewNotFoundErr(err, "coupon not found")
		default:
			return "", errs.NewUnexpectedError(fmt.Errorf("failed to validate coupon: %w", err))
		}
	}

	now := time.Now().UTC().Truncate(time.Microsecond)

	// Convert restrictions if provided
	var restrictions *domain.PromotionCodeRestrictions
	if request.Restrictions != nil {
		restrictions = &domain.PromotionCodeRestrictions{
			CurrencyOptions: request.Restrictions.CurrencyOptions,
		}
	}

	promotionCode := &domain.PromotionCode{
		ID:                    uuid.New(),
		StripePromotionID:     fmt.Sprintf("promo_%s", uuid.New().String()[:12]), // Generate temp Stripe ID
		CouponID:              couponID,
		Code:                  strings.ToUpper(strings.TrimSpace(request.Code)),
		Active:                true,
		MaxRedemptions:        request.MaxRedemptions,
		TimesRedeemed:         0,
		ExpiresAt:             request.ExpiresAt,
		FirstTimeTransaction:  request.FirstTimeTransaction,
		MinimumAmount:         request.MinimumAmount,
		MinimumAmountCurrency: request.MinimumAmountCurrency,
		Restrictions:          restrictions,
		CreatedAt:             now,
		UpdatedAt:             now,
		Metadata:              request.Metadata,
	}

	promotionCodeID, err := s.repo.CreatePromotionCode(ctx, promotionCode)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			return "", errs.NewInvalidValueErr(fmt.Sprintf("promotion code data: %v", err))
		case errors.Is(err, errs.ErrUniqueViolation):
			return "", errs.NewAlreadyExistsError(err, "promotion code with this code")
		case errors.Is(err, errs.ErrNotNullViolation):
			return "", errs.NewInvalidValueErr(fmt.Sprintf("missing required data for promotion code: %v", err))
		case errors.Is(err, errs.ErrForeignKeyViolation):
			return "", errs.NewInvalidValueErr(fmt.Sprintf("invalid coupon reference: %v", err))
		case errors.Is(err, errs.ErrCheckViolation):
			return "", errs.NewInvalidValueErr(fmt.Sprintf("promotion code data failed check constraint: %v", err))
		case errors.Is(err, errs.ErrDBQuery):
			return "", errs.NewQueryFailedErr(fmt.Errorf("repository query failed for promotion code: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return "", errs.NewUnexpectedError(fmt.Errorf("database connection error for promotion code: %w", err))
		case errors.Is(err, errs.ErrContext):
			return "", errs.NewUnexpectedError(fmt.Errorf("context error during promotion code creation: %w", err))
		default:
			return "", errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during promotion code creation: %w", err))
		}
	}

	return promotionCodeID, nil
}

