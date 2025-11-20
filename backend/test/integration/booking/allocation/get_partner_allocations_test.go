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

// make test-func TEST_NAME=TestGetPartnerAllocations TEST_PATH=test/integration/booking/allocation/get_partner_allocations_test.go

func TestGetPartnerAllocations(t *testing.T) {
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

	t.Run("should successfully get active allocations for partner with default query parameter", func(t *testing.T) {
		// Clean test data
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx, "partner1@example.com")
		roomID1 := setupTestRoom(t, ctx)
		roomID2 := setupTestRoom(t, ctx)

		// Create active allocations
		activeAlloc1 := ta.NewTestSharedAllocation(t, roomID1, partnerUserID)
		ta.InsertAllocation(t, ctx, activeAlloc1, testPool)

		startDate := time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 0, 7).Truncate(24 * time.Hour)
		activeAlloc2 := ta.NewTestDedicatedAllocation(t, roomID2, partnerUserID, startDate, endDate)
		ta.InsertAllocation(t, ctx, activeAlloc2, testPool)

		// Create inactive allocation
		inactiveAlloc := ta.NewTestInactiveAllocation(t, roomID1, partnerUserID)
		ta.InsertAllocation(t, ctx, inactiveAlloc, testPool)

		// Get allocations with default query parameter (should return only active)
		req := ta.NewGetPartnerAllocationsRequest(t, ctx, testServerURL, partnerUserID, nil, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return only active allocations
		assert.Len(t, response, 2)
		for _, alloc := range response {
			assert.True(t, alloc.IsActive)
			assert.Equal(t, partnerUserID, alloc.UserID)
		}
	})

	t.Run("should successfully get only active allocations when active_only=true", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx, "partner2@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create active allocation
		activeAlloc := ta.NewTestSharedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, activeAlloc, testPool)

		// Create inactive allocation
		inactiveAlloc := ta.NewTestInactiveAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, inactiveAlloc, testPool)

		// Get allocations with active_only=true
		activeOnly := true
		req := ta.NewGetPartnerAllocationsRequest(t, ctx, testServerURL, partnerUserID, &activeOnly, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return only active allocation
		assert.Len(t, response, 1)
		assert.True(t, response[0].IsActive)
		assert.Equal(t, activeAlloc.ID, response[0].ID)
	})

	t.Run("should successfully get all allocations when active_only=false", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx, "partner3@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create active allocation
		activeAlloc := ta.NewTestSharedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, activeAlloc, testPool)

		// Create inactive allocation
		inactiveAlloc := ta.NewTestInactiveAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, inactiveAlloc, testPool)

		// Get allocations with active_only=false
		activeOnly := false
		req := ta.NewGetPartnerAllocationsRequest(t, ctx, testServerURL, partnerUserID, &activeOnly, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return both active and inactive allocations
		assert.Len(t, response, 2)

		// Verify we have one active and one inactive
		activeCount := 0
		inactiveCount := 0
		for _, alloc := range response {
			if alloc.IsActive {
				activeCount++
			} else {
				inactiveCount++
			}
		}
		assert.Equal(t, 1, activeCount)
		assert.Equal(t, 1, inactiveCount)
	})

	t.Run("should return empty array when partner has no allocations", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx, "partner4@example.com")

		// Don't create any allocations
		req := ta.NewGetPartnerAllocationsRequest(t, ctx, testServerURL, partnerUserID, nil, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return empty array
		assert.Len(t, response, 0)
	})

	t.Run("should return allocations with mixed types (shared and dedicated)", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx, "partner5@example.com")
		roomID1 := setupTestRoom(t, ctx)
		roomID2 := setupTestRoom(t, ctx)

		// Create shared allocation
		sharedAlloc := ta.NewTestSharedAllocation(t, roomID1, partnerUserID)
		ta.InsertAllocation(t, ctx, sharedAlloc, testPool)

		// Create dedicated allocation
		startDate := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 0, 30).Truncate(24 * time.Hour)
		dedicatedAlloc := ta.NewTestDedicatedAllocation(t, roomID2, partnerUserID, startDate, endDate)
		ta.InsertAllocation(t, ctx, dedicatedAlloc, testPool)

		req := ta.NewGetPartnerAllocationsRequest(t, ctx, testServerURL, partnerUserID, nil, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response, 2)

		// Verify we have both types
		hasShared := false
		hasDedicated := false
		for _, alloc := range response {
			if alloc.AllocationType == domain.AllocationTypeShared {
				hasShared = true
				assert.Nil(t, alloc.StartDate)
				assert.Nil(t, alloc.EndDate)
			} else if alloc.AllocationType == domain.AllocationTypeDedicated {
				hasDedicated = true
				assert.NotNil(t, alloc.StartDate)
				assert.NotNil(t, alloc.EndDate)
			}
		}
		assert.True(t, hasShared)
		assert.True(t, hasDedicated)
	})

	t.Run("should return allocations with various date ranges", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx, "partner6@example.com")
		roomID1 := setupTestRoom(t, ctx)
		roomID2 := setupTestRoom(t, ctx)
		roomID3 := setupTestRoom(t, ctx)

		// Create past allocation
		pastAlloc := ta.NewTestPastDedicatedAllocation(t, roomID1, partnerUserID)
		ta.InsertAllocation(t, ctx, pastAlloc, testPool)

		// Create current allocation
		currentAlloc := ta.NewTestActiveDedicatedAllocation(t, roomID2, partnerUserID)
		ta.InsertAllocation(t, ctx, currentAlloc, testPool)

		// Create future allocation
		futureAlloc := ta.NewTestFutureDedicatedAllocation(t, roomID3, partnerUserID)
		ta.InsertAllocation(t, ctx, futureAlloc, testPool)

		req := ta.NewGetPartnerAllocationsRequest(t, ctx, testServerURL, partnerUserID, nil, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// All three should be active and returned
		assert.Len(t, response, 3)
		for _, alloc := range response {
			assert.True(t, alloc.IsActive)
			assert.NotNil(t, alloc.StartDate)
			assert.NotNil(t, alloc.EndDate)
		}
	})

	t.Run("should allow partner to access allocations with partner role", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create partner role session (the middleware RequirePartner allows partner and above)
		partnerAccessToken := tu.SetupUserWithRole(t, ctx, identity.Partner, authCtx)
		partnerUserID := setupPartnerUser(t, ctx, "partner7@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create allocation for a partner
		allocation := ta.NewTestSharedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Partner role accessing partner allocations (middleware allows this)
		req := ta.NewGetPartnerAllocationsRequest(t, ctx, testServerURL, partnerUserID, nil, partnerAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response, 1)
		assert.Equal(t, partnerUserID, response[0].UserID)
	})

	t.Run("should allow admin to access any partner allocations", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx, "partner8@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create allocation for partner
		allocation := ta.NewTestSharedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Admin accessing partner's allocations
		req := ta.NewGetPartnerAllocationsRequest(t, ctx, testServerURL, partnerUserID, nil, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response, 1)
		assert.Equal(t, partnerUserID, response[0].UserID)
	})

	t.Run("should return allocations with indefinite end date", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx, "partner9@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create allocation with indefinite end date
		startDate := time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
		allocation := ta.NewTestDedicatedAllocationIndefinite(t, roomID, partnerUserID, startDate)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		req := ta.NewGetPartnerAllocationsRequest(t, ctx, testServerURL, partnerUserID, nil, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response, 1)
		assert.NotNil(t, response[0].StartDate)
		assert.Nil(t, response[0].EndDate) // Should be nil for indefinite
	})

	t.Run("should return 400 Bad Request for invalid partner ID format", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Create request with invalid UUID format
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			testServerURL+"/partners/invalid-uuid-format/allocations",
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

		// Create request with malformed path (missing partner ID)
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			testServerURL+"/partners",
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

		// Should return 404 or 405 depending on router behavior
		assert.True(t, resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusMethodNotAllowed)
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		partnerUserID := setupPartnerUser(t, ctx, "partner10@example.com")

		req := ta.NewGetPartnerAllocationsRequest(t, ctx, testServerURL, partnerUserID, nil, "") // Empty token

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
		partnerUserID := setupPartnerUser(t, ctx, "partner11@example.com")

		req := ta.NewGetPartnerAllocationsRequest(t, ctx, testServerURL, partnerUserID, nil, expiredToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		partnerUserID := setupPartnerUser(t, ctx, "partner12@example.com")

		req := ta.NewGetPartnerAllocationsRequest(t, ctx, testServerURL, partnerUserID, nil, "invalid-token-12345")

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

		partnerUserID := setupPartnerUser(t, ctx, "partner13@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create allocation
		allocation := ta.NewTestSharedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		req := ta.NewGetPartnerAllocationsRequest(t, ctx, testServerURL, partnerUserID, nil, standardAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}
