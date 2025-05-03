package stripeService

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/price"
)

// CreatePrice creates a new price for a given product in the Stripe payment system.
//
// Parameters:
//   - ctx: A context.Context instance to manage request lifecycle and cancellation.
//   - productID: A string representing the ID of the product for which the price is being created.
//   - priceValue: An int64 representing the value of the price in the smallest currency unit (e.g., cents for EUR).
//
// Returns:
//   - string: The ID of the newly created price in the Stripe system.
//   - error: An error if the price could not be created.
func (s *service) CreatePrice(ctx context.Context, productID string, priceValue int64) (string, error) {
	price_params := &stripe.PriceParams{
		Currency:   stripe.String(string(stripe.CurrencyEUR)),
		Product:    stripe.String(productID),
		UnitAmount: stripe.Int64(priceValue),
	}
	price, err := price.New(price_params)
	if err != nil {
		return "", domain.NewNotCreatedErr(err)
	}
	return price.ID, nil
}
