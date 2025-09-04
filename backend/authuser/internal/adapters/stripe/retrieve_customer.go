package authPayment

import (
	"context"

	"github.com/Leviosa-care/core/errs"
	"github.com/stripe/stripe-go/v82"
)

// RetrieveCustomer retrieves a Stripe customer by ID
func (s *service) RetrieveCustomer(ctx context.Context, customerID string) (*stripe.Customer, error) {
	customer, err := s.client.V1Customers.Retrieve(ctx, customerID, &stripe.CustomerRetrieveParams{})
	if err != nil {
		return nil, errs.ClassifyStripeError("retrieve customer", err)
	}

	return customer, nil
}
