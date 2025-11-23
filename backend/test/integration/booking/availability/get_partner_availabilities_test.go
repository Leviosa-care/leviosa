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
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	ta "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetPartnerAvailabilities TEST_PATH=test/integration/booking/availability/get_partner_availabilities_test.go

func TestGetPartnerAvailabilities(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Helper to setup test environment with building and room
	setupTestEnvironment := func(t *testing.T, ctx context.Context) (uuid.UUID, uuid.UUID) {
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

		return building.ID, room.ID
	}

	// Helper to create and insert availability
	createAvailability := func(t *testing.T, ctx context.Context, roomID, partnerID uuid.UUID, startOffset, duration time.Duration, serviceType string, priceCents *int, status domain.AvailabilityStatus) uuid.UUID {
		t.Helper()

		startTime := time.Now().Add(startOffset).Truncate(time.Second)
		endTime := startTime.Add(duration)

		availability := ta.NewTestAvailabilityWithParams(
			t,
			partnerID,
			roomID,
			startTime,
			endTime,
			serviceType,
			priceCents,
			1,
			status,
		)

		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

		return availability.ID
	}

	t.Run("should successfully retrieve all availabilities for a partner", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)
		_, roomID := setupTestEnvironment(t, ctx)
		partnerID := uuid.New()

		// Create multiple availabilities for the partner
		createAvailability(t, ctx, roomID, partnerID, 24*time.Hour, 1*time.Hour, "Consultation", intPtr(15000), domain.AvailabilityStatusAvailable)
		createAvailability(t, ctx, roomID, partnerID, 48*time.Hour, 2*time.Hour, "Therapy", intPtr(20000), domain.AvailabilityStatusAvailable)
		createAvailability(t, ctx, roomID, partnerID, 72*time.Hour, 1*time.Hour, "Check-up", intPtr(10000), domain.AvailabilityStatusAvailable)

		// Make request
		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partnerID.String(), nil, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responses []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		require.NoError(t, err)

		// Verify we got all 3 availabilities
		assert.Len(t, responses, 3)

		// Verify all belong to the partner
		for _, response := range responses {
			assert.Equal(t, partnerID, response.UserID)
			assert.Equal(t, roomID, response.RoomID)
		}
	})

	t.Run("should filter availabilities by single status", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)
		_, roomID := setupTestEnvironment(t, ctx)
		partnerID := uuid.New()

		// Create availabilities with different statuses
		createAvailability(t, ctx, roomID, partnerID, 24*time.Hour, 1*time.Hour, "Available", intPtr(15000), domain.AvailabilityStatusAvailable)
		createAvailability(t, ctx, roomID, partnerID, 48*time.Hour, 1*time.Hour, "Booked", intPtr(15000), domain.AvailabilityStatusBooked)
		createAvailability(t, ctx, roomID, partnerID, 72*time.Hour, 1*time.Hour, "Cancelled", intPtr(15000), domain.AvailabilityStatusCancelled)

		// Request only available status
		queryParams := map[string]string{
			"status": string(domain.AvailabilityStatusAvailable),
		}
		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partnerID.String(), queryParams, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responses []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		require.NoError(t, err)

		// Should only return 1 available availability
		assert.Len(t, responses, 1)
		assert.Equal(t, domain.AvailabilityStatusAvailable, responses[0].Status)
	})

	t.Run("should filter availabilities by multiple statuses", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)
		_, roomID := setupTestEnvironment(t, ctx)
		partnerID := uuid.New()

		// Create availabilities with different statuses
		createAvailability(t, ctx, roomID, partnerID, 24*time.Hour, 1*time.Hour, "Available", intPtr(15000), domain.AvailabilityStatusAvailable)
		createAvailability(t, ctx, roomID, partnerID, 48*time.Hour, 1*time.Hour, "Booked", intPtr(15000), domain.AvailabilityStatusBooked)
		createAvailability(t, ctx, roomID, partnerID, 72*time.Hour, 1*time.Hour, "Cancelled", intPtr(15000), domain.AvailabilityStatusCancelled)

		// Request available and booked statuses using url.Values for multiple params
		baseURL := testServerURL + "/partners/" + partnerID.String() + "/availabilities"
		params := url.Values{}
		params.Add("status", string(domain.AvailabilityStatusAvailable))
		params.Add("status", string(domain.AvailabilityStatusBooked))
		fullURL := baseURL + "?" + params.Encode()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
		require.NoError(t, err)

		// Add auth cookie
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responses []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		require.NoError(t, err)

		// Should return 2 availabilities (available and booked)
		assert.Len(t, responses, 2)
		statuses := []domain.AvailabilityStatus{responses[0].Status, responses[1].Status}
		assert.Contains(t, statuses, domain.AvailabilityStatusAvailable)
		assert.Contains(t, statuses, domain.AvailabilityStatusBooked)
		assert.NotContains(t, statuses, domain.AvailabilityStatusCancelled)
	})

	t.Run("should filter availabilities by start time", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)
		_, roomID := setupTestEnvironment(t, ctx)
		partnerID := uuid.New()

		// Create availabilities at different times
		createAvailability(t, ctx, roomID, partnerID, 24*time.Hour, 1*time.Hour, "Day 1", intPtr(15000), domain.AvailabilityStatusAvailable)
		createAvailability(t, ctx, roomID, partnerID, 72*time.Hour, 1*time.Hour, "Day 3", intPtr(15000), domain.AvailabilityStatusAvailable)
		createAvailability(t, ctx, roomID, partnerID, 120*time.Hour, 1*time.Hour, "Day 5", intPtr(15000), domain.AvailabilityStatusAvailable)

		// Request availabilities starting after 48 hours from now
		filterTime := time.Now().Add(48 * time.Hour)
		queryParams := map[string]string{
			"start_time": filterTime.Format(time.RFC3339),
		}
		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partnerID.String(), queryParams, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responses []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		require.NoError(t, err)

		// Should return 2 availabilities (Day 3 and Day 5)
		assert.Len(t, responses, 2)
		for _, response := range responses {
			assert.True(t, response.StartTime.After(filterTime) || response.StartTime.Equal(filterTime))
		}
	})

	t.Run("should filter availabilities by end time", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)
		_, roomID := setupTestEnvironment(t, ctx)
		partnerID := uuid.New()

		// Create availabilities at different times
		createAvailability(t, ctx, roomID, partnerID, 24*time.Hour, 1*time.Hour, "Day 1", intPtr(15000), domain.AvailabilityStatusAvailable)
		createAvailability(t, ctx, roomID, partnerID, 72*time.Hour, 1*time.Hour, "Day 3", intPtr(15000), domain.AvailabilityStatusAvailable)
		createAvailability(t, ctx, roomID, partnerID, 120*time.Hour, 1*time.Hour, "Day 5", intPtr(15000), domain.AvailabilityStatusAvailable)

		// Request availabilities ending before 96 hours from now
		filterTime := time.Now().Add(96 * time.Hour)
		queryParams := map[string]string{
			"end_time": filterTime.Format(time.RFC3339),
		}
		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partnerID.String(), queryParams, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responses []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		require.NoError(t, err)

		// Should return 2 availabilities (Day 1 and Day 3)
		assert.Len(t, responses, 2)
		for _, response := range responses {
			assert.True(t, response.EndTime.Before(filterTime) || response.EndTime.Equal(filterTime))
		}
	})

	t.Run("should filter availabilities by time range", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)
		_, roomID := setupTestEnvironment(t, ctx)
		partnerID := uuid.New()

		// Create availabilities at different times
		createAvailability(t, ctx, roomID, partnerID, 24*time.Hour, 1*time.Hour, "Day 1", intPtr(15000), domain.AvailabilityStatusAvailable)
		createAvailability(t, ctx, roomID, partnerID, 72*time.Hour, 1*time.Hour, "Day 3", intPtr(15000), domain.AvailabilityStatusAvailable)
		createAvailability(t, ctx, roomID, partnerID, 120*time.Hour, 1*time.Hour, "Day 5", intPtr(15000), domain.AvailabilityStatusAvailable)

		// Request availabilities in time range (48 hours to 96 hours from now)
		startFilter := time.Now().Add(48 * time.Hour)
		endFilter := time.Now().Add(96 * time.Hour)
		queryParams := map[string]string{
			"start_time": startFilter.Format(time.RFC3339),
			"end_time":   endFilter.Format(time.RFC3339),
		}
		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partnerID.String(), queryParams, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responses []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		require.NoError(t, err)

		// Should return 1 availability (Day 3)
		assert.Len(t, responses, 1)
		assert.Equal(t, "Day 3", responses[0].ServiceType)
	})

	t.Run("should respect limit parameter", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)
		_, roomID := setupTestEnvironment(t, ctx)
		partnerID := uuid.New()

		// Create 5 availabilities
		for i := 0; i < 5; i++ {
			createAvailability(t, ctx, roomID, partnerID, time.Duration(24*(i+1))*time.Hour, 1*time.Hour, "Session", intPtr(15000), domain.AvailabilityStatusAvailable)
		}

		// Request with limit of 2
		queryParams := map[string]string{
			"limit": "2",
		}
		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partnerID.String(), queryParams, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responses []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		require.NoError(t, err)

		// Should return exactly 2 availabilities
		assert.Len(t, responses, 2)
	})

	t.Run("should combine multiple filters", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)
		_, roomID := setupTestEnvironment(t, ctx)
		partnerID := uuid.New()

		// Create various availabilities
		createAvailability(t, ctx, roomID, partnerID, 24*time.Hour, 1*time.Hour, "Available 1", intPtr(15000), domain.AvailabilityStatusAvailable)
		createAvailability(t, ctx, roomID, partnerID, 48*time.Hour, 1*time.Hour, "Booked 1", intPtr(15000), domain.AvailabilityStatusBooked)
		createAvailability(t, ctx, roomID, partnerID, 72*time.Hour, 1*time.Hour, "Available 2", intPtr(15000), domain.AvailabilityStatusAvailable)
		createAvailability(t, ctx, roomID, partnerID, 96*time.Hour, 1*time.Hour, "Available 3", intPtr(15000), domain.AvailabilityStatusAvailable)

		// Request available status within time range with limit
		startFilter := time.Now().Add(48 * time.Hour)
		endFilter := time.Now().Add(120 * time.Hour)
		queryParams := map[string]string{
			"status":     string(domain.AvailabilityStatusAvailable),
			"start_time": startFilter.Format(time.RFC3339),
			"end_time":   endFilter.Format(time.RFC3339),
			"limit":      "1",
		}
		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partnerID.String(), queryParams, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responses []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		require.NoError(t, err)

		// Should return 1 available availability in the time range
		assert.Len(t, responses, 1)
		assert.Equal(t, domain.AvailabilityStatusAvailable, responses[0].Status)
		assert.True(t, responses[0].StartTime.After(startFilter) || responses[0].StartTime.Equal(startFilter))
		assert.True(t, responses[0].EndTime.Before(endFilter) || responses[0].EndTime.Equal(endFilter))
	})

	t.Run("should return empty array when no availabilities exist for partner", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)
		partnerID := uuid.New()

		// Make request for partner with no availabilities
		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partnerID.String(), nil, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responses []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		require.NoError(t, err)

		// Should return empty array
		assert.Empty(t, responses)
	})

	t.Run("should return empty array when filters match no availabilities", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)
		_, roomID := setupTestEnvironment(t, ctx)
		partnerID := uuid.New()

		// Create available availabilities
		createAvailability(t, ctx, roomID, partnerID, 24*time.Hour, 1*time.Hour, "Available", intPtr(15000), domain.AvailabilityStatusAvailable)

		// Request booked status (which doesn't exist)
		queryParams := map[string]string{
			"status": string(domain.AvailabilityStatusBooked),
		}
		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partnerID.String(), queryParams, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responses []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		require.NoError(t, err)

		// Should return empty array
		assert.Empty(t, responses)
	})

	t.Run("should only return availabilities for specified partner", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)
		_, roomID := setupTestEnvironment(t, ctx)
		partner1ID := uuid.New()
		partner2ID := uuid.New()

		// Create availabilities for partner 1
		createAvailability(t, ctx, roomID, partner1ID, 24*time.Hour, 1*time.Hour, "Partner 1 - Slot 1", intPtr(15000), domain.AvailabilityStatusAvailable)
		createAvailability(t, ctx, roomID, partner1ID, 48*time.Hour, 1*time.Hour, "Partner 1 - Slot 2", intPtr(15000), domain.AvailabilityStatusAvailable)

		// Create availabilities for partner 2
		createAvailability(t, ctx, roomID, partner2ID, 24*time.Hour, 1*time.Hour, "Partner 2 - Slot 1", intPtr(15000), domain.AvailabilityStatusAvailable)
		createAvailability(t, ctx, roomID, partner2ID, 48*time.Hour, 1*time.Hour, "Partner 2 - Slot 2", intPtr(15000), domain.AvailabilityStatusAvailable)

		// Request availabilities for partner 1 only
		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partner1ID.String(), nil, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responses []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		require.NoError(t, err)

		// Should return only partner 1's availabilities
		assert.Len(t, responses, 2)
		for _, response := range responses {
			assert.Equal(t, partner1ID, response.UserID)
		}
	})

	t.Run("should work with partner role", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Partner, authCtx)
		_, roomID := setupTestEnvironment(t, ctx)
		partnerID := uuid.New()

		createAvailability(t, ctx, roomID, partnerID, 24*time.Hour, 1*time.Hour, "Consultation", intPtr(15000), domain.AvailabilityStatusAvailable)

		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partnerID.String(), nil, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responses []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		require.NoError(t, err)

		assert.Len(t, responses, 1)
	})

	t.Run("should work with admin role", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Administrator, authCtx)
		_, roomID := setupTestEnvironment(t, ctx)
		partnerID := uuid.New()

		createAvailability(t, ctx, roomID, partnerID, 24*time.Hour, 1*time.Hour, "Consultation", intPtr(15000), domain.AvailabilityStatusAvailable)

		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partnerID.String(), nil, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responses []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		require.NoError(t, err)

		assert.Len(t, responses, 1)
	})

	t.Run("should properly decrypt encrypted fields", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)
		_, roomID := setupTestEnvironment(t, ctx)
		partnerID := uuid.New()

		// Create availability with specific encrypted values
		startTime := time.Now().Add(24 * time.Hour).Truncate(time.Second)
		endTime := startTime.Add(2 * time.Hour)
		availability := &domain.Availability{
			ID:          uuid.New(),
			UserID:      partnerID,
			RoomID:      roomID,
			StartTime:   startTime,
			EndTime:     endTime,
			ServiceType: "Specialized Therapy Session",
			PriceCents:  intPtr(25000),
			MaxCapacity: 1,
			Notes:       "Bring comfortable clothing and arrive 10 minutes early",
			IsRecurring: false,
			Status:      domain.AvailabilityStatusAvailable,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		ta.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

		// Request availabilities
		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partnerID.String(), nil, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responses []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		require.NoError(t, err)

		require.Len(t, responses, 1)
		// Verify encrypted fields are properly decrypted
		assert.Equal(t, "Specialized Therapy Session", responses[0].ServiceType)
		assert.Equal(t, "Bring comfortable clothing and arrive 10 minutes early", responses[0].Notes)
		assert.Equal(t, 25000, *responses[0].PriceCents)
	})

	t.Run("should return 400 Bad Request for invalid partner ID format", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)

		// Use invalid UUID format
		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, "invalid-uuid-format", nil, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "invalid partner ID format")
	})

	t.Run("should return 400 Bad Request for missing partner ID in path", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)

		// Make request without partner ID (this will result in path parsing error)
		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, "", nil, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should get 404 or 400 depending on router behavior
		assert.True(t, resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound)
	})

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		partnerID := uuid.New()

		// Make request without token
		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partnerID.String(), nil, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create expired user session
		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Standard, authCtx)
		partnerID := uuid.New()

		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partnerID.String(), nil, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		partnerID := uuid.New()

		// Use invalid token
		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partnerID.String(), nil, "invalid-token-12345")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should handle invalid time format gracefully", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)
		_, roomID := setupTestEnvironment(t, ctx)
		partnerID := uuid.New()

		createAvailability(t, ctx, roomID, partnerID, 24*time.Hour, 1*time.Hour, "Test", intPtr(15000), domain.AvailabilityStatusAvailable)

		// Request with invalid time format (should be ignored)
		queryParams := map[string]string{
			"start_time": "invalid-time-format",
		}
		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partnerID.String(), queryParams, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should succeed and ignore invalid time filter
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responses []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		require.NoError(t, err)

		// Should return all availabilities (invalid filter ignored)
		assert.Len(t, responses, 1)
	})

	t.Run("should handle invalid limit format gracefully", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)
		_, roomID := setupTestEnvironment(t, ctx)
		partnerID := uuid.New()

		createAvailability(t, ctx, roomID, partnerID, 24*time.Hour, 1*time.Hour, "Test", intPtr(15000), domain.AvailabilityStatusAvailable)

		// Request with invalid limit format (should be ignored)
		queryParams := map[string]string{
			"limit": "not-a-number",
		}
		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partnerID.String(), queryParams, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should succeed and ignore invalid limit
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responses []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		require.NoError(t, err)

		// Should return all availabilities (invalid limit ignored)
		assert.Len(t, responses, 1)
	})

	t.Run("should handle availabilities with nil price", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)
		_, roomID := setupTestEnvironment(t, ctx)
		partnerID := uuid.New()

		// Create availability with nil price (free)
		createAvailability(t, ctx, roomID, partnerID, 24*time.Hour, 1*time.Hour, "Free Consultation", nil, domain.AvailabilityStatusAvailable)

		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partnerID.String(), nil, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responses []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		require.NoError(t, err)

		require.Len(t, responses, 1)
		// Verify price is nil
		assert.Nil(t, responses[0].PriceCents)
		assert.Equal(t, "Free Consultation", responses[0].ServiceType)
	})

	t.Run("should handle empty status filter parameter", func(t *testing.T) {
		ta.ClearAvailabilityTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupUserWithRole(t, ctx, identity.Standard, authCtx)
		_, roomID := setupTestEnvironment(t, ctx)
		partnerID := uuid.New()

		createAvailability(t, ctx, roomID, partnerID, 24*time.Hour, 1*time.Hour, "Test", intPtr(15000), domain.AvailabilityStatusAvailable)

		// Request with empty status parameter (should be ignored)
		queryParams := map[string]string{
			"status": "",
		}
		req := ta.NewGetPartnerAvailabilitiesRequest(t, ctx, testServerURL, partnerID.String(), queryParams, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responses []domain.AvailabilityResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		require.NoError(t, err)

		// Should return all availabilities (empty status ignored)
		assert.Len(t, responses, 1)
	})
}
