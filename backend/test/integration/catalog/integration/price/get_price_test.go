package price_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPrice_Integration(t *testing.T) {
	ctx := context.Background()

	t.Run("should get price successfully with valid ID", func(t *testing.T) {
		// Setup
		productID := setupTestProduct(t, ctx)

		// Create a price first
		createRequest := td.NewValidCreatePriceRequest()
		createResp, createBody := createPriceViaAPI(t, productID, createRequest)
		assertSuccessResponse(t, createResp, 201)

		var createdPriceID string
		err := json.Unmarshal(createBody, &createdPriceID)
		require.NoError(t, err)

		// Execute
		resp, body := getPriceViaAPI(t, createdPriceID)

		// Assert
		assertSuccessResponse(t, resp, 200)

		var retrievedPrice domain.Price
		err = json.Unmarshal(body, &retrievedPrice)
		require.NoError(t, err)

		// Verify response structure
		assert.Equal(t, createdPriceID, retrievedPrice.ID.String())
		assert.Equal(t, createRequest.Amount, retrievedPrice.Amount)
		assert.Equal(t, createRequest.Currency, retrievedPrice.Currency)
		assert.Equal(t, createRequest.Interval, string(retrievedPrice.Interval))
		assert.True(t, retrievedPrice.IsActive) // Should be active by default
		assert.NotZero(t, retrievedPrice.CreatedAt)
		assert.NotZero(t, retrievedPrice.UpdatedAt)
	})

	t.Run("should return 400 when price ID is missing", func(t *testing.T) {
		// Execute - using empty price ID results in 404 (route not found)
		resp, _ := getPriceViaAPI(t, "")

		// Assert
		assertErrorResponse(t, resp, 404)
	})

	t.Run("should return 404 when price does not exist", func(t *testing.T) {
		// Setup
		td.ClearPricesTable(t, ctx, testPool)
		nonExistentPriceID := "550e8400-e29b-41d4-a716-446655440000"

		// Execute
		resp, _ := getPriceViaAPI(t, nonExistentPriceID)

		// Assert
		assertErrorResponse(t, resp, 404)
	})

	t.Run("should return 400 with invalid price ID format", func(t *testing.T) {
		// Execute - using invalid UUID format
		resp, _ := getPriceViaAPI(t, "invalid-uuid-format")

		// Assert
		assertErrorResponse(t, resp, 400)
	})

	t.Run("should get price from database correctly", func(t *testing.T) {
		// Setup - Insert price directly to database
		productID := setupTestProduct(t, ctx)

		// Create price via API first to have Stripe integration
		createRequest := td.NewValidCreatePriceRequest()
		createRequest.Amount = 1500
		createRequest.Currency = "EUR"
		createRequest.Interval = "year"
		createRequest.Nickname = td.StrPtr("Annual Plan")

		createResp, createBody := createPriceViaAPI(t, productID, createRequest)
		assertSuccessResponse(t, createResp, 201)

		var createdPriceID string
		err := json.Unmarshal(createBody, &createdPriceID)
		require.NoError(t, err)

		// Execute - Get the price
		resp, body := getPriceViaAPI(t, createdPriceID)

		// Assert
		assertSuccessResponse(t, resp, 200)

		var retrievedPrice domain.Price
		err = json.Unmarshal(body, &retrievedPrice)
		require.NoError(t, err)

		// Verify specific fields
		assert.Equal(t, 1500, retrievedPrice.Amount)
		assert.Equal(t, "EUR", retrievedPrice.Currency)
		assert.Equal(t, "year", retrievedPrice.Interval)
		assert.True(t, retrievedPrice.IsActive)

		// Verify against database record
		dbPrice := td.GetPriceByID(t, ctx, createdPrice.ID, testPool)
		assert.Equal(t, retrievedPrice.Amount, dbPrice.Amount)
		assert.Equal(t, retrievedPrice.Currency, dbPrice.Currency)
		assert.Equal(t, string(retrievedPrice.Interval), string(dbPrice.Interval))
		assert.Equal(t, retrievedPrice.IsActive, dbPrice.IsActive)
	})

	t.Run("should get inactive price correctly", func(t *testing.T) {
		// Setup
		productID := setupTestProduct(t, ctx)

		// Create a price
		createRequest := td.NewValidCreatePriceRequest()
		createResp, createBody := createPriceViaAPI(t, productID, createRequest)
		assertSuccessResponse(t, createResp, 201)

		var createdPriceID string
		err := json.Unmarshal(createBody, &createdPriceID)
		require.NoError(t, err)

		// Deactivate the price
		updateRequest := &domain.UpdatePriceRequest{
			Active: &[]bool{false}[0], // Create pointer to false
		}
		updateResp, _ := updatePriceViaAPI(t, createdPriceID, updateRequest)
		assertSuccessResponse(t, updateResp, 200)

		// Execute - Get the deactivated price
		resp, body := getPriceViaAPI(t, createdPriceID)

		// Assert
		assertSuccessResponse(t, resp, 200)

		var retrievedPrice domain.Price
		err = json.Unmarshal(body, &retrievedPrice)
		require.NoError(t, err)

		// Verify it's inactive
		assert.False(t, retrievedPrice.IsActive)
		assert.Equal(t, createdPriceID, retrievedPrice.ID.String())
	})
}
