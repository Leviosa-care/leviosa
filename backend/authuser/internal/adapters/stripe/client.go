package authPayment

import "github.com/stripe/stripe-go/v82"

// NewClient creates a new Stripe client
func NewClient(apiKey string) *stripe.Client {
	return stripe.NewClient(apiKey)
}
