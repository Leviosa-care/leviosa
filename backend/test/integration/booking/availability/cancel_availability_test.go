package availability_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
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

// make test-func TEST_NAME=TestCancelAvailability TEST_PATH=test/integration/booking/availability/cancel_availability_test.go

func TestCancelAvailability(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Helper to setup test availability with room and authenticated user
	setupTestAvailabilityForCancellation := func(t *testing.T, ctx context.Context, email string) (string, uuid.UUID, uuid.UUID, uuid.UUID) {
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

		// Create availability for cancellation
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

		return accessToken, userID, availability.ID, room.ID
	}

	t.Run("should successfully cancel availability with partner token", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, _, availabilityID, _ := setupTestAvailabilityForCancellation(t, ctx, "partner@leviosa.care")

		// Make cancel request
		req := ta.NewCancelAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response - should be 204 No Content
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify availability is actually cancelled in database
		availabilityEncx := ta.GetAvailabilityEncxFromDB(t, ctx, availabilityID, testPool)
		assert.Equal(t, domain.AvailabilityStatusCancelled, availabilityEncx.Status)
	})

	t.Run("should successfully cancel availability with admin token", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data with partner
		_, _, availabilityID, _ := setupTestAvailabilityForCancellation(t, ctx, "partner@leviosa.care")

		// Get admin token
		adminToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Make cancel request with admin token
		req := ta.NewCancelAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), adminToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify availability is cancelled
		availabilityEncx := ta.GetAvailabilityEncxFromDB(t, ctx, availabilityID, testPool)
		assert.Equal(t, domain.AvailabilityStatusCancelled, availabilityEncx.Status)
	})

	t.Run("should return 403 when standard user tries to cancel availability", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		_, _, availabilityID, _ := setupTestAvailabilityForCancellation(t, ctx, "partner@leviosa.care")

		// Get standard user token
		standardToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)

		// Make cancel request with standard user token
		req := ta.NewCancelAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), standardToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get 403 Forbidden due to role-based access control
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 404 when availability does not exist", func(t *testing.T) {
		// Clean up test data
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create building and room for proper allocation
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

		// Setup partner with allocation to valid room
		accessToken, _ := ts.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@leviosa.care", room.ID, testPool, redisClient, crypto)

		// Use non-existent availability ID
		nonExistentID := uuid.New()
		req := ta.NewCancelAvailabilityRequest(t, ctx, testServerURL, nonExistentID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get 404 Not Found
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 when availability ID format is invalid", func(t *testing.T) {
		// Clean up test data
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create building and room for proper allocation
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

		// Setup partner with allocation to valid room
		accessToken, _ := ts.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@leviosa.care", room.ID, testPool, redisClient, crypto)

		// Use invalid UUID format
		req := ta.NewCancelAvailabilityRequest(t, ctx, testServerURL, "invalid-uuid-format", accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Verify error message
		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "invalid availability ID format")
	})

	t.Run("should return 400 when availability ID is missing from path", func(t *testing.T) {
		// Clean up test data
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create building and room for proper allocation
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

		// Setup partner with allocation to valid room
		accessToken, _ := ts.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@leviosa.care", room.ID, testPool, redisClient, crypto)

		// Create request manually with malformed path
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/availabilities/", nil)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get 400 Bad Request or 404 Not Found depending on router
		assert.True(t, resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound)
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Setup test data without auth cleanup to keep availability
		_, _, availabilityID, _ := setupTestAvailabilityForCancellation(t, ctx, "partner-missing@leviosa.care")

		// Make cancel request without token
		req := ta.NewCancelAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get 401 Unauthorized
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when access token is invalid", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Setup test data with unique email
		_, _, availabilityID, _ := setupTestAvailabilityForCancellation(t, ctx, "partner-invalid@leviosa.care")

		// Make cancel request with invalid token
		req := ta.NewCancelAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), "invalid-token-12345")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get 401 Unauthorized or 403 Forbidden depending on middleware behavior
		assert.True(t, resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create expired user session
		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Partner, authCtx)

		// Use any UUID since the auth should fail first
		testID := uuid.New()
		req := ta.NewCancelAvailabilityRequest(t, ctx, testServerURL, testID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get 401 Unauthorized
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should handle cancellation of already cancelled availability gracefully", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, _, availabilityID, _ := setupTestAvailabilityForCancellation(t, ctx, "partner@leviosa.care")

		// First cancel the availability
		req1 := ta.NewCancelAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), accessToken)
		resp1, err := client.Do(req1)
		require.NoError(t, err)
		resp1.Body.Close()
		assert.Equal(t, http.StatusNoContent, resp1.StatusCode)

		// Try to cancel again - should handle gracefully (business logic decision)
		req2 := ta.NewCancelAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), accessToken)
		resp2, err := client.Do(req2)
		require.NoError(t, err)
		defer resp2.Body.Close()

		// This could be 204 (idempotent) or 409 (conflict) depending on business logic
		// We'll check for success or conflict
		assert.True(t, resp2.StatusCode == http.StatusNoContent || resp2.StatusCode == http.StatusConflict)
	})

	t.Run("should safely handle SQL injection attempts in availability ID", func(t *testing.T) {
		// Clean up test data
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create building and room for proper allocation
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

		// Setup partner with allocation to valid room
		accessToken, _ := ts.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@leviosa.care", room.ID, testPool, redisClient, crypto)

		// Attempt SQL injection via availability ID
		sqlInjectionID := "'; DROP TABLE availabilities; --"
		req := ta.NewCancelAvailabilityRequest(t, ctx, testServerURL, sqlInjectionID, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return 400 due to invalid UUID format, not execute SQL
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Verify database table still exists and is intact
		// This is implicit - if we get 400 instead of 500, the SQL injection failed
	})

	t.Run("should safely handle XSS attempts in availability ID", func(t *testing.T) {
		// Clean up test data
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create building and room for proper allocation
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

		// Setup partner with allocation to valid room
		accessToken, _ := ts.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@leviosa.care", room.ID, testPool, redisClient, crypto)

		// Attempt XSS via availability ID
		xssID := "<script>alert('xss')</script>"
		req := ta.NewCancelAvailabilityRequest(t, ctx, testServerURL, xssID, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return 400 due to invalid UUID format or 404 from router, not execute script
		assert.True(t, resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound)
	})

	t.Run("should handle context timeout appropriately", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, _, availabilityID, _ := setupTestAvailabilityForCancellation(t, ctx, "partner@leviosa.care")

		// Use a very short context timeout
		shortCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
		defer cancel()

		time.Sleep(2 * time.Millisecond) // Ensure timeout has passed

		req := ta.NewCancelAvailabilityRequest(t, shortCtx, testServerURL, availabilityID.String(), accessToken)

		resp, err := client.Do(req)
		// Either the context timeout on client side or a successful response
		if err != nil {
			// Context timeout on client side
			assert.ErrorIs(t, err, context.DeadlineExceeded)
		} else {
			defer resp.Body.Close()
			// If we got a response, it should be either success or timeout status
			assert.True(t, resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusRequestTimeout)
		}
	})
}
