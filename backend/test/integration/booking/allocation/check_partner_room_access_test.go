package allocation_test

import (
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

// make test-func TEST_NAME=TestCheckPartnerRoomAccess TEST_PATH=test/integration/booking/allocation/check_partner_room_access_test.go

func TestCheckPartnerRoomAccess(t *testing.T) {
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

	setupPartnerUser := func(t *testing.T, ctx context.Context, email string) uuid.UUID {
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

	type AccessResponse struct {
		HasAccess bool      `json:"has_access"`
		CheckedAt time.Time `json:"checked_at"`
	}

	t.Run("should return true when partner has active shared allocation", func(t *testing.T) {
		// Clean test data
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		partnerAccessToken := tu.SetupUserWithRole(t, ctx, identity.Partner, authCtx)
		partnerID := setupPartnerUser(t, ctx, "partner1@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create active shared allocation
		allocation := ta.NewTestSharedAllocation(t, roomID, partnerID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Check access without time parameter (defaults to now)
		req := ta.NewCheckPartnerRoomAccessRequest(t, ctx, testServerURL, partnerID, roomID, nil, partnerAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response AccessResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.HasAccess)
		assert.False(t, response.CheckedAt.IsZero())
	})

	t.Run("should return true when partner has active dedicated allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerID := setupPartnerUser(t, ctx, "partner2@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create active dedicated allocation (starts 7 days ago, ends 7 days from now)
		allocation := ta.NewTestActiveDedicatedAllocation(t, roomID, partnerID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		req := ta.NewCheckPartnerRoomAccessRequest(t, ctx, testServerURL, partnerID, roomID, nil, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response AccessResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.HasAccess)
	})

	t.Run("should return false when partner has no allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerID := setupPartnerUser(t, ctx, "partner3@example.com")
		roomID := setupTestRoom(t, ctx)

		// Don't create any allocation
		req := ta.NewCheckPartnerRoomAccessRequest(t, ctx, testServerURL, partnerID, roomID, nil, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response AccessResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.False(t, response.HasAccess)
	})

	t.Run("should return false when partner has inactive allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerID := setupPartnerUser(t, ctx, "partner4@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create inactive allocation
		allocation := ta.NewTestInactiveAllocation(t, roomID, partnerID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		req := ta.NewCheckPartnerRoomAccessRequest(t, ctx, testServerURL, partnerID, roomID, nil, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response AccessResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.False(t, response.HasAccess)
	})

	t.Run("should check access at specific past time using at parameter", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerID := setupPartnerUser(t, ctx, "partner5@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create allocation that was active 20 days ago
		startDate := time.Now().AddDate(0, 0, -30).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 0, -10).Truncate(24 * time.Hour)
		allocation := ta.NewTestDedicatedAllocation(t, roomID, partnerID, startDate, endDate)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Check access 20 days ago (when allocation was active)
		checkTime := time.Now().AddDate(0, 0, -20).Format(time.RFC3339)
		req := ta.NewCheckPartnerRoomAccessRequest(t, ctx, testServerURL, partnerID, roomID, &checkTime, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response AccessResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.HasAccess)
		// Verify the checked_at time matches what we requested
		assert.WithinDuration(t, time.Now().AddDate(0, 0, -20), response.CheckedAt, 24*time.Hour)
	})

	t.Run("should return false when checking past time before allocation started", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerID := setupPartnerUser(t, ctx, "partner6@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create allocation starting 10 days ago
		startDate := time.Now().AddDate(0, 0, -10).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		allocation := ta.NewTestDedicatedAllocation(t, roomID, partnerID, startDate, endDate)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Check access 20 days ago (before allocation started)
		checkTime := time.Now().AddDate(0, 0, -20).Format(time.RFC3339)
		req := ta.NewCheckPartnerRoomAccessRequest(t, ctx, testServerURL, partnerID, roomID, &checkTime, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response AccessResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.False(t, response.HasAccess)
	})

	t.Run("should check access at specific future time using at parameter", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerID := setupPartnerUser(t, ctx, "partner7@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create future allocation
		allocation := ta.NewTestFutureDedicatedAllocation(t, roomID, partnerID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Check access 20 days from now (when allocation will be active)
		checkTime := time.Now().AddDate(0, 0, 20).Format(time.RFC3339)
		req := ta.NewCheckPartnerRoomAccessRequest(t, ctx, testServerURL, partnerID, roomID, &checkTime, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response AccessResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.HasAccess)
	})

	t.Run("should return false when checking current time for future allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerID := setupPartnerUser(t, ctx, "partner8@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create future allocation (starts 15 days from now)
		allocation := ta.NewTestFutureDedicatedAllocation(t, roomID, partnerID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Check access now (allocation hasn't started yet)
		req := ta.NewCheckPartnerRoomAccessRequest(t, ctx, testServerURL, partnerID, roomID, nil, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response AccessResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.False(t, response.HasAccess)
	})

	t.Run("should handle allocation with indefinite end date", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerID := setupPartnerUser(t, ctx, "partner9@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create allocation with indefinite end date
		startDate := time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
		allocation := ta.NewTestDedicatedAllocationIndefinite(t, roomID, partnerID, startDate)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		req := ta.NewCheckPartnerRoomAccessRequest(t, ctx, testServerURL, partnerID, roomID, nil, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response AccessResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.HasAccess)
	})

	t.Run("should ignore invalid at parameter and default to now", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerID := setupPartnerUser(t, ctx, "partner10@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create current allocation
		allocation := ta.NewTestActiveDedicatedAllocation(t, roomID, partnerID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Use invalid time format (should be ignored, defaults to now)
		invalidTime := "not-a-valid-timestamp"
		req := ta.NewCheckPartnerRoomAccessRequest(t, ctx, testServerURL, partnerID, roomID, &invalidTime, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response AccessResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should have access because it defaults to now and allocation is active
		assert.True(t, response.HasAccess)
	})

	t.Run("should return 400 Bad Request for invalid partner ID format", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		roomID := setupTestRoom(t, ctx)

		// Create request with invalid UUID format
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			testServerURL+"/partners/invalid-uuid/rooms/"+roomID.String()+"/access",
			nil,
		)
		require.NoError(t, err)

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
		assert.Contains(t, respBody.Error, errs.ErrInvalidValue.Error())
	})

	t.Run("should return 400 Bad Request for invalid room ID format", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerID := setupPartnerUser(t, ctx, "partner11@example.com")

		// Create request with invalid UUID format
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			testServerURL+"/partners/"+partnerID.String()+"/rooms/invalid-uuid/access",
			nil,
		)
		require.NoError(t, err)

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
		assert.Contains(t, respBody.Error, errs.ErrInvalidValue.Error())
	})

	t.Run("should return 400 Bad Request for invalid path format", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Create request with malformed path (missing IDs)
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			testServerURL+"/partners/access",
			nil,
		)
		require.NoError(t, err)

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

		// Should return 404 or 400 depending on router behavior
		assert.True(t, resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest)
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		partnerID := setupPartnerUser(t, ctx, "partner12@example.com")
		roomID := setupTestRoom(t, ctx)

		req := ta.NewCheckPartnerRoomAccessRequest(t, ctx, testServerURL, partnerID, roomID, nil, "") // Empty token

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
		expiredToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Partner, authCtx)
		partnerID := setupPartnerUser(t, ctx, "partner13@example.com")
		roomID := setupTestRoom(t, ctx)

		req := ta.NewCheckPartnerRoomAccessRequest(t, ctx, testServerURL, partnerID, roomID, nil, expiredToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		partnerID := setupPartnerUser(t, ctx, "partner14@example.com")
		roomID := setupTestRoom(t, ctx)

		req := ta.NewCheckPartnerRoomAccessRequest(t, ctx, testServerURL, partnerID, roomID, nil, "invalid-token-12345")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when standard user tries to access", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create standard user (not partner or admin)
		standardAccessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)

		partnerID := setupPartnerUser(t, ctx, "partner15@example.com")
		roomID := setupTestRoom(t, ctx)

		req := ta.NewCheckPartnerRoomAccessRequest(t, ctx, testServerURL, partnerID, roomID, nil, standardAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}
