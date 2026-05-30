package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
)

// StripeService defines operations for Stripe customer management
type StripeService interface {
	CreateCustomer(ctx context.Context, userID uuid.UUID, email, firstName, lastName string) (*stripe.Customer, error)
	RetrieveCustomer(ctx context.Context, customerID string) (*stripe.Customer, error)
	UpdateCustomer(ctx context.Context, customerID string, params *stripe.CustomerUpdateParams) (*stripe.Customer, error)
	DeleteCustomer(ctx context.Context, customerID string) (*stripe.Customer, error)
	FindCustomerByUserID(ctx context.Context, userID uuid.UUID) (*stripe.Customer, error)

	// CreateConnectedAccount creates a Stripe Connect Express account for a partner.
	// Returns the Stripe account ID on success.
	CreateConnectedAccount(ctx context.Context, userID uuid.UUID) (string, error)

	// CreateAccountLink creates a Stripe Account Link for onboarding.
	// Returns the URL the partner should be redirected to.
	CreateAccountLink(ctx context.Context, accountID, returnType, returnURL, refreshURL string) (string, error)

	// VerifyConnectWebhookSignature verifies the Stripe webhook signature using the
	// Connect-specific webhook secret and returns the parsed account.updated event data.
	// Returns the Stripe account ID and whether charges/payouts are enabled.
	VerifyConnectWebhookSignature(payload []byte, signature string) (accountID string, chargesEnabled bool, payoutsEnabled bool, err error)
}

