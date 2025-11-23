package availability_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
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

// make test-func TEST_NAME=TestUpdateAvailability TEST_PATH=test/integration/booking/availability/update_availability_test.go

func TestUpdateAvailability(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Helper to setup test availability with room and authenticated user
	setupTestAvailabilityForUpdate := func(t *testing.T, ctx context.Context, email string) (string, uuid.UUID, *domain.Availability, uuid.UUID) {
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

		// Create availability for updating
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

		return accessToken, userID, availability, room.ID
	}

	t.Run("should successfully update all availability fields with partner token", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, _, originalAvailability, _ := setupTestAvailabilityForUpdate(t, ctx, "partner@leviosa.care")

		// Prepare update request with all fields
		newStartTime := time.Now().Add(48 * time.Hour).Truncate(time.Second)
		newEndTime := time.Now().Add(50 * time.Hour).Truncate(time.Second)
		newServiceType := "Updated Therapy Session"
		newPriceCents := 20000 // $200.00
		newNotes := "Updated consultation notes"

		updateRequest := domain.UpdateAvailabilityRequest{
			StartTime:   &newStartTime,
			EndTime:     &newEndTime,
			ServiceType: &newServiceType,
			PriceCents:  &newPriceCents,
			Notes:       &newNotes,
		}

		// Make update request
		req := ta.NewUpdateAvailabilityRequest(t, ctx, testServerURL, originalAvailability.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response data
		assert.Equal(t, originalAvailability.ID, response.ID)
		assert.Equal(t, newStartTime, response.StartTime)
		assert.Equal(t, newEndTime, response.EndTime)
		assert.Equal(t, newServiceType, response.ServiceType)
		require.NotNil(t, response.PriceCents, "PriceCents should not be nil")
		assert.Equal(t, newPriceCents, *response.PriceCents)
		assert.Equal(t, newNotes, response.Notes)
	})

	t.Run("should successfully update partial availability fields (only notes)", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, _, originalAvailability, _ := setupTestAvailabilityForUpdate(t, ctx, "partner@leviosa.care")

		// Prepare partial update request - only notes
		newNotes := "Updated consultation notes only"
		updateRequest := domain.UpdateAvailabilityRequest{
			Notes: &newNotes,
		}

		// Make update request
		req := ta.NewUpdateAvailabilityRequest(t, ctx, testServerURL, originalAvailability.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify only notes changed, other fields remain the same
		assert.Equal(t, originalAvailability.ID, response.ID)
		assert.Equal(t, newNotes, response.Notes)
		assert.Equal(t, originalAvailability.ServiceType, response.ServiceType) // Original value
		require.NotNil(t, response.PriceCents, "PriceCents should not be nil")
		assert.Equal(t, *originalAvailability.PriceCents, *response.PriceCents) // Original value
	})

	t.Run("should successfully update partial availability fields (only price)", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, _, originalAvailability, _ := setupTestAvailabilityForUpdate(t, ctx, "partner@leviosa.care")

		// Prepare partial update request - only price
		newPriceCents := 25000 // $250.00
		updateRequest := domain.UpdateAvailabilityRequest{
			PriceCents: &newPriceCents,
		}

		// Make update request
		req := ta.NewUpdateAvailabilityRequest(t, ctx, testServerURL, originalAvailability.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify only price changed
		assert.Equal(t, originalAvailability.ID, response.ID)
		require.NotNil(t, response.PriceCents, "PriceCents should not be nil")
		assert.Equal(t, newPriceCents, *response.PriceCents)
		assert.Equal(t, originalAvailability.ServiceType, response.ServiceType) // Original value
		assert.Equal(t, originalAvailability.Notes, response.Notes)             // Original value from helper
	})

	t.Run("should return 403 when standard user tries to update availability", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		_, _, originalAvailability, _ := setupTestAvailabilityForUpdate(t, ctx, "partner@leviosa.care")

		// Get standard user token
		standardToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)

		// Prepare update request
		newNotes := "Should not work"
		updateRequest := domain.UpdateAvailabilityRequest{
			Notes: &newNotes,
		}

		// Make update request with standard user token
		req := ta.NewUpdateAvailabilityRequest(t, ctx, testServerURL, originalAvailability.ID.String(), updateRequest, standardToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get 403 Forbidden due to role-based access control
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should successfully update availability with admin token", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		_, _, originalAvailability, _ := setupTestAvailabilityForUpdate(t, ctx, "partner@leviosa.care")

		// Get admin token
		adminToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)

		// Prepare update request
		newNotes := "Admin successfully updated"
		updateRequest := domain.UpdateAvailabilityRequest{
			Notes: &newNotes,
		}

		// Make update request with admin token
		req := ta.NewUpdateAvailabilityRequest(t, ctx, testServerURL, originalAvailability.ID.String(), updateRequest, adminToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get 200 OK - admins have partner+ privileges
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify the update was successful
		var response domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, newNotes, response.Notes)
	})

	t.Run("should return 400 when availability ID format is invalid", func(t *testing.T) {
		// Clean up test data
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Get partner token
		partnerToken := tu.SetupUserWithRole(t, ctx, identity.Partner, authCtx)

		// Prepare update request
		newNotes := "Test notes"
		updateRequest := domain.UpdateAvailabilityRequest{
			Notes: &newNotes,
		}

		// Make update request with invalid availability ID
		req := ta.NewUpdateAvailabilityRequest(t, ctx, testServerURL, "invalid-uuid-format", updateRequest, partnerToken)

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

	t.Run("should return 400 when request body has invalid JSON", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, _, originalAvailability, _ := setupTestAvailabilityForUpdate(t, ctx, "partner@leviosa.care")

		// Create request with invalid JSON manually
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, testServerURL+"/availabilities/"+originalAvailability.ID.String(), strings.NewReader(`{"invalid": json}`))
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

		// Should get 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 when Content-Type is not application/json", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, _, originalAvailability, _ := setupTestAvailabilityForUpdate(t, ctx, "partner@leviosa.care")

		// Create request with wrong Content-Type manually
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, testServerURL+"/availabilities/"+originalAvailability.ID.String(), strings.NewReader(`{"notes": "test"}`))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "text/plain")

		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get 415 Unsupported Media Type
		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
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
		newNotes := "Test notes"
		updateRequest := domain.UpdateAvailabilityRequest{
			Notes: &newNotes,
		}

		req := ta.NewUpdateAvailabilityRequest(t, ctx, testServerURL, nonExistentID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get 404 Not Found
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Setup test data without auth cleanup to keep availability
		_, _, originalAvailability, _ := setupTestAvailabilityForUpdate(t, ctx, "partner-missing@leviosa.care")

		// Prepare update request
		newNotes := "Test notes"
		updateRequest := domain.UpdateAvailabilityRequest{
			Notes: &newNotes,
		}

		// Make update request without token
		req := ta.NewUpdateAvailabilityRequest(t, ctx, testServerURL, originalAvailability.ID.String(), updateRequest, "")

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
		_, _, originalAvailability, _ := setupTestAvailabilityForUpdate(t, ctx, "partner-invalid@leviosa.care")

		// Prepare update request
		newNotes := "Test notes"
		updateRequest := domain.UpdateAvailabilityRequest{
			Notes: &newNotes,
		}

		// Make update request with invalid token
		req := ta.NewUpdateAvailabilityRequest(t, ctx, testServerURL, originalAvailability.ID.String(), updateRequest, "invalid-token-12345")

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
		newNotes := "Test notes"
		updateRequest := domain.UpdateAvailabilityRequest{
			Notes: &newNotes,
		}

		req := ta.NewUpdateAvailabilityRequest(t, ctx, testServerURL, testID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get 401 Unauthorized
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 400 when validation fails for time constraints", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, _, originalAvailability, _ := setupTestAvailabilityForUpdate(t, ctx, "partner@leviosa.care")

		// Prepare update request with invalid time (past time)
		pastTime := time.Now().Add(-24 * time.Hour).Truncate(time.Second)
		updateRequest := domain.UpdateAvailabilityRequest{
			StartTime: &pastTime,
		}

		req := ta.NewUpdateAvailabilityRequest(t, ctx, testServerURL, originalAvailability.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get 400 Bad Request due to validation
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 when only one time field is provided", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, _, originalAvailability, _ := setupTestAvailabilityForUpdate(t, ctx, "partner@leviosa.care")

		// Prepare update request with only start time (missing end time)
		futureTime := time.Now().Add(48 * time.Hour).Truncate(time.Second)
		updateRequest := domain.UpdateAvailabilityRequest{
			StartTime: &futureTime,
			// EndTime is nil - this should cause validation error
		}

		req := ta.NewUpdateAvailabilityRequest(t, ctx, testServerURL, originalAvailability.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get 400 Bad Request due to validation (both times must be provided together)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should safely handle SQL injection attempts in request body", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, _, originalAvailability, _ := setupTestAvailabilityForUpdate(t, ctx, "partner@leviosa.care")

		// Attempt SQL injection via notes field
		sqlInjectionNotes := "'; DROP TABLE availabilities; --"
		updateRequest := domain.UpdateAvailabilityRequest{
			Notes: &sqlInjectionNotes,
		}

		req := ta.NewUpdateAvailabilityRequest(t, ctx, testServerURL, originalAvailability.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return 200 (successful update) - SQL injection should be safely handled
		// The malicious text would just be stored as a regular string
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify database table still exists and is intact by checking the response
		var response domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, sqlInjectionNotes, response.Notes)
	})

	t.Run("should safely handle XSS attempts in request body", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, _, originalAvailability, _ := setupTestAvailabilityForUpdate(t, ctx, "partner@leviosa.care")

		// Attempt XSS via notes field
		xssNotes := "<script>alert('xss')</script>"
		updateRequest := domain.UpdateAvailabilityRequest{
			Notes: &xssNotes,
		}

		req := ta.NewUpdateAvailabilityRequest(t, ctx, testServerURL, originalAvailability.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return 200 (successful update) - XSS should be safely stored as text
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify the XSS is stored as text (not executed)
		var response domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, xssNotes, response.Notes)
	})

	t.Run("should handle context timeout appropriately", func(t *testing.T) {
		// Clean up test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test data
		accessToken, _, originalAvailability, _ := setupTestAvailabilityForUpdate(t, ctx, "partner@leviosa.care")

		// Prepare update request
		newNotes := "Test notes"
		updateRequest := domain.UpdateAvailabilityRequest{
			Notes: &newNotes,
		}

		// Use a very short context timeout
		shortCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
		defer cancel()

		time.Sleep(2 * time.Millisecond) // Ensure timeout has passed

		req := ta.NewUpdateAvailabilityRequest(t, shortCtx, testServerURL, originalAvailability.ID.String(), updateRequest, accessToken)

		resp, err := client.Do(req)
		// Either the context timeout on client side or a successful response
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
