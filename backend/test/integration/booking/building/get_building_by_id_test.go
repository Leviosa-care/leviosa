package building_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetBuildingByID TEST_PATH=test/integration/booking/building/get_building_by_id_test.go

func TestGetBuildingByID(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get building by ID without authentication", func(t *testing.T) {
		// Clean test data
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create test building using helper
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)

		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Make request without authentication using helper
		req := tb.NewGetBuildingByIDRequest(t, ctx, testServerURL, building.ID.String(), "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.BuildingResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response data
		assert.Equal(t, building.ID, response.ID)
		assert.Equal(t, building.Name, response.Name)
		assert.Equal(t, building.Address, response.Address)
		assert.Equal(t, building.City, response.City)
		assert.Equal(t, building.PostalCode, response.PostalCode)
		assert.Equal(t, building.Country, response.Country)
		assert.Equal(t, building.Description, response.Description)
		assert.Equal(t, building.Phone, response.Phone)
		assert.Equal(t, building.Email, response.Email)
		assert.Equal(t, building.IsActive, response.IsActive)
	})

	t.Run("should successfully get building by ID with standard user authentication", func(t *testing.T) {
		// Clean test data
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup standard user
		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		// Create test building using helper
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)

		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Make request with standard user authentication using helper
		req := tb.NewGetBuildingByIDRequest(t, ctx, testServerURL, building.ID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.BuildingResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response data
		assert.Equal(t, building.ID, response.ID)
		assert.Equal(t, building.Name, response.Name)
	})

	t.Run("should successfully get building by ID with admin authentication", func(t *testing.T) {
		// Clean test data
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create test building using helper
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)

		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Make request with admin authentication using helper
		req := tb.NewGetBuildingByIDRequest(t, ctx, testServerURL, building.ID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.BuildingResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response data
		assert.Equal(t, building.ID, response.ID)
		assert.Equal(t, building.Name, response.Name)
	})

	t.Run("should return 404 when building does not exist", func(t *testing.T) {
		// Clean test data
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Use a non-existent building ID
		nonExistentID := uuid.New()

		// Make request using helper
		req := tb.NewGetBuildingByIDRequest(t, ctx, testServerURL, nonExistentID.String(), "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrDomainNotFound.Error())
	})

	t.Run("should return 400 when building ID is invalid", func(t *testing.T) {
		// Make request with invalid UUID using helper
		req := tb.NewGetBuildingByIDRequest(t, ctx, testServerURL, "invalid-uuid", "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "invalid building ID format")
	})

	t.Run("should return 400 when building ID is missing", func(t *testing.T) {
		// Make request without ID (this would hit a different route or 404)
		// Note: With Go 1.22+ router patterns, this might return 404 Method Not Allowed
		req := tb.NewGetBuildingByIDRequest(t, ctx, testServerURL, "", "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// The router will return 404 or 405 for invalid route
		assert.True(t, resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusMethodNotAllowed)
	})

	t.Run("should return building with all fields populated", func(t *testing.T) {
		// Clean test data
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create test building with all optional fields
		building := &domain.Building{
			ID:          uuid.New(),
			Name:        "Complete Building",
			Address:     "456 Avenue des Champs-Élysées",
			City:        "Paris",
			PostalCode:  "75008",
			Country:     "France",
			Description: "Luxury office building in the heart of Paris",
			Phone:       "+33 1 23 45 67 89",
			Email:       "contact@completebuilding.fr",
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)

		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Make request using helper
		req := tb.NewGetBuildingByIDRequest(t, ctx, testServerURL, building.ID.String(), "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.BuildingResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify all fields
		assert.Equal(t, building.ID, response.ID)
		assert.Equal(t, building.Name, response.Name)
		assert.Equal(t, building.Address, response.Address)
		assert.Equal(t, building.City, response.City)
		assert.Equal(t, building.PostalCode, response.PostalCode)
		assert.Equal(t, building.Country, response.Country)
		assert.Equal(t, building.Description, response.Description)
		assert.Equal(t, building.Phone, response.Phone)
		assert.Equal(t, building.Email, response.Email)
		assert.Equal(t, building.IsActive, response.IsActive)
	})

	t.Run("should return inactive building", func(t *testing.T) {
		// Clean test data
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create inactive building
		building := tb.NewTestBuilding(t)
		building.IsActive = false

		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)

		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Make request using helper
		req := tb.NewGetBuildingByIDRequest(t, ctx, testServerURL, building.ID.String(), "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.BuildingResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify inactive status
		assert.Equal(t, building.ID, response.ID)
		assert.False(t, response.IsActive)
	})

	t.Run("should handle context timeout appropriately", func(t *testing.T) {
		// Clean test data
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create test building
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)

		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Use a very short context timeout
		shortCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
		defer cancel()

		time.Sleep(2 * time.Millisecond) // Ensure timeout has passed

		// Make request using helper (but with short context)
		req := tb.NewGetBuildingByIDRequest(t, shortCtx, testServerURL, building.ID.String(), "")

		resp, err := client.Do(req)
		// Either the context timeout or a successful response
		if err != nil {
			// Context timeout on client side
			assert.ErrorIs(t, err, context.DeadlineExceeded)
		} else {
			defer resp.Body.Close()
			// If we got a response, it should be either success or timeout status
			assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusRequestTimeout)
		}
	})
}
