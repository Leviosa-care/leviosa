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

// make test-func TEST_NAME=TestCreateSharedAllocation TEST_PATH=test/integration/booking/allocation/create_shared_allocation_test.go

func TestCreateSharedAllocation(t *testing.T) {
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
	createValidRequest := func(roomID, userID uuid.UUID) domain.CreateSharedAllocationRequest {
		return domain.CreateSharedAllocationRequest{
			RoomID: roomID,
			UserID: userID,
		}
	}

	setupPartnerUser := func(t *testing.T, ctx context.Context, email string) uuid.UUID {
		t.Helper()

		user := th.NewTestUser(t, email, "John", "DOE")
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

	t.Run("should successfully create shared allocation with valid partner and admin token", func(t *testing.T) {
		// Clean test data
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx, "partner1@example.com")

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Prepare request
		request := createValidRequest(roomID, partnerUserID)

		req := ta.NewCreateSharedAllocationRequest(t, ctx, testServerURL, request, adminAccessToken)

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
		assert.Equal(t, domain.AllocationTypeShared, response.AllocationType)
		assert.Nil(t, response.StartDate) // Shared allocations have no start date
		assert.Nil(t, response.EndDate)   // Shared allocations have no end date
		assert.True(t, response.IsActive)

		// Verify allocation exists in database
		allocation, err := ta.GetAllocationByID(t, ctx, response.ID, testPool, crypto)
		require.NoError(t, err)
		require.NotNil(t, allocation)

		assert.Equal(t, request.RoomID, allocation.RoomID)
		assert.Equal(t, request.UserID, allocation.UserID)
		assert.Equal(t, domain.AllocationTypeShared, allocation.AllocationType)
		assert.Nil(t, allocation.StartDate)
		assert.Nil(t, allocation.EndDate)
		assert.True(t, allocation.IsActive)
	})

	t.Run("should successfully create shared allocation with admin token", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx, "partner2@example.com")
		roomID := setupTestRoom(t, ctx)

		request := createValidRequest(roomID, partnerUserID)

		req := ta.NewCreateSharedAllocationRequest(t, ctx, testServerURL, request, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotNil(t, response.ID)
		assert.Equal(t, domain.AllocationTypeShared, response.AllocationType)
	})

	t.Run("should return 400 Bad Request for missing room_id", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx, "partner3@example.com")

		request := createValidRequest(uuid.Nil, partnerUserID) // Invalid room ID

		req := ta.NewCreateSharedAllocationRequest(t, ctx, testServerURL, request, adminAccessToken)

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

		req := ta.NewCreateSharedAllocationRequest(t, ctx, testServerURL, request, adminAccessToken)

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

	t.Run("should return 400 Bad Request for non-existent room", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx, "partner4@example.com")

		request := createValidRequest(uuid.New(), partnerUserID) // Non-existent room ID

		req := ta.NewCreateSharedAllocationRequest(t, ctx, testServerURL, request, adminAccessToken)

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
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/allocations/shared", bytes.NewBuffer([]byte("{invalid json")))
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
		partnerUserID := setupPartnerUser(t, ctx, "partner5@example.com")
		roomID := setupTestRoom(t, ctx)

		request := createValidRequest(roomID, partnerUserID)
		jsonBody, err := json.Marshal(request)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/allocations/shared", bytes.NewReader(jsonBody))
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
		partnerUserID := setupPartnerUser(t, ctx, "partner6@example.com")
		request := createValidRequest(roomID, partnerUserID)

		req := ta.NewCreateSharedAllocationRequest(t, ctx, testServerURL, request, "") // Empty token

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create expired session
		expiredToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx, "partner7@example.com")
		roomID := setupTestRoom(t, ctx)

		request := createValidRequest(roomID, partnerUserID)

		req := ta.NewCreateSharedAllocationRequest(t, ctx, testServerURL, request, expiredToken)

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
		partnerUserID := setupPartnerUser(t, ctx, "partner8@example.com")
		request := createValidRequest(roomID, partnerUserID)

		req := ta.NewCreateSharedAllocationRequest(t, ctx, testServerURL, request, "invalid-token-12345")

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

		partnerUserID := setupPartnerUser(t, ctx, "partner9@example.com")
		roomID := setupTestRoom(t, ctx)

		request := createValidRequest(roomID, partnerUserID)

		req := ta.NewCreateSharedAllocationRequest(t, ctx, testServerURL, request, standardAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should allow multiple shared allocations for the same room with different partners", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		roomID := setupTestRoom(t, ctx)

		// Create first partner and allocation
		partner1ID := setupPartnerUser(t, ctx, "partner10@example.com")
		allocation1 := ta.NewTestSharedAllocation(t, roomID, partner1ID)
		ta.InsertAllocation(t, ctx, allocation1, testPool, crypto)

		// Create second partner and allocation for the same room
		partner2ID := setupPartnerUser(t, ctx, "partner11@example.com")
		request := createValidRequest(roomID, partner2ID)

		req := ta.NewCreateSharedAllocationRequest(t, ctx, testServerURL, request, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should succeed - multiple partners can share the same room
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotNil(t, response.ID)
		assert.Equal(t, roomID, response.RoomID)
		assert.Equal(t, partner2ID, response.UserID)
		assert.Equal(t, domain.AllocationTypeShared, response.AllocationType)
	})

	t.Run("should allow multiple shared allocations for the same partner with different rooms", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx, "partner12@example.com")

		// Create first room and allocation
		room1ID := setupTestRoom(t, ctx)
		allocation1 := ta.NewTestSharedAllocation(t, room1ID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation1, testPool, crypto)

		// Create second room and allocation for the same partner
		room2ID := setupTestRoom(t, ctx)
		request := createValidRequest(room2ID, partnerUserID)

		req := ta.NewCreateSharedAllocationRequest(t, ctx, testServerURL, request, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should succeed - same partner can have shared access to multiple rooms
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotNil(t, response.ID)
		assert.Equal(t, room2ID, response.RoomID)
		assert.Equal(t, partnerUserID, response.UserID)
		assert.Equal(t, domain.AllocationTypeShared, response.AllocationType)
	})

	t.Run("should return 400 Bad Request for non-existent partner", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		roomID := setupTestRoom(t, ctx)

		// Use a non-existent user ID
		request := createValidRequest(roomID, uuid.New())

		req := ta.NewCreateSharedAllocationRequest(t, ctx, testServerURL, request, adminAccessToken)

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
}
