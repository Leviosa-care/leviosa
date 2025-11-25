package availability_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	ts "github.com/Leviosa-care/leviosa/backend/test/helpers/booking"
	ta "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetAvailableSlots TEST_PATH=test/integration/booking/availability/get_available_slots_test.go

func TestGetAvailableSlots(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully retrieve available slots with no filters", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create test data
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)

		room := tr.NewTestRoomWithBuilding(t, building.ID)
		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		tr.InsertRoomEncx(t, ctx, testPool, roomEncx)

		// Setup authenticated partners with allocation to the room
		accessToken1, userID1 := ts.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner1@leviosa.care", room.ID, testPool, redisClient, crypto)
		_, userID2 := ts.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner2@leviosa.care", room.ID, testPool, redisClient, crypto)

		// Create multiple availabilities with properly allocated users
		availability1 := ta.NewTestAvailabilityWithParams(t, userID1, room.ID, time.Now().Add(24*time.Hour), time.Now().Add(26*time.Hour), "Consultation", intPtr(15000), 1, domain.AvailabilityStatusAvailable)
		availability1Encx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability1)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availability1Encx, testPool)

		availability2 := ta.NewTestAvailabilityWithParams(t, userID2, room.ID, time.Now().Add(48*time.Hour), time.Now().Add(50*time.Hour), "Consultation", intPtr(15000), 1, domain.AvailabilityStatusAvailable)
		availability2Encx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability2)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availability2Encx, testPool)

		// Make request with one of the tokens
		req := ta.NewGetAvailableSlotsRequest(t, ctx, testServerURL, nil, accessToken1)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return 2 available slots
		assert.Len(t, response, 2)

		// Verify response structure
		for _, slot := range response {
			assert.NotEmpty(t, slot.ID)
			assert.NotEmpty(t, slot.UserID)
			assert.NotEmpty(t, slot.RoomID)
			assert.NotZero(t, slot.StartTime)
			assert.NotZero(t, slot.EndTime)
			assert.Equal(t, domain.AvailabilityStatusAvailable, slot.Status)
		}
	})

	t.Run("should filter availabilities by time range", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create test data
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)

		room := tr.NewTestRoomWithBuilding(t, building.ID)
		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		tr.InsertRoomEncx(t, ctx, testPool, roomEncx)

		now := time.Now().Truncate(time.Second)

		// Setup authenticated partners with allocation to the room
		accessToken1, userID1 := ts.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner1@leviosa.care", room.ID, testPool, redisClient, crypto)
		_, userID2 := ts.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner2@leviosa.care", room.ID, testPool, redisClient, crypto)

		// Create availability within the filter range
		availabilityInRange := ta.NewTestAvailabilityWithParams(t, userID1, room.ID, now.Add(25*time.Hour), now.Add(27*time.Hour), "Consultation", intPtr(15000), 1, domain.AvailabilityStatusAvailable)
		availabilityInRangeEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availabilityInRange)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availabilityInRangeEncx, testPool)

		// Create availability outside the filter range
		availabilityOutOfRange := ta.NewTestAvailabilityWithParams(t, userID2, room.ID, now.Add(72*time.Hour), now.Add(74*time.Hour), "Consultation", intPtr(15000), 1, domain.AvailabilityStatusAvailable)
		availabilityOutOfRangeEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availabilityOutOfRange)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availabilityOutOfRangeEncx, testPool)

		// Filter by time range
		queryParams := map[string]string{
			"start_time": now.Add(24 * time.Hour).Format(time.RFC3339),
			"end_time":   now.Add(48 * time.Hour).Format(time.RFC3339),
		}

		req := ta.NewGetAvailableSlotsRequest(t, ctx, testServerURL, queryParams, accessToken1)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should only return 1 slot within the time range
		assert.Len(t, response, 1)
		assert.True(t, response[0].StartTime.After(now.Add(24*time.Hour)))
		assert.True(t, response[0].EndTime.Before(now.Add(48*time.Hour)))
	})

	t.Run("should filter availabilities by room_id", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create test data with multiple rooms
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)

		room1 := tr.NewTestRoomWithBuilding(t, building.ID)
		room1Encx, err := domain.ProcessRoomEncx(ctx, crypto, room1)
		require.NoError(t, err)
		tr.InsertRoomEncx(t, ctx, testPool, room1Encx)

		room2 := tr.NewTestRoomWithBuilding(t, building.ID)
		room2.ID = uuid.New()
		room2Encx, err := domain.ProcessRoomEncx(ctx, crypto, room2)
		require.NoError(t, err)
		tr.InsertRoomEncx(t, ctx, testPool, room2Encx)

		// Setup authenticated partners with allocation to different rooms
		accessToken1, userID1 := ts.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner1@leviosa.care", room1.ID, testPool, redisClient, crypto)
		_, userID2 := ts.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner2@leviosa.care", room2.ID, testPool, redisClient, crypto)

		// Create availabilities for different rooms
		availabilityRoom1 := ta.NewTestAvailabilityWithParams(t, userID1, room1.ID, time.Now().Add(24*time.Hour), time.Now().Add(26*time.Hour), "Consultation", intPtr(15000), 1, domain.AvailabilityStatusAvailable)
		availabilityRoom1Encx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availabilityRoom1)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availabilityRoom1Encx, testPool)

		availabilityRoom2 := ta.NewTestAvailabilityWithParams(t, userID2, room2.ID, time.Now().Add(24*time.Hour), time.Now().Add(26*time.Hour), "Consultation", intPtr(15000), 1, domain.AvailabilityStatusAvailable)
		availabilityRoom2Encx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availabilityRoom2)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availabilityRoom2Encx, testPool)

		// Filter by room1
		queryParams := map[string]string{
			"room_id": room1.ID.String(),
		}

		req := ta.NewGetAvailableSlotsRequest(t, ctx, testServerURL, queryParams, accessToken1)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should only return 1 slot for room1
		assert.Len(t, response, 1)
		assert.Equal(t, room1.ID, response[0].RoomID)
	})

	t.Run("should filter availabilities by partner_id", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create test data
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)

		room := tr.NewTestRoomWithBuilding(t, building.ID)
		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		tr.InsertRoomEncx(t, ctx, testPool, roomEncx)

		// Setup authenticated partners with allocation to the room
		accessToken1, userID1 := ts.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner1@leviosa.care", room.ID, testPool, redisClient, crypto)
		_, userID2 := ts.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner2@leviosa.care", room.ID, testPool, redisClient, crypto)

		// Create availabilities for different partners
		availabilityPartner1 := ta.NewTestAvailabilityWithParams(t, userID1, room.ID, time.Now().Add(24*time.Hour), time.Now().Add(26*time.Hour), "Consultation", intPtr(15000), 1, domain.AvailabilityStatusAvailable)
		availabilityPartner1Encx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availabilityPartner1)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availabilityPartner1Encx, testPool)

		availabilityPartner2 := ta.NewTestAvailabilityWithParams(t, userID2, room.ID, time.Now().Add(24*time.Hour), time.Now().Add(26*time.Hour), "Consultation", intPtr(15000), 1, domain.AvailabilityStatusAvailable)
		availabilityPartner2Encx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availabilityPartner2)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availabilityPartner2Encx, testPool)

		// Filter by partner1
		queryParams := map[string]string{
			"partner_id": userID1.String(),
		}

		req := ta.NewGetAvailableSlotsRequest(t, ctx, testServerURL, queryParams, accessToken1)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should only return 1 slot for partner1
		assert.Len(t, response, 1)
		assert.Equal(t, userID1, response[0].UserID)
	})

	t.Run("should apply limit parameter", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create test data
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)

		room := tr.NewTestRoomWithBuilding(t, building.ID)
		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		tr.InsertRoomEncx(t, ctx, testPool, roomEncx)

		// Setup authenticated partner with allocation to the room
		accessToken, userID := ts.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@leviosa.care", room.ID, testPool, redisClient, crypto)

		// Create multiple availabilities
		for i := 0; i < 5; i++ {
			availability := ta.NewTestAvailabilityWithParams(t, userID, room.ID, time.Now().Add(time.Duration(i+1)*24*time.Hour), time.Now().Add(time.Duration(i+1)*24*time.Hour+2*time.Hour), "Consultation", intPtr(15000), 1, domain.AvailabilityStatusAvailable)
			availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
			require.NoError(t, err)
			ta.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)
		}

		// Apply limit of 3
		queryParams := map[string]string{
			"limit": "3",
		}

		req := ta.NewGetAvailableSlotsRequest(t, ctx, testServerURL, queryParams, accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return only 3 slots
		assert.Len(t, response, 3)
	})

	t.Run("should return empty results when no availabilities match filters", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)

		// Filter by non-existent room
		nonExistentRoomID := uuid.New()
		queryParams := map[string]string{
			"room_id": nonExistentRoomID.String(),
		}

		req := ta.NewGetAvailableSlotsRequest(t, ctx, testServerURL, queryParams, accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return empty array
		assert.Len(t, response, 0)
	})

	t.Run("should combine multiple filters correctly", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create test data
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)

		room1 := tr.NewTestRoomWithBuilding(t, building.ID)
		room1Encx, err := domain.ProcessRoomEncx(ctx, crypto, room1)
		require.NoError(t, err)
		tr.InsertRoomEncx(t, ctx, testPool, room1Encx)

		room2 := tr.NewTestRoomWithBuilding(t, building.ID)
		room2.ID = uuid.New()
		room2Encx, err := domain.ProcessRoomEncx(ctx, crypto, room2)
		require.NoError(t, err)
		tr.InsertRoomEncx(t, ctx, testPool, room2Encx)

		now := time.Now().Truncate(time.Second)

		// Setup authenticated partners with allocation to different rooms
		accessToken1, userID1 := ts.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner1@leviosa.care", room1.ID, testPool, redisClient, crypto)
		_, userID2 := ts.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner2@leviosa.care", room2.ID, testPool, redisClient, crypto)

		// Create matching availability (correct room, partner, and time range)
		matchingAvailability := ta.NewTestAvailabilityWithParams(t, userID1, room1.ID, now.Add(25*time.Hour), now.Add(27*time.Hour), "Consultation", intPtr(15000), 1, domain.AvailabilityStatusAvailable)
		matchingAvailabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, matchingAvailability)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, matchingAvailabilityEncx, testPool)

		// Create non-matching availabilities
		wrongRoomAvailability := ta.NewTestAvailabilityWithParams(t, userID2, room2.ID, now.Add(25*time.Hour), now.Add(27*time.Hour), "Consultation", intPtr(15000), 1, domain.AvailabilityStatusAvailable)
		wrongRoomAvailabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, wrongRoomAvailability)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, wrongRoomAvailabilityEncx, testPool)

		wrongTimeAvailability := ta.NewTestAvailabilityWithParams(t, userID1, room1.ID, now.Add(72*time.Hour), now.Add(74*time.Hour), "Consultation", intPtr(15000), 1, domain.AvailabilityStatusAvailable)
		wrongTimeAvailabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, wrongTimeAvailability)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, wrongTimeAvailabilityEncx, testPool)

		// Apply combined filters
		queryParams := map[string]string{
			"room_id":    room1.ID.String(),
			"partner_id": userID1.String(),
			"start_time": now.Add(24 * time.Hour).Format(time.RFC3339),
			"end_time":   now.Add(48 * time.Hour).Format(time.RFC3339),
			"limit":      "10",
		}

		req := ta.NewGetAvailableSlotsRequest(t, ctx, testServerURL, queryParams, accessToken1)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should only return 1 matching slot
		assert.Len(t, response, 1)
		assert.Equal(t, room1.ID, response[0].RoomID)
		assert.Equal(t, userID1, response[0].UserID)
	})

	// Authentication Tests
	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)

		// Make request without access token
		req := ta.NewGetAvailableSlotsRequest(t, ctx, testServerURL, nil, "")
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when access token is invalid", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)

		// Make request with invalid token
		req := ta.NewGetAvailableSlotsRequest(t, ctx, testServerURL, nil, "invalid_token")
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	// Validation Tests
	t.Run("should handle invalid time format gracefully", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create test data for proper setup
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)

		room := tr.NewTestRoomWithBuilding(t, building.ID)
		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		tr.InsertRoomEncx(t, ctx, testPool, roomEncx)

		accessToken, userID := ts.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@leviosa.care", room.ID, testPool, redisClient, crypto)

		// Create a test availability
		availability := ta.NewTestAvailabilityWithParams(t, userID, room.ID, time.Now().Add(24*time.Hour), time.Now().Add(26*time.Hour), "Consultation", intPtr(15000), 1, domain.AvailabilityStatusAvailable)
		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

		// Test with invalid time format - should be ignored and return all availabilities
		queryParams := map[string]string{
			"start_time": "invalid-time-format",
		}

		req := ta.NewGetAvailableSlotsRequest(t, ctx, testServerURL, queryParams, accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Handler gracefully ignores invalid parameters and returns 200
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return all availabilities since invalid time filter is ignored
		assert.Len(t, response, 1)
	})

	t.Run("should handle invalid room_id format gracefully", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create test data for proper setup
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)

		room := tr.NewTestRoomWithBuilding(t, building.ID)
		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		tr.InsertRoomEncx(t, ctx, testPool, roomEncx)

		accessToken, userID := ts.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@leviosa.care", room.ID, testPool, redisClient, crypto)

		// Create a test availability
		availability := ta.NewTestAvailabilityWithParams(t, userID, room.ID, time.Now().Add(24*time.Hour), time.Now().Add(26*time.Hour), "Consultation", intPtr(15000), 1, domain.AvailabilityStatusAvailable)
		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

		// Test with invalid room_id format - should be ignored and return all availabilities
		queryParams := map[string]string{
			"room_id": "invalid-uuid-format",
		}

		req := ta.NewGetAvailableSlotsRequest(t, ctx, testServerURL, queryParams, accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Handler gracefully ignores invalid parameters and returns 200
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return all availabilities since invalid room filter is ignored
		assert.Len(t, response, 1)
	})

	// Security Tests
	t.Run("should safely handle SQL injection attempts via query parameters", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)

		// Attempt SQL injection via room_id parameter
		maliciousParams := map[string]string{
			"room_id": "'; DROP TABLE availabilities; --",
		}

		req := ta.NewGetAvailableSlotsRequest(t, ctx, testServerURL, maliciousParams, accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should handle gracefully - either bad request for invalid UUID or ok with no results
		assert.True(t, resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusOK)

		// Verify table still exists and is functional
		ta.ClearAvailabilityTable(t, ctx, testPool) // This would fail if table was dropped
	})

	t.Run("should safely handle XSS attempts in query parameters", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)

		// Attempt XSS via query parameters
		maliciousParams := map[string]string{
			"start_time": "<script>alert('xss')</script>",
		}

		req := ta.NewGetAvailableSlotsRequest(t, ctx, testServerURL, maliciousParams, accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should handle gracefully - either bad request for invalid time or ok with no results
		assert.True(t, resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusOK)
	})
}
