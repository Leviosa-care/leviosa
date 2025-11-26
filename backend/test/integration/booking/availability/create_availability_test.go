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

// make test-func TEST_NAME=TestCreateAvailability TEST_PATH=test/integration/booking/availability/create_availability_test.go

func TestCreateAvailability(t *testing.T) {
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

		// Create room schedule for all days of the week (8:00 AM - 8:00 PM)
		for dayOfWeek := 0; dayOfWeek < 7; dayOfWeek++ {
			schedule := ta.NewTestRoomScheduleRecurring(room.ID, dayOfWeek, "08:00", "20:00")
			ta.InsertRoomSchedule(t, ctx, schedule, testPool)
		}

		return room.ID
	}

	// Create valid request helper
	createValidRequest := func(roomID uuid.UUID) domain.CreateAvailabilityRequest {
		// Create start time tomorrow at 9:00 AM (within operating hours 08:00-20:00)
		tomorrow := time.Now().AddDate(0, 0, 1)
		startTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 9, 0, 0, 0, tomorrow.Location())
		endTime := startTime.Add(2 * time.Hour) // 9:00 AM - 11:00 AM
		return domain.CreateAvailabilityRequest{
			RoomID:      roomID,
			StartTime:   startTime,
			EndTime:     endTime,
			MaxCapacity: 10,
		}
	}

	t.Run("should successfully create availability with valid partner token", func(t *testing.T) {
		// Clean test data
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID // userID available if needed for assertions

		// Prepare request
		request := createValidRequest(roomID)

		req := ta.NewCreateAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

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
	})

	t.Run("should successfully create availability with service details", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner2@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID // userID available if needed for assertions

		// Create request with service details
		request := createValidRequest(roomID)
		priceCents := 5000 // $50.00
		request.ServiceType = "Massage Therapy"
		request.PriceCents = &priceCents
		request.Notes = "Bring comfortable clothing"

		req := ta.NewCreateAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify service details in response
		assert.Equal(t, request.ServiceType, response.ServiceType)
		require.NotNil(t, response.PriceCents)
		assert.Equal(t, *request.PriceCents, *response.PriceCents)
		assert.Equal(t, request.Notes, response.Notes)

		// Verify in database with encryption
		availabilityEncx := ta.GetAvailabilityEncxFromDB(t, ctx, response.ID, testPool)
		availability, err := domain.DecryptAvailabilityEncx(ctx, crypto, availabilityEncx)
		require.NoError(t, err)

		assert.Equal(t, request.ServiceType, availability.ServiceType)
		assert.Equal(t, *request.PriceCents, *availability.PriceCents)
		assert.Equal(t, request.Notes, availability.Notes)
	})

	t.Run("should successfully create availability with admin token", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		ta.ClearRoomSchedulesTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup admin with room allocation
		accessToken, userID := tsetup.SetupAdminWithAllocation(t, ctx, roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		request := createValidRequest(roomID)

		req := ta.NewCreateAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

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
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner3@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		request := createValidRequest(roomID)
		request.RoomID = uuid.Nil // Invalid room ID

		req := ta.NewCreateAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

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

	t.Run("should return 400 Bad Request for start time in the past", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner4@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		request := createValidRequest(roomID)
		request.StartTime = time.Now().Add(-24 * time.Hour) // Past time
		request.EndTime = request.StartTime.Add(2 * time.Hour)

		req := ta.NewCreateAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

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

	t.Run("should return 400 Bad Request for start time after end time", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner5@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		request := createValidRequest(roomID)
		request.StartTime = time.Now().Add(48 * time.Hour)
		request.EndTime = request.StartTime.Add(-2 * time.Hour) // End before start

		req := ta.NewCreateAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

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

	t.Run("should return 400 Bad Request for duration less than 15 minutes", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner6@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		request := createValidRequest(roomID)
		request.StartTime = time.Now().Add(24 * time.Hour)
		request.EndTime = request.StartTime.Add(10 * time.Minute) // Less than 15 minutes

		req := ta.NewCreateAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "15 minutes")
	})

	t.Run("should return 400 Bad Request for duration more than 12 hours", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner7@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		request := createValidRequest(roomID)
		request.StartTime = time.Now().Add(24 * time.Hour)
		request.EndTime = request.StartTime.Add(13 * time.Hour) // More than 12 hours

		req := ta.NewCreateAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "12 hours")
	})

	t.Run("should return 400 Bad Request for max capacity less than 1", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner8@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		request := createValidRequest(roomID)
		request.MaxCapacity = 0 // Invalid capacity

		req := ta.NewCreateAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "max_capacity")
	})

	t.Run("should return 400 Bad Request for max capacity greater than 50", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner9@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		request := createValidRequest(roomID)
		request.MaxCapacity = 51 // Exceeds limit

		req := ta.NewCreateAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "max_capacity")
	})

	t.Run("should return 400 Bad Request for service type exceeding 255 characters", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner10@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		request := createValidRequest(roomID)
		request.ServiceType = string(make([]byte, 256)) // Exceeds limit

		req := ta.NewCreateAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "service_type")
	})

	t.Run("should return 400 Bad Request for notes exceeding 1000 characters", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner11@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		request := createValidRequest(roomID)
		request.Notes = string(make([]byte, 1001)) // Exceeds limit

		req := ta.NewCreateAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "notes")
	})

	t.Run("should return 400 Bad Request for negative price", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner12@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		request := createValidRequest(roomID)
		negativePriceCents := -100
		request.PriceCents = &negativePriceCents

		req := ta.NewCreateAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "price")
	})

	t.Run("should return 400 Bad Request for price exceeding maximum", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner13@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		request := createValidRequest(roomID)
		excessivePriceCents := 1000000 // More than $9,999.99
		request.PriceCents = &excessivePriceCents

		req := ta.NewCreateAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "price")
	})

	t.Run("should return 400 Bad Request for invalid JSON body", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Partner, authCtx)

		// Create malformed JSON
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/availabilities", bytes.NewBuffer([]byte("{invalid json")))
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

	t.Run("should return 415 Unsupported Media Type when Content-Type is not application/json", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner14@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		request := createValidRequest(roomID)
		jsonBody, err := json.Marshal(request)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/availabilities", bytes.NewReader(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "text/plain") // Wrong content type

		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t, ctx)
		request := createValidRequest(roomID)

		req := ta.NewCreateAvailabilityRequest(t, ctx, testServerURL, request, "") // Empty token

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

		request := createValidRequest(roomID)

		req := ta.NewCreateAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t, ctx)
		request := createValidRequest(roomID)

		req := ta.NewCreateAvailabilityRequest(t, ctx, testServerURL, request, "invalid-token-12345")

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

		request := createValidRequest(roomID)

		req := ta.NewCreateAvailabilityRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should handle context timeout appropriately", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup test room
		roomID := setupTestRoom(t, ctx)

		// Setup authenticated partner with room allocation
		accessToken, userID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner15@test.com", roomID, testPool, authCtx.Redis, crypto)
		_ = userID

		request := createValidRequest(roomID)

		// Use a very short context timeout to potentially trigger timeout
		shortCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
		defer cancel()

		time.Sleep(2 * time.Millisecond) // Ensure timeout has passed

		req := ta.NewCreateAvailabilityRequest(t, shortCtx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		// Either the context timeout or a successful response (if operation was fast enough)
		if err != nil {
			// Context timeout on client side
			assert.ErrorIs(t, err, context.DeadlineExceeded)
		} else {
			defer resp.Body.Close()
			// If we got a response, it should be either success or timeout status
			assert.True(t, resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusRequestTimeout)
		}
	})
}
