package productPayment

import (
	"context"
	"fmt"

	"github.com/stripe/stripe-go/v82"
)

func (c *service) DeactivateProduct(ctx context.Context, stripeProductID string) error {
	params := &stripe.ProductUpdateParams{
		Active: stripe.Bool(false),
	}
	_, err := c.V1Products.Update(ctx, stripeProductID, params)
	if err != nil {
		return fmt.Errorf("failed to deactivate product %s: %w", stripeProductID, err)
	}
	return nil
}
