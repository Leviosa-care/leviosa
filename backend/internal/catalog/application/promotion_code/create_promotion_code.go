package promotionCode

import (
	"context"
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
		return "", fmt.Errorf("get coupon by ID: %w", err)
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
		return "", fmt.Errorf("create promotion code: %w", err)
	}

	return promotionCodeID, nil
}
