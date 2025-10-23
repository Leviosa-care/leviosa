package productPayment_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeactivateProduct(t *testing.T) {
	ctx := context.Background()

	t.Run("successful product deactivation", func(t *testing.T) {
		// Create an active product
		createReq := domain.CreateStripeProductRequest{
			Name:        "Product to Deactivate",
			Description: "This product will be deactivated",
		}

		createdProduct, err := stripeService.CreateProduct(ctx, createReq)
		require.NoError(t, err)
		require.True(t, createdProduct.Active)

		// Deactivate the product
		err = stripeService.DeactivateProduct(ctx, createdProduct.ID)

		require.NoError(t, err)

		// Verify the product is deactivated
		retrievedProduct, err := stripeService.GetProduct(ctx, createdProduct.ID)
		require.NoError(t, err)
		assert.False(t, retrievedProduct.Active)
	})

	t.Run("deactivation fails with non-existent product ID", func(t *testing.T) {
		nonExistentID := "prod_nonexistent123456789"

		err := stripeService.DeactivateProduct(ctx, nonExistentID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to deactivate product")
	})

	t.Run("deactivation of already inactive product succeeds", func(t *testing.T) {
		// Create a product
		createReq := domain.CreateStripeProductRequest{
			Name: "Product to Deactivate Twice",
		}

		createdProduct, err := stripeService.CreateProduct(ctx, createReq)
		require.NoError(t, err)

		// Deactivate once
		err = stripeService.DeactivateProduct(ctx, createdProduct.ID)
		require.NoError(t, err)

		// Deactivate again - should succeed
		err = stripeService.DeactivateProduct(ctx, createdProduct.ID)

		assert.NoError(t, err)

		// Verify still inactive
		retrievedProduct, err := stripeService.GetProduct(ctx, createdProduct.ID)
		require.NoError(t, err)
		assert.False(t, retrievedProduct.Active)
	})
}

func TestReactivateProduct(t *testing.T) {
	ctx := context.Background()

	t.Run("successful product reactivation", func(t *testing.T) {
		// Create a product
		createReq := domain.CreateStripeProductRequest{
			Name:        "Product to Reactivate",
			Description: "This product will be deactivated then reactivated",
		}

		createdProduct, err := stripeService.CreateProduct(ctx, createReq)
		require.NoError(t, err)

		// Deactivate first
		err = stripeService.DeactivateProduct(ctx, createdProduct.ID)
		require.NoError(t, err)

		// Verify it's deactivated
		deactivatedProduct, err := stripeService.GetProduct(ctx, createdProduct.ID)
		require.NoError(t, err)
		require.False(t, deactivatedProduct.Active)

		// Reactivate the product
		err = stripeService.ReactivateProduct(ctx, createdProduct.ID)

		require.NoError(t, err)

		// Verify the product is reactivated
		retrievedProduct, err := stripeService.GetProduct(ctx, createdProduct.ID)
		require.NoError(t, err)
		assert.True(t, retrievedProduct.Active)
	})

	t.Run("reactivation fails with non-existent product ID", func(t *testing.T) {
		nonExistentID := "prod_nonexistent123456789"

		err := stripeService.ReactivateProduct(ctx, nonExistentID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to reactivate product")
	})

	t.Run("reactivation of already active product succeeds", func(t *testing.T) {
		// Create a product (already active by default)
		createReq := domain.CreateStripeProductRequest{
			Name: "Product to Reactivate Twice",
		}

		createdProduct, err := stripeService.CreateProduct(ctx, createReq)
		require.NoError(t, err)
		require.True(t, createdProduct.Active)

		// Reactivate already active product - should succeed
		err = stripeService.ReactivateProduct(ctx, createdProduct.ID)

		assert.NoError(t, err)

		// Verify still active
		retrievedProduct, err := stripeService.GetProduct(ctx, createdProduct.ID)
		require.NoError(t, err)
		assert.True(t, retrievedProduct.Active)
	})
}

func TestProductActivationWorkflow(t *testing.T) {
	ctx := context.Background()

	t.Run("full activation workflow", func(t *testing.T) {
		// Create product
		createReq := domain.CreateStripeProductRequest{
			Name:        "Workflow Test Product",
			Description: "Testing activation workflow",
		}

		createdProduct, err := stripeService.CreateProduct(ctx, createReq)
		require.NoError(t, err)
		assert.True(t, createdProduct.Active)

		// Deactivate
		err = stripeService.DeactivateProduct(ctx, createdProduct.ID)
		require.NoError(t, err)

		deactivatedProduct, err := stripeService.GetProduct(ctx, createdProduct.ID)
		require.NoError(t, err)
		assert.False(t, deactivatedProduct.Active)

		// Reactivate
		err = stripeService.ReactivateProduct(ctx, createdProduct.ID)
		require.NoError(t, err)

		reactivatedProduct, err := stripeService.GetProduct(ctx, createdProduct.ID)
		require.NoError(t, err)
		assert.True(t, reactivatedProduct.Active)

		// Deactivate again
		err = stripeService.DeactivateProduct(ctx, createdProduct.ID)
		require.NoError(t, err)

		finalProduct, err := stripeService.GetProduct(ctx, createdProduct.ID)
		require.NoError(t, err)
		assert.False(t, finalProduct.Active)
	})
}