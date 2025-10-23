package productPayment_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProduct(t *testing.T) {
	ctx := context.Background()

	t.Run("successful product retrieval", func(t *testing.T) {
		// First create a product
		createReq := domain.CreateStripeProductRequest{
			Name:        "Test Product for Get",
			Description: "Product to test retrieval",
			Metadata: map[string]string{
				"test": "get_product",
			},
		}

		createdProduct, err := stripeService.CreateProduct(ctx, createReq)
		require.NoError(t, err)
		require.NotNil(t, createdProduct)

		// Now retrieve the product
		retrievedProduct, err := stripeService.GetProduct(ctx, createdProduct.ID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedProduct)
		assert.Equal(t, createdProduct.ID, retrievedProduct.ID)
		assert.Equal(t, createdProduct.Name, retrievedProduct.Name)
		assert.Equal(t, createdProduct.Description, retrievedProduct.Description)
		assert.Equal(t, createdProduct.Active, retrievedProduct.Active)
		assert.Equal(t, createdProduct.Metadata, retrievedProduct.Metadata)
		assert.False(t, retrievedProduct.CreatedAt.IsZero())
	})

	t.Run("retrieval fails with empty product ID", func(t *testing.T) {
		product, err := stripeService.GetProduct(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Contains(t, err.Error(), "stripeProductID cannot be empty")
	})

	t.Run("retrieval fails with non-existent product ID", func(t *testing.T) {
		nonExistentID := "prod_nonexistent123456789"

		product, err := stripeService.GetProduct(ctx, nonExistentID)

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Contains(t, err.Error(), "failed to retrieve product")
	})

	t.Run("retrieve product with minimal data", func(t *testing.T) {
		// Create a minimal product
		createReq := domain.CreateStripeProductRequest{
			Name: "Minimal Product for Get",
		}

		createdProduct, err := stripeService.CreateProduct(ctx, createReq)
		require.NoError(t, err)

		// Retrieve it
		retrievedProduct, err := stripeService.GetProduct(ctx, createdProduct.ID)

		require.NoError(t, err)
		assert.NotNil(t, retrievedProduct)
		assert.Equal(t, createdProduct.ID, retrievedProduct.ID)
		assert.Equal(t, createReq.Name, retrievedProduct.Name)
		assert.True(t, retrievedProduct.Active)
	})
}