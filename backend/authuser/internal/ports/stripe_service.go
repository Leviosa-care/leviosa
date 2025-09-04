package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
)

// StripeService defines operations for Stripe customer management
type StripeService interface {
	CreateCustomer(ctx context.Context, userID uuid.UUID, email, firstName, lastName string) (*stripe.Customer, error)
}

