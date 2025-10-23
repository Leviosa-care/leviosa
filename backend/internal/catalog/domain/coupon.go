package domain

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

// Coupon represents a discount coupon that can be applied to subscriptions or invoices
type Coupon struct {
	ID               uuid.UUID         `json:"id"`
	StripeCouponID   string            `json:"stripeCouponId"`             // Coupon ID from Stripe
	Name             string            `json:"name"`                       // Human-readable name
	PercentOff       *float64          `json:"percentOff,omitempty"`       // Percentage discount (e.g., 25.5 for 25.5%)
	AmountOff        *int              `json:"amountOff,omitempty"`        // Fixed amount off in cents
	Currency         *string           `json:"currency,omitempty"`         // Currency for amount_off (required if amount_off is set)
	Duration         CouponDuration    `json:"duration"`                   // How long the discount applies
	DurationInMonths *int              `json:"durationInMonths,omitempty"` // Required if duration is "repeating"
	MaxRedemptions   *int              `json:"maxRedemptions,omitempty"`   // Maximum number of times this coupon can be redeemed
	TimesRedeemed    int               `json:"timesRedeemed"`              // Number of times redeemed so far
	IsValid          bool              `json:"isValid"`                    // Whether the coupon is still valid
	RedeemBy         *time.Time        `json:"redeemBy,omitempty"`         // Last time coupon can be redeemed
	CreatedAt        time.Time         `json:"createdAt"`
	UpdatedAt        time.Time         `json:"updatedAt"`
	Metadata         map[string]string `json:"metadata,omitempty"` // Additional metadata
}

// CouponDuration represents how long a coupon discount applies
type CouponDuration string

const (
	CouponDurationOnce      CouponDuration = "once"      // Applies to one charge only
	CouponDurationRepeating CouponDuration = "repeating" // Applies to multiple charges for a specified number of months
	CouponDurationForever   CouponDuration = "forever"   // Applies to all future charges
)

// IsValid checks if the CouponDuration is one of the defined constants
func (cd CouponDuration) IsValid() bool {
	switch cd {
	case CouponDurationOnce, CouponDurationRepeating, CouponDurationForever:
		return true
	default:
		return false
	}
}

var currencyCodeRegex = regexp.MustCompile("^[A-Z]{3}$")

// Valid validates the coupon fields
func (c Coupon) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Name is required
	if c.Name == "" {
		errs.Set("name", "Coupon name cannot be empty")
	}

	// Must have either percent_off or amount_off, but not both
	if c.PercentOff == nil && c.AmountOff == nil {
		errs.Set("discount", "Must specify either percentOff or amountOff")
	}
	if c.PercentOff != nil && c.AmountOff != nil {
		errs.Set("discount", "Cannot specify both percentOff and amountOff")
	}

	// Validate percent_off range
	if c.PercentOff != nil {
		if *c.PercentOff <= 0 || *c.PercentOff > 100 {
			errs.Set("percentOff", "Must be between 0 and 100")
		}
	}

	// Validate amount_off and currency
	if c.AmountOff != nil {
		if *c.AmountOff <= 0 {
			errs.Set("amountOff", "Must be a positive value")
		}
		if c.Currency == nil || *c.Currency == "" {
			errs.Set("currency", "Currency is required when amountOff is specified")
		} else if !currencyCodeRegex.MatchString(*c.Currency) {
			errs.Set("currency", "Must be a 3-letter uppercase ISO currency code (e.g., 'USD', 'EUR')")
		}
	}

	// Validate duration
	if !c.Duration.IsValid() {
		errs.Set("duration", fmt.Sprintf("Invalid duration: '%s'. Must be 'once', 'repeating', or 'forever'", c.Duration))
	}

	// Validate duration_in_months
	if c.Duration == CouponDurationRepeating {
		if c.DurationInMonths == nil || *c.DurationInMonths <= 0 {
			errs.Set("durationInMonths", "Duration in months is required and must be positive when duration is 'repeating'")
		}
	} else if c.DurationInMonths != nil {
		errs.Set("durationInMonths", "Duration in months should only be specified when duration is 'repeating'")
	}

	// Validate max_redemptions
	if c.MaxRedemptions != nil && *c.MaxRedemptions <= 0 {
		errs.Set("maxRedemptions", "Max redemptions must be a positive value")
	}

	// Validate times_redeemed
	if c.TimesRedeemed < 0 {
		errs.Set("timesRedeemed", "Times redeemed cannot be negative")
	}

	// Check if max redemptions exceeded
	if c.MaxRedemptions != nil && c.TimesRedeemed > *c.MaxRedemptions {
		errs.Set("timesRedeemed", "Times redeemed cannot exceed max redemptions")
	}

	// Validate redeem_by is in the future (if set)
	if c.RedeemBy != nil && c.RedeemBy.Before(time.Now()) {
		errs.Set("redeemBy", "Redeem by date must be in the future")
	}

	return errs.AsError()
}

// ToResponse converts a domain.Coupon to a domain.CouponResponse for API responses
func (c *Coupon) ToResponse() *CouponResponse {
	return &CouponResponse{
		ID:               c.ID.String(),
		StripeCouponID:   c.StripeCouponID,
		Name:             c.Name,
		PercentOff:       c.PercentOff,
		AmountOff:        c.AmountOff,
		Currency:         c.Currency,
		Duration:         string(c.Duration),
		DurationInMonths: c.DurationInMonths,
		MaxRedemptions:   c.MaxRedemptions,
		TimesRedeemed:    c.TimesRedeemed,
		Valid:            c.IsValid,
		RedeemBy:         c.RedeemBy,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
		Metadata:         c.Metadata,
	}
}
