package availability_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	ts "github.com/Leviosa-care/leviosa/backend/test/helpers/booking"
	ta "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCheckAvailabilityConflict TEST_PATH=test/integration/booking/availability/check_availability_conflict_test.go

func TestCheckAvailabilityConflict(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Helper to setup test availability with room and authenticated user
	setupTestAvailabilityForConflictCheck := func(t *testing.T, ctx context.Context, email string) (string, uuid.UUID, *domain.Availability, uuid.UUID) {
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

		// Setup authenticated partner with allocation to the room
		accessToken, userID := ts.SetupAuthenticatedPartnerWithAllocation(t, ctx, email, room.ID, testPool, redisClient, crypto)

		// Create existing availability for conflict testing
		startTime := time.Now().Add(24 * time.Hour).Truncate(time.Second)
		endTime := time.Now().Add(26 * time.Hour).Truncate(time.Second)
		availability := ta.NewTestAvailabilityWithParams(
			t,
			userID,
			room.ID,
			startTime,
			endTime,
			"Consultation",
			intPtr(15000), // $150.00
			1,
			domain.AvailabilityStatusAvailable,
		)

		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

		return accessToken, userID, availability, room.ID
	}

	// Helper to create check conflict request
	createCheckConflictRequest := func(t *testing.T, ctx context.Context, baseURL, partnerID, startTime, endTime, excludeID string, accessToken string) *http.Request {
		t.Helper()

		requestURL, err := url.Parse(baseURL + "/partners/" + partnerID + "/availabilities/conflict")
		require.NoError(t, err)

		// Add query parameters
		q := requestURL.Query()
		q.Set("start_time", startTime)
		q.Set("end_time", endTime)
		if excludeID != "" {
			q.Set("exclude_id", excludeID)
		}
		requestURL.RawQuery = q.Encode()

		req, err := http.NewRequestWithContext(ctx, "GET", requestURL.String(), nil)
		require.NoError(t, err)

		req.AddCookie(&http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
			Path:  "/",
		})

		return req
	}

	t.Run("should detect no conflict when time slot is completely different", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, userID, _, _ := setupTestAvailabilityForConflictCheck(t, ctx, "partner@leviosa.care")

		// Check for conflict with non-overlapping time slot (2 days later)
		newStartTime := time.Now().Add(72 * time.Hour).Truncate(time.Second) // 3 days from now
		newEndTime := time.Now().Add(74 * time.Hour).Truncate(time.Second)   // 3 days + 2 hours from now

		req := createCheckConflictRequest(
			t,
			ctx,
			testServerURL,
			userID.String(),
			newStartTime.Format(time.RFC3339),
			newEndTime.Format(time.RFC3339),
			"",
			accessToken,
		)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			HasConflict bool `json:"has_conflict"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.False(t, response.HasConflict, "Expected no conflict for non-overlapping time slot")
	})

	t.Run("should detect conflict when time slot overlaps partially", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, userID, existingAvailability, _ := setupTestAvailabilityForConflictCheck(t, ctx, "partner@leviosa.care")

		// Check for conflict with overlapping time slot (starts 1 hour into existing availability)
		newStartTime := existingAvailability.StartTime.Add(1 * time.Hour)
		newEndTime := existingAvailability.EndTime.Add(1 * time.Hour)

		req := createCheckConflictRequest(
			t,
			ctx,
			testServerURL,
			userID.String(),
			newStartTime.Format(time.RFC3339),
			newEndTime.Format(time.RFC3339),
			"",
			accessToken,
		)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			HasConflict bool `json:"has_conflict"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.HasConflict, "Expected conflict for overlapping time slot")
	})

	t.Run("should detect conflict when time slot is completely within existing availability", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, userID, existingAvailability, _ := setupTestAvailabilityForConflictCheck(t, ctx, "partner@leviosa.care")

		// Check for conflict with time slot completely inside existing availability
		newStartTime := existingAvailability.StartTime.Add(30 * time.Minute)
		newEndTime := existingAvailability.StartTime.Add(90 * time.Minute)

		req := createCheckConflictRequest(
			t,
			ctx,
			testServerURL,
			userID.String(),
			newStartTime.Format(time.RFC3339),
			newEndTime.Format(time.RFC3339),
			"",
			accessToken,
		)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			HasConflict bool `json:"has_conflict"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.HasConflict, "Expected conflict for time slot within existing availability")
	})

	t.Run("should detect conflict when time slot completely encompasses existing availability", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, userID, existingAvailability, _ := setupTestAvailabilityForConflictCheck(t, ctx, "partner@leviosa.care")

		// Check for conflict with time slot that completely encompasses existing availability
		newStartTime := existingAvailability.StartTime.Add(-30 * time.Minute)
		newEndTime := existingAvailability.EndTime.Add(30 * time.Minute)

		req := createCheckConflictRequest(
			t,
			ctx,
			testServerURL,
			userID.String(),
			newStartTime.Format(time.RFC3339),
			newEndTime.Format(time.RFC3339),
			"",
			accessToken,
		)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			HasConflict bool `json:"has_conflict"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.HasConflict, "Expected conflict for time slot encompassing existing availability")
	})

	t.Run("should exclude specific availability from conflict check", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, userID, existingAvailability, _ := setupTestAvailabilityForConflictCheck(t, ctx, "partner@leviosa.care")

		// Check for conflict with same time slot but excluding the existing availability
		req := createCheckConflictRequest(
			t,
			ctx,
			testServerURL,
			userID.String(),
			existingAvailability.StartTime.Format(time.RFC3339),
			existingAvailability.EndTime.Format(time.RFC3339),
			existingAvailability.ID.String(),
			accessToken,
		)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			HasConflict bool `json:"has_conflict"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.False(t, response.HasConflict, "Expected no conflict when excluding the existing availability")
	})

	t.Run("should return 401 for unauthenticated request", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		_, userID, _, _ := setupTestAvailabilityForConflictCheck(t, ctx, "partner@leviosa.care")

		newStartTime := time.Now().Add(48 * time.Hour).Truncate(time.Second)
		newEndTime := time.Now().Add(50 * time.Hour).Truncate(time.Second)

		// Create request without authentication
		requestURL, err := url.Parse(testServerURL + "/partners/" + userID.String() + "/availabilities/conflict")
		require.NoError(t, err)

		q := requestURL.Query()
		q.Set("start_time", newStartTime.Format(time.RFC3339))
		q.Set("end_time", newEndTime.Format(time.RFC3339))
		requestURL.RawQuery = q.Encode()

		req, err := http.NewRequestWithContext(ctx, "GET", requestURL.String(), nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 for standard user trying to check another partner's conflicts", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data with partner
		_, partnerUserID, _, roomID := setupTestAvailabilityForConflictCheck(t, ctx, "partner@leviosa.care")

		// Create standard user
		standardAccessToken, _ := ts.SetupStandardUser(t, ctx, "standarduser@leviosa.care", roomID, testPool, redisClient, crypto)

		newStartTime := time.Now().Add(48 * time.Hour).Truncate(time.Second)
		newEndTime := time.Now().Add(50 * time.Hour).Truncate(time.Second)

		// Standard user trying to check partner's conflicts
		req := createCheckConflictRequest(
			t,
			ctx,
			testServerURL,
			partnerUserID.String(),
			newStartTime.Format(time.RFC3339),
			newEndTime.Format(time.RFC3339),
			"",
			standardAccessToken,
		)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 400 for invalid partner ID format", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authenticated partner
		accessToken, _, _, _ := setupTestAvailabilityForConflictCheck(t, ctx, "partner@leviosa.care")

		newStartTime := time.Now().Add(48 * time.Hour).Truncate(time.Second)
		newEndTime := time.Now().Add(50 * time.Hour).Truncate(time.Second)

		// Create request with invalid partner ID
		requestURL, err := url.Parse(testServerURL + "/partners/invalid-uuid/availabilities/conflict")
		require.NoError(t, err)

		q := requestURL.Query()
		q.Set("start_time", newStartTime.Format(time.RFC3339))
		q.Set("end_time", newEndTime.Format(time.RFC3339))
		requestURL.RawQuery = q.Encode()

		req, err := http.NewRequestWithContext(ctx, "GET", requestURL.String(), nil)
		require.NoError(t, err)

		req.AddCookie(&http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
			Path:  "/",
		})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for missing start_time parameter", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authenticated partner
		accessToken, userID, _, _ := setupTestAvailabilityForConflictCheck(t, ctx, "partner@leviosa.care")

		newEndTime := time.Now().Add(50 * time.Hour).Truncate(time.Second)

		// Create request without start_time
		requestURL, err := url.Parse(testServerURL + "/partners/" + userID.String() + "/availabilities/conflict")
		require.NoError(t, err)

		q := requestURL.Query()
		q.Set("end_time", newEndTime.Format(time.RFC3339))
		requestURL.RawQuery = q.Encode()

		req, err := http.NewRequestWithContext(ctx, "GET", requestURL.String(), nil)
		require.NoError(t, err)

		req.AddCookie(&http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
			Path:  "/",
		})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for missing end_time parameter", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authenticated partner
		accessToken, userID, _, _ := setupTestAvailabilityForConflictCheck(t, ctx, "partner@leviosa.care")

		newStartTime := time.Now().Add(48 * time.Hour).Truncate(time.Second)

		// Create request without end_time
		requestURL, err := url.Parse(testServerURL + "/partners/" + userID.String() + "/availabilities/conflict")
		require.NoError(t, err)

		q := requestURL.Query()
		q.Set("start_time", newStartTime.Format(time.RFC3339))
		requestURL.RawQuery = q.Encode()

		req, err := http.NewRequestWithContext(ctx, "GET", requestURL.String(), nil)
		require.NoError(t, err)

		req.AddCookie(&http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
			Path:  "/",
		})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid start_time format", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authenticated partner
		accessToken, userID, _, _ := setupTestAvailabilityForConflictCheck(t, ctx, "partner@leviosa.care")

		newEndTime := time.Now().Add(50 * time.Hour).Truncate(time.Second)

		// Create request with invalid start_time format
		requestURL, err := url.Parse(testServerURL + "/partners/" + userID.String() + "/availabilities/conflict")
		require.NoError(t, err)

		q := requestURL.Query()
		q.Set("start_time", "invalid-time")
		q.Set("end_time", newEndTime.Format(time.RFC3339))
		requestURL.RawQuery = q.Encode()

		req, err := http.NewRequestWithContext(ctx, "GET", requestURL.String(), nil)
		require.NoError(t, err)

		req.AddCookie(&http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
			Path:  "/",
		})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid end_time format", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup authenticated partner
		accessToken, userID, _, _ := setupTestAvailabilityForConflictCheck(t, ctx, "partner@leviosa.care")

		newStartTime := time.Now().Add(48 * time.Hour).Truncate(time.Second)

		// Create request with invalid end_time format
		requestURL, err := url.Parse(testServerURL + "/partners/" + userID.String() + "/availabilities/conflict")
		require.NoError(t, err)

		q := requestURL.Query()
		q.Set("start_time", newStartTime.Format(time.RFC3339))
		q.Set("end_time", "invalid-time")
		requestURL.RawQuery = q.Encode()

		req, err := http.NewRequestWithContext(ctx, "GET", requestURL.String(), nil)
		require.NoError(t, err)

		req.AddCookie(&http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
			Path:  "/",
		})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should allow admin to check any partner's conflicts", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data with partner
		_, partnerUserID, existingAvailability, roomID := setupTestAvailabilityForConflictCheck(t, ctx, "partner@leviosa.care")

		// Create admin user
		adminAccessToken, _ := ts.SetupAdminWithAllocation(t, ctx, roomID, testPool, redisClient, crypto)

		// Admin checking partner's conflicts with overlapping time slot
		newStartTime := existingAvailability.StartTime.Add(1 * time.Hour)
		newEndTime := existingAvailability.EndTime.Add(1 * time.Hour)

		req := createCheckConflictRequest(
			t,
			ctx,
			testServerURL,
			partnerUserID.String(),
			newStartTime.Format(time.RFC3339),
			newEndTime.Format(time.RFC3339),
			"",
			adminAccessToken,
		)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			HasConflict bool `json:"has_conflict"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.HasConflict, "Expected conflict for overlapping time slot")
	})

	t.Run("should allow partner to check their own conflicts", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, userID, _, _ := setupTestAvailabilityForConflictCheck(t, ctx, "partner@leviosa.care")

		// Partner checking their own conflicts with non-overlapping time slot
		newStartTime := time.Now().Add(72 * time.Hour).Truncate(time.Second)
		newEndTime := time.Now().Add(74 * time.Hour).Truncate(time.Second)

		req := createCheckConflictRequest(
			t,
			ctx,
			testServerURL,
			userID.String(),
			newStartTime.Format(time.RFC3339),
			newEndTime.Format(time.RFC3339),
			"",
			accessToken,
		)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			HasConflict bool `json:"has_conflict"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.False(t, response.HasConflict, "Expected no conflict for non-overlapping time slot")
	})

	t.Run("should handle invalid exclude_id gracefully", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, userID, existingAvailability, _ := setupTestAvailabilityForConflictCheck(t, ctx, "partner@leviosa.care")

		// Check for conflict with same time slot but with invalid exclude_id (should be ignored)
		newStartTime := existingAvailability.StartTime
		newEndTime := existingAvailability.EndTime

		req := createCheckConflictRequest(
			t,
			ctx,
			testServerURL,
			userID.String(),
			newStartTime.Format(time.RFC3339),
			newEndTime.Format(time.RFC3339),
			"invalid-uuid",
			accessToken,
		)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			HasConflict bool `json:"has_conflict"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.HasConflict, "Expected conflict when invalid exclude_id is ignored")
	})

	t.Run("should detect edge conflict when new time slot starts exactly when existing ends", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, userID, existingAvailability, _ := setupTestAvailabilityForConflictCheck(t, ctx, "partner@leviosa.care")

		// Check for conflict with time slot that starts exactly when existing ends
		newStartTime := existingAvailability.EndTime
		newEndTime := existingAvailability.EndTime.Add(2 * time.Hour)

		req := createCheckConflictRequest(
			t,
			ctx,
			testServerURL,
			userID.String(),
			newStartTime.Format(time.RFC3339),
			newEndTime.Format(time.RFC3339),
			"",
			accessToken,
		)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			HasConflict bool `json:"has_conflict"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// This depends on the implementation - typically, end == start is NOT a conflict
		assert.False(t, response.HasConflict, "Expected no conflict when new slot starts exactly when existing ends")
	})

	t.Run("should detect edge conflict when new time slot ends exactly when existing starts", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, userID, existingAvailability, _ := setupTestAvailabilityForConflictCheck(t, ctx, "partner@leviosa.care")

		// Check for conflict with time slot that ends exactly when existing starts
		newStartTime := existingAvailability.StartTime.Add(-2 * time.Hour)
		newEndTime := existingAvailability.StartTime

		req := createCheckConflictRequest(
			t,
			ctx,
			testServerURL,
			userID.String(),
			newStartTime.Format(time.RFC3339),
			newEndTime.Format(time.RFC3339),
			"",
			accessToken,
		)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			HasConflict bool `json:"has_conflict"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// This depends on the implementation - typically, end == start is NOT a conflict
		assert.False(t, response.HasConflict, "Expected no conflict when new slot ends exactly when existing starts")
	})
}

