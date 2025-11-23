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

// make test-func TEST_NAME=TestBlockAvailability TEST_PATH=test/integration/booking/availability/block_availability_test.go

func TestBlockAvailability(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Helper to setup test availability with room and authenticated user
	setupTestAvailabilityForBlocking := func(t *testing.T, ctx context.Context, email string) (string, uuid.UUID, uuid.UUID, uuid.UUID) {
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

		// Create availability for blocking
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

	t.Run("should successfully block availability with admin token", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data with partner
		_, _, availabilityID, _ := setupTestAvailabilityForBlocking(t, ctx, "partner@leviosa.care")

		// Get admin token
		adminToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Make block request with admin token
		req := ta.NewBlockAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), adminToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response - should be 204 No Content
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify availability is actually blocked in database
		availabilityEncx := ta.GetAvailabilityEncxFromDB(t, ctx, availabilityID, testPool)
		assert.Equal(t, domain.AvailabilityStatusBlocked, availabilityEncx.Status)
	})

	t.Run("should return 403 when partner tries to block availability", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, _, availabilityID, _ := setupTestAvailabilityForBlocking(t, ctx, "partner@leviosa.care")

		// Make block request with partner token
		req := ta.NewBlockAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get 403 Forbidden due to role-based access control
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 403 when standard user tries to block availability", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		_, _, availabilityID, _ := setupTestAvailabilityForBlocking(t, ctx, "partner@leviosa.care")

		// Get standard user token
		standardToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)

		// Make block request with standard user token
		req := ta.NewBlockAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), standardToken)

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

		// Get admin token
		adminToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Use non-existent availability ID
		nonExistentID := uuid.New()
		req := ta.NewBlockAvailabilityRequest(t, ctx, testServerURL, nonExistentID.String(), adminToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get 404 Not Found
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 when availability ID format is invalid", func(t *testing.T) {
		// Clean up test data
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Get admin token
		adminToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Use invalid UUID format
		req := ta.NewBlockAvailabilityRequest(t, ctx, testServerURL, "invalid-uuid-format", adminToken)

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
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Get admin token
		adminToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Create request manually with malformed path
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/availabilities/", nil)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: adminToken,
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
		_, _, availabilityID, _ := setupTestAvailabilityForBlocking(t, ctx, "partner-missing@leviosa.care")

		// Make block request without token
		req := ta.NewBlockAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), "")

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
		_, _, availabilityID, _ := setupTestAvailabilityForBlocking(t, ctx, "partner-invalid@leviosa.care")

		// Make block request with invalid token
		req := ta.NewBlockAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), "invalid-token-12345")

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
		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Use any UUID since the auth should fail first
		testID := uuid.New()
		req := ta.NewBlockAvailabilityRequest(t, ctx, testServerURL, testID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get 401 Unauthorized
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should handle blocking of already blocked availability gracefully", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		_, _, availabilityID, _ := setupTestAvailabilityForBlocking(t, ctx, "partner@leviosa.care")

		// Get admin token
		adminToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// First block the availability
		req1 := ta.NewBlockAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), adminToken)
		resp1, err := client.Do(req1)
		require.NoError(t, err)
		resp1.Body.Close()
		assert.Equal(t, http.StatusNoContent, resp1.StatusCode)

		// Try to block again - should handle gracefully (business logic decision)
		req2 := ta.NewBlockAvailabilityRequest(t, ctx, testServerURL, availabilityID.String(), adminToken)
		resp2, err := client.Do(req2)
		require.NoError(t, err)
		defer resp2.Body.Close()

		// This could be 204 (idempotent) or 409 (conflict) depending on business logic
		// We'll check for success or conflict
		assert.True(t, resp2.StatusCode == http.StatusNoContent || resp2.StatusCode == http.StatusConflict)
	})

	t.Run("should safely handle SQL injection attempts in availability ID", func(t *testing.T) {
		// Clean up test data
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Get admin token
		adminToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Attempt SQL injection via availability ID
		sqlInjectionID := "'; DROP TABLE availabilities; --"
		req := ta.NewBlockAvailabilityRequest(t, ctx, testServerURL, sqlInjectionID, adminToken)

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
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Get admin token
		adminToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Attempt XSS via availability ID
		xssID := "<script>alert('xss')</script>"
		req := ta.NewBlockAvailabilityRequest(t, ctx, testServerURL, xssID, adminToken)

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
		_, _, availabilityID, _ := setupTestAvailabilityForBlocking(t, ctx, "partner@leviosa.care")

		// Get admin token
		adminToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Use a very short context timeout
		shortCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
		defer cancel()

		time.Sleep(2 * time.Millisecond) // Ensure timeout has passed

		req := ta.NewBlockAvailabilityRequest(t, shortCtx, testServerURL, availabilityID.String(), adminToken)

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
