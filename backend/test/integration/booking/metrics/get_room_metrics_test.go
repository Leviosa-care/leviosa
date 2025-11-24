package metrics_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	availabilityHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	buildingHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	metricsHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/metrics"
	roomHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRoomMetrics(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should retrieve metrics for room with availabilities", func(t *testing.T) {
		// Clean state
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		roomHelpers.ClearRoomTables(t, ctx, testPool)
		buildingHelpers.ClearBuildingTable(t, ctx, testPool)

		// Create test data
		partnerUser := authCtx.CreateTestUser(t, ctx, "partner", "partner@metrics.com")
		sessionToken := authCtx.CreateTestSession(t, ctx, partnerUser)

		building := buildingHelpers.CreateTestBuilding(t, ctx, testPool, crypto)

		// Room with 8-hour operating day (9 AM - 5 PM)
		operatingStart := time.Date(2000, 1, 1, 9, 0, 0, 0, time.UTC)
		operatingEnd := time.Date(2000, 1, 1, 17, 0, 0, 0, time.UTC)
		room := roomHelpers.CreateTestRoomWithOperatingHours(t, ctx, testPool, crypto, building.ID, operatingStart, operatingEnd)

		// Create availabilities for specific dates
		testDate1 := time.Now().Add(-2 * 24 * time.Hour).Truncate(24 * time.Hour)
		testDate2 := time.Now().Add(-1 * 24 * time.Hour).Truncate(24 * time.Hour)

		// Day 1: 3 hours booked (180 min)
		booking1Start := time.Date(testDate1.Year(), testDate1.Month(), testDate1.Day(), 10, 0, 0, 0, time.UTC)
		booking1End := booking1Start.Add(90 * time.Minute)
		availability1 := availabilityHelpers.CreateTestAvailability(t, ctx, testPool, crypto, partnerUser.ID, room.ID, booking1Start, booking1End, "available")
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, availability1, testPool)

		booking2Start := time.Date(testDate1.Year(), testDate1.Month(), testDate1.Day(), 14, 0, 0, 0, time.UTC)
		booking2End := booking2Start.Add(90 * time.Minute)
		availability2 := availabilityHelpers.CreateTestAvailability(t, ctx, testPool, crypto, partnerUser.ID, room.ID, booking2Start, booking2End, "available")
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, availability2, testPool)

		// Day 2: 2 hours booked (120 min)
		booking3Start := time.Date(testDate2.Year(), testDate2.Month(), testDate2.Day(), 11, 0, 0, 0, time.UTC)
		booking3End := booking3Start.Add(120 * time.Minute)
		availability3 := availabilityHelpers.CreateTestAvailability(t, ctx, testPool, crypto, partnerUser.ID, room.ID, booking3Start, booking3End, "available")
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, availability3, testPool)

		// Refresh materialized view to calculate metrics
		metricsHelpers.RefreshMetricsMaterializedView(t, ctx, testPool)

		// Make request
		startDateStr := testDate1.Format("2006-01-02")
		endDateStr := testDate2.Format("2006-01-02")
		req := metricsHelpers.NewGetRoomMetricsRequest(t, ctx, testServerURL, room.ID.String(), startDateStr, endDateStr, sessionToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.GetRoomMetricsResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response structure
		assert.Equal(t, room.ID, response.RoomID)
		assert.Len(t, response.DailyMetrics, 2, "Should have 2 days of metrics")

		// Verify daily metrics
		for _, daily := range response.DailyMetrics {
			assert.Equal(t, 480, daily.TotalMinutesOpen, "8-hour day = 480 minutes")
			assert.Greater(t, daily.TotalMinutesBooked, 0, "Should have booked minutes")
			assert.Greater(t, daily.UtilizationPercent, 0.0, "Should have utilization percentage")
		}

		// Verify summary
		assert.Equal(t, 2, response.Summary.DaysAnalyzed)
		assert.Greater(t, response.Summary.AverageUtilization, 0.0)
	})

	t.Run("should return empty metrics for room with no bookings", func(t *testing.T) {
		// Clean state
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		roomHelpers.ClearRoomTables(t, ctx, testPool)
		buildingHelpers.ClearBuildingTable(t, ctx, testPool)

		// Create test data
		partnerUser := authCtx.CreateTestUser(t, ctx, "partner", "partner2@metrics.com")
		sessionToken := authCtx.CreateTestSession(t, ctx, partnerUser)

		building := buildingHelpers.CreateTestBuilding(t, ctx, testPool, crypto)
		room := roomHelpers.CreateTestRoom(t, ctx, testPool, crypto, building.ID)

		// No bookings - refresh view
		metricsHelpers.RefreshMetricsMaterializedView(t, ctx, testPool)

		// Make request for date range
		startDate := time.Now().Add(-7 * 24 * time.Hour)
		endDate := time.Now()
		req := metricsHelpers.NewGetRoomMetricsRequest(t, ctx, testServerURL, room.ID.String(), startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), sessionToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.GetRoomMetricsResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return empty metrics
		assert.Empty(t, response.DailyMetrics)
		assert.Equal(t, 0, response.Summary.DaysAnalyzed)
	})

	t.Run("should return 400 for invalid date range", func(t *testing.T) {
		partnerUser := authCtx.CreateTestUser(t, ctx, "partner", "partner3@metrics.com")
		sessionToken := authCtx.CreateTestSession(t, ctx, partnerUser)

		roomID := uuid.New()
		// End date before start date
		req := metricsHelpers.NewGetRoomMetricsRequest(t, ctx, testServerURL, roomID.String(), "2025-01-10", "2025-01-05", sessionToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid date format", func(t *testing.T) {
		partnerUser := authCtx.CreateTestUser(t, ctx, "partner", "partner4@metrics.com")
		sessionToken := authCtx.CreateTestSession(t, ctx, partnerUser)

		roomID := uuid.New()
		req := metricsHelpers.NewGetRoomMetricsRequest(t, ctx, testServerURL, roomID.String(), "invalid-date", "2025-01-10", sessionToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid room ID format", func(t *testing.T) {
		partnerUser := authCtx.CreateTestUser(t, ctx, "partner", "partner5@metrics.com")
		sessionToken := authCtx.CreateTestSession(t, ctx, partnerUser)

		req := metricsHelpers.NewGetRoomMetricsRequest(t, ctx, testServerURL, "invalid-uuid", "2025-01-01", "2025-01-10", sessionToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
