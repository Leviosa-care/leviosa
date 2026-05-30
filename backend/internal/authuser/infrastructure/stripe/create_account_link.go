package authPayment

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/stripe/stripe-go/v82"
)

// CreateAccountLink creates a Stripe Account Link for partner onboarding.
// The partner will be redirected to returnURL after completing onboarding,
// and to refreshURL if the link expires and needs regeneration.
func (s *service) CreateAccountLink(ctx context.Context, accountID, returnType, returnURL, refreshURL string) (string, error) {
	params := &stripe.AccountLinkCreateParams{
		Account:    stripe.String(accountID),
		ReturnURL:  stripe.String(returnURL),
		RefreshURL: stripe.String(refreshURL),
		Type:       stripe.String(returnType),
	}

	accountLink, err := s.client.V1AccountLinks.Create(ctx, params)
	if err != nil {
		return "", errs.ClassifyStripeError("create account link", err)
	}

	return accountLink.URL, nil
}
