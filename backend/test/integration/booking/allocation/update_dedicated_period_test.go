package allocation_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
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

// make test-func TEST_NAME=TestUpdateDedicatedPeriod TEST_PATH=test/integration/booking/allocation/update_dedicated_period_test.go

func TestUpdateDedicatedPeriod(t *testing.T) {
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

	t.Run("should successfully update dedicated period with both start and end date", func(t *testing.T) {
		// Clean test data
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authentication
		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Setup test entities
		partnerUserID := setupPartnerUser(t, ctx, "partner1@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create and insert dedicated allocation
		allocation := ta.NewTestActiveDedicatedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Create update request with new dates
		newStartDate := time.Now().Add(30 * 24 * time.Hour).Truncate(24 * time.Hour)
		newEndDate := time.Now().Add(60 * 24 * time.Hour).Truncate(24 * time.Hour)
		request := domain.UpdateDedicatedAllocationRequest{
			StartDate: &newStartDate,
			EndDate:   &newEndDate,
		}

		// Make HTTP request
		req := ta.NewUpdateDedicatedAllocationRequest(t, ctx, testServerURL, allocation.ID, request, adminAccessToken)
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
		assert.NotNil(t, response.EndDate)
		assert.WithinDuration(t, newStartDate, *response.StartDate, time.Second)
		assert.WithinDuration(t, newEndDate, *response.EndDate, time.Second)

		// Verify database persistence
		updated, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)

		assert.WithinDuration(t, newStartDate, *updated.StartDate, time.Second)
		assert.WithinDuration(t, newEndDate, *updated.EndDate, time.Second)
	})

	t.Run("should successfully update dedicated period to indefinite (no end date)", func(t *testing.T) {
		// Clean test data
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authentication
		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Setup test entities
		partnerUserID := setupPartnerUser(t, ctx, "partner2@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create and insert dedicated allocation with end date
		allocation := ta.NewTestActiveDedicatedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Create update request with only start date (indefinite)
		newStartDate := time.Now().Add(15 * 24 * time.Hour).Truncate(24 * time.Hour)
		request := domain.UpdateDedicatedAllocationRequest{
			StartDate: &newStartDate,
			EndDate:   nil,
		}

		// Make HTTP request
		req := ta.NewUpdateDedicatedAllocationRequest(t, ctx, testServerURL, allocation.ID, request, adminAccessToken)
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
		assert.NotNil(t, response.StartDate)
		assert.Nil(t, response.EndDate)

		// Verify database persistence
		updated, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.Nil(t, updated.EndDate)
	})

	t.Run("should successfully extend existing dedicated period", func(t *testing.T) {
		// Clean test data
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authentication
		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Setup test entities
		partnerUserID := setupPartnerUser(t, ctx, "partner3@example.com")
		roomID := setupTestRoom(t, ctx)

		startDate := time.Now().AddDate(0, 0, 10)
		endDate := time.Now().AddDate(0, 0, 20)

		// Create and insert dedicated allocation
		allocation := ta.NewTestDedicatedAllocation(t, roomID, partnerUserID, startDate, endDate)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Create update request extending the end date
		extendedEndDate := allocation.EndDate.Add(30 * 24 * time.Hour)
		request := domain.UpdateDedicatedAllocationRequest{
			StartDate: allocation.StartDate,
			EndDate:   &extendedEndDate,
		}

		// Make HTTP request
		req := ta.NewUpdateDedicatedAllocationRequest(t, ctx, testServerURL, allocation.ID, request, adminAccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response body
		var response domain.RoomAllocationResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, extendedEndDate.Unix(), response.EndDate.Unix())
	})

	t.Run("should return 400 for invalid UUID format", func(t *testing.T) {
		// Clean test data
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authentication
		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Create request with invalid UUID
		invalidID := "not-a-valid-uuid"
		newStartDate := time.Now().Add(15 * 24 * time.Hour)
		newEndDate := time.Now().Add(45 * 24 * time.Hour)
		request := domain.UpdateDedicatedAllocationRequest{
			StartDate: &newStartDate,
			EndDate:   &newEndDate,
		}

		// Create HTTP request manually with invalid ID
		jsonBody, err := json.Marshal(request)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, http.MethodPut, testServerURL+"/allocations/"+invalidID+"/period",
			bytes.NewReader(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: adminAccessToken,
		})

		// Make HTTP request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid JSON body", func(t *testing.T) {
		// Clean test data
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authentication
		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Setup test entities
		partnerUserID := setupPartnerUser(t, ctx, "partner4@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create and insert dedicated allocation
		allocation := ta.NewTestActiveDedicatedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Create request with invalid JSON
		invalidJSON := `{"start_date": "not-a-date", "end_date": 12345}`

		req, err := http.NewRequestWithContext(ctx, http.MethodPut, testServerURL+"/allocations/"+allocation.ID.String()+"/period",
			strings.NewReader(invalidJSON))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: adminAccessToken,
		})

		// Make HTTP request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 415 for unsupported media type", func(t *testing.T) {
		// Clean test data
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authentication
		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Setup test entities
		partnerUserID := setupPartnerUser(t, ctx, "partner5@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create and insert dedicated allocation
		allocation := ta.NewTestActiveDedicatedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Create request without Content-Type header
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, testServerURL+"/allocations/"+allocation.ID.String()+"/period",
			strings.NewReader(`{}`))
		require.NoError(t, err)
		// Deliberately not setting Content-Type
		req.AddCookie(&http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: adminAccessToken,
		})

		// Make HTTP request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
	})

	t.Run("should return 400 when allocation does not exist", func(t *testing.T) {
		// Clean test data
		ta.ClearAllocationTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authentication
		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Use a valid UUID that doesn't exist in database
		nonExistentID := uuid.New()

		// Create update request
		newStartDate := time.Now().Add(15 * 24 * time.Hour)
		newEndDate := time.Now().Add(45 * 24 * time.Hour)
		request := domain.UpdateDedicatedAllocationRequest{
			StartDate: &newStartDate,
			EndDate:   &newEndDate,
		}

		// Make HTTP request
		req := ta.NewUpdateDedicatedAllocationRequest(t, ctx, testServerURL, nonExistentID, request, adminAccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 401 when not authenticated", func(t *testing.T) {
		// Clean test data
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Setup test entities (without authentication)
		partnerUserID := uuid.New()
		roomID := setupTestRoom(t, ctx)

		// Create and insert dedicated allocation
		allocation := ta.NewTestActiveDedicatedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Create update request
		newStartDate := time.Now().Add(15 * 24 * time.Hour)
		newEndDate := time.Now().Add(45 * 24 * time.Hour)
		request := domain.UpdateDedicatedAllocationRequest{
			StartDate: &newStartDate,
			EndDate:   &newEndDate,
		}

		// Create request without authentication token
		jsonBody, err := json.Marshal(request)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, http.MethodPut, testServerURL+"/allocations/"+allocation.ID.String()+"/period",
			bytes.NewReader(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		// No authentication cookie

		// Make HTTP request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 400 when trying to update shared allocation", func(t *testing.T) {
		// Clean test data
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authentication
		adminAccessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Setup test entities
		partnerUserID := setupPartnerUser(t, ctx, "partner6@example.com")
		roomID := setupTestRoom(t, ctx)

		// Create and insert SHARED allocation (not dedicated)
		allocation := ta.NewTestSharedAllocation(t, roomID, partnerUserID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Try to update shared allocation's period (should fail)
		newStartDate := time.Now().Add(15 * 24 * time.Hour)
		newEndDate := time.Now().Add(45 * 24 * time.Hour)
		request := domain.UpdateDedicatedAllocationRequest{
			StartDate: &newStartDate,
			EndDate:   &newEndDate,
		}

		// Make HTTP request
		req := ta.NewUpdateDedicatedAllocationRequest(t, ctx, testServerURL, allocation.ID, request, adminAccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response - should fail because it's a shared allocation
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
