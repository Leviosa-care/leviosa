package productPayment

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stripe/stripe-go/v82"
)

func (c *service) CreateProduct(ctx context.Context, req domain.CreateStripeProductRequest) (*domain.PaymentProduct, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	// Include more fields in idempotency key or use UUID
	idempotencyKey := fmt.Sprintf("prod_%s_%s_%d",
		strings.ReplaceAll(req.Name, " ", "_"),
		strings.ReplaceAll(req.Description, " ", "_"),
		time.Now().Unix(), // Or use a hash of metadata
	)

	params := &stripe.ProductCreateParams{
		Name:        stripe.String(req.Name),
		Description: stripe.String(req.Description),
		Active:      stripe.Bool(true),
		Metadata:    req.Metadata,
	}
	params.SetIdempotencyKey(idempotencyKey)

	stripeProduct, err := c.V1Products.Create(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("stripe product creation failed: %w", err)
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
