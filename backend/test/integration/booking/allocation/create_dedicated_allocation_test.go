package allocation_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	userDomain "github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"

	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"
	ta "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/allocation"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreateDedicatedAllocation TEST_PATH=test/integration/booking/allocation/create_dedicated_allocation_test.go

func TestCreateDedicatedAllocation(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Setup test room for allocation tests
	setupTestRoom := func(t *testing.T, ctx context.Context) uuid.UUID {
		t.Helper()

		// Create building
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Create room
		room := tr.NewTestRoomWithBuilding(t, building.ID)
		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		return room.ID
	}

	// Create valid request helper
	createValidRequest := func(roomID, userID uuid.UUID) domain.CreateDedicatedAllocationRequest {
		startDate := time.Now().AddDate(0, 0, 7).Truncate(24 * time.Hour) // 7 days from now
		endDate := time.Now().AddDate(0, 0, 37).Truncate(24 * time.Hour)  // 37 days from now

		return domain.CreateDedicatedAllocationRequest{
			RoomID:    roomID,
			UserID:    userID,
			StartDate: &startDate,
			EndDate:   &endDate,
		}
	}

	setupPartnerUser := func(t *testing.T, ctx context.Context) uuid.UUID {
		user := th.NewTestUser(t, "john.doe@example.com", "John", "DOE")
		user.Role = identity.Partner.String()
		userEncx, err := userDomain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = th.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)
		partner := th.NewTestPartner(t, user.ID)
		partner.StripeAccountStatus = userDomain.StripeAccountStatusActive
		partner.StripeOnboardingComplete = true
		partnerEncx, err := userDomain.ProcessPartnerEncx(ctx, crypto, partner)
		require.NoError(t, err)
		err = th.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		return user.ID
	}

	t.Run("should successfully create dedicated allocation with valid partner token", func(t *testing.T) {
		// Clean test data
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Prepare request
		request := createValidRequest(roomID, partnerUserID)

		req := ta.NewCreateDedicatedAllocationRequest(t, ctx, testServerURL, request, adminAccessToken)

		// Execute request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response data
		assert.NotNil(t, response.ID)
		assert.NotEqual(t, uuid.Nil, response.ID)
		assert.Equal(t, request.RoomID, response.RoomID)
		assert.Equal(t, request.UserID, response.UserID)
		assert.Equal(t, domain.AllocationTypeDedicated, response.AllocationType)
		assert.WithinDuration(t, *request.StartDate, *response.StartDate, time.Second)
		assert.WithinDuration(t, *request.EndDate, *response.EndDate, time.Second)
		assert.True(t, response.IsActive)

		// Verify allocation exists in database
		allocation, err := ta.GetAllocationByID(t, ctx, response.ID, testPool, crypto)
		require.NoError(t, err)
		require.NotNil(t, allocation)

		assert.Equal(t, request.RoomID, allocation.RoomID)
		assert.Equal(t, request.UserID, allocation.UserID)
		assert.Equal(t, domain.AllocationTypeDedicated, allocation.AllocationType)
		assert.WithinDuration(t, *request.StartDate, *allocation.StartDate, time.Second)
		assert.WithinDuration(t, *request.EndDate, *allocation.EndDate, time.Second)
		assert.True(t, allocation.IsActive)
	})

	t.Run("should successfully create dedicated allocation with indefinite end date", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup partner user and get access token
		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		partnerUserID := setupPartnerUser(t, ctx)
		roomID := setupTestRoom(t, ctx)

		// Create request with nil end date (indefinite allocation)
		startDate := time.Now().AddDate(0, 0, 7).Truncate(24 * time.Hour)
		request := domain.CreateDedicatedAllocationRequest{
			RoomID: roomID,
			// UserID:    userID,
			UserID:    partnerUserID,
			StartDate: &startDate,
			EndDate:   nil, // Indefinite
		}

		// req := ta.NewCreateDedicatedAllocationRequest(t, ctx, testServerURL, request, accessToken)
		req := ta.NewCreateDedicatedAllocationRequest(t, ctx, testServerURL, request, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotNil(t, response.ID)
		assert.Nil(t, response.EndDate) // Should be nil for indefinite allocation

		// Verify in database
		allocation, err := ta.GetAllocationByID(t, ctx, response.ID, testPool, crypto)
		require.NoError(t, err)
		assert.Nil(t, allocation.EndDate)
	})

	t.Run("should successfully create dedicated allocation with admin token", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx)
		roomID := setupTestRoom(t, ctx)

		request := createValidRequest(roomID, partnerUserID)

		req := ta.NewCreateDedicatedAllocationRequest(t, ctx, testServerURL, request, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotNil(t, response.ID)
	})

	t.Run("should return 400 Bad Request for missing room_id", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx)

		request := createValidRequest(uuid.Nil, partnerUserID) // Invalid room ID

		req := ta.NewCreateDedicatedAllocationRequest(t, ctx, testServerURL, request, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "room")
	})

	t.Run("should return 400 Bad Request for missing user_id", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		roomID := setupTestRoom(t, ctx)

		request := createValidRequest(roomID, uuid.Nil) // Invalid user ID

		req := ta.NewCreateDedicatedAllocationRequest(t, ctx, testServerURL, request, adminAccessToken)

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

	t.Run("should return 400 Bad Request for missing start_date", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx)
		roomID := setupTestRoom(t, ctx)

		request := domain.CreateDedicatedAllocationRequest{
			RoomID:    roomID,
			UserID:    partnerUserID,
			StartDate: nil, // Missing start date
			EndDate:   nil,
		}

		req := ta.NewCreateDedicatedAllocationRequest(t, ctx, testServerURL, request, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "start date")
	})

	t.Run("should return 400 Bad Request for end date before start date", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx)
		roomID := setupTestRoom(t, ctx)

		startDate := time.Now().AddDate(0, 0, 30).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 0, 7).Truncate(24 * time.Hour) // Before start date

		request := domain.CreateDedicatedAllocationRequest{
			RoomID:    roomID,
			UserID:    partnerUserID,
			StartDate: &startDate,
			EndDate:   &endDate,
		}

		req := ta.NewCreateDedicatedAllocationRequest(t, ctx, testServerURL, request, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "end date")
	})

	t.Run("should return 409 Conflict for overlapping allocation on same room and user", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx)
		roomID := setupTestRoom(t, ctx)

		// Create existing allocation
		existingStartDate := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		existingEndDate := time.Now().AddDate(0, 0, 30).Truncate(24 * time.Hour)
		existingAllocation := ta.NewTestDedicatedAllocation(t, roomID, partnerUserID, existingStartDate, existingEndDate)
		ta.InsertAllocation(t, ctx, existingAllocation, testPool, crypto)

		// Try to create overlapping allocation
		newStartDate := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour) // Overlaps with existing
		newEndDate := time.Now().AddDate(0, 0, 40).Truncate(24 * time.Hour)

		request := domain.CreateDedicatedAllocationRequest{
			RoomID:    roomID,
			UserID:    partnerUserID,
			StartDate: &newStartDate,
			EndDate:   &newEndDate,
		}

		req := ta.NewCreateDedicatedAllocationRequest(t, ctx, testServerURL, request, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrAlreadyExists.Error())
	})

	t.Run("should return 400 Bad Request for non-existent room", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx)

		request := createValidRequest(uuid.New(), partnerUserID) // Non-existent room ID

		req := ta.NewCreateDedicatedAllocationRequest(t, ctx, testServerURL, request, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrInvalidInput.Error())
	})

	t.Run("should return 400 Bad Request for invalid JSON body", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Create malformed JSON
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/allocations/dedicated", bytes.NewBuffer([]byte("{invalid json")))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		if adminAccessToken != "" {
			cookie := &http.Cookie{
				Name:  ck.AccessTokenCookieName,
				Value: adminAccessToken,
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
		assert.NotEmpty(t, respBody.Error)
	})

	t.Run("should return 415 Unsupported Media Type when Content-Type is not application/json", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx)
		roomID := setupTestRoom(t, ctx)

		request := createValidRequest(roomID, partnerUserID)
		jsonBody, err := json.Marshal(request)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/allocations/dedicated", bytes.NewReader(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "text/plain") // Wrong content type

		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: adminAccessToken,
		}
		req.AddCookie(cookie)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		roomID := setupTestRoom(t, ctx)
		partnerUserID := setupPartnerUser(t, ctx)
		request := createValidRequest(roomID, partnerUserID)

		req := ta.NewCreateDedicatedAllocationRequest(t, ctx, testServerURL, request, "") // Empty token

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create expired partner session
		expiredToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx)
		roomID := setupTestRoom(t, ctx)

		request := createValidRequest(roomID, partnerUserID)

		req := ta.NewCreateDedicatedAllocationRequest(t, ctx, testServerURL, request, expiredToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		roomID := setupTestRoom(t, ctx)
		partnerUserID := setupPartnerUser(t, ctx)
		request := createValidRequest(roomID, partnerUserID)

		req := ta.NewCreateDedicatedAllocationRequest(t, ctx, testServerURL, request, "invalid-token-12345")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role (standard user)", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create standard user (not admin)
		standardAccessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)

		partnerUserID := setupPartnerUser(t, ctx)
		roomID := setupTestRoom(t, ctx)

		request := createValidRequest(roomID, partnerUserID)

		req := ta.NewCreateDedicatedAllocationRequest(t, ctx, testServerURL, request, standardAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should allow non-overlapping allocation for same room and user", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx)
		roomID := setupTestRoom(t, ctx)

		// Create first allocation
		firstStartDate := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		firstEndDate := time.Now().AddDate(0, 0, 30).Truncate(24 * time.Hour)
		firstAllocation := ta.NewTestDedicatedAllocation(t, roomID, partnerUserID, firstStartDate, firstEndDate)
		ta.InsertAllocation(t, ctx, firstAllocation, testPool, crypto)

		// Create second allocation that doesn't overlap
		secondStartDate := time.Now().AddDate(0, 0, 40).Truncate(24 * time.Hour) // After first ends
		secondEndDate := time.Now().AddDate(0, 0, 60).Truncate(24 * time.Hour)

		request := domain.CreateDedicatedAllocationRequest{
			RoomID:    roomID,
			UserID:    partnerUserID,
			StartDate: &secondStartDate,
			EndDate:   &secondEndDate,
		}

		req := ta.NewCreateDedicatedAllocationRequest(t, ctx, testServerURL, request, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotNil(t, response.ID)
	})
}
