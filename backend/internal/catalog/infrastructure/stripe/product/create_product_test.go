package productPayment_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateProduct(t *testing.T) {
	ctx := context.Background()

	t.Run("successful product creation", func(t *testing.T) {
		req := domain.CreateStripeProductRequest{
			Name:        "Test Product",
			Description: "Test product description",
			Metadata: map[string]string{
				"catalog_id": "test-catalog-123",
				"category":   "electronics",
			},
		}

		product, err := stripeService.CreateProduct(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, product)
		assert.NotEmpty(t, product.ID)
		assert.Equal(t, req.Name, product.Name)
		assert.Equal(t, req.Description, product.Description)
		assert.True(t, product.Active)
		assert.Equal(t, req.Metadata, product.Metadata)
		assert.False(t, product.CreatedAt.IsZero())
	})

	t.Run("creation with minimal required fields", func(t *testing.T) {
		req := domain.CreateStripeProductRequest{
			Name: "Minimal Product",
		}

		product, err := stripeService.CreateProduct(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, product)
		assert.NotEmpty(t, product.ID)
		assert.Equal(t, req.Name, product.Name)
		assert.True(t, product.Active)
		assert.False(t, product.CreatedAt.IsZero())
	})

	t.Run("creation fails with empty name", func(t *testing.T) {
		req := domain.CreateStripeProductRequest{
			Name:        "",
			Description: "Product without name",
		}

		product, err := stripeService.CreateProduct(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Contains(t, err.Error(), "name cannot be empty")
	})

	t.Run("creation with metadata", func(t *testing.T) {
		metadata := map[string]string{
			"internal_id": "internal-123",
			"category":    "books",
			"brand":       "test-brand",
		}

		req := domain.CreateStripeProductRequest{
			Name:        "Product with Metadata",
			Description: "Product with custom metadata",
			Metadata:    metadata,
		}

		product, err := stripeService.CreateProduct(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, metadata, product.Metadata)
	})
}