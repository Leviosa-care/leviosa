package authPayment

import (
	"context"

	"github.com/Leviosa-care/core/errs"
	"github.com/stripe/stripe-go/v82"
)

// DeleteCustomer deletes a Stripe customer
func (s *service) DeleteCustomer(ctx context.Context, customerID string) (*stripe.Customer, error) {
	customer, err := s.client.V1Customers.Delete(ctx, customerID, &stripe.CustomerDeleteParams{})
	if err != nil {
		return nil, errs.ClassifyStripeError("delete customer", err)
	}

	return customer, nil
}
