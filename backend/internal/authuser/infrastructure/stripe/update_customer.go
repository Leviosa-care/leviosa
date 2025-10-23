package authPayment

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/stripe/stripe-go/v82"
)

// UpdateCustomer updates a Stripe customer
func (s *service) UpdateCustomer(ctx context.Context, customerID string, params *stripe.CustomerUpdateParams) (*stripe.Customer, error) {
	customer, err := s.client.V1Customers.Update(ctx, customerID, params)
	if err != nil {
		return nil, errs.ClassifyStripeError("update customer", err)
	}

	return customer, nil
}
