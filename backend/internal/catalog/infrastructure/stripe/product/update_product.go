package productPayment

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stripe/stripe-go/v82"
)

func (c *service) UpdateProduct(ctx context.Context, productID string, req *domain.UpdateStripeProductRequest) (*domain.PaymentProduct, error) {
	params := &stripe.ProductUpdateParams{}

	// Apply updates only if provided
	if req.Name != nil {
		params.Name = stripe.String(*req.Name)
	}
	if req.Description != nil {
		params.Description = stripe.String(*req.Description)
	}
	if req.Metadata != nil {
		params.Metadata = req.Metadata
	}

	stripeProduct, err := c.V1Products.Update(ctx, productID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update product %s: %w", productID, err) // ← Fixed
	}

	var createdAt time.Time
	if stripeProduct.Created > 0 {
		t := time.Unix(stripeProduct.Created, 0)
		createdAt = t
	}

	var updatedAt time.Time
	if stripeProduct.Updated > 0 {
		t := time.Unix(stripeProduct.Updated, 0)
		createdAt = t
	}
	return &domain.PaymentProduct{
		ID:          stripeProduct.ID,
		Name:        stripeProduct.Name,
		Description: stripeProduct.Description,
		Active:      stripeProduct.Active,
		Metadata:    stripeProduct.Metadata,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}
