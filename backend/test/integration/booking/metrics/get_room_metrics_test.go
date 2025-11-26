package metrics_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	ta "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tm "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/metrics"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetRoomMetrics TEST_PATH=test/integration/booking/metrics/get_room_metrics_test.go

func TestGetRoomMetrics(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should retrieve metrics for room with availabilities", func(t *testing.T) {
		// Clean state
		ta.ClearRoomSchedulesTable(t, ctx, testPool)
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authenticated partner
		accessToken := tu.SetupPartnerUser(t, ctx, authCtx)

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

		// Set up operating hours: 9 AM - 5 PM (8 hours) for all weekdays
		for day := 0; day <= 6; day++ {
			schedule := ta.NewTestRoomScheduleRecurring(room.ID, day, "09:00", "17:00")
			ta.InsertRoomSchedule(t, ctx, schedule, testPool)
		}

		// Create availabilities for specific dates (future dates to pass CHECK constraint)
		testDate1 := time.Now().Add(2 * 24 * time.Hour).Truncate(24 * time.Hour)
		testDate2 := time.Now().Add(3 * 24 * time.Hour).Truncate(24 * time.Hour)

		// Day 1: 3 hours booked (180 min) - use 'booked' status so they count as bookings
		booking1Start := time.Date(testDate1.Year(), testDate1.Month(), testDate1.Day(), 10, 0, 0, 0, time.UTC)
		booking1End := booking1Start.Add(90 * time.Minute)
		availability1 := ta.NewTestAvailabilityWithParams(t, uuid.New(), room.ID, booking1Start, booking1End, "Massage", nil, 1, domain.AvailabilityStatusBooked)
		availability1Encx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability1)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availability1Encx, testPool)

		booking2Start := time.Date(testDate1.Year(), testDate1.Month(), testDate1.Day(), 14, 0, 0, 0, time.UTC)
		booking2End := booking2Start.Add(90 * time.Minute)
		availability2 := ta.NewTestAvailabilityWithParams(t, uuid.New(), room.ID, booking2Start, booking2End, "Massage", nil, 1, domain.AvailabilityStatusBooked)
		availability2Encx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability2)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availability2Encx, testPool)

		// Day 2: 2 hours booked (120 min)
		booking3Start := time.Date(testDate2.Year(), testDate2.Month(), testDate2.Day(), 11, 0, 0, 0, time.UTC)
		booking3End := booking3Start.Add(120 * time.Minute)
		availability3 := ta.NewTestAvailabilityWithParams(t, uuid.New(), room.ID, booking3Start, booking3End, "Massage", nil, 1, domain.AvailabilityStatusBooked)
		availability3Encx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability3)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availability3Encx, testPool)

		// Refresh materialized view to calculate metrics
		tm.RefreshMetricsMaterializedView(t, ctx, testPool)

		// Make request
		startDateStr := testDate1.Format("2006-01-02")
		endDateStr := testDate2.Format("2006-01-02")
		req := tm.NewGetRoomMetricsRequest(t, ctx, testServerURL, room.ID.String(), startDateStr, endDateStr, accessToken)
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

		// Verify daily metrics exist
		for _, daily := range response.DailyMetrics {
			assert.Greater(t, daily.TotalMinutesBooked, 0, "Should have booked minutes")
			assert.Greater(t, daily.UtilizationPercent, 0.0, "Should have utilization percentage")
		}

		// Verify summary
		assert.Equal(t, 2, response.Summary.DaysAnalyzed)
		assert.Greater(t, response.Summary.AverageUtilization, 0.0)
	})

	t.Run("should return empty metrics for room with no bookings", func(t *testing.T) {
		// Clean state
		ta.ClearRoomSchedulesTable(t, ctx, testPool)
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authenticated partner
		accessToken := tu.SetupPartnerUser(t, ctx, authCtx)

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

		// Set up operating hours: 9 AM - 5 PM for all days
		for day := 0; day <= 6; day++ {
			schedule := ta.NewTestRoomScheduleRecurring(room.ID, day, "09:00", "17:00")
			ta.InsertRoomSchedule(t, ctx, schedule, testPool)
		}

		// No bookings - refresh view
		tm.RefreshMetricsMaterializedView(t, ctx, testPool)

		// Make request for date range
		startDate := time.Now().Add(-7 * 24 * time.Hour)
		endDate := time.Now()
		req := tm.NewGetRoomMetricsRequest(t, ctx, testServerURL, room.ID.String(), startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), accessToken)
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
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupPartnerUser(t, ctx, authCtx)
		roomID := uuid.New()
		// End date before start date
		req := tm.NewGetRoomMetricsRequest(t, ctx, testServerURL, roomID.String(), "2025-01-10", "2025-01-05", accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid date format", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupPartnerUser(t, ctx, authCtx)
		roomID := uuid.New()
		req := tm.NewGetRoomMetricsRequest(t, ctx, testServerURL, roomID.String(), "invalid-date", "2025-01-10", accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid room ID format", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupPartnerUser(t, ctx, authCtx)
		req := tm.NewGetRoomMetricsRequest(t, ctx, testServerURL, "invalid-uuid", "2025-01-01", "2025-01-10", accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
