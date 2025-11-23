package availability_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	tsetup "github.com/Leviosa-care/leviosa/backend/test/helpers/booking"
	ta "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreateRecurringAvailability TEST_PATH=test/integration/booking/availability/create_recurring_availability_test.go

func TestCreateRecurringAvailability(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Setup test building and room for availability tests
	setupTestRoom := func(t *testing.T, ctx context.Context) uuid.UUID {
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

		return room.ID
	}

	// Create valid daily recurrence request helper
	createValidDailyRequest := func(roomID uuid.UUID) domain.CreateRecurringAvailabilityRequest {
		startTime := time.Now().Add(24 * time.Hour).Truncate(time.Second)
		endTime := startTime.Add(2 * time.Hour)
		until := startTime.Add(30 * 24 * time.Hour) // 30 days

		return domain.CreateRecurringAvailabilityRequest{
			RoomID:      roomID,
			StartTime:   startTime,
			EndTime:     endTime,
			MaxCapacity: 10,
			Pattern: domain.RecurrencePattern{
				Type:     "daily",
				Interval: 1,
				Until:    &until,
			},
		}
	}

	// Create valid weekly recurrence request helper
	createValidWeeklyRequest := func(roomID uuid.UUID) domain.CreateRecurringAvailabilityRequest {
		startTime := time.Now().Add(24 * time.Hour).Truncate(time.Second)
		endTime := startTime.Add(2 * time.Hour)
		until := startTime.Add(60 * 24 * time.Hour) // 60 days

		return domain.CreateRecurringAvailabilityRequest{
			RoomID:      roomID,
			StartTime:   startTime,
			EndTime:     endTime,
			MaxCapacity: 10,
			Pattern: domain.RecurrencePattern{
				Type:       "weekly",
				Interval:   1,
				Until:      &until,
				DaysOfWeek: []int{1, 3, 5}, // Monday, Wednesday, Friday
			},
		}
	}

	t.Run("should successfully create daily recurring availability", func(t *testing.T) {
		// Clean test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		// Prepare request
		request := createValidDailyRequest(roomID)

		req := ta.NewCreateRecurringAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		// Execute request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response data
		assert.NotNil(t, response.ID)
		assert.NotEqual(t, uuid.Nil, response.UserID)
		assert.Equal(t, request.RoomID, response.RoomID)
		assert.WithinDuration(t, request.StartTime, response.StartTime, time.Second)
		assert.WithinDuration(t, request.EndTime, response.EndTime, time.Second)
		assert.Equal(t, request.MaxCapacity, response.MaxCapacity)
		assert.Equal(t, domain.AvailabilityStatusAvailable, response.Status)

		// Verify availability exists in database
		availabilityEncx := ta.GetAvailabilityEncxFromDB(t, ctx, response.ID, testPool)
		require.NotNil(t, availabilityEncx)

		// Decrypt and verify
		availability, err := domain.DecryptAvailabilityEncx(ctx, crypto, availabilityEncx)
		require.NoError(t, err)

		assert.Equal(t, request.RoomID, availability.RoomID)
		assert.WithinDuration(t, request.StartTime, availability.StartTime, time.Second)
		assert.WithinDuration(t, request.EndTime, availability.EndTime, time.Second)
		assert.Equal(t, request.MaxCapacity, availability.MaxCapacity)
		assert.Equal(t, domain.AvailabilityStatusAvailable, availability.Status)
		assert.True(t, availability.IsRecurring)
		assert.NotNil(t, availability.RecurrencePattern)
		assert.Equal(t, "daily", availability.RecurrencePattern.Type)
		assert.Equal(t, 1, availability.RecurrencePattern.Interval)
	})

	t.Run("should successfully create weekly recurring availability", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner2@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		// Prepare request
		request := createValidWeeklyRequest(roomID)

		req := ta.NewCreateRecurringAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		// Execute request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response data
		assert.NotNil(t, response.ID)
		assert.Equal(t, request.RoomID, response.RoomID)

		// Verify in database
		availabilityEncx := ta.GetAvailabilityEncxFromDB(t, ctx, response.ID, testPool)
		availability, err := domain.DecryptAvailabilityEncx(ctx, crypto, availabilityEncx)
		require.NoError(t, err)

		assert.True(t, availability.IsRecurring)
		assert.NotNil(t, availability.RecurrencePattern)
		assert.Equal(t, "weekly", availability.RecurrencePattern.Type)
		assert.Equal(t, 1, availability.RecurrencePattern.Interval)
		assert.Equal(t, []int{1, 3, 5}, availability.RecurrencePattern.DaysOfWeek)
	})

	t.Run("should successfully create recurring availability with service details", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner3@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		// Create request with service details
		request := createValidDailyRequest(roomID)
		priceCents := 5000 // $50.00
		request.ServiceType = "Massage Therapy"
		request.PriceCents = &priceCents
		request.Notes = "Bring comfortable clothing"

		req := ta.NewCreateRecurringAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify service details in response
		assert.Equal(t, request.ServiceType, response.ServiceType)
		assert.Equal(t, *request.PriceCents, *response.PriceCents)
		assert.Equal(t, request.Notes, response.Notes)

		// Verify in database
		availabilityEncx := ta.GetAvailabilityEncxFromDB(t, ctx, response.ID, testPool)
		availability, err := domain.DecryptAvailabilityEncx(ctx, crypto, availabilityEncx)
		require.NoError(t, err)

		assert.Equal(t, request.ServiceType, availability.ServiceType)
		assert.Equal(t, *request.PriceCents, *availability.PriceCents)
		assert.Equal(t, request.Notes, availability.Notes)
	})

	t.Run("should successfully create recurring availability with admin token", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup admin with room allocation
		accessToken, userID := tsetup.SetupAdminWithAllocation(t, ctx, roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		request := createValidDailyRequest(roomID)

		req := ta.NewCreateRecurringAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotNil(t, response.ID)
	})

	t.Run("should return 400 Bad Request for missing room_id", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner4@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		request := createValidDailyRequest(roomID)
		request.RoomID = uuid.Nil // Invalid room ID

		req := ta.NewCreateRecurringAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "room_id")
	})

	t.Run("should return 400 Bad Request for invalid recurrence type", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner5@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		request := createValidDailyRequest(roomID)
		request.Pattern.Type = "invalid_type"

		req := ta.NewCreateRecurringAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 Bad Request for invalid interval", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner6@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		request := createValidDailyRequest(roomID)
		request.Pattern.Interval = 0 // Invalid interval

		req := ta.NewCreateRecurringAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 Bad Request for start time in the past", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner7@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		request := createValidDailyRequest(roomID)
		request.StartTime = time.Now().Add(-24 * time.Hour) // Past time
		request.EndTime = request.StartTime.Add(2 * time.Hour)

		req := ta.NewCreateRecurringAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrInvalidValue.Error())
	})

	t.Run("should return 400 Bad Request for invalid JSON body", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Partner, authCtx)

		// Create malformed JSON
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/availabilities/recurring", bytes.NewBuffer([]byte("{invalid json")))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		if accessToken != "" {
			cookie := &http.Cookie{
				Name:  ck.AccessTokenCookieName,
				Value: accessToken,
			}
			req.AddCookie(cookie)
		}

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrInvalidValue.Error())
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t, ctx)
		request := createValidDailyRequest(roomID)

		req := ta.NewCreateRecurringAvailabilityRequest(t, ctx, testServerURL, request, "") // Empty token

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create expired partner session
		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Partner, authCtx)
		roomID := setupTestRoom(t, ctx)

		request := createValidDailyRequest(roomID)

		req := ta.NewCreateRecurringAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t, ctx)
		request := createValidDailyRequest(roomID)

		req := ta.NewCreateRecurringAvailabilityRequest(t, ctx, testServerURL, request, "invalid-token-12345")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role (standard user)", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create standard user (not partner or admin)
		accessToken := tu.SetupStandardUser(t, ctx, authCtx)
		roomID := setupTestRoom(t, ctx)

		request := createValidDailyRequest(roomID)

		req := ta.NewCreateRecurringAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}
