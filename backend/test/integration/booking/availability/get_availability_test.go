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
	ta "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetAvailability TEST_PATH=test/integration/booking/availability/get_availability_test.go

func TestGetAvailability(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Helper to setup test availability with room
	setupTestAvailability := func(t *testing.T, ctx context.Context, userID uuid.UUID) (uuid.UUID, uuid.UUID) {
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

		// Create availability
		availability := ta.NewTestAvailabilityWithParams(
			t,
			userID,
			room.ID,
			time.Now().Add(24*time.Hour).Truncate(time.Second),
			time.Now().Add(26*time.Hour).Truncate(time.Second),
			"Consultation",
			intPtr(15000), // $150.00
			1,
			domain.AvailabilityStatusAvailable,
		)

		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

		return availability.ID, room.ID
	}

	t.Run("should successfully retrieve availability with standard user token", func(t *testing.T) { // Clean test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup user and availability
		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)
		availabilityID, _ := setupTestAvailability(t, ctx, uuid.New())

		// Make request
		req := ta.NewGetAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response data
		assert.Equal(t, availabilityID, response.ID)
		assert.NotEqual(t, uuid.Nil, response.UserID)
		assert.NotEqual(t, uuid.Nil, response.RoomID)
		assert.Equal(t, domain.AvailabilityStatusAvailable, response.Status)
		assert.Equal(t, "Consultation", response.ServiceType)
		assert.NotNil(t, response.PriceCents)
		assert.Equal(t, 15000, *response.PriceCents)
		assert.Equal(t, 1, response.MaxCapacity)
	})

	t.Run("should successfully retrieve availability with partner token", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup partner and availability
		accessToken := tu.SetupUserWithRole(t, ctx, identity.Partner, authCtx)
		availabilityID, _ := setupTestAvailability(t, ctx, uuid.New())

		req := ta.NewGetAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, availabilityID, response.ID)
	})

	t.Run("should successfully retrieve availability with admin token", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin and availability
		accessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		availabilityID, _ := setupTestAvailability(t, ctx, uuid.New())

		req := ta.NewGetAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, availabilityID, response.ID)
	})

	t.Run("should retrieve availability with encrypted fields properly decrypted", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)

		// Create building and room
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		room := tr.NewTestRoomWithBuilding(t, building.ID)
		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		// Create availability with specific encrypted values
		availability := ta.NewTestAvailabilityWithParams(
			t,
			uuid.New(),
			room.ID,
			time.Now().Add(24*time.Hour).Truncate(time.Second),
			time.Now().Add(26*time.Hour).Truncate(time.Second),
			"Specialized Therapy Session",
			intPtr(25000), // $250.00
			2,
			domain.AvailabilityStatusAvailable,
		)
		availability.Notes = "Bring comfortable clothing and arrive 10 minutes early"

		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

		// Make request
		req := ta.NewGetAvailabilityRequest(t, ctx, testServerURL, availability.ID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify encrypted fields are properly decrypted
		assert.Equal(t, "Specialized Therapy Session", response.ServiceType)
		assert.Equal(t, "Bring comfortable clothing and arrive 10 minutes early", response.Notes)
		assert.Equal(t, 25000, *response.PriceCents)
	})

	t.Run("should retrieve recurring availability with pattern", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)

		// Create building and room
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		room := tr.NewTestRoomWithBuilding(t, building.ID)
		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		// Create recurring availability
		startTime := time.Now().Add(24 * time.Hour).Truncate(time.Second)
		endTime := startTime.Add(1 * time.Hour)
		until := time.Now().Add(60 * 24 * time.Hour) // 60 days from now

		availability := &domain.Availability{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			RoomID:      room.ID,
			StartTime:   startTime,
			EndTime:     endTime,
			ServiceType: "Weekly Therapy",
			PriceCents:  intPtr(10000),
			MaxCapacity: 1,
			Notes:       "Recurring weekly session",
			IsRecurring: true,
			RecurrencePattern: &domain.RecurrencePattern{
				Type:       "weekly",
				Interval:   1,
				Until:      &until,
				DaysOfWeek: []int{1, 3, 5}, // Monday, Wednesday, Friday
			},
			Status:    domain.AvailabilityStatusAvailable,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

		// Make request
		req := ta.NewGetAvailabilityRequest(t, ctx, testServerURL, availability.ID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify recurrence pattern
		assert.NotNil(t, response.RecurrencePattern)
		assert.Equal(t, "weekly", response.RecurrencePattern.Type)
		assert.Equal(t, 1, response.RecurrencePattern.Interval)
		assert.Equal(t, []int{1, 3, 5}, response.RecurrencePattern.DaysOfWeek)
		assert.NotNil(t, response.RecurrencePattern.Until)
	})

	t.Run("should return 400 Bad Request for invalid availability ID format", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)

		// Use invalid UUID format
		req := ta.NewGetAvailabilityRequest(t, ctx, testServerURL, "invalid-uuid-format", accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "invalid availability ID format")
	})

	t.Run("should return 400 Bad Request for missing availability ID in path", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)

		// Make request without ID (this will result in path parsing error)
		req := ta.NewGetAvailabilityRequest(t, ctx, testServerURL, "", accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get 404 or 400 depending on router behavior
		assert.True(t, resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound)
	})

	t.Run("should return 404 Not Found for non-existent availability", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)

		// Use a valid UUID that doesn't exist in database
		nonExistentID := uuid.New()
		req := ta.NewGetAvailabilityRequest(t, ctx, testServerURL, nonExistentID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create availability without auth
		availabilityID, _ := setupTestAvailability(t, ctx, uuid.New())

		// Make request without token
		req := ta.NewGetAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create expired user session
		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Standard, authCtx)
		availabilityID, _ := setupTestAvailability(t, ctx, uuid.New())

		req := ta.NewGetAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		availabilityID, _ := setupTestAvailability(t, ctx, uuid.New())

		// Use invalid token
		req := ta.NewGetAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), "invalid-token-12345")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should retrieve availability with nil price", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)

		// Create building and room
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		room := tr.NewTestRoomWithBuilding(t, building.ID)
		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		// Create availability with nil price (free)
		availability := ta.NewTestAvailabilityWithParams(
			t,
			uuid.New(),
			room.ID,
			time.Now().Add(24*time.Hour).Truncate(time.Second),
			time.Now().Add(25*time.Hour).Truncate(time.Second),
			"Free Consultation",
			nil, // No price
			1,
			domain.AvailabilityStatusAvailable,
		)

		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

		// Make request
		req := ta.NewGetAvailabilityRequest(t, ctx, testServerURL, availability.ID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify price is nil
		assert.Nil(t, response.PriceCents)
		assert.Equal(t, "Free Consultation", response.ServiceType)
	})

	t.Run("should retrieve booked availability", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)

		// Create building and room
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		room := tr.NewTestRoomWithBuilding(t, building.ID)
		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		// Create booked availability
		availability := ta.NewTestAvailabilityWithParams(
			t,
			uuid.New(),
			room.ID,
			time.Now().Add(24*time.Hour).Truncate(time.Second),
			time.Now().Add(26*time.Hour).Truncate(time.Second),
			"Consultation",
			intPtr(15000),
			1,
			domain.AvailabilityStatusBooked, // Booked status
		)

		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

		// Make request
		req := ta.NewGetAvailabilityRequest(t, ctx, testServerURL, availability.ID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify booked status
		assert.Equal(t, domain.AvailabilityStatusBooked, response.Status)
	})

	t.Run("should retrieve cancelled availability", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)

		// Create building and room
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		room := tr.NewTestRoomWithBuilding(t, building.ID)
		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		// Create cancelled availability
		availability := ta.NewTestAvailabilityWithParams(
			t,
			uuid.New(),
			room.ID,
			time.Now().Add(24*time.Hour).Truncate(time.Second),
			time.Now().Add(26*time.Hour).Truncate(time.Second),
			"Consultation",
			intPtr(15000),
			1,
			domain.AvailabilityStatusCancelled, // Cancelled status
		)

		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

		// Make request
		req := ta.NewGetAvailabilityRequest(t, ctx, testServerURL, availability.ID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify cancelled status
		assert.Equal(t, domain.AvailabilityStatusCancelled, response.Status)
	})

	t.Run("should handle context timeout appropriately", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)
		availabilityID, _ := setupTestAvailability(t, ctx, uuid.New())

		// Use a very short context timeout
		shortCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
		defer cancel()

		time.Sleep(2 * time.Millisecond) // Ensure timeout has passed

		req := ta.NewGetAvailabilityRequest(t, shortCtx, testServerURL, availabilityID.String(), accessToken)

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

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}
