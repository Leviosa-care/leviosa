package domain

import "time"

// TODO: split the file into product, price and coupon eventually

type Metadata = map[string]string

type PaymentProduct struct {
	ID          string
	Name        string
	Description string
	Active      bool
	Metadata    Metadata
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PaymentPrice struct {
	ID        string
	Product   string // Stripe Product ID
	ProductID string
	Amount    int64  // in cents
	Currency  string // "usd", "eur", etc.
	Interval  string // "month", "year", "one-time"
	Active    bool
	Nickname  string
	Metadata  Metadata
	CreatedAt time.Time
}

type PaymentCoupon struct {
	ID         string
	PercentOff *float64 // percentage discount (0-100)
	AmountOff  *int64   // fixed amount discount in cents
	Duration   string   // "once", "repeating", "forever"
	Valid      bool
	Metadata   Metadata
	CreatedAt  time.Time
}

type PaymentPromotionCode struct {
	ID             string
	CouponID       string
	Code           string
	MaxRedemptions *int
	TimesRedeemed  int
	ExpiresAt      *time.Time
	Active         bool
	Metadata       Metadata
	CreatedAt      time.Time
}

// Request/Response structs
type CreateStripeProductRequest struct {
	Name        string
	Description string
	Metadata    Metadata
}

type UpdateStripeProductRequest struct {
	Name        *string
	Description *string
	Metadata    Metadata
}

type CreateStripePriceRequest struct {
	ProductID string
	Amount    int64
	Currency  string
	Interval  string
	Active    bool
	Nickname  string
	Metadata  Metadata
}

// UpdateStripePriceRequest is the input for updating a price in Stripe.
// Note: Only 'Active', 'Metadata', 'Nickname', 'LookupKey' can be updated on Stripe Prices.
type UpdateStripePriceRequest struct {
	Active   *bool    // Use pointer to distinguish between false and not provided
	Metadata Metadata // Will overwrite existing metadata if provided
	Nickname *string
	// LookupKey - if you use this, add it here too
}

type PriceListOptions struct {
	Active *bool // nil = all, true = active only, false = inactive only
	Limit  int64
}
