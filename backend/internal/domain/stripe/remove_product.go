package stripeService

import (
	"context"
	"fmt"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/product"
)

func (s *service) RemoveProduct(ctx context.Context, productID string) (string, error) {
	product_params := &stripe.ProductParams{
		Active: stripe.Bool(false),
		ID:     &productID,
	}
	// NOTE: that thing does not work if the product has a price associated with it
	_, err := product.Del(productID, product_params)
	if err != nil {
		return "", fmt.Errorf("Failed to delete product with ID %s: %w", productID, err)
	}
	return productID, nil
}
