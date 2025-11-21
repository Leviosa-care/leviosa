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
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"
	ta "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/allocation"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetAllocation TEST_PATH=test/integration/booking/allocation/get_allocation_test.go

func TestGetAllocation(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Setup helper for creating test room
	setupTestRoom := func(t *testing.T, ctx context.Context) uuid.UUID {
		// Create building with encryption
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Create room with encryption
		room := tr.NewTestRoomWithBuilding(t, building.ID)
		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		return room.ID
	}

	// Setup helper for creating partner user
	setupPartnerUser := func(t *testing.T, ctx context.Context, email string) uuid.UUID {
		// Create user with encryption
		user := th.NewTestUser(t, email, "John", "DOE")
		user.Role = identity.Partner.String()
		userEncx, err := userDomain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = th.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Create partner
		partner := th.NewTestPartner(t, user.ID)
		partner.StripeAccountStatus = userDomain.StripeAccountStatusActive
		partner.StripeOnboardingComplete = true
		partnerEncx, err := userDomain.ProcessPartnerEncx(ctx, crypto, partner)
		require.NoError(t, err)
		err = th.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		return user.ID
	}

	t.Run("should successfully get shared allocation by ID", func(t *testing.T) {
		// Clean test data
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authentication
		partnerAccessToken := tu.SetupUserWithRole(t, ctx, identity.Partner, authCtx)

		// Setup test entities
		partnerUserID := setupPartnerUser(t, ctx, "partner1@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create and insert test allocation
		allocation := ta.NewTestSharedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool, crypto)

		// Make HTTP request
		req := ta.NewGetAllocationRequest(t, ctx, testServerURL, allocation.ID, partnerAccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response body
		var response domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, allocation.ID, response.ID)
		assert.Equal(t, allocation.RoomID, response.RoomID)
		assert.Equal(t, allocation.UserID, response.UserID)
		assert.Equal(t, allocation.AllocationType, response.AllocationType)
		assert.Equal(t, allocation.IsActive, response.IsActive)
		assert.Nil(t, response.StartDate)
		assert.Nil(t, response.EndDate)
	})

	t.Run("should successfully get dedicated allocation by ID with time bounds", func(t *testing.T) {
		// Clean test data
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authentication
		partnerAccessToken := tu.SetupUserWithRole(t, ctx, identity.Partner, authCtx)

		// Setup test entities
		partnerUserID := setupPartnerUser(t, ctx, "partner2@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create and insert dedicated allocation with time bounds
		allocation := ta.NewTestActiveDedicatedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool, crypto)

		// Make HTTP request
		req := ta.NewGetAllocationRequest(t, ctx, testServerURL, allocation.ID, partnerAccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response body
		var response domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, allocation.ID, response.ID)
		assert.Equal(t, allocation.RoomID, response.RoomID)
		assert.Equal(t, allocation.UserID, response.UserID)
		assert.Equal(t, allocation.AllocationType, response.AllocationType)
		assert.Equal(t, allocation.IsActive, response.IsActive)
		assert.NotNil(t, response.StartDate)
		assert.NotNil(t, response.EndDate)
		assert.Equal(t, allocation.StartDate.Unix(), response.StartDate.Unix())
		assert.Equal(t, allocation.EndDate.Unix(), response.EndDate.Unix())
	})

	t.Run("should successfully get dedicated allocation without end date", func(t *testing.T) {
		// Clean test data
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authentication
		partnerAccessToken := tu.SetupUserWithRole(t, ctx, identity.Partner, authCtx)

		// Setup test entities
		partnerUserID := setupPartnerUser(t, ctx, "partner3@example.com")
		roomID := setupTestRoom(t, ctx)

		startDate := time.Now().AddDate(0, 0, 10)

		// Create and insert dedicated allocation without end date (indefinite)
		allocation := ta.NewTestDedicatedAllocationIndefinite(t, roomID, partnerUserID, startDate)
		ta.InsertAllocation(t, ctx, allocation, testPool, crypto)

		// Make HTTP request
		req := ta.NewGetAllocationRequest(t, ctx, testServerURL, allocation.ID, partnerAccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response body
		var response domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, allocation.ID, response.ID)
		assert.Equal(t, domain.AllocationTypeDedicated, response.AllocationType)
		assert.NotNil(t, response.StartDate)
		assert.Nil(t, response.EndDate)
	})

	t.Run("should successfully get inactive allocation", func(t *testing.T) {
		// Clean test data
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authentication
		partnerAccessToken := tu.SetupUserWithRole(t, ctx, identity.Partner, authCtx)

		// Setup test entities
		partnerUserID := setupPartnerUser(t, ctx, "partner4@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create and insert inactive allocation
		allocation := ta.NewTestInactiveAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool, crypto)

		// Make HTTP request
		req := ta.NewGetAllocationRequest(t, ctx, testServerURL, allocation.ID, partnerAccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response body
		var response domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, allocation.ID, response.ID)
		assert.False(t, response.IsActive)
	})

	t.Run("should return 400 for invalid UUID format", func(t *testing.T) {
		// Clean test data
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authentication
		partnerAccessToken := tu.SetupUserWithRole(t, ctx, identity.Partner, authCtx)

		// Create request with invalid UUID
		invalidID := "not-a-valid-uuid"
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+"/allocations/"+invalidID, nil)
		require.NoError(t, err)
		if partnerAccessToken != "" {
			req.AddCookie(&http.Cookie{
				Name:  ck.AccessTokenCookieName,
				Value: partnerAccessToken,
			})
		}

		// Make HTTP request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 404 when allocation does not exist", func(t *testing.T) {
		// Clean test data
		ta.ClearAllocationTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authentication
		partnerAccessToken := tu.SetupUserWithRole(t, ctx, identity.Partner, authCtx)

		// Use a valid UUID that doesn't exist in database
		nonExistentID := uuid.New()

		// Make HTTP request
		req := ta.NewGetAllocationRequest(t, ctx, testServerURL, nonExistentID, partnerAccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 401 when not authenticated", func(t *testing.T) {
		// Clean test data
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Setup test entities (without authentication)
		partnerUserID := uuid.New()
		roomID := setupTestRoom(t, ctx)

		// Create and insert test allocation
		allocation := ta.NewTestSharedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool, crypto)

		// Create request without authentication token
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+"/allocations/"+allocation.ID.String(), nil)
		require.NoError(t, err)

		// Make HTTP request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
