package productPayment

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stripe/stripe-go/v82"
)

func (c *service) GetProduct(ctx context.Context, productID string) (*domain.PaymentProduct, error) {
	if productID == "" {
		return nil, fmt.Errorf("stripeProductID cannot be empty")
	}

	params := &stripe.ProductRetrieveParams{}
	stripeProduct, err := c.V1Products.Retrieve(ctx, productID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve product %s: %w", productID, err) // ← Fixed
	}

	var createdAt time.Time
	if stripeProduct.Created > 0 {
		t := time.Unix(stripeProduct.Created, 0)
		createdAt = t

	}

	return &domain.PaymentProduct{
		ID:          stripeProduct.ID,
		Name:        stripeProduct.Name,
		Description: stripeProduct.Description,
		Active:      stripeProduct.Active,
		Metadata:    stripeProduct.Metadata,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}, nil
}
