package building_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestUpdateBuilding TEST_PATH=test/integration/booking/building/update_building_test.go

func TestUpdateBuilding(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully update building name", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create initial building
		initialBuilding := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, initialBuilding)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Update name
		newName := "Updated Building Name"
		updateRequest := domain.UpdateBuildingRequest{
			ID:   initialBuilding.ID,
			Name: &newName,
		}

		req := tb.NewUpdateBuildingRequest(t, ctx, testServerURL, initialBuilding.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.BuildingResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, initialBuilding.ID, response.ID)
		assert.Equal(t, newName, response.Name)
		assert.Equal(t, initialBuilding.Address, response.Address)

		// Verify in database
		updatedBuildingEncx, err := tb.GetBuildingEncxByID(t, ctx, testPool, initialBuilding.ID)
		require.NoError(t, err)

		updatedBuilding, err := domain.DecryptBuildingEncx(ctx, crypto, updatedBuildingEncx)
		require.NoError(t, err)

		assert.Equal(t, newName, updatedBuilding.Name)
		assert.Equal(t, initialBuilding.Address, updatedBuilding.Address)
	})

	t.Run("should successfully update building address", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create initial building
		initialBuilding := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, initialBuilding)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Update address
		newAddress := "456 New Street"
		updateRequest := domain.UpdateBuildingRequest{
			ID:      initialBuilding.ID,
			Address: &newAddress,
		}

		req := tb.NewUpdateBuildingRequest(t, ctx, testServerURL, initialBuilding.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.BuildingResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, newAddress, response.Address)
		assert.Equal(t, initialBuilding.Name, response.Name)

		// Verify in database
		updatedBuildingEncx, err := tb.GetBuildingEncxByID(t, ctx, testPool, initialBuilding.ID)
		require.NoError(t, err)

		updatedBuilding, err := domain.DecryptBuildingEncx(ctx, crypto, updatedBuildingEncx)
		require.NoError(t, err)

		assert.Equal(t, newAddress, updatedBuilding.Address)
	})

	t.Run("should successfully update multiple fields", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create initial building
		initialBuilding := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, initialBuilding)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Update multiple fields
		newName := "Updated Name"
		newCity := "Lyon"
		newCountry := "France"
		newPostalCode := "69001"
		updateRequest := domain.UpdateBuildingRequest{
			ID:         initialBuilding.ID,
			Name:       &newName,
			City:       &newCity,
			Country:    &newCountry,
			PostalCode: &newPostalCode,
		}

		req := tb.NewUpdateBuildingRequest(t, ctx, testServerURL, initialBuilding.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.BuildingResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, newName, response.Name)
		assert.Equal(t, newCity, response.City)
		assert.Equal(t, newCountry, response.Country)
		assert.Equal(t, newPostalCode, response.PostalCode)

		// Verify in database
		updatedBuildingEncx, err := tb.GetBuildingEncxByID(t, ctx, testPool, initialBuilding.ID)
		require.NoError(t, err)

		updatedBuilding, err := domain.DecryptBuildingEncx(ctx, crypto, updatedBuildingEncx)
		require.NoError(t, err)

		assert.Equal(t, newName, updatedBuilding.Name)
		assert.Equal(t, newCity, updatedBuilding.City)
		assert.Equal(t, newCountry, updatedBuilding.Country)
		assert.Equal(t, newPostalCode, updatedBuilding.PostalCode)
	})

	t.Run("should successfully update optional fields", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create initial building
		initialBuilding := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, initialBuilding)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Update optional fields
		newDescription := "Updated description"
		newPhone := "0612345679"
		newEmail := "updated@example.com"
		updateRequest := domain.UpdateBuildingRequest{
			ID:          initialBuilding.ID,
			Description: &newDescription,
			Phone:       &newPhone,
			Email:       &newEmail,
		}

		req := tb.NewUpdateBuildingRequest(t, ctx, testServerURL, initialBuilding.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.BuildingResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, newDescription, response.Description)
		assert.Equal(t, newPhone, response.Phone)
		assert.Equal(t, newEmail, response.Email)

		// Verify in database
		updatedBuildingEncx, err := tb.GetBuildingEncxByID(t, ctx, testPool, initialBuilding.ID)
		require.NoError(t, err)

		updatedBuilding, err := domain.DecryptBuildingEncx(ctx, crypto, updatedBuildingEncx)
		require.NoError(t, err)

		assert.Equal(t, newDescription, updatedBuilding.Description)
		assert.Equal(t, newPhone, updatedBuilding.Phone)
		assert.Equal(t, newEmail, updatedBuilding.Email)
	})

	t.Run("should preserve unchanged fields", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create initial building
		initialBuilding := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, initialBuilding)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Update only name, other fields should remain unchanged
		newName := "Only Name Updated"
		updateRequest := domain.UpdateBuildingRequest{
			ID:   initialBuilding.ID,
			Name: &newName,
		}

		req := tb.NewUpdateBuildingRequest(t, ctx, testServerURL, initialBuilding.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify in database - all other fields should be unchanged
		updatedBuildingEncx, err := tb.GetBuildingEncxByID(t, ctx, testPool, initialBuilding.ID)
		require.NoError(t, err)

		updatedBuilding, err := domain.DecryptBuildingEncx(ctx, crypto, updatedBuildingEncx)
		require.NoError(t, err)

		assert.Equal(t, newName, updatedBuilding.Name)
		assert.Equal(t, initialBuilding.Address, updatedBuilding.Address)
		assert.Equal(t, initialBuilding.City, updatedBuilding.City)
		assert.Equal(t, initialBuilding.PostalCode, updatedBuilding.PostalCode)
		assert.Equal(t, initialBuilding.Country, updatedBuilding.Country)
		assert.Equal(t, initialBuilding.Description, updatedBuilding.Description)
		assert.Equal(t, initialBuilding.Phone, updatedBuilding.Phone)
		assert.Equal(t, initialBuilding.Email, updatedBuilding.Email)
	})

	t.Run("should return 404 Not Found when building does not exist", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		nonExistentID := uuid.New()
		newName := "Updated Name"
		updateRequest := domain.UpdateBuildingRequest{
			ID:   nonExistentID,
			Name: &newName,
		}

		req := tb.NewUpdateBuildingRequest(t, ctx, testServerURL, nonExistentID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrDomainNotFound.Error())
	})

	t.Run("should return 400 Bad Request for invalid building ID format", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		invalidID := "not-a-valid-uuid"
		newName := "Updated Name"
		updateRequest := domain.UpdateBuildingRequest{
			Name: &newName,
		}

		req := tb.NewUpdateBuildingRequest(t, ctx, testServerURL, invalidID, updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrInvalidValue.Error())
	})

	t.Run("should return 400 Bad Request for invalid JSON body", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create initial building
		initialBuilding := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, initialBuilding)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Create malformed JSON
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, testServerURL+"/buildings/"+initialBuilding.ID.String(), bytes.NewBuffer([]byte("{invalid json")))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		if accessToken != "" {
			cookie := &http.Cookie{
				Name:  ck.AccessTokenCookieName,
				Value: accessToken,
			}
			req.AddCookie(cookie)
		}

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrInvalidValue.Error())
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create initial building
		initialBuilding := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, initialBuilding)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		newName := "Updated Name"
		updateRequest := domain.UpdateBuildingRequest{
			ID:   initialBuilding.ID,
			Name: &newName,
		}

		req := tb.NewUpdateBuildingRequest(t, ctx, testServerURL, initialBuilding.ID.String(), updateRequest, "") // Empty token

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create expired admin session
		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Create initial building
		initialBuilding := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, initialBuilding)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		newName := "Updated Name"
		updateRequest := domain.UpdateBuildingRequest{
			ID:   initialBuilding.ID,
			Name: &newName,
		}

		req := tb.NewUpdateBuildingRequest(t, ctx, testServerURL, initialBuilding.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role (standard user)", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create standard user (not admin)
		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		// Create initial building
		initialBuilding := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, initialBuilding)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		newName := "Updated Name"
		updateRequest := domain.UpdateBuildingRequest{
			ID:   initialBuilding.ID,
			Name: &newName,
		}

		req := tb.NewUpdateBuildingRequest(t, ctx, testServerURL, initialBuilding.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role (partner)", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create partner user (not admin)
		accessToken := tu.SetupUserWithRole(t, ctx, identity.Partner, authCtx)

		// Create initial building
		initialBuilding := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, initialBuilding)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		newName := "Updated Name"
		updateRequest := domain.UpdateBuildingRequest{
			ID:   initialBuilding.ID,
			Name: &newName,
		}

		req := tb.NewUpdateBuildingRequest(t, ctx, testServerURL, initialBuilding.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create initial building
		initialBuilding := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, initialBuilding)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		newName := "Updated Name"
		updateRequest := domain.UpdateBuildingRequest{
			ID:   initialBuilding.ID,
			Name: &newName,
		}

		req := tb.NewUpdateBuildingRequest(t, ctx, testServerURL, initialBuilding.ID.String(), updateRequest, "invalid-token-12345")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should handle building deleted between GET and UPDATE", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create initial building
		initialBuilding := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, initialBuilding)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Delete the building to simulate race condition
		_, err = testPool.Exec(ctx, "DELETE FROM booking.buildings WHERE id = $1", initialBuilding.ID)
		require.NoError(t, err)

		// Try to update deleted building
		newName := "Updated Name"
		updateRequest := domain.UpdateBuildingRequest{
			ID:   initialBuilding.ID,
			Name: &newName,
		}

		req := tb.NewUpdateBuildingRequest(t, ctx, testServerURL, initialBuilding.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return 404 since building doesn't exist at GET time
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrDomainNotFound.Error())
	})

	t.Run("should handle context timeout appropriately", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create initial building
		initialBuilding := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, initialBuilding)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		newName := "Updated Name"
		updateRequest := domain.UpdateBuildingRequest{
			ID:   initialBuilding.ID,
			Name: &newName,
		}

		// Use a very short context timeout to potentially trigger timeout
		shortCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
		defer cancel()

		time.Sleep(2 * time.Millisecond) // Ensure timeout has passed

		req := tb.NewUpdateBuildingRequest(t, shortCtx, testServerURL, initialBuilding.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		// Either the context timeout or a successful response (if operation was fast enough)
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
