package stripeService

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/product"
)

// CreateProduct creates a Stripe product with a generated price using Stripe's Product API.
//
// Parameters:
//   - ctx: A context.Context instance to manage request lifecycle and cancellation.
//   - object: A Payment instance containing information about the product to be created.
//
// Returns:
//   - string: The ID of the newly created product in the Stripe system.
//   - error: An error if the product could not be created.
func (s *service) CreateProduct(ctx context.Context, object Payment) (string, error) {
	paymentInfo := object.GetPaymentInfo(ctx)
	product_params := &stripe.ProductParams{
		ID:          stripe.String(paymentInfo.ID),
		Name:        stripe.String(paymentInfo.Name),
		Description: stripe.String(paymentInfo.Description),
	}
	product, err := product.New(product_params)
	if err != nil {
		return "", domain.NewNotCreatedErr(err)
	}
	return product.ID, nil
}
