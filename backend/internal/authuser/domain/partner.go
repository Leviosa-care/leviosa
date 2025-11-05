package domain

import (
	"time"

	"github.com/google/uuid"
)

type Partner struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"-"`
	Bio        string    `json:"bio"`
	Experience string    `json:"experience"`
	// Certifications []string    `json:"certifications" encx:"encrypt"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	CategoryIDs []uuid.UUID `json:"category_ids"`
	ProductIDs  []uuid.UUID `json:"product_ids"`

	// Stripe Connect fields for Option 2
	StripeConnectedAccountID string              `json:"stripe_connected_account_id" encx:"encrypt"`
	StripeAccountStatus      StripeAccountStatus `json:"stripe_account_status"`
	StripeOnboardingComplete bool                `json:"stripe_onboarding_complete"`
}
