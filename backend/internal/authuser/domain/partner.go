package domain

import (
	"time"

	"github.com/google/uuid"
)

type Partner struct {
	UserID           uuid.UUID   `json:"user_id"`
	Bio              string      `json:"bio" encx:"encrypt"`
	Experience       string      `json:"experience" encx:"encrypt"`
	Certifications   []string    `json:"certifications" encx:"encrypt"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	CategoryIDs      []uuid.UUID `json:"category_ids" encx:"encrypt"`
	ProductIDs       []uuid.UUID `json:"product_ids" encx:"encrypt"`

	// Stripe Connect fields for Option 2
	StripeConnectedAccountID   string              `json:"stripe_connected_account_id" encx:"encrypt"`
	StripeAccountStatus        StripeAccountStatus `json:"stripe_account_status"`
	StripeOnboardingComplete   bool                `json:"stripe_onboarding_complete"`
}