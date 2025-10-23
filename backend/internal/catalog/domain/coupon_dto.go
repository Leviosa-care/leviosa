package domain

import (
	"context"
	"time"

	"github.com/hengadev/errsx"
)

// CreateCouponRequest represents the data required to create a new coupon.
// This is used for inbound API requests (from client to handler/service).
type CreateCouponRequest struct {
	Name             string            `json:"name" validate:"required"`                                                             // Human-readable name
	PercentOff       *float64          `json:"percentOff,omitempty" validate:"omitempty,min=0.1,max=100"`                            // Percentage discount
	AmountOff        *int              `json:"amountOff,omitempty" validate:"omitempty,min=1"`                                       // Fixed amount off in cents
	Currency         *string           `json:"currency,omitempty" validate:"required_with=AmountOff,len=3"`                          // Currency for amount_off
	Duration         string            `json:"duration" validate:"required,oneof=once repeating forever"`                            // How long the discount applies
	DurationInMonths *int              `json:"durationInMonths,omitempty" validate:"required_if=Duration repeating,omitempty,min=1"` // Required if duration is "repeating"
	MaxRedemptions   *int              `json:"maxRedemptions,omitempty" validate:"omitempty,min=1"`                                  // Maximum number of times this coupon can be redeemed
	RedeemBy         *time.Time        `json:"redeemBy,omitempty"`                                                                   // Last time coupon can be redeemed
	Metadata         map[string]string `json:"metadata,omitempty"`                                                                   // Additional metadata
}

func (r CreateCouponRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	if r.Name == "" {
		errs.Set("name", "name is required")
	}

	if r.Duration == "" {
		errs.Set("duration", "duration is required")
	} else {
		switch r.Duration {
		case "once", "repeating", "forever":
		default:
			errs.Set("duration", "duration must be one of: once, repeating, forever")
		}
	}

	if r.Duration == "repeating" && r.DurationInMonths == nil {
		errs.Set("durationInMonths", "durationInMonths is required when duration is 'repeating'")
	}

	if r.PercentOff == nil && r.AmountOff == nil {
		errs.Set("discount", "either percentOff or amountOff must be provided")
	}

	if r.PercentOff != nil && r.AmountOff != nil {
		errs.Set("discount", "cannot provide both percentOff and amountOff")
	}

	if r.PercentOff != nil {
		if *r.PercentOff <= 0 || *r.PercentOff > 100 {
			errs.Set("percentOff", "percentOff must be between 0.1 and 100")
		}
	}

	if r.AmountOff != nil {
		if *r.AmountOff < 1 {
			errs.Set("amountOff", "amountOff must be at least 1")
		}
		if r.Currency == nil {
			errs.Set("currency", "currency is required when amountOff is provided")
		} else if len(*r.Currency) != 3 {
			errs.Set("currency", "currency must be a 3-character ISO code")
		}
	}

	if r.MaxRedemptions != nil && *r.MaxRedemptions < 1 {
		errs.Set("maxRedemptions", "maxRedemptions must be at least 1")
	}

	if r.DurationInMonths != nil && *r.DurationInMonths < 1 {
		errs.Set("durationInMonths", "durationInMonths must be at least 1")
	}

	return errs.AsError()
}

// UpdateCouponRequest represents the fields that can be updated for an existing coupon.
// This is used for inbound API requests (from client to handler/service) for PATCH operations.
type UpdateCouponRequest struct {
	Name     *string           `json:"name,omitempty" validate:"omitempty,min=1"` // Human-readable name
	Metadata map[string]string `json:"metadata,omitempty"`                        // Full map replacement if provided
}

func (r UpdateCouponRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	if r.Name != nil && *r.Name == "" {
		errs.Set("name", "name cannot be empty if provided")
	}

	return errs.AsError()
}

// CouponResponse represents the coupon data returned to clients.
// This is used for outbound API responses (from handler to client).
type CouponResponse struct {
	ID               string            `json:"id"`
	StripeCouponID   string            `json:"stripeCouponId"`
	Name             string            `json:"name"`
	PercentOff       *float64          `json:"percentOff,omitempty"`
	AmountOff        *int              `json:"amountOff,omitempty"`
	Currency         *string           `json:"currency,omitempty"`
	Duration         string            `json:"duration"`
	DurationInMonths *int              `json:"durationInMonths,omitempty"`
	MaxRedemptions   *int              `json:"maxRedemptions,omitempty"`
	TimesRedeemed    int               `json:"timesRedeemed"`
	Valid            bool              `json:"valid"`
	RedeemBy         *time.Time        `json:"redeemBy,omitempty"`
	CreatedAt        time.Time         `json:"createdAt"`
	UpdatedAt        time.Time         `json:"updatedAt"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}
