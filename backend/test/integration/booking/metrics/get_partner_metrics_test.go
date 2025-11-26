package metrics_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	talloc "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/allocation"
	ta "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tm "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/metrics"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetPartnerMetrics TEST_PATH=test/integration/booking/metrics/get_partner_metrics_test.go

func TestGetPartnerMetrics(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should aggregate metrics across multiple rooms", func(t *testing.T) {
		// Clean state
		ta.ClearRoomSchedulesTable(t, ctx, testPool)
		talloc.ClearAllocationTable(t, ctx, testPool)
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authenticated partner - this creates the user and session
		accessToken := tu.SetupPartnerUser(t, ctx, authCtx)

		// Get the user ID from the session
		// We need to extract it from authCtx or create a known user
		// For simplicity, let's create a specific partner user
		partnerUserID := uuid.New()

		// Create building
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Create 2 rooms
		room1 := tr.NewTestRoomWithBuilding(t, building.ID)
		room1Encx, err := domain.ProcessRoomEncx(ctx, crypto, room1)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)

		room2 := tr.NewTestRoomWithBuilding(t, building.ID)
		room2Encx, err := domain.ProcessRoomEncx(ctx, crypto, room2)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)

		// Set up operating hours for both rooms: 9 AM - 5 PM
		for day := 0; day <= 6; day++ {
			schedule1 := ta.NewTestRoomScheduleRecurring(room1.ID, day, "09:00", "17:00")
			ta.InsertRoomSchedule(t, ctx, schedule1, testPool)
			schedule2 := ta.NewTestRoomScheduleRecurring(room2.ID, day, "09:00", "17:00")
			ta.InsertRoomSchedule(t, ctx, schedule2, testPool)
		}

		// Create allocations for both rooms
		allocation1 := talloc.NewTestSharedAllocation(t, room1.ID, partnerUserID)
		talloc.InsertAllocation(t, ctx, allocation1, testPool, crypto)

		allocation2 := talloc.NewTestSharedAllocation(t, room2.ID, partnerUserID)
		talloc.InsertAllocation(t, ctx, allocation2, testPool, crypto)

		// Create bookings for both rooms (future dates to pass CHECK constraint)
		testDate := time.Now().Add(2 * 24 * time.Hour).Truncate(24 * time.Hour)

		// Room 1: 2 hours booked - use 'booked' status
		booking1Start := time.Date(testDate.Year(), testDate.Month(), testDate.Day(), 10, 0, 0, 0, time.UTC)
		booking1End := booking1Start.Add(120 * time.Minute)
		availability1 := ta.NewTestAvailabilityWithParams(t, partnerUserID, room1.ID, booking1Start, booking1End, "Massage", nil, 1, domain.AvailabilityStatusBooked)
		availability1Encx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability1)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availability1Encx, testPool)

		// Room 2: 3 hours booked
		booking2Start := time.Date(testDate.Year(), testDate.Month(), testDate.Day(), 11, 0, 0, 0, time.UTC)
		booking2End := booking2Start.Add(180 * time.Minute)
		availability2 := ta.NewTestAvailabilityWithParams(t, partnerUserID, room2.ID, booking2Start, booking2End, "Massage", nil, 1, domain.AvailabilityStatusBooked)
		availability2Encx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability2)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availability2Encx, testPool)

		// Refresh materialized view
		tm.RefreshMetricsMaterializedView(t, ctx, testPool)

		// Make request
		dateStr := testDate.Format("2006-01-02")
		req := tm.NewGetPartnerMetricsRequest(t, ctx, testServerURL, partnerUserID.String(), dateStr, dateStr, accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.GetPartnerMetricsResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response structure
		assert.Equal(t, partnerUserID, response.PartnerID)
		assert.Len(t, response.RoomMetrics, 2, "Should have metrics for 2 rooms")

		// Verify each room has metrics
		for _, roomMetric := range response.RoomMetrics {
			assert.NotEmpty(t, roomMetric.DailyMetrics)
			assert.Greater(t, roomMetric.Summary.AverageUtilization, 0.0)
		}

		// Verify aggregated summary
		assert.Greater(t, response.Summary.AverageUtilization, 0.0)
	})

	t.Run("should handle partner with no rooms", func(t *testing.T) {
		// Clean state
		talloc.ClearAllocationTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup partner with no room allocations
		accessToken := tu.SetupPartnerUser(t, ctx, authCtx)
		partnerUserID := uuid.New()

		// Make request
		startDate := time.Now().Add(-7 * 24 * time.Hour)
		endDate := time.Now()
		req := tm.NewGetPartnerMetricsRequest(t, ctx, testServerURL, partnerUserID.String(), startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.GetPartnerMetricsResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return empty metrics
		assert.Empty(t, response.RoomMetrics)
		assert.Equal(t, 0, response.Summary.DaysAnalyzed)
	})

	t.Run("should return 400 for invalid date range", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupPartnerUser(t, ctx, authCtx)
		partnerUserID := uuid.New()

		// End date before start date
		req := tm.NewGetPartnerMetricsRequest(t, ctx, testServerURL, partnerUserID.String(), "2025-01-10", "2025-01-05", accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid partner ID format", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupPartnerUser(t, ctx, authCtx)

		req := tm.NewGetPartnerMetricsRequest(t, ctx, testServerURL, "invalid-uuid", "2025-01-01", "2025-01-10", accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
