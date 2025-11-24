package availability_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	availabilityHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	buildingHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	roomHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"
	catalogHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/catalog"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRoomGaps(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully retrieve gaps for room with bookings", func(t *testing.T) {
		// Clean state
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		roomHelpers.ClearRoomTables(t, ctx, testPool)
		buildingHelpers.ClearBuildingTable(t, ctx, testPool)
		catalogHelpers.ClearProductsTable(t, ctx, testPool)

		// Create test data
		partnerUser := authCtx.CreateTestUser(t, ctx, "partner", "partner@test.com")
		sessionToken := authCtx.CreateTestSession(t, ctx, partnerUser)

		building := buildingHelpers.CreateTestBuilding(t, ctx, testPool, crypto)

		// Create room with operating hours 9 AM to 5 PM
		operatingStart := time.Date(2000, 1, 1, 9, 0, 0, 0, time.UTC)
		operatingEnd := time.Date(2000, 1, 1, 17, 0, 0, 0, time.UTC)
		room := roomHelpers.CreateTestRoomWithOperatingHours(t, ctx, testPool, crypto, building.ID, operatingStart, operatingEnd)

		// Create products
		categoryID := catalogHelpers.CreateTestCategory(t, ctx, testPool)
		catalogHelpers.CreateDefaultTestProducts(t, ctx, testPool, categoryID)

		// Create availabilities with gaps
		testDate := time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour)

		// Booking 1: 10:00 - 11:15 (60min + 15min buffer)
		booking1Start := time.Date(testDate.Year(), testDate.Month(), testDate.Day(), 10, 0, 0, 0, time.UTC)
		booking1End := booking1Start.Add(75 * time.Minute)
		availability1 := availabilityHelpers.CreateTestAvailability(t, ctx, testPool, crypto, partnerUser.ID, room.ID, booking1Start, booking1End, "available")
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, availability1, testPool)

		// Booking 2: 13:00 - 14:45 (90min + 15min buffer)
		booking2Start := time.Date(testDate.Year(), testDate.Month(), testDate.Day(), 13, 0, 0, 0, time.UTC)
		booking2End := booking2Start.Add(105 * time.Minute)
		availability2 := availabilityHelpers.CreateTestAvailability(t, ctx, testPool, crypto, partnerUser.ID, room.ID, booking2Start, booking2End, "available")
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, availability2, testPool)

		// Expected gaps:
		// 1. 09:00 - 10:00 (60 min) - before first booking
		// 2. 11:15 - 13:00 (105 min) - between bookings
		// 3. 14:45 - 17:00 (135 min) - after last booking

		// Make request
		dateStr := testDate.Format("2006-01-02")
		req := availabilityHelpers.NewGetRoomGapsRequest(t, ctx, testServerURL, room.ID.String(), dateStr, sessionToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.GetRoomGapsResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify gaps
		assert.Equal(t, room.ID, response.RoomID)
		assert.Equal(t, testDate.Format("2006-01-02"), response.Date.Format("2006-01-02"))
		assert.Len(t, response.Gaps, 3, "Should have 3 gaps")

		// Verify gap 1: Before first booking (60 min)
		gap1 := response.Gaps[0]
		assert.Equal(t, 60, gap1.DurationMinutes)
		assert.True(t, gap1.IsBookable, "60-min gap should be bookable")
		assert.NotEmpty(t, gap1.SuggestedProducts, "Should have product suggestions")

		// Verify gap 2: Between bookings (105 min)
		gap2 := response.Gaps[1]
		assert.Equal(t, 105, gap2.DurationMinutes)
		assert.True(t, gap2.IsBookable)

		// Verify gap 3: After last booking (135 min)
		gap3 := response.Gaps[2]
		assert.Equal(t, 135, gap3.DurationMinutes)
		assert.True(t, gap3.IsBookable)

		// Verify total gap minutes
		expectedTotal := 60 + 105 + 135
		assert.Equal(t, expectedTotal, response.TotalGapMinutes)
	})

	t.Run("should return entire day as gap when no bookings exist", func(t *testing.T) {
		// Clean state
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		roomHelpers.ClearRoomTables(t, ctx, testPool)
		buildingHelpers.ClearBuildingTable(t, ctx, testPool)

		// Create test data
		partnerUser := authCtx.CreateTestUser(t, ctx, "partner", "partner2@test.com")
		sessionToken := authCtx.CreateTestSession(t, ctx, partnerUser)

		building := buildingHelpers.CreateTestBuilding(t, ctx, testPool, crypto)

		// Room with 8-hour operating day (9 AM - 5 PM = 480 minutes)
		operatingStart := time.Date(2000, 1, 1, 9, 0, 0, 0, time.UTC)
		operatingEnd := time.Date(2000, 1, 1, 17, 0, 0, 0, time.UTC)
		room := roomHelpers.CreateTestRoomWithOperatingHours(t, ctx, testPool, crypto, building.ID, operatingStart, operatingEnd)

		// No bookings - don't create any availabilities

		// Make request
		testDate := time.Now().Add(24 * time.Hour)
		dateStr := testDate.Format("2006-01-02")
		req := availabilityHelpers.NewGetRoomGapsRequest(t, ctx, testServerURL, room.ID.String(), dateStr, sessionToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.GetRoomGapsResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify single gap for entire day
		assert.Len(t, response.Gaps, 1, "Should have 1 gap (entire day)")
		gap := response.Gaps[0]
		assert.Equal(t, 480, gap.DurationMinutes, "Should be 8 hours (480 minutes)")
		assert.True(t, gap.IsBookable)
		assert.Equal(t, 480, response.TotalGapMinutes)
	})

	t.Run("should return empty gaps when room is fully booked", func(t *testing.T) {
		// Clean state
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		roomHelpers.ClearRoomTables(t, ctx, testPool)
		buildingHelpers.ClearBuildingTable(t, ctx, testPool)

		// Create test data
		partnerUser := authCtx.CreateTestUser(t, ctx, "partner", "partner3@test.com")
		sessionToken := authCtx.CreateTestSession(t, ctx, partnerUser)

		building := buildingHelpers.CreateTestBuilding(t, ctx, testPool, crypto)

		operatingStart := time.Date(2000, 1, 1, 9, 0, 0, 0, time.UTC)
		operatingEnd := time.Date(2000, 1, 1, 17, 0, 0, 0, time.UTC)
		room := roomHelpers.CreateTestRoomWithOperatingHours(t, ctx, testPool, crypto, building.ID, operatingStart, operatingEnd)

		// Create booking that covers entire operating hours
		testDate := time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour)
		bookingStart := time.Date(testDate.Year(), testDate.Month(), testDate.Day(), 9, 0, 0, 0, time.UTC)
		bookingEnd := time.Date(testDate.Year(), testDate.Month(), testDate.Day(), 17, 0, 0, 0, time.UTC)

		availability := availabilityHelpers.CreateTestAvailability(t, ctx, testPool, crypto, partnerUser.ID, room.ID, bookingStart, bookingEnd, "available")
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, availability, testPool)

		// Make request
		dateStr := testDate.Format("2006-01-02")
		req := availabilityHelpers.NewGetRoomGapsRequest(t, ctx, testServerURL, room.ID.String(), dateStr, sessionToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.GetRoomGapsResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify no gaps
		assert.Empty(t, response.Gaps, "Should have no gaps")
		assert.Equal(t, 0, response.TotalGapMinutes)
	})

	t.Run("should return 400 for invalid date format", func(t *testing.T) {
		partnerUser := authCtx.CreateTestUser(t, ctx, "partner", "partner4@test.com")
		sessionToken := authCtx.CreateTestSession(t, ctx, partnerUser)

		roomID := uuid.New()
		req := availabilityHelpers.NewGetRoomGapsRequest(t, ctx, testServerURL, roomID.String(), "invalid-date", sessionToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid room ID format", func(t *testing.T) {
		partnerUser := authCtx.CreateTestUser(t, ctx, "partner", "partner5@test.com")
		sessionToken := authCtx.CreateTestSession(t, ctx, partnerUser)

		testDate := time.Now().Add(24 * time.Hour)
		dateStr := testDate.Format("2006-01-02")
		req := availabilityHelpers.NewGetRoomGapsRequest(t, ctx, testServerURL, "invalid-uuid", dateStr, sessionToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
