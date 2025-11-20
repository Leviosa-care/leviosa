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

// make test-func TEST_NAME=TestDeactivateAllocation TEST_PATH=test/integration/booking/allocation/deactivate_allocation_test.go

func TestDeactivateAllocation(t *testing.T) {
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

	t.Run("should successfully deactivate an active allocation", func(t *testing.T) {
		// Clean test data
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx)
		roomID := setupTestRoom(t, ctx)

		// Create active allocation
		startDate := time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 0, 7).Truncate(24 * time.Hour)
		allocation := ta.NewTestDedicatedAllocation(t, roomID, partnerUserID, startDate, endDate)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Verify allocation is active before deactivation
		activeCountBefore := ta.CountActiveAllocations(t, ctx, testPool)
		assert.Equal(t, 1, activeCountBefore)

		// Deactivate allocation
		req := ta.NewDeactivateAllocationRequest(t, ctx, testServerURL, allocation.ID, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response data
		assert.Equal(t, allocation.ID, response.ID)
		assert.Equal(t, allocation.RoomID, response.RoomID)
		assert.Equal(t, allocation.UserID, response.UserID)
		assert.Equal(t, domain.AllocationTypeDedicated, response.AllocationType)
		assert.False(t, response.IsActive)

		// Verify allocation is deactivated in database
		deactivatedAllocation, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.False(t, deactivatedAllocation.IsActive)

		// Verify active allocation count decreased
		activeCountAfter := ta.CountActiveAllocations(t, ctx, testPool)
		assert.Equal(t, 0, activeCountAfter)

		inactiveCountAfter := ta.CountInactiveAllocations(t, ctx, testPool)
		assert.Equal(t, 1, inactiveCountAfter)
	})

	t.Run("should successfully deactivate a shared allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx)
		roomID := setupTestRoom(t, ctx)

		// Create shared allocation
		allocation := ta.NewTestSharedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Deactivate allocation
		req := ta.NewDeactivateAllocationRequest(t, ctx, testServerURL, allocation.ID, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, allocation.ID, response.ID)
		assert.Equal(t, domain.AllocationTypeShared, response.AllocationType)
		assert.False(t, response.IsActive)

		// Verify in database
		deactivatedAllocation, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.False(t, deactivatedAllocation.IsActive)
	})

	t.Run("should be idempotent when deactivating already inactive allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx)
		roomID := setupTestRoom(t, ctx)

		// Create already inactive allocation
		allocation := ta.NewTestInactiveAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Verify allocation is already inactive
		assert.False(t, allocation.IsActive)

		// Deactivate allocation again (should be idempotent)
		req := ta.NewDeactivateAllocationRequest(t, ctx, testServerURL, allocation.ID, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.False(t, response.IsActive)

		// Verify still inactive in database
		deactivatedAllocation, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.False(t, deactivatedAllocation.IsActive)
	})

	t.Run("should return 400 Bad Request for non-existent allocation ID", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Use a random UUID that doesn't exist in database
		nonExistentID := uuid.New()

		req := ta.NewDeactivateAllocationRequest(t, ctx, testServerURL, nonExistentID, adminAccessToken)

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

	t.Run("should return 400 Bad Request for invalid allocation ID format", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Create request with invalid UUID format
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			testServerURL+"/allocations/invalid-uuid-format/deactivate",
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

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		partnerUserID := setupPartnerUser(t, ctx)
		roomID := setupTestRoom(t, ctx)

		// Create allocation
		allocation := ta.NewTestSharedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		req := ta.NewDeactivateAllocationRequest(t, ctx, testServerURL, allocation.ID, "") // Empty token

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
		partnerUserID := setupPartnerUser(t, ctx)
		roomID := setupTestRoom(t, ctx)

		// Create allocation
		allocation := ta.NewTestSharedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		req := ta.NewDeactivateAllocationRequest(t, ctx, testServerURL, allocation.ID, expiredToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		partnerUserID := setupPartnerUser(t, ctx)
		roomID := setupTestRoom(t, ctx)

		// Create allocation
		allocation := ta.NewTestSharedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		req := ta.NewDeactivateAllocationRequest(t, ctx, testServerURL, allocation.ID, "invalid-token-12345")

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

		// Create allocation
		allocation := ta.NewTestSharedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		req := ta.NewDeactivateAllocationRequest(t, ctx, testServerURL, allocation.ID, standardAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role (partner user)", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create partner user (not admin)
		partnerAccessToken := tu.SetupUserWithRole(t, ctx, identity.Partner, authCtx)

		partnerUserID := setupPartnerUser(t, ctx)
		roomID := setupTestRoom(t, ctx)

		// Create allocation
		allocation := ta.NewTestSharedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		req := ta.NewDeactivateAllocationRequest(t, ctx, testServerURL, allocation.ID, partnerAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should successfully deactivate allocation with indefinite end date", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx)
		roomID := setupTestRoom(t, ctx)

		// Create allocation with indefinite end date
		startDate := time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
		allocation := ta.NewTestDedicatedAllocationIndefinite(t, roomID, partnerUserID, startDate)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Deactivate allocation
		req := ta.NewDeactivateAllocationRequest(t, ctx, testServerURL, allocation.ID, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.False(t, response.IsActive)
		assert.Nil(t, response.EndDate) // Should still be nil

		// Verify in database
		deactivatedAllocation, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.False(t, deactivatedAllocation.IsActive)
		assert.Nil(t, deactivatedAllocation.EndDate)
	})

	t.Run("should successfully deactivate past allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx)
		roomID := setupTestRoom(t, ctx)

		// Create past allocation (already ended)
		allocation := ta.NewTestPastDedicatedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Deactivate allocation
		req := ta.NewDeactivateAllocationRequest(t, ctx, testServerURL, allocation.ID, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.False(t, response.IsActive)

		// Verify in database
		deactivatedAllocation, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.False(t, deactivatedAllocation.IsActive)
	})

	t.Run("should successfully deactivate future allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		partnerUserID := setupPartnerUser(t, ctx)
		roomID := setupTestRoom(t, ctx)

		// Create future allocation (hasn't started yet)
		allocation := ta.NewTestFutureDedicatedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Deactivate allocation
		req := ta.NewDeactivateAllocationRequest(t, ctx, testServerURL, allocation.ID, adminAccessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.False(t, response.IsActive)

		// Verify in database
		deactivatedAllocation, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.False(t, deactivatedAllocation.IsActive)
	})
}
