package building_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBuilding(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Clear test data before each test
	clearBuildingsTable(t, ctx)

	t.Run("should successfully create a building", func(t *testing.T) {
		// Prepare request
		request := domain.CreateBuildingRequest{
			Name:       "Test Building",
			Address:    "123 Test Street",
			City:       "Test City",
			PostalCode: "12345",
			Country:    "Test Country",
		}

		reqBody, err := json.Marshal(request)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/buildings", bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Execute request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response domain.BuildingResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response data
		assert.NotNil(t, response.ID)
		assert.Equal(t, "Test Building", response.Name)
		assert.Equal(t, "123 Test Street", response.Address)
		assert.Equal(t, "Test City", response.City)
		assert.Equal(t, "12345", response.PostalCode)
		assert.Equal(t, "Test Country", response.Country)
		assert.True(t, response.IsActive)
		assert.NotZero(t, response.CreatedAt)
		assert.NotZero(t, response.UpdatedAt)
	})

	t.Run("should return validation error for invalid data", func(t *testing.T) {
		// Prepare request with missing required fields
		request := domain.CreateBuildingRequest{
			Name: "", // Missing name
		}

		reqBody, err := json.Marshal(request)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/buildings", bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Execute request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

// clearBuildingsTable removes all test data from the buildings table
func clearBuildingsTable(t *testing.T, ctx context.Context) {
	t.Helper()
	_, err := testPool.Exec(ctx, "DELETE FROM booking.buildings WHERE name LIKE 'Test%'")
	require.NoError(t, err)
}
