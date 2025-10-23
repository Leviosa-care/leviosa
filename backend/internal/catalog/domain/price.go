package domain

import (
	"context"
	"fmt"
	"regexp" // For currency code validation
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type Price struct {
	ID            uuid.UUID
	ProductID     uuid.UUID `json:"productId"`     // FK to your internal product
	StripePriceID string    `json:"stripePriceId"` // Price ID from Stripe
	Amount        int       `json:"amount"`        // in cents
	Currency      string    `json:"currency"`      // "usd", "eur", etc.
	Interval      Interval  `json:"interval"`      // "month", "year", "one_time"
	IsActive      bool      `json:"isActive"`      // Track which one is current
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

var currencyRegex = regexp.MustCompile("^[A-Z]{3}$")

func (p Price) Valid(ctx context.Context) error {
	var errs errsx.Map
	// Amount validation
	if p.Amount <= 0 {
		errs.Set("amount", "must be a positive value")
	}

	// Currency validation (3-letter uppercase ISO code)
	if !currencyRegex.MatchString(p.Currency) {
		errs.Set("currency", "must be a 3-letter uppercase ISO currency code (e.g., 'USD', 'EUR')")
	}

	// Interval validation (must be one of the defined ENUM values)
	if !p.Interval.IsValid() {
		errs.Set("interval", fmt.Sprintf("Invalid interval type: '%s'. Must be 'month', 'year', or 'one_time'.", p.Interval))
	}

	return errs.AsError()
}

type Interval string

const (
	OneTime Interval = "one_time"
	Month   Interval = "month"
	Year    Interval = "year"
)

// IsValid checks if the AvailabilityType is one of the defined constants.
func (it Interval) IsValid() bool {
	switch it {
	case Month, Year, OneTime:
		return true
	default:
		return false
	}
}
