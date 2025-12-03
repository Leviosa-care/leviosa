package booking

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	tsetup "github.com/Leviosa-care/leviosa/backend/test/helpers/booking"
	tavail "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	tbooking "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/booking"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"
	tcatalog "github.com/Leviosa-care/leviosa/backend/test/helpers/catalog"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCompleteBooking TEST_PATH=test/integration/booking/booking/complete_booking_test.go

func TestCompleteBooking(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Helper function to create a complete booking setup
	setupBooking := func(t *testing.T) (uuid.UUID, string, string) {
		t.Helper()

		// Clean state
		tbooking.ClearBookingsTable(t, ctx, testPool)
		tavail.ClearAvailabilityTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tcatalog.ClearProductsTable(t, ctx, testPool)

		// Setup test category and products
		categoryID := tcatalog.CreateTestCategory(t, ctx, testPool)
		products := tcatalog.CreateDefaultTestProducts(t, ctx, testPool, categoryID)

		// Setup test building and room
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

		// Setup authenticated partner and client
		partnerToken, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@test.com", room.ID, testPool, authCtx.Redis, crypto)
		clientToken, clientID := tsetup.SetupStandardUser(t, ctx, "client@test.com", room.ID, testPool, authCtx.Redis, crypto)

		// Create availability
		startTime := time.Now().Add(24 * time.Hour).Truncate(time.Hour)
		endTime := startTime.Add(2 * time.Hour)
		availability := tavail.NewTestAvailability(t, partnerID, room.ID, startTime, endTime, 5000)
		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

		// Create booking via HTTP using first test product (60-minute massage)
		productID := products[0].ID
		slotStartTime := startTime.Add(30 * time.Minute).Truncate(10 * time.Minute)
		requestBody := map[string]interface{}{
			"availability_id": availability.ID.String(),
			"client_id":       clientID,
			"product_id":      productID.String(),
			"slot_start_time": slotStartTime.Format(time.RFC3339),
		}
		bodyBytes, err := json.Marshal(requestBody)
		require.NoError(t, err)

		reqCreate, err := http.NewRequest("POST", testServerURL+"/bookings", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		reqCreate.Header.Set("Content-Type", "application/json")
		reqCreate.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: clientToken})

		respCreate, err := client.Do(reqCreate)
		require.NoError(t, err)
		defer respCreate.Body.Close()
		require.Equal(t, http.StatusCreated, respCreate.StatusCode)

		var createdBooking domain.BookingResponse
		err = json.NewDecoder(respCreate.Body).Decode(&createdBooking)
		require.NoError(t, err)

		return createdBooking.ID, clientToken, partnerToken
	}

	t.Run("should successfully complete a booking", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup booking
		bookingID, _, partnerToken := setupBooking(t)

		// Complete the booking (as partner)
		req, err := http.NewRequest("POST", testServerURL+"/bookings/"+bookingID.String()+"/complete", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: partnerToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		var completedBooking domain.BookingResponse
		err = json.NewDecoder(resp.Body).Decode(&completedBooking)
		require.NoError(t, err)

		// Assert booking is completed
		assert.Equal(t, domain.BookingStatusCompleted, completedBooking.Status)
		assert.NotNil(t, completedBooking.CompletedAt)

		// Verify in database
		dbBooking, err := bookingRepo.GetByID(ctx, bookingID)
		require.NoError(t, err)
		assert.Equal(t, domain.BookingStatusCompleted, dbBooking.Status)
		assert.NotNil(t, dbBooking.CompletedAt)
	})

	t.Run("should return 404 for non-existent booking", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup partner user
		partnerToken := tu.SetupPartnerUser(t, ctx, authCtx)

		// Try to complete non-existent booking
		nonExistentID := uuid.New()
		req, err := http.NewRequest("POST", testServerURL+"/bookings/"+nonExistentID.String()+"/complete", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: partnerToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert 404
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 when completing cancelled booking", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup and cancel a booking
		bookingID, clientToken, partnerToken := setupBooking(t)

		// Cancel the booking first
		cancelRequest := domain.CancelBookingRequest{
			Reason: "Cancelled by client",
		}
		cancelBytes, err := json.Marshal(cancelRequest)
		require.NoError(t, err)

		reqCancel, _ := http.NewRequest("POST", testServerURL+"/bookings/"+bookingID.String()+"/cancel", bytes.NewReader(cancelBytes))
		reqCancel.Header.Set("Content-Type", "application/json")
		reqCancel.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: clientToken})

		respCancel, err := client.Do(reqCancel)
		require.NoError(t, err)
		respCancel.Body.Close()
		require.Equal(t, http.StatusOK, respCancel.StatusCode)

		// Try to complete the cancelled booking
		req, err := http.NewRequest("POST", testServerURL+"/bookings/"+bookingID.String()+"/complete", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: partnerToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 when completing already completed booking", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup and complete booking
		bookingID, _, partnerToken := setupBooking(t)

		// Complete the booking first
		reqComplete1, _ := http.NewRequest("POST", testServerURL+"/bookings/"+bookingID.String()+"/complete", nil)
		reqComplete1.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: partnerToken})

		respComplete1, err := client.Do(reqComplete1)
		require.NoError(t, err)
		respComplete1.Body.Close()
		require.Equal(t, http.StatusOK, respComplete1.StatusCode)

		// Try to complete the already completed booking
		reqComplete2, err := http.NewRequest("POST", testServerURL+"/bookings/"+bookingID.String()+"/complete", nil)
		require.NoError(t, err)
		reqComplete2.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: partnerToken})

		respComplete2, err := client.Do(reqComplete2)
		require.NoError(t, err)
		defer respComplete2.Body.Close()

		// Assert 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, respComplete2.StatusCode)
	})
}
