package authPayment

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
)

// CreateCustomer creates a new Stripe customer
func (s *service) CreateCustomer(ctx context.Context, userID uuid.UUID, email, firstName, lastName string) (*stripe.Customer, error) {
	params := &stripe.CustomerCreateParams{
		Email:       stripe.String(email),
		Name:        stripe.String(fmt.Sprintf("%s %s", firstName, lastName)),
		Description: stripe.String(fmt.Sprintf("User ID: %s", userID.String())),
		Metadata: map[string]string{
			"user_id": userID.String(),
		},
	}

	customer, err := s.client.V1Customers.Create(ctx, params)
	if err != nil {
		return nil, errs.ClassifyStripeError("create customer", err)
	}

	return customer, nil
}
