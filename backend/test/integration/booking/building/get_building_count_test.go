package building_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetBuildingCount TEST_PATH=test/integration/booking/building/get_building_count_test.go

func TestGetBuildingCount(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 15 * time.Second}

	t.Run("should successfully get building count without authentication", func(t *testing.T) {
		// Clean test data
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create 3 test buildings
		for i := 0; i < 3; i++ {
			building := tb.NewTestBuilding(t)
			buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
			require.NoError(t, err)
			err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
			require.NoError(t, err)
		}

		// Make request without authentication using helper
		req := tb.NewGetBuildingCountRequest(t, ctx, testServerURL, nil, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]int
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify we got count of 3
		assert.Equal(t, 3, response["count"])
	})

	t.Run("should successfully get building count with standard user authentication", func(t *testing.T) {
		// Clean test data
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup standard user
		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		// Create 5 test buildings
		for i := 0; i < 5; i++ {
			building := tb.NewTestBuilding(t)
			buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
			require.NoError(t, err)
			err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
			require.NoError(t, err)
		}

		// Make request with standard user authentication using helper
		req := tb.NewGetBuildingCountRequest(t, ctx, testServerURL, nil, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]int
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, 5, response["count"])
	})

	t.Run("should successfully get building count with admin authentication", func(t *testing.T) {
		// Clean test data
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create 2 test buildings
		for i := 0; i < 2; i++ {
			building := tb.NewTestBuilding(t)
			buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
			require.NoError(t, err)
			err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
			require.NoError(t, err)
		}

		// Make request with admin authentication using helper
		req := tb.NewGetBuildingCountRequest(t, ctx, testServerURL, nil, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]int
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, 2, response["count"])
	})

	t.Run("should return 0 when no buildings exist", func(t *testing.T) {
		// Clean test data
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Make request using helper
		req := tb.NewGetBuildingCountRequest(t, ctx, testServerURL, nil, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]int
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, 0, response["count"])
	})

	t.Run("should filter count by is_active=true", func(t *testing.T) {
		// Clean test data
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create 2 active and 1 inactive building
		activeBuilding1 := tb.NewTestBuildingWithParams(t, "Active Building 1", "Paris", "France", true)
		activeBuilding2 := tb.NewTestBuildingWithParams(t, "Active Building 2", "Lyon", "France", true)
		inactiveBuilding := tb.NewTestBuildingWithParams(t, "Inactive Building", "Marseille", "France", false)

		for _, building := range []*domain.Building{activeBuilding1, activeBuilding2, inactiveBuilding} {
			buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
			require.NoError(t, err)
			err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
			require.NoError(t, err)
		}

		// Make request with is_active=true filter using helper
		req := tb.NewGetBuildingCountRequest(t, ctx, testServerURL, map[string]string{
			"is_active": "true",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]int
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should only count 2 active buildings
		assert.Equal(t, 2, response["count"])
	})

	t.Run("should filter count by is_active=false", func(t *testing.T) {
		// Clean test data
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create 1 active and 2 inactive buildings
		activeBuilding := tb.NewTestBuildingWithParams(t, "Active Building", "Paris", "France", true)
		inactiveBuilding1 := tb.NewTestBuildingWithParams(t, "Inactive Building 1", "Lyon", "France", false)
		inactiveBuilding2 := tb.NewTestBuildingWithParams(t, "Inactive Building 2", "Marseille", "France", false)

		for _, building := range []*domain.Building{activeBuilding, inactiveBuilding1, inactiveBuilding2} {
			buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
			require.NoError(t, err)
			err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
			require.NoError(t, err)
		}

		// Make request with is_active=false filter using helper
		req := tb.NewGetBuildingCountRequest(t, ctx, testServerURL, map[string]string{
			"is_active": "false",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]int
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should only count 2 inactive buildings
		assert.Equal(t, 2, response["count"])
	})

	t.Run("should filter count by city", func(t *testing.T) {
		// Clean test data
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create buildings in different cities
		parisBuilding1 := tb.NewTestBuildingWithParams(t, "Paris Building 1", "Paris", "France", true)
		parisBuilding2 := tb.NewTestBuildingWithParams(t, "Paris Building 2", "Paris", "France", true)
		lyonBuilding := tb.NewTestBuildingWithParams(t, "Lyon Building", "Lyon", "France", true)

		for _, building := range []*domain.Building{parisBuilding1, parisBuilding2, lyonBuilding} {
			buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
			require.NoError(t, err)
			err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
			require.NoError(t, err)
		}

		// Make request with city filter using helper
		req := tb.NewGetBuildingCountRequest(t, ctx, testServerURL, map[string]string{
			"city": "Paris",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]int
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should only count 2 Paris buildings
		assert.Equal(t, 2, response["count"])
	})

	t.Run("should filter count by country", func(t *testing.T) {
		// Clean test data
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create buildings in different countries
		franceBuilding1 := tb.NewTestBuildingWithParams(t, "France Building 1", "Paris", "France", true)
		franceBuilding2 := tb.NewTestBuildingWithParams(t, "France Building 2", "Lyon", "France", true)
		spainBuilding := tb.NewTestBuildingWithParams(t, "Spain Building", "Madrid", "Spain", true)

		for _, building := range []*domain.Building{franceBuilding1, franceBuilding2, spainBuilding} {
			buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
			require.NoError(t, err)
			err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
			require.NoError(t, err)
		}

		// Make request with country filter using helper
		req := tb.NewGetBuildingCountRequest(t, ctx, testServerURL, map[string]string{
			"country": "France",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]int
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should only count 2 France buildings
		assert.Equal(t, 2, response["count"])
	})

	t.Run("should combine multiple filters (city and is_active)", func(t *testing.T) {
		// Clean test data
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create various buildings
		parisActive := tb.NewTestBuildingWithParams(t, "Paris Active", "Paris", "France", true)
		parisInactive := tb.NewTestBuildingWithParams(t, "Paris Inactive", "Paris", "France", false)
		lyonActive := tb.NewTestBuildingWithParams(t, "Lyon Active", "Lyon", "France", true)

		for _, building := range []*domain.Building{parisActive, parisInactive, lyonActive} {
			buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
			require.NoError(t, err)
			err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
			require.NoError(t, err)
		}

		// Make request with city and is_active filters using helper
		req := tb.NewGetBuildingCountRequest(t, ctx, testServerURL, map[string]string{
			"city":      "Paris",
			"is_active": "true",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]int
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should only count 1 building (Paris and active)
		assert.Equal(t, 1, response["count"])
	})

	t.Run("should return 400 for invalid is_active parameter", func(t *testing.T) {
		// Make request with invalid is_active value using helper
		req := tb.NewGetBuildingCountRequest(t, ctx, testServerURL, map[string]string{
			"is_active": "invalid",
		}, "")

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
		assert.Contains(t, respBody.Error, "is_active must be a boolean")
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
		req := tb.NewGetBuildingCountRequest(t, shortCtx, testServerURL, nil, "")

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

	t.Run("should return 0 count for filters with no matches", func(t *testing.T) {
		// Clean test data
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create some test buildings
		building1 := tb.NewTestBuildingWithParams(t, "Building 1", "Paris", "France", true)
		building2 := tb.NewTestBuildingWithParams(t, "Building 2", "Lyon", "France", true)

		for _, building := range []*domain.Building{building1, building2} {
			buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
			require.NoError(t, err)
			err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
			require.NoError(t, err)
		}

		// Make request with filter that won't match any building
		req := tb.NewGetBuildingCountRequest(t, ctx, testServerURL, map[string]string{
			"city": "NonExistentCity",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]int
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, 0, response["count"])
	})
}
