package productPayment

import (
	"context"
	"fmt"

	"github.com/stripe/stripe-go/v82"
)

func (c *service) ReactivateProduct(ctx context.Context, stripeProductID string) error {
	params := &stripe.ProductUpdateParams{
		Active: stripe.Bool(true),
	}
	_, err := c.V1Products.Update(ctx, stripeProductID, params)
	if err != nil {
		return fmt.Errorf("failed to reactivate product with ID %s: %w", stripeProductID, err)
	}
	return nil
}
