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

func (s *PromotionCodeService) ValidatePromotionCode(ctx context.Context, req *domain.ValidatePromotionCodeRequest) (*domain.ValidatePromotionCodeResponse, error) {
	if err := req.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	code := strings.ToUpper(strings.TrimSpace(req.Code))

	// Get promotion code
	promotionCode, err := s.repo.GetPromotionCodeByCode(ctx, code)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return &domain.ValidatePromotionCodeResponse{
				Valid:  false,
				Reason: "promotion code not found",
			}, nil
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("failed to get promotion code: %w", err))
		}
	}

	// Get associated coupon
	coupon, err := s.couponRepo.GetCouponByID(ctx, promotionCode.CouponID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return &domain.ValidatePromotionCodeResponse{
				Valid:  false,
				Reason: "associated coupon not found",
			}, nil
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("failed to get associated coupon: %w", err))
		}
	}

	// Validate promotion code
	validationResult := s.validatePromotionCodeRules(promotionCode, coupon, req)
	if !validationResult.Valid {
		return validationResult, nil
	}

	// Build successful response with full details
	promotionCodeResponse := s.buildPromotionCodeResponse(promotionCode)
	couponResponse := s.buildCouponResponse(coupon)

	return &domain.ValidatePromotionCodeResponse{
		Valid: true,
		PromotionCode: &domain.PromotionCodeWithCouponResponse{
			PromotionCode: *promotionCodeResponse,
			Coupon:        *couponResponse,
		},
	}, nil
}

func (s *PromotionCodeService) IncrementRedemptionCount(ctx context.Context, promotionCodeID string) error {
	id, err := uuid.Parse(promotionCodeID)
	if err != nil {
		return errs.NewInvalidValueErr("invalid promotion code ID format")
	}

	err = s.repo.IncrementRedemptionCount(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "promotion code not found")
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

func (s *PromotionCodeService) CheckRedemptionLimit(ctx context.Context, promotionCodeID string) (bool, error) {
	id, err := uuid.Parse(promotionCodeID)
	if err != nil {
		return false, errs.NewInvalidValueErr("invalid promotion code ID format")
	}

	promotionCode, err := s.repo.GetPromotionCodeByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return false, errs.NewNotFoundErr(err, "promotion code not found")
		default:
			return false, errs.NewUnexpectedError(fmt.Errorf("failed to get promotion code: %w", err))
		}
	}

	// Check if there's a redemption limit
	if promotionCode.MaxRedemptions == nil {
		return true, nil // No limit means it can be redeemed
	}

	// Check if current redemptions exceed or equal the limit
	return promotionCode.TimesRedeemed < *promotionCode.MaxRedemptions, nil
}

func (s *PromotionCodeService) GetPromotionCodeWithCoupon(ctx context.Context, code string) (*domain.PromotionCodeWithCouponResponse, error) {
	if code == "" {
		return nil, errs.NewInvalidValueErr("promotion code cannot be empty")
	}

	code = strings.ToUpper(strings.TrimSpace(code))

	// Get promotion code
	promotionCode, err := s.repo.GetPromotionCodeByCode(ctx, code)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "promotion code not found")
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("failed to get promotion code: %w", err))
		}
	}

	// Get associated coupon
	coupon, err := s.couponRepo.GetCouponByID(ctx, promotionCode.CouponID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return nil, errs.NewNotFoundErr(err, "associated coupon not found")
		default:
			return nil, errs.NewUnexpectedError(fmt.Errorf("failed to get associated coupon: %w", err))
		}
	}

	promotionCodeResponse := s.buildPromotionCodeResponse(promotionCode)
	couponResponse := s.buildCouponResponse(coupon)

	return &domain.PromotionCodeWithCouponResponse{
		PromotionCode: *promotionCodeResponse,
		Coupon:        *couponResponse,
	}, nil
}

// Helper functions

func (s *PromotionCodeService) validatePromotionCodeRules(promotionCode *domain.PromotionCode, coupon *domain.Coupon, req *domain.ValidatePromotionCodeRequest) *domain.ValidatePromotionCodeResponse {
	// Check if promotion code is active
	if !promotionCode.Active {
		return &domain.ValidatePromotionCodeResponse{
			Valid:  false,
			Reason: "promotion code is not active",
		}
	}

	// Check if coupon is valid
	if !coupon.IsValid {
		return &domain.ValidatePromotionCodeResponse{
			Valid:  false,
			Reason: "associated coupon is not valid",
		}
	}

	// Check if promotion code is expired
	if promotionCode.ExpiresAt != nil && time.Now().UTC().After(*promotionCode.ExpiresAt) {
		return &domain.ValidatePromotionCodeResponse{
			Valid:  false,
			Reason: "promotion code has expired",
		}
	}

	// Check if coupon is expired
	if coupon.RedeemBy != nil && time.Now().UTC().After(*coupon.RedeemBy) {
		return &domain.ValidatePromotionCodeResponse{
			Valid:  false,
			Reason: "associated coupon has expired",
		}
	}

	// Check redemption limits for promotion code
	if promotionCode.MaxRedemptions != nil && promotionCode.TimesRedeemed >= *promotionCode.MaxRedemptions {
		return &domain.ValidatePromotionCodeResponse{
			Valid:  false,
			Reason: "promotion code has reached its redemption limit",
		}
	}

	// Check redemption limits for coupon
	if coupon.MaxRedemptions != nil && coupon.TimesRedeemed >= *coupon.MaxRedemptions {
		return &domain.ValidatePromotionCodeResponse{
			Valid:  false,
			Reason: "associated coupon has reached its redemption limit",
		}
	}

	// Check minimum amount requirements
	if promotionCode.MinimumAmount != nil && req.OrderAmount != nil {
		if promotionCode.MinimumAmountCurrency != nil && req.OrderCurrency != nil {
			if strings.ToUpper(*promotionCode.MinimumAmountCurrency) == strings.ToUpper(*req.OrderCurrency) {
				if *req.OrderAmount < *promotionCode.MinimumAmount {
					return &domain.ValidatePromotionCodeResponse{
						Valid:  false,
						Reason: fmt.Sprintf("order amount must be at least %d %s", *promotionCode.MinimumAmount, *promotionCode.MinimumAmountCurrency),
					}
				}
			}
		}
	}

	// Check currency restrictions
	if promotionCode.Restrictions != nil && len(promotionCode.Restrictions.CurrencyOptions) > 0 && req.OrderCurrency != nil {
		currencyAllowed := false
		orderCurrency := strings.ToUpper(*req.OrderCurrency)
		for _, allowedCurrency := range promotionCode.Restrictions.CurrencyOptions {
			if strings.ToUpper(allowedCurrency) == orderCurrency {
				currencyAllowed = true
				break
			}
		}
		if !currencyAllowed {
			return &domain.ValidatePromotionCodeResponse{
				Valid:  false,
				Reason: fmt.Sprintf("promotion code is not valid for currency %s", *req.OrderCurrency),
			}
		}
	}

	return &domain.ValidatePromotionCodeResponse{
		Valid: true,
	}
}

func (s *PromotionCodeService) buildPromotionCodeResponse(promotionCode *domain.PromotionCode) *domain.PromotionCodeResponse {
	var restrictions *domain.PromotionCodeRestrictionsResponse
	if promotionCode.Restrictions != nil {
		restrictions = &domain.PromotionCodeRestrictionsResponse{
			CurrencyOptions: promotionCode.Restrictions.CurrencyOptions,
		}
	}

	return &domain.PromotionCodeResponse{
		ID:                    promotionCode.ID.String(),
		StripePromotionID:     promotionCode.StripePromotionID,
		CouponID:              promotionCode.CouponID.String(),
		Code:                  promotionCode.Code,
		Active:                promotionCode.Active,
		MaxRedemptions:        promotionCode.MaxRedemptions,
		TimesRedeemed:         promotionCode.TimesRedeemed,
		ExpiresAt:             promotionCode.ExpiresAt,
		FirstTimeTransaction:  promotionCode.FirstTimeTransaction,
		MinimumAmount:         promotionCode.MinimumAmount,
		MinimumAmountCurrency: promotionCode.MinimumAmountCurrency,
		Restrictions:          restrictions,
		CreatedAt:             promotionCode.CreatedAt,
		UpdatedAt:             promotionCode.UpdatedAt,
		Metadata:              promotionCode.Metadata,
	}
}

func (s *PromotionCodeService) buildCouponResponse(coupon *domain.Coupon) *domain.CouponResponse {
	return &domain.CouponResponse{
		ID:               coupon.ID.String(),
		StripeCouponID:   coupon.StripeCouponID,
		Name:             coupon.Name,
		PercentOff:       coupon.PercentOff,
		AmountOff:        coupon.AmountOff,
		Currency:         coupon.Currency,
		Duration:         string(coupon.Duration),
		DurationInMonths: coupon.DurationInMonths,
		MaxRedemptions:   coupon.MaxRedemptions,
		TimesRedeemed:    coupon.TimesRedeemed,
		Valid:            coupon.IsValid,
		RedeemBy:         coupon.RedeemBy,
		CreatedAt:        coupon.CreatedAt,
		UpdatedAt:        coupon.UpdatedAt,
		Metadata:         coupon.Metadata,
	}
}
