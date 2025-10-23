package domain

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

// CreatePromotionCodeRequest represents the data required to create a new promotion code.
// This is used for inbound API requests (from client to handler/service).
type CreatePromotionCodeRequest struct {
	CouponID              string                            `json:"couponId" validate:"required,uuid"`                                            // FK to Coupon
	Code                  string                            `json:"code" validate:"required,min=3,max=50"`                                        // The actual code customers enter
	MaxRedemptions        *int                              `json:"maxRedemptions,omitempty" validate:"omitempty,min=1"`                          // Max times this code can be used
	ExpiresAt             *time.Time                        `json:"expiresAt,omitempty"`                                                          // When this code expires
	FirstTimeTransaction  bool                              `json:"firstTimeTransaction"`                                                         // Whether code can only be used by new customers
	MinimumAmount         *int                              `json:"minimumAmount,omitempty" validate:"omitempty,min=1"`                           // Minimum order amount in cents
	MinimumAmountCurrency *string                           `json:"minimumAmountCurrency,omitempty" validate:"required_with=MinimumAmount,len=3"` // Currency for minimum amount
	Restrictions          *PromotionCodeRestrictionsRequest `json:"restrictions,omitempty"`                                                       // Additional restrictions
	Metadata              map[string]string                 `json:"metadata,omitempty"`                                                           // Additional metadata
}

func (r CreatePromotionCodeRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Validate CouponID
	if r.CouponID == "" {
		errs.Set("couponId", "coupon ID cannot be empty.")
	} else if err := uuid.Validate(r.CouponID); err != nil {
		errs.Set("couponId", "coupon ID must be a valid UUID.")
	}

	// Validate Code
	if r.Code == "" {
		errs.Set("code", "promotion code cannot be empty.")
	} else {
		r.Code = strings.ToUpper(strings.TrimSpace(r.Code))
		if len(r.Code) < 3 {
			errs.Set("code", "promotion code must be at least 3 characters long.")
		}
		if len(r.Code) > 50 {
			errs.Set("code", "promotion code cannot exceed 50 characters.")
		}
		// Check if code matches required format (uppercase alphanumeric with underscores/hyphens)
		for _, char := range r.Code {
			if !((char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_' || char == '-') {
				errs.Set("code", "promotion code can only contain uppercase letters, numbers, underscores, and hyphens.")
				break
			}
		}
	}

	// Validate MaxRedemptions
	if r.MaxRedemptions != nil && *r.MaxRedemptions <= 0 {
		errs.Set("maxRedemptions", "max redemptions must be greater than 0.")
	}

	// Validate MinimumAmount and Currency together
	if r.MinimumAmount != nil && r.MinimumAmountCurrency == nil {
		errs.Set("minimumAmountCurrency", "minimum amount currency is required when minimum amount is provided.")
	}
	if r.MinimumAmountCurrency != nil && r.MinimumAmount == nil {
		errs.Set("minimumAmount", "minimum amount is required when minimum amount currency is provided.")
	}
	if r.MinimumAmount != nil && *r.MinimumAmount <= 0 {
		errs.Set("minimumAmount", "minimum amount must be greater than 0.")
	}
	if r.MinimumAmountCurrency != nil && len(*r.MinimumAmountCurrency) != 3 {
		errs.Set("minimumAmountCurrency", "minimum amount currency must be a 3-letter ISO currency code.")
	}

	// Validate Restrictions
	if r.Restrictions != nil {
		if err := r.Restrictions.Valid(ctx); err != nil {
			errs.Set("restrictions", err.Error())
		}
	}

	return errs.AsError()
}

// PromotionCodeRestrictionsRequest represents restrictions for promotion code creation
type PromotionCodeRestrictionsRequest struct {
	CurrencyOptions []string `json:"currencyOptions,omitempty" validate:"dive,len=3"` // Currencies this code applies to
}

func (r PromotionCodeRestrictionsRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Validate currency options
	for i, currency := range r.CurrencyOptions {
		if len(currency) != 3 {
			errs.Set("currencyOptions", "all currency codes must be exactly 3 characters long.")
			break
		}
		// Convert to uppercase for consistency
		r.CurrencyOptions[i] = strings.ToUpper(currency)
	}

	return errs.AsError()
}

// UpdatePromotionCodeRequest represents the fields that can be updated for an existing promotion code.
// This is used for inbound API requests (from client to handler/service) for PATCH operations.
type UpdatePromotionCodeRequest struct {
	Active   *bool             `json:"active,omitempty"`   // Whether the promotion code is active
	Metadata map[string]string `json:"metadata,omitempty"` // Full map replacement if provided
}

func (r UpdatePromotionCodeRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// For update requests, we mainly validate that if fields are provided, they're reasonable
	// Most validation is optional since this is a PATCH operation

	// No specific validation needed for Active (it's just a boolean pointer)
	// No specific validation needed for Metadata (map[string]string is flexible)

	return errs.AsError()
}

// PromotionCodeResponse represents the promotion code data returned to clients.
// This is used for outbound API responses (from handler to client).
type PromotionCodeResponse struct {
	ID                    string                             `json:"id"`
	StripePromotionID     string                             `json:"stripePromotionId"`
	CouponID              string                             `json:"couponId"`
	Code                  string                             `json:"code"`
	Active                bool                               `json:"active"`
	MaxRedemptions        *int                               `json:"maxRedemptions,omitempty"`
	TimesRedeemed         int                                `json:"timesRedeemed"`
	ExpiresAt             *time.Time                         `json:"expiresAt,omitempty"`
	FirstTimeTransaction  bool                               `json:"firstTimeTransaction"`
	MinimumAmount         *int                               `json:"minimumAmount,omitempty"`
	MinimumAmountCurrency *string                            `json:"minimumAmountCurrency,omitempty"`
	Restrictions          *PromotionCodeRestrictionsResponse `json:"restrictions,omitempty"`
	CreatedAt             time.Time                          `json:"createdAt"`
	UpdatedAt             time.Time                          `json:"updatedAt"`
	Metadata              map[string]string                  `json:"metadata,omitempty"`
}

// PromotionCodeRestrictionsResponse represents restrictions in the response
type PromotionCodeRestrictionsResponse struct {
	CurrencyOptions []string `json:"currencyOptions,omitempty"`
}

// PromotionCodeWithCouponResponse represents a promotion code with its associated coupon details.
// This is useful for API responses that need to show both promotion code and coupon information.
type PromotionCodeWithCouponResponse struct {
	PromotionCode PromotionCodeResponse `json:"promotionCode"`
	Coupon        CouponResponse        `json:"coupon"`
}

// ValidatePromotionCodeRequest represents a request to validate if a promotion code can be used.
// This is used for checking code validity before applying it to a purchase.
type ValidatePromotionCodeRequest struct {
	Code          string  `json:"code" validate:"required,min=3,max=50"`                              // The promotion code to validate
	OrderAmount   *int    `json:"orderAmount,omitempty" validate:"omitempty,min=1"`                   // Order amount in cents
	OrderCurrency *string `json:"orderCurrency,omitempty" validate:"required_with=OrderAmount,len=3"` // Order currency
	CustomerID    *string `json:"customerId,omitempty"`                                               // Customer ID for first-time transaction checks
}

func (r ValidatePromotionCodeRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Validate Code
	if r.Code == "" {
		errs.Set("code", "promotion code cannot be empty.")
	} else {
		code := strings.ToUpper(strings.TrimSpace(r.Code))
		if len(code) < 3 {
			errs.Set("code", "promotion code must be at least 3 characters long.")
		}
		if len(code) > 50 {
			errs.Set("code", "promotion code cannot exceed 50 characters.")
		}
	}

	// Validate OrderAmount and Currency together
	if r.OrderAmount != nil && r.OrderCurrency == nil {
		errs.Set("orderCurrency", "order currency is required when order amount is provided.")
	}
	if r.OrderCurrency != nil && r.OrderAmount == nil {
		errs.Set("orderAmount", "order amount is required when order currency is provided.")
	}
	if r.OrderAmount != nil && *r.OrderAmount <= 0 {
		errs.Set("orderAmount", "order amount must be greater than 0.")
	}
	if r.OrderCurrency != nil && len(*r.OrderCurrency) != 3 {
		errs.Set("orderCurrency", "order currency must be a 3-letter ISO currency code.")
	}

	return errs.AsError()
}

// ValidatePromotionCodeResponse represents the response to a promotion code validation request.
type ValidatePromotionCodeResponse struct {
	Valid         bool                             `json:"valid"`                   // Whether the code is valid
	Reason        string                           `json:"reason,omitempty"`        // Reason why code is invalid (if applicable)
	PromotionCode *PromotionCodeWithCouponResponse `json:"promotionCode,omitempty"` // Full promotion code details if valid
}
