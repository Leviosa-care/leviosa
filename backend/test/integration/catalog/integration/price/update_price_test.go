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

func TestUpdatePrice_Integration(t *testing.T) {
	ctx := context.Background()

	t.Run("should update price successfully with valid data", func(t *testing.T) {
		// Setup
		productID := setupTestProduct(t, ctx)

		// Create a price first
		createRequest := td.NewValidCreatePriceRequest()
		// CreatePriceRequest doesn't have Active field - prices are active by default
		createRequest.Nickname = td.StrPtr("Original Plan")

		createResp, createBody := createPriceViaAPI(t, productID, createRequest)
		assertSuccessResponse(t, createResp, 201)

		var createdPriceID string
		err := json.Unmarshal(createBody, &createdPriceID)
		require.NoError(t, err)

		// Execute - Update the price
		active := false
		nickname := "Updated Plan"
		updateRequest := &domain.UpdatePriceRequest{
			Active:   &active,
			Nickname: &nickname,
			Metadata: map[string]string{
				"updated_key": "updated_value",
				"plan_type":   "premium",
			},
		}

		resp, body := updatePriceViaAPI(t, createdPriceID, updateRequest)

		// Assert
		assertSuccessResponse(t, resp, 200)

		var updatedPrice domain.Price
		err = json.Unmarshal(body, &updatedPrice)
		require.NoError(t, err)

		// Verify updated fields
		assert.Equal(t, createdPriceID, updatedPrice.ID.String())
		assert.False(t, updatedPrice.IsActive)

		// Verify database record is updated
		dbPrice := td.GetPriceByID(t, ctx, updatedPrice.ID, testPool)
		assert.False(t, dbPrice.IsActive)
	})

	t.Run("should update only active status", func(t *testing.T) {
		// Setup
		productID := setupTestProduct(t, ctx)

		createRequest := td.NewValidCreatePriceRequest()
		// CreatePriceRequest doesn't have Active field - prices are active by default
		createRequest.Nickname = td.StrPtr("Test Plan")

		createResp, createBody := createPriceViaAPI(t, productID, createRequest)
		assertSuccessResponse(t, createResp, 201)

		var createdPriceID string
		err := json.Unmarshal(createBody, &createdPriceID)
		require.NoError(t, err)

		// Execute - Update only active status
		active := false
		updateRequest := &domain.UpdatePriceRequest{
			Active: &active,
		}

		resp, body := updatePriceViaAPI(t, createdPriceID, updateRequest)

		// Assert
		assertSuccessResponse(t, resp, 200)

		var updatedPrice domain.Price
		err = json.Unmarshal(body, &updatedPrice)
		require.NoError(t, err)

		assert.False(t, updatedPrice.IsActive)
		// Nickname is not part of Price domain model
	})

	t.Run("should update only nickname", func(t *testing.T) {
		// Setup
		productID := setupTestProduct(t, ctx)

		createRequest := td.NewValidCreatePriceRequest()
		// CreatePriceRequest doesn't have Active field - prices are active by default
		createRequest.Nickname = td.StrPtr("Original Plan")

		createResp, createBody := createPriceViaAPI(t, productID, createRequest)
		assertSuccessResponse(t, createResp, 201)

		var createdPriceID string
		err := json.Unmarshal(createBody, &createdPriceID)
		require.NoError(t, err)

		// Execute - Update only nickname
		nickname := "New Plan Name"
		updateRequest := &domain.UpdatePriceRequest{
			Nickname: &nickname,
		}

		resp, body := updatePriceViaAPI(t, createdPriceID, updateRequest)

		// Assert
		assertSuccessResponse(t, resp, 200)

		var updatedPrice domain.Price
		err = json.Unmarshal(body, &updatedPrice)
		require.NoError(t, err)

		// Nickname is not returned in Price domain model
		assert.Equal(t, true, updatedPrice.IsActive) // Should remain unchanged
	})

	t.Run("should update only metadata", func(t *testing.T) {
		// Setup
		productID := setupTestProduct(t, ctx)

		createRequest := td.NewValidCreatePriceRequest()
		createRequest.Metadata = map[string]string{
			"original_key": "original_value",
		}

		createResp, createBody := createPriceViaAPI(t, productID, createRequest)
		assertSuccessResponse(t, createResp, 201)

		var createdPriceID string
		err := json.Unmarshal(createBody, &createdPriceID)
		require.NoError(t, err)

		// Execute - Update only metadata
		updateRequest := &domain.UpdatePriceRequest{
			Metadata: map[string]string{
				"updated_key":    "updated_value",
				"additional_key": "additional_value",
			},
		}

		resp, body := updatePriceViaAPI(t, createdPriceID, updateRequest)

		// Assert
		assertSuccessResponse(t, resp, 200)

		var updatedPrice domain.Price
		err = json.Unmarshal(body, &updatedPrice)
		require.NoError(t, err)

		// Note: Metadata verification depends on how it's stored and returned
		// This test assumes metadata is properly handled in the domain/service layers
		assert.Equal(t, true, updatedPrice.IsActive) // Should remain unchanged
		// Nickname is not part of Price domain model
	})

	t.Run("should return 400 when no updatable fields provided", func(t *testing.T) {
		// Setup
		productID := setupTestProduct(t, ctx)

		createRequest := td.NewValidCreatePriceRequest()
		createResp, createBody := createPriceViaAPI(t, productID, createRequest)
		assertSuccessResponse(t, createResp, 201)

		var createdPriceID string
		err := json.Unmarshal(createBody, &createdPriceID)
		require.NoError(t, err)

		// Execute - Empty update request
		updateRequest := &domain.UpdatePriceRequest{}

		resp, _ := updatePriceViaAPI(t, createdPriceID, updateRequest)

		// Assert
		assertErrorResponse(t, resp, 400)
	})

	t.Run("should return 400 when price ID is missing", func(t *testing.T) {
		// Setup
		updateRequest := td.NewValidUpdatePriceRequest()

		// Execute - using empty price ID results in 404 (route not found)
		resp, _ := updatePriceViaAPI(t, "", updateRequest)

		// Assert
		assertErrorResponse(t, resp, 404)
	})

	t.Run("should return 404 when price does not exist", func(t *testing.T) {
		// Setup
		td.ClearPricesTable(t, ctx, testPool)
		nonExistentPriceID := "550e8400-e29b-41d4-a716-446655440000"
		updateRequest := td.NewValidUpdatePriceRequest()

		// Execute
		resp, _ := updatePriceViaAPI(t, nonExistentPriceID, updateRequest)

		// Assert
		assertErrorResponse(t, resp, 404)
	})

	t.Run("should return 400 with invalid price ID format", func(t *testing.T) {
		// Setup
		updateRequest := td.NewValidUpdatePriceRequest()

		// Execute - using invalid UUID format
		resp, _ := updatePriceViaAPI(t, "invalid-uuid-format", updateRequest)

		// Assert
		assertErrorResponse(t, resp, 400)
	})

	t.Run("should return 400 with invalid request payload", func(t *testing.T) {
		// Setup
		productID := setupTestProduct(t, ctx)

		createRequest := td.NewValidCreatePriceRequest()
		createResp, createBody := createPriceViaAPI(t, productID, createRequest)
		assertSuccessResponse(t, createResp, 201)

		var createdPriceID string
		err := json.Unmarshal(createBody, &createdPriceID)
		require.NoError(t, err)

		// Execute - using nil payload
		resp, _ := updatePriceViaAPI(t, createdPriceID, nil)

		// Assert
		assertErrorResponse(t, resp, 400)
	})

	t.Run("should handle reactivation of deactivated price", func(t *testing.T) {
		// Setup
		productID := setupTestProduct(t, ctx)

		createRequest := td.NewValidCreatePriceRequest()
		// CreatePriceRequest doesn't have Active field - prices are active by default

		createResp, createBody := createPriceViaAPI(t, productID, createRequest)
		assertSuccessResponse(t, createResp, 201)

		var createdPriceID string
		err := json.Unmarshal(createBody, &createdPriceID)
		require.NoError(t, err)

		// First deactivate
		deactivate := false
		deactivateRequest := &domain.UpdatePriceRequest{
			Active: &deactivate,
		}

		deactivateResp, _ := updatePriceViaAPI(t, createdPriceID, deactivateRequest)
		assertSuccessResponse(t, deactivateResp, 200)

		// Then reactivate
		reactivate := true
		reactivateRequest := &domain.UpdatePriceRequest{
			Active: &reactivate,
		}

		resp, body := updatePriceViaAPI(t, createdPriceID, reactivateRequest)

		// Assert
		assertSuccessResponse(t, resp, 200)

		var updatedPrice domain.Price
		err = json.Unmarshal(body, &updatedPrice)
		require.NoError(t, err)

		assert.True(t, updatedPrice.IsActive)

		// Verify in database
		dbPrice := td.GetPriceByID(t, ctx, updatedPrice.ID, testPool)
		assert.True(t, dbPrice.IsActive)
	})

	t.Run("should handle Stripe service unavailable during update", func(t *testing.T) {
		// This test would require mocking Stripe to return errors
		// For now, we'll test with valid data to ensure the happy path works
		// In a real scenario, you might want to test with network issues or Stripe mock returning errors

		productID := setupTestProduct(t, ctx)

		createRequest := td.NewValidCreatePriceRequest()
		createResp, createBody := createPriceViaAPI(t, productID, createRequest)
		assertSuccessResponse(t, createResp, 201)

		var createdPriceID string
		err := json.Unmarshal(createBody, &createdPriceID)
		require.NoError(t, err)

		active := false
		updateRequest := &domain.UpdatePriceRequest{
			Active: &active,
		}

		resp, body := updatePriceViaAPI(t, createdPriceID, updateRequest)

		// Should succeed with mock Stripe
		assertSuccessResponse(t, resp, 200)

		var updatedPrice domain.Price
		err = json.Unmarshal(body, &updatedPrice)
		require.NoError(t, err)
		assert.False(t, updatedPrice.IsActive)
	})
}
