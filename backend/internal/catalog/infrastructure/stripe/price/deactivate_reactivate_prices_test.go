package pricePayment_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeactivatePrices(t *testing.T) {
	ctx := context.Background()

	createMockProduct := func() string {
		return "prod_test123456789"
	}

	t.Run("successful single price deactivation", func(t *testing.T) {
		// Create an active price
		productID := createMockProduct()
		createReq := domain.CreateStripePriceRequest{
			ProductID: productID,
			Amount:    1000,
			Currency:  "USD",
			Interval:  "month",
			Active:    true,
			Nickname:  "Price to Deactivate",
		}

		createdPrice, err := stripeService.CreatePrice(ctx, createReq)
		require.NoError(t, err)
		require.True(t, createdPrice.Active)

		// Deactivate the price
		err = stripeService.DeactivatePrices(ctx, []string{createdPrice.ID})

		require.NoError(t, err)

		// Verify the price is deactivated
		retrievedPrice, err := stripeService.GetPrice(ctx, createdPrice.ID)
		require.NoError(t, err)
		assert.False(t, retrievedPrice.Active)
	})

	t.Run("successful multiple prices deactivation", func(t *testing.T) {
		// Create multiple active prices
		productID := createMockProduct()
		var priceIDs []string

		for i := 0; i < 3; i++ {
			createReq := domain.CreateStripePriceRequest{
				ProductID: productID,
				Amount:    int64(1000 + (i * 100)), // Different amounts
				Currency:  "USD",
				Interval:  "month",
				Active:    true,
			}

			createdPrice, err := stripeService.CreatePrice(ctx, createReq)
			require.NoError(t, err)
			priceIDs = append(priceIDs, createdPrice.ID)
		}

		// Deactivate all prices
		err := stripeService.DeactivatePrices(ctx, priceIDs)

		require.NoError(t, err)

		// Verify all prices are deactivated
		for _, priceID := range priceIDs {
			retrievedPrice, err := stripeService.GetPrice(ctx, priceID)
			require.NoError(t, err)
			assert.False(t, retrievedPrice.Active, "Price %s should be deactivated", priceID)
		}
	})

	t.Run("deactivation fails with non-existent price ID", func(t *testing.T) {
		nonExistentID := "price_nonexistent123456789"

		err := stripeService.DeactivatePrices(ctx, []string{nonExistentID})

		assert.Error(t, err)
	})

	t.Run("deactivation with empty slice succeeds", func(t *testing.T) {
		err := stripeService.DeactivatePrices(ctx, []string{})

		assert.NoError(t, err)
	})
}

func TestReactivatePrices(t *testing.T) {
	ctx := context.Background()

	createMockProduct := func() string {
		return "prod_test123456789"
	}

	t.Run("successful single price reactivation", func(t *testing.T) {
		// Create a price
		productID := createMockProduct()
		createReq := domain.CreateStripePriceRequest{
			ProductID: productID,
			Amount:    1200,
			Currency:  "USD",
			Interval:  "month",
			Active:    true,
			Nickname:  "Price to Reactivate",
		}

		createdPrice, err := stripeService.CreatePrice(ctx, createReq)
		require.NoError(t, err)

		// Deactivate first
		err = stripeService.DeactivatePrices(ctx, []string{createdPrice.ID})
		require.NoError(t, err)

		// Verify it's deactivated
		deactivatedPrice, err := stripeService.GetPrice(ctx, createdPrice.ID)
		require.NoError(t, err)
		require.False(t, deactivatedPrice.Active)

		// Reactivate the price
		err = stripeService.ReactivatePrices(ctx, []string{createdPrice.ID})

		require.NoError(t, err)

		// Verify the price is reactivated
		retrievedPrice, err := stripeService.GetPrice(ctx, createdPrice.ID)
		require.NoError(t, err)
		assert.True(t, retrievedPrice.Active)
	})

	t.Run("successful multiple prices reactivation", func(t *testing.T) {
		// Create multiple prices
		productID := createMockProduct()
		var priceIDs []string

		for i := 0; i < 3; i++ {
			createReq := domain.CreateStripePriceRequest{
				ProductID: productID,
				Amount:    int64(1500 + (i * 200)), // Different amounts
				Currency:  "USD",
				Interval:  "month",
				Active:    true,
			}

			createdPrice, err := stripeService.CreatePrice(ctx, createReq)
			require.NoError(t, err)
			priceIDs = append(priceIDs, createdPrice.ID)
		}

		// Deactivate all first
		err := stripeService.DeactivatePrices(ctx, priceIDs)
		require.NoError(t, err)

		// Reactivate all prices
		err = stripeService.ReactivatePrices(ctx, priceIDs)

		require.NoError(t, err)

		// Verify all prices are reactivated
		for _, priceID := range priceIDs {
			retrievedPrice, err := stripeService.GetPrice(ctx, priceID)
			require.NoError(t, err)
			assert.True(t, retrievedPrice.Active, "Price %s should be reactivated", priceID)
		}
	})

	t.Run("reactivation fails with non-existent price ID", func(t *testing.T) {
		nonExistentID := "price_nonexistent123456789"

		err := stripeService.ReactivatePrices(ctx, []string{nonExistentID})

		assert.Error(t, err)
	})

	t.Run("reactivation with empty slice succeeds", func(t *testing.T) {
		err := stripeService.ReactivatePrices(ctx, []string{})

		assert.NoError(t, err)
	})
}

func TestPriceActivationWorkflow(t *testing.T) {
	ctx := context.Background()

	createMockProduct := func() string {
		return "prod_test123456789"
	}

	t.Run("full activation workflow with multiple prices", func(t *testing.T) {
		// Create multiple prices
		productID := createMockProduct()
		var priceIDs []string

		for i := 0; i < 2; i++ {
			createReq := domain.CreateStripePriceRequest{
				ProductID: productID,
				Amount:    int64(1000 + (i * 500)),
				Currency:  "USD",
				Interval:  "month",
				Active:    true,
			}

			createdPrice, err := stripeService.CreatePrice(ctx, createReq)
			require.NoError(t, err)
			assert.True(t, createdPrice.Active)
			priceIDs = append(priceIDs, createdPrice.ID)
		}

		// Deactivate all
		err := stripeService.DeactivatePrices(ctx, priceIDs)
		require.NoError(t, err)

		// Verify deactivation
		for _, priceID := range priceIDs {
			price, err := stripeService.GetPrice(ctx, priceID)
			require.NoError(t, err)
			assert.False(t, price.Active)
		}

		// Reactivate all
		err = stripeService.ReactivatePrices(ctx, priceIDs)
		require.NoError(t, err)

		// Verify reactivation
		for _, priceID := range priceIDs {
			price, err := stripeService.GetPrice(ctx, priceID)
			require.NoError(t, err)
			assert.True(t, price.Active)
		}

		// Deactivate again
		err = stripeService.DeactivatePrices(ctx, priceIDs)
		require.NoError(t, err)

		// Final verification
		for _, priceID := range priceIDs {
			price, err := stripeService.GetPrice(ctx, priceID)
			require.NoError(t, err)
			assert.False(t, price.Active)
		}
	})
}

