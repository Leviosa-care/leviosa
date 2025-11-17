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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreateBuilding TEST_PATH=test/integration/booking/building/create_building_test.go

func TestCreateBuilding(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	validBuilding := tb.NewTestBuilding(t)

	validRequest := domain.CreateBuildingRequest{
		Name:        validBuilding.Name,
		Address:     validBuilding.Address,
		City:        validBuilding.City,
		PostalCode:  validBuilding.PostalCode,
		Country:     validBuilding.Country,
		Description: validBuilding.Description,
		Phone:       validBuilding.Phone,
		Email:       validBuilding.Email,
		IsActive:    true,
	}

	t.Run("should successfully create a building with valid admin token", func(t *testing.T) {
		// Clean test data
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and get access token
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Prepare request
		request := validRequest

		req := tb.NewCreateBuildingRequest(t, ctx, testServerURL, request, accessToken)

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
		assert.Equal(t, request.Name, response.Name)
		assert.Equal(t, request.Address, response.Address)
		assert.Equal(t, request.City, response.City)
		assert.Equal(t, request.PostalCode, response.PostalCode)
		assert.Equal(t, request.Country, response.Country)
		assert.True(t, response.IsActive)

		// Verify building exists in database
		buildingEncx, err := tb.GetBuildingEncxByID(t, ctx, testPool, response.ID)
		require.NoError(t, err)

		building, err := domain.DecryptBuildingEncx(ctx, crypto, buildingEncx)
		require.NoError(t, err)

		assert.Equal(t, request.Name, building.Name)
		assert.Equal(t, request.Address, building.Address)
		assert.Equal(t, request.City, building.City)
		assert.Equal(t, request.PostalCode, building.PostalCode)
		assert.Equal(t, request.Country, building.Country)
		assert.True(t, building.IsActive)
	})

	t.Run("should successfully create inactive building", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := validRequest
		request.IsActive = false

		req := tb.NewCreateBuildingRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response domain.BuildingResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.False(t, response.IsActive)

		// Verify in database
		buildingEncx, err := tb.GetBuildingEncxByID(t, ctx, testPool, response.ID)
		require.NoError(t, err)

		building, err := domain.DecryptBuildingEncx(ctx, crypto, buildingEncx)
		require.NoError(t, err)

		assert.False(t, building.IsActive)
	})

	t.Run("should return 400 Bad Request for empty name", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := validRequest
		request.Name = ""

		req := tb.NewCreateBuildingRequest(t, ctx, testServerURL, request, accessToken)

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

	t.Run("should return 400 Bad Request for empty address", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := validRequest
		request.Address = ""

		req := tb.NewCreateBuildingRequest(t, ctx, testServerURL, request, accessToken)

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

	t.Run("should return 400 Bad Request for empty city", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := validRequest
		request.City = ""

		req := tb.NewCreateBuildingRequest(t, ctx, testServerURL, request, accessToken)

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

	t.Run("should return 400 Bad Request for empty postal code", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := validRequest
		request.PostalCode = ""

		req := tb.NewCreateBuildingRequest(t, ctx, testServerURL, request, accessToken)

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

	t.Run("should return 400 Bad Request for empty country", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := validRequest
		request.Country = ""

		req := tb.NewCreateBuildingRequest(t, ctx, testServerURL, request, accessToken)

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

		// Create malformed JSON
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/buildings", bytes.NewBuffer([]byte("{invalid json")))
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

	t.Run("should return 409 Conflict when building with same name or address already exists", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create first building
		request := validRequest

		req := tb.NewCreateBuildingRequest(t, ctx, testServerURL, request, accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		resp.Body.Close()
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		// Try to create duplicate
		req2 := tb.NewCreateBuildingRequest(t, ctx, testServerURL, request, accessToken)
		resp2, err := client.Do(req2)
		require.NoError(t, err)
		defer resp2.Body.Close()

		assert.Equal(t, http.StatusConflict, resp2.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp2.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrAlreadyExists.Error())
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)

		request := validRequest

		req := tb.NewCreateBuildingRequest(t, ctx, testServerURL, request, "") // Empty token

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

		request := validRequest

		req := tb.NewCreateBuildingRequest(t, ctx, testServerURL, request, accessToken)

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

		request := validRequest

		req := tb.NewCreateBuildingRequest(t, ctx, testServerURL, request, accessToken)

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

		request := validRequest

		req := tb.NewCreateBuildingRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)

		request := validRequest

		req := tb.NewCreateBuildingRequest(t, ctx, testServerURL, request, "invalid-token-12345")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should handle context timeout appropriately", func(t *testing.T) {
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := validRequest

		// Use a very short context timeout to potentially trigger timeout
		shortCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
		defer cancel()

		time.Sleep(2 * time.Millisecond) // Ensure timeout has passed

		req := tb.NewCreateBuildingRequest(t, shortCtx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		// Either the context timeout or a successful response (if operation was fast enough)
		if err != nil {
			// Context timeout on client side
			assert.ErrorIs(t, err, context.DeadlineExceeded)
		} else {
			defer resp.Body.Close()
			// If we got a response, it should be either success or timeout status
			assert.True(t, resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusRequestTimeout)
		}
	})
}
