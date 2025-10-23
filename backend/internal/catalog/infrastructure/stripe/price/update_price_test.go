package pricePayment_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdatePrice(t *testing.T) {
	ctx := context.Background()

	createMockProduct := func() string {
		return "prod_test123456789"
	}

	t.Run("successful price update with all fields", func(t *testing.T) {
		// Create a price first
		productID := createMockProduct()
		createReq := domain.CreateStripePriceRequest{
			ProductID: productID,
			Amount:    1000,
			Currency:  "USD",
			Interval:  "month",
			Active:    true,
			Nickname:  "Original Price",
			Metadata: map[string]string{
				"version": "1.0",
			},
		}

		createdPrice, err := stripeService.CreatePrice(ctx, createReq)
		require.NoError(t, err)

		// Update the price
		newNickname := "Updated Price"
		newActive := false
		newMetadata := map[string]string{
			"version": "2.0",
			"updated": "true",
		}

		updateReq := domain.UpdateStripePriceRequest{
			Active:   &newActive,
			Nickname: &newNickname,
			Metadata: newMetadata,
		}

		updatedPrice, err := stripeService.UpdatePrice(ctx, createdPrice.ID, updateReq)

		require.NoError(t, err)
		assert.NotNil(t, updatedPrice)
		assert.Equal(t, createdPrice.ID, updatedPrice.ID)
		assert.Equal(t, newNickname, updatedPrice.Nickname)
		assert.Equal(t, newActive, updatedPrice.Active)
		assert.Equal(t, newMetadata, updatedPrice.Metadata)
	})

	t.Run("successful partial update - active only", func(t *testing.T) {
		// Create a price
		productID := createMockProduct()
		createReq := domain.CreateStripePriceRequest{
			ProductID: productID,
			Amount:    1500,
			Currency:  "USD",
			Interval:  "month",
			Active:    true,
			Nickname:  "Test Price",
		}

		createdPrice, err := stripeService.CreatePrice(ctx, createReq)
		require.NoError(t, err)

		// Update only the active status
		newActive := false
		updateReq := domain.UpdateStripePriceRequest{
			Active: &newActive,
		}

		updatedPrice, err := stripeService.UpdatePrice(ctx, createdPrice.ID, updateReq)

		require.NoError(t, err)
		assert.False(t, updatedPrice.Active)
		assert.Equal(t, createdPrice.Nickname, updatedPrice.Nickname) // Should remain unchanged
	})

	t.Run("successful partial update - nickname only", func(t *testing.T) {
		// Create a price
		productID := createMockProduct()
		createReq := domain.CreateStripePriceRequest{
			ProductID: productID,
			Amount:    2000,
			Currency:  "USD",
			Interval:  "month",
			Active:    true,
			Nickname:  "Original Nickname",
		}

		createdPrice, err := stripeService.CreatePrice(ctx, createReq)
		require.NoError(t, err)

		// Update only the nickname
		newNickname := "Updated Nickname Only"
		updateReq := domain.UpdateStripePriceRequest{
			Nickname: &newNickname,
		}

		updatedPrice, err := stripeService.UpdatePrice(ctx, createdPrice.ID, updateReq)

		require.NoError(t, err)
		assert.Equal(t, newNickname, updatedPrice.Nickname)
		assert.Equal(t, createdPrice.Active, updatedPrice.Active) // Should remain unchanged
	})

	t.Run("successful partial update - metadata only", func(t *testing.T) {
		// Create a price
		productID := createMockProduct()
		createReq := domain.CreateStripePriceRequest{
			ProductID: productID,
			Amount:    1200,
			Currency:  "USD",
			Interval:  "month",
			Active:    true,
			Metadata: map[string]string{
				"original": "value",
			},
		}

		createdPrice, err := stripeService.CreatePrice(ctx, createReq)
		require.NoError(t, err)

		// Update only the metadata
		newMetadata := map[string]string{
			"updated":  "metadata",
			"category": "premium",
		}
		updateReq := domain.UpdateStripePriceRequest{
			Metadata: newMetadata,
		}

		updatedPrice, err := stripeService.UpdatePrice(ctx, createdPrice.ID, updateReq)

		require.NoError(t, err)
		assert.Equal(t, newMetadata, updatedPrice.Metadata)
		assert.Equal(t, createdPrice.Active, updatedPrice.Active)
		assert.Equal(t, createdPrice.Nickname, updatedPrice.Nickname)
	})

	t.Run("update fails with non-existent price ID", func(t *testing.T) {
		nonExistentID := "price_nonexistent123456789"
		newActive := false

		updateReq := domain.UpdateStripePriceRequest{
			Active: &newActive,
		}

		updatedPrice, err := stripeService.UpdatePrice(ctx, nonExistentID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, updatedPrice)
	})

	t.Run("activate deactivated price", func(t *testing.T) {
		// Create an inactive price
		productID := createMockProduct()
		createReq := domain.CreateStripePriceRequest{
			ProductID: productID,
			Amount:    800,
			Currency:  "USD",
			Interval:  "month",
			Active:    false,
		}

		createdPrice, err := stripeService.CreatePrice(ctx, createReq)
		require.NoError(t, err)
		require.False(t, createdPrice.Active)

		// Activate it
		newActive := true
		updateReq := domain.UpdateStripePriceRequest{
			Active: &newActive,
		}

		updatedPrice, err := stripeService.UpdatePrice(ctx, createdPrice.ID, updateReq)

		require.NoError(t, err)
		assert.True(t, updatedPrice.Active)
	})
}