package pricePayment

import (
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
	"github.com/stripe/stripe-go/v82"
)

type service struct {
	*stripe.Client
}

func NewPrice(apiKey, baseURL string) ports.PricePaymentGateway {
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
