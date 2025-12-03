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

// make test-func TEST_NAME=TestMarkNoShow TEST_PATH=test/integration/booking/booking/mark_no_show_test.go

func TestMarkNoShow(t *testing.T) {
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
		require.NoError(t, err)

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

	t.Run("should successfully mark booking as no-show", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup booking
		bookingID, _, partnerToken := setupBooking(t)

		// Mark as no-show (as partner)
		req, err := http.NewRequest("POST", testServerURL+"/bookings/"+bookingID.String()+"/no-show", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: partnerToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		var noShowBooking domain.BookingResponse
		err = json.NewDecoder(resp.Body).Decode(&noShowBooking)
		require.NoError(t, err)

		// Assert booking is marked as no-show
		assert.Equal(t, domain.BookingStatusNoShow, noShowBooking.Status)

		// Verify in database
		dbBooking, err := bookingRepo.GetByID(ctx, bookingID)
		require.NoError(t, err)
		assert.Equal(t, domain.BookingStatusNoShow, dbBooking.Status)
	})

	t.Run("should return 404 for non-existent booking", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup partner user
		partnerToken := tu.SetupPartnerUser(t, ctx, authCtx)

		// Try to mark non-existent booking as no-show
		nonExistentID := uuid.New()
		req, err := http.NewRequest("POST", testServerURL+"/bookings/"+nonExistentID.String()+"/no-show", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: partnerToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert 404
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 when marking cancelled booking as no-show", func(t *testing.T) {
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

		// Try to mark the cancelled booking as no-show
		req, err := http.NewRequest("POST", testServerURL+"/bookings/"+bookingID.String()+"/no-show", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: partnerToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 when marking completed booking as no-show", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup and complete a booking
		bookingID, _, partnerToken := setupBooking(t)

		// Complete the booking first
		reqComplete, _ := http.NewRequest("POST", testServerURL+"/bookings/"+bookingID.String()+"/complete", nil)
		reqComplete.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: partnerToken})

		respComplete, err := client.Do(reqComplete)
		require.NoError(t, err)
		respComplete.Body.Close()
		require.Equal(t, http.StatusOK, respComplete.StatusCode)

		// Try to mark the completed booking as no-show
		req, err := http.NewRequest("POST", testServerURL+"/bookings/"+bookingID.String()+"/no-show", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: partnerToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
