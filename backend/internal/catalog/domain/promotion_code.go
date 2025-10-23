package domain

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

// PromotionCode represents a customer-facing code that references a coupon
type PromotionCode struct {
	ID                    uuid.UUID                  `json:"id"`
	StripePromotionID     string                     `json:"stripePromotionId"`               // Promotion Code ID from Stripe
	CouponID              uuid.UUID                  `json:"couponId"`                        // FK to Coupon
	Code                  string                     `json:"code"`                            // The actual code customers enter
	Active                bool                       `json:"active"`                          // Whether the promotion code is active
	MaxRedemptions        *int                       `json:"maxRedemptions,omitempty"`        // Max times this code can be used
	TimesRedeemed         int                        `json:"timesRedeemed"`                   // Number of times this code has been used
	ExpiresAt             *time.Time                 `json:"expiresAt,omitempty"`             // When this code expires
	FirstTimeTransaction  bool                       `json:"firstTimeTransaction"`            // Whether code can only be used by new customers
	MinimumAmount         *int                       `json:"minimumAmount,omitempty"`         // Minimum order amount in cents
	MinimumAmountCurrency *string                    `json:"minimumAmountCurrency,omitempty"` // Currency for minimum amount
	Restrictions          *PromotionCodeRestrictions `json:"restrictions,omitempty"`          // Additional restrictions
	CreatedAt             time.Time                  `json:"createdAt"`
	UpdatedAt             time.Time                  `json:"updatedAt"`
	Metadata              map[string]string          `json:"metadata,omitempty"` // Additional metadata
}

// PromotionCodeRestrictions represents additional restrictions for promotion codes
type PromotionCodeRestrictions struct {
	// Stripe allows more complex restrictions, but these are the most common
	CurrencyOptions []string `json:"currencyOptions,omitempty"` // Currencies this code applies to
}

var promotionCodeRegex = regexp.MustCompile("^[A-Z0-9_-]+$") // Alphanumeric, underscore, hyphen only

// Valid validates the promotion code fields
func (pc PromotionCode) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Code validation
	if pc.Code == "" {
		errs.Set("code", "Promotion code cannot be empty")
	} else {
		// Normalize code to uppercase for validation
		normalizedCode := strings.ToUpper(pc.Code)
		if !promotionCodeRegex.MatchString(normalizedCode) {
			errs.Set("code", "Code must contain only uppercase letters, numbers, hyphens, and underscores")
		}
		if len(normalizedCode) < 3 {
			errs.Set("code", "Code must be at least 3 characters long")
		}
		if len(normalizedCode) > 50 {
			errs.Set("code", "Code cannot exceed 50 characters")
		}
	}

	// CouponID is required (FK constraint will be enforced at DB level)
	if pc.CouponID == uuid.Nil {
		errs.Set("couponId", "Coupon ID is required")
	}

	// Max redemptions validation
	if pc.MaxRedemptions != nil && *pc.MaxRedemptions <= 0 {
		errs.Set("maxRedemptions", "Max redemptions must be a positive value")
	}

	// Times redeemed validation
	if pc.TimesRedeemed < 0 {
		errs.Set("timesRedeemed", "Times redeemed cannot be negative")
	}

	// Check if max redemptions exceeded
	if pc.MaxRedemptions != nil && pc.TimesRedeemed > *pc.MaxRedemptions {
		errs.Set("timesRedeemed", "Times redeemed cannot exceed max redemptions")
	}

	// Expiry date validation
	if pc.ExpiresAt != nil && pc.ExpiresAt.Before(time.Now()) {
		errs.Set("expiresAt", "Expiry date must be in the future")
	}

	// Minimum amount validation
	if pc.MinimumAmount != nil {
		if *pc.MinimumAmount <= 0 {
			errs.Set("minimumAmount", "Minimum amount must be a positive value")
		}
		if pc.MinimumAmountCurrency == nil || *pc.MinimumAmountCurrency == "" {
			errs.Set("minimumAmountCurrency", "Currency is required when minimum amount is specified")
		} else if !currencyCodeRegex.MatchString(*pc.MinimumAmountCurrency) {
			errs.Set("minimumAmountCurrency", "Must be a 3-letter uppercase ISO currency code (e.g., 'USD', 'EUR')")
		}
	}

	// Restrictions validation
	if pc.Restrictions != nil {
		for i, currency := range pc.Restrictions.CurrencyOptions {
			if !currencyCodeRegex.MatchString(currency) {
				errs.Set("restrictions.currencyOptions", fmt.Sprintf("All currency codes must be 3-letter uppercase ISO codes (e.g., 'USD', 'EUR') at index %d", i))
			}
		}
	}

	return errs.AsError()
}

// IsExpired checks if the promotion code has expired
func (pc PromotionCode) IsExpired() bool {
	return pc.ExpiresAt != nil && pc.ExpiresAt.Before(time.Now())
}

// IsRedemptionLimitReached checks if the promotion code has reached its redemption limit
func (pc PromotionCode) IsRedemptionLimitReached() bool {
	return pc.MaxRedemptions != nil && pc.TimesRedeemed >= *pc.MaxRedemptions
}

// CanBeUsed checks if the promotion code can currently be used
func (pc PromotionCode) CanBeUsed() bool {
	return pc.Active && !pc.IsExpired() && !pc.IsRedemptionLimitReached()
}
