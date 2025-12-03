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

// make test-func TEST_NAME=TestCancelBooking TEST_PATH=test/integration/booking/booking/cancel_booking_test.go

func TestCancelBooking(t *testing.T) {
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

	t.Run("should successfully cancel a booking", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup booking
		bookingID, clientToken, _ := setupBooking(t)

		// Cancel the booking
		cancelRequest := domain.CancelBookingRequest{
			Reason: "Client needs to reschedule",
		}
		cancelBytes, err := json.Marshal(cancelRequest)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", testServerURL+"/bookings/"+bookingID.String()+"/cancel", bytes.NewReader(cancelBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: clientToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		var cancelledBooking domain.BookingResponse
		err = json.NewDecoder(resp.Body).Decode(&cancelledBooking)
		require.NoError(t, err)

		// Assert booking is cancelled
		assert.Equal(t, domain.BookingStatusCancelled, cancelledBooking.Status)
		assert.NotNil(t, cancelledBooking.CancelledAt)
		assert.NotNil(t, cancelledBooking.CancellationReason)
		assert.Equal(t, "Client needs to reschedule", *cancelledBooking.CancellationReason)

		// Verify in database
		dbBooking, err := bookingRepo.GetByID(ctx, bookingID)
		require.NoError(t, err)
		assert.Equal(t, domain.BookingStatusCancelled, dbBooking.Status)
		assert.NotNil(t, dbBooking.CancelledAt)
	})

	t.Run("should return 404 for non-existent booking", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup user
		userToken := tu.SetupStandardUser(t, ctx, authCtx)

		// Try to cancel non-existent booking
		nonExistentID := uuid.New()
		cancelRequest := domain.CancelBookingRequest{
			Reason: "Some reason",
		}
		cancelBytes, err := json.Marshal(cancelRequest)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", testServerURL+"/bookings/"+nonExistentID.String()+"/cancel", bytes.NewReader(cancelBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: userToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert 404
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 when cancelling already cancelled booking", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup and cancel booking
		bookingID, clientToken, _ := setupBooking(t)

		cancelRequest := domain.CancelBookingRequest{
			Reason: "First cancellation",
		}
		cancelBytes, err := json.Marshal(cancelRequest)
		require.NoError(t, err)

		// Cancel once
		req1, _ := http.NewRequest("POST", testServerURL+"/bookings/"+bookingID.String()+"/cancel", bytes.NewReader(cancelBytes))
		req1.Header.Set("Content-Type", "application/json")
		req1.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: clientToken})

		resp1, err := client.Do(req1)
		require.NoError(t, err)
		resp1.Body.Close()
		require.Equal(t, http.StatusOK, resp1.StatusCode)

		// Try to cancel again
		req2, err := http.NewRequest("POST", testServerURL+"/bookings/"+bookingID.String()+"/cancel", bytes.NewReader(cancelBytes))
		require.NoError(t, err)
		req2.Header.Set("Content-Type", "application/json")
		req2.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: clientToken})

		resp2, err := client.Do(req2)
		require.NoError(t, err)
		defer resp2.Body.Close()

		// Assert 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp2.StatusCode)
	})

	t.Run("should return 415 for missing content-type header", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		userToken := tu.SetupStandardUser(t, ctx, authCtx)

		req, err := http.NewRequest("POST", testServerURL+"/bookings/"+uuid.New().String()+"/cancel", bytes.NewReader([]byte(`{"reason":"test"}`)))
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: userToken})
		// Intentionally not setting Content-Type

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert 415 Unsupported Media Type
		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
	})
}
