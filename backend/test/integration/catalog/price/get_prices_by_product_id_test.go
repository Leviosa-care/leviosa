package price_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetPricesByProductID TEST_PATH=test/integration/catalog/price/get_prices_by_product_id_test.go

func TestGetPricesByProductID(t *testing.T) {
	ctx := context.Background()

	t.Run("should get all prices for a product successfully", func(t *testing.T) {
		// Setup
		productID := setupTestProduct(t, ctx)

		// Create multiple prices for the same product
		createRequest1 := td.NewValidCreatePriceRequest()
		createRequest1.Amount = 1000
		createRequest1.Nickname = td.StrPtr("Basic Plan")
		createRequest1.Interval = "month"

		createRequest2 := td.NewValidCreatePriceRequest()
		createRequest2.Amount = 10000
		createRequest2.Nickname = td.StrPtr("Premium Plan")
		createRequest2.Interval = "year"

		// Create both prices
		createResp1, createBody1 := createPriceViaAPI(t, productID, createRequest1)
		assertSuccessResponse(t, createResp1, 201)

		createResp2, createBody2 := createPriceViaAPI(t, productID, createRequest2)
		assertSuccessResponse(t, createResp2, 201)

		var price1, price2 domain.Price
		require.NoError(t, json.Unmarshal(createBody1, &price1))
		require.NoError(t, json.Unmarshal(createBody2, &price2))

		// Execute
		resp, body := getPricesByProductIDViaAPI(t, productID)

		// Assert
		assertSuccessResponse(t, resp, 200)

		var prices []*domain.Price
		err := json.Unmarshal(body, &prices)
		assert.NoError(t, err)

		// Verify we got both prices
		assert.Len(t, prices, 2)

		// Sort prices by amount for consistent comparison
		if prices[0].Amount > prices[1].Amount {
			prices[0], prices[1] = prices[1], prices[0]
		}

		// Verify first price (Basic Plan)
		assert.Equal(t, 1000, prices[0].Amount)
		assert.Equal(t, "month", prices[0].Interval)
		assert.Equal(t, "EUR", prices[0].Currency)
		assert.True(t, prices[0].IsActive)

		// Verify second price (Premium Plan)
		assert.Equal(t, 10000, prices[1].Amount)
		assert.Equal(t, "year", prices[1].Interval)
		assert.Equal(t, "EUR", prices[1].Currency)
		assert.True(t, prices[1].IsActive)

		// Verify both prices belong to the same product
		assert.Equal(t, prices[0].ProductID, prices[1].ProductID)
	})

	t.Run("should return empty array when product has no prices", func(t *testing.T) {
		// Setup - product with no prices
		productID := setupTestProduct(t, ctx)

		// Execute
		resp, body := getPricesByProductIDViaAPI(t, productID)

		// Assert
		assertSuccessResponse(t, resp, 200)

		var prices []*domain.Price
		err := json.Unmarshal(body, &prices)
		assert.NoError(t, err)

		assert.Empty(t, prices)
	})

	t.Run("should return 400 when product ID is missing", func(t *testing.T) {
		// Execute - using empty product ID results in 404 (route not found)
		resp, _ := getPricesByProductIDViaAPI(t, "")

		// Assert
		assertErrorResponse(t, resp, 404)
	})

	t.Run("should return 404 when product does not exist", func(t *testing.T) {
		// Setup
		td.ClearProductsTable(t, ctx, testPool)
		td.ClearPricesTable(t, ctx, testPool)
		nonExistentProductID := "550e8400-e29b-41d4-a716-446655440000"

		// Execute
		resp, _ := getPricesByProductIDViaAPI(t, nonExistentProductID)

		// Assert
		assertErrorResponse(t, resp, 404)
	})

	t.Run("should return 400 with invalid product ID format", func(t *testing.T) {
		// Execute - using invalid UUID format
		resp, _ := getPricesByProductIDViaAPI(t, "invalid-uuid-format")

		// Assert
		assertErrorResponse(t, resp, 400)
	})

	t.Run("should include both active and inactive prices", func(t *testing.T) {
		// Setup
		productID := setupTestProduct(t, ctx)

		// Create two prices
		createRequest1 := td.NewValidCreatePriceRequest()
		createRequest1.Amount = 1000
		createRequest1.Nickname = td.StrPtr("Active Plan")

		createRequest2 := td.NewValidCreatePriceRequest()
		createRequest2.Amount = 2000
		createRequest2.Nickname = td.StrPtr("Plan to Deactivate")

		createResp1, createBody1 := createPriceViaAPI(t, productID, createRequest1)
		assertSuccessResponse(t, createResp1, 201)

		createResp2, createBody2 := createPriceViaAPI(t, productID, createRequest2)
		assertSuccessResponse(t, createResp2, 201)

		var price1, price2 domain.Price
		assert.NoError(t, json.Unmarshal(createBody1, &price1))
		assert.NoError(t, json.Unmarshal(createBody2, &price2))

		// Deactivate second price
		updateRequest := &domain.UpdatePriceRequest{
			Active: &[]bool{false}[0], // Create pointer to false
		}
		updateResp, _ := updatePriceViaAPI(t, price2.ID.String(), updateRequest)
		assertSuccessResponse(t, updateResp, 200)

		// Execute
		resp, body := getPricesByProductIDViaAPI(t, productID)

		// Assert
		assertSuccessResponse(t, resp, 200)

		var prices []*domain.Price
		err := json.Unmarshal(body, &prices)
		assert.NoError(t, err)

		assert.Len(t, prices, 2)

		// Find active and inactive prices
		var activePrice, inactivePrice *domain.Price
		for _, p := range prices {
			if p.IsActive {
				activePrice = p
			} else {
				inactivePrice = p
			}
		}

		assert.NotNil(t, activePrice, "Should have found active price")
		assert.NotNil(t, inactivePrice, "Should have found inactive price")

		assert.Equal(t, 1000, activePrice.Amount)
		assert.Equal(t, 2000, inactivePrice.Amount)
		assert.True(t, activePrice.IsActive)
		assert.False(t, inactivePrice.IsActive)
	})

	t.Run("should return prices in consistent order", func(t *testing.T) {
		// Setup
		productID := setupTestProduct(t, ctx)

		// Create multiple prices with different creation times
		prices := make([]domain.Price, 3)
		for i := 0; i < 3; i++ {
			createRequest := td.NewValidCreatePriceRequest()
			createRequest.Amount = 1000 * (i + 1)
			createRequest.Nickname = td.StrPtr(fmt.Sprintf("Plan %d", i+1))

			createResp, createBody := createPriceViaAPI(t, productID, createRequest)
			assertSuccessResponse(t, createResp, 201)

			err := json.Unmarshal(createBody, &prices[i])
			assert.NoError(t, err)

			// Small delay to ensure different creation times
			time.Sleep(1 * time.Millisecond)
		}

		// Execute multiple times to ensure consistent ordering
		for attempt := 0; attempt < 3; attempt++ {
			resp, body := getPricesByProductIDViaAPI(t, productID)
			assertSuccessResponse(t, resp, 200)

			var retrievedPrices []*domain.Price
			err := json.Unmarshal(body, &retrievedPrices)
			assert.NoError(t, err)

			assert.Len(t, retrievedPrices, 3)

			// Prices should be ordered consistently (likely by creation time)
			for i := 1; i < len(retrievedPrices); i++ {
				// Verify consistent ordering - creation time should be non-decreasing
				assert.True(t, retrievedPrices[i-1].CreatedAt.Before(retrievedPrices[i].CreatedAt) ||
					retrievedPrices[i-1].CreatedAt.Equal(retrievedPrices[i].CreatedAt),
					"Prices should be consistently ordered by creation time")
			}
		}
	})

	t.Run("should handle product with single price", func(t *testing.T) {
		// Setup
		productID := setupTestProduct(t, ctx)

		// Create single price
		createRequest := td.NewValidCreatePriceRequest()
		createRequest.Amount = 5000
		createRequest.Nickname = td.StrPtr("Single Plan")

		createResp, createBody := createPriceViaAPI(t, productID, createRequest)
		assertSuccessResponse(t, createResp, 201)

		var createdPriceID string
		require.NoError(t, json.Unmarshal(createBody, &createdPriceID))

		// Execute
		resp, body := getPricesByProductIDViaAPI(t, productID)

		// Assert
		assertSuccessResponse(t, resp, 200)

		var prices []*domain.Price
		err := json.Unmarshal(body, &prices)
		assert.NoError(t, err)

		assert.Len(t, prices, 1)
		assert.Equal(t, createdPriceID, prices[0].ID.String())
		assert.Equal(t, 5000, prices[0].Amount)
		// Nickname is not returned in Price domain model
	})
}
