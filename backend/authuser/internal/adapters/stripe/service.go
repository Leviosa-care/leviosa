package authPayment

import (
	"github.com/Leviosa-care/authuser/internal/ports"
	"github.com/stripe/stripe-go/v82"
)

// service handles Stripe operations
type service struct {
	client *stripe.Client
}

// Compile-time check to ensure Service implements ports.StripeService
var _ ports.StripeService = (*service)(nil)

// NewService creates a new Stripe service
func NewService(apiKey, baseURL string) *service {
	var sc *stripe.Client

	if baseURL != "" {
		backends := stripe.NewBackendsWithConfig(&stripe.BackendConfig{
			URL: &baseURL,
		})
		sc = stripe.NewClient(apiKey, stripe.WithBackends(backends))
	} else {
		// Set default backend for production if not using a custom base URL
		sc = stripe.NewClient(apiKey)
	}

	return &service{sc}
}
