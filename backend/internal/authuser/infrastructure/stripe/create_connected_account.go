package authPayment

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
)

// CreateConnectedAccount creates a Stripe Connect Express account for a partner.
// The account is created with the partner's user ID in metadata for traceability.
// Partners will complete onboarding separately via Stripe Account Links.
func (s *service) CreateConnectedAccount(ctx context.Context, userID uuid.UUID) (string, error) {
	params := &stripe.AccountCreateParams{
		Type: stripe.String("express"),
		Metadata: map[string]string{
			"user_id": userID.String(),
		},
	}

	account, err := s.client.V1Accounts.Create(ctx, params)
	if err != nil {
		return "", errs.ClassifyStripeError("create connected account", err)
	}

	return account.ID, nil
}
