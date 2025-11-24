package metrics_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	allocationHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/allocation"
	availabilityHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	buildingHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	metricsHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/metrics"
	roomHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPartnerMetrics(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should aggregate metrics across multiple rooms", func(t *testing.T) {
		// Clean state
		allocationHelpers.ClearAllocationTables(t, ctx, testPool)
		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
		roomHelpers.ClearRoomTables(t, ctx, testPool)
		buildingHelpers.ClearBuildingTable(t, ctx, testPool)

		// Create test data
		partnerUser := authCtx.CreateTestUser(t, ctx, "partner", "partner@partner-metrics.com")
		sessionToken := authCtx.CreateTestSession(t, ctx, partnerUser)

		building := buildingHelpers.CreateTestBuilding(t, ctx, testPool, crypto)

		// Create 2 rooms
		operatingStart := time.Date(2000, 1, 1, 9, 0, 0, 0, time.UTC)
		operatingEnd := time.Date(2000, 1, 1, 17, 0, 0, 0, time.UTC)
		room1 := roomHelpers.CreateTestRoomWithOperatingHours(t, ctx, testPool, crypto, building.ID, operatingStart, operatingEnd)
		room2 := roomHelpers.CreateTestRoomWithOperatingHours(t, ctx, testPool, crypto, building.ID, operatingStart, operatingEnd)

		// Create allocations for both rooms
		startDate := time.Now().Add(-30 * 24 * time.Hour)
		endDate := time.Now().Add(365 * 24 * time.Hour)
		allocation1 := allocationHelpers.CreateTestAllocation(t, ctx, testPool, crypto, partnerUser.ID, room1.ID, domain.AllocationTypeDedicated, startDate, endDate)
		allocationHelpers.InsertAllocationEncx(t, ctx, allocation1, testPool)
		allocation2 := allocationHelpers.CreateTestAllocation(t, ctx, testPool, crypto, partnerUser.ID, room2.ID, domain.AllocationTypeDedicated, startDate, endDate)
		allocationHelpers.InsertAllocationEncx(t, ctx, allocation2, testPool)

		// Create bookings for both rooms
		testDate := time.Now().Add(-1 * 24 * time.Hour).Truncate(24 * time.Hour)

		// Room 1: 2 hours booked
		booking1Start := time.Date(testDate.Year(), testDate.Month(), testDate.Day(), 10, 0, 0, 0, time.UTC)
		booking1End := booking1Start.Add(120 * time.Minute)
		availability1 := availabilityHelpers.CreateTestAvailability(t, ctx, testPool, crypto, partnerUser.ID, room1.ID, booking1Start, booking1End, "available")
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, availability1, testPool)

		// Room 2: 3 hours booked
		booking2Start := time.Date(testDate.Year(), testDate.Month(), testDate.Day(), 11, 0, 0, 0, time.UTC)
		booking2End := booking2Start.Add(180 * time.Minute)
		availability2 := availabilityHelpers.CreateTestAvailability(t, ctx, testPool, crypto, partnerUser.ID, room2.ID, booking2Start, booking2End, "available")
		availabilityHelpers.InsertAvailabilityEncx(t, ctx, availability2, testPool)

		// Refresh materialized view
		metricsHelpers.RefreshMetricsMaterializedView(t, ctx, testPool)

		// Make request
		dateStr := testDate.Format("2006-01-02")
		req := metricsHelpers.NewGetPartnerMetricsRequest(t, ctx, testServerURL, partnerUser.ID.String(), dateStr, dateStr, sessionToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.GetPartnerMetricsResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response structure
		assert.Equal(t, partnerUser.ID, response.PartnerID)
		assert.Len(t, response.RoomMetrics, 2, "Should have metrics for 2 rooms")

		// Verify each room has metrics
		for _, roomMetric := range response.RoomMetrics {
			assert.NotEmpty(t, roomMetric.DailyMetrics)
			assert.Greater(t, roomMetric.Summary.AverageUtilization, 0.0)
		}

		// Verify aggregated summary
		assert.Equal(t, 2, response.Summary.DaysAnalyzed, "Should analyze 2 room-days")
		assert.Greater(t, response.Summary.AverageUtilization, 0.0)
		assert.Greater(t, response.Summary.TotalIdleMinutes, 0)
	})

	t.Run("should handle partner with no rooms", func(t *testing.T) {
		// Clean state
		allocationHelpers.ClearAllocationTables(t, ctx, testPool)

		// Create partner with no room allocations
		partnerUser := authCtx.CreateTestUser(t, ctx, "partner", "partner2@partner-metrics.com")
		sessionToken := authCtx.CreateTestSession(t, ctx, partnerUser)

		// Make request
		startDate := time.Now().Add(-7 * 24 * time.Hour)
		endDate := time.Now()
		req := metricsHelpers.NewGetPartnerMetricsRequest(t, ctx, testServerURL, partnerUser.ID.String(), startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), sessionToken)
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
		partnerUser := authCtx.CreateTestUser(t, ctx, "partner", "partner3@partner-metrics.com")
		sessionToken := authCtx.CreateTestSession(t, ctx, partnerUser)

		// End date before start date
		req := metricsHelpers.NewGetPartnerMetricsRequest(t, ctx, testServerURL, partnerUser.ID.String(), "2025-01-10", "2025-01-05", sessionToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid partner ID format", func(t *testing.T) {
		partnerUser := authCtx.CreateTestUser(t, ctx, "partner", "partner4@partner-metrics.com")
		sessionToken := authCtx.CreateTestSession(t, ctx, partnerUser)

		req := metricsHelpers.NewGetPartnerMetricsRequest(t, ctx, testServerURL, "invalid-uuid", "2025-01-01", "2025-01-10", sessionToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
