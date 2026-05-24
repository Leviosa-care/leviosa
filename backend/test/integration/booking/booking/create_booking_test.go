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

// make test-func TEST_NAME=TestCreateBooking TEST_PATH=test/integration/booking/booking/create_booking_test.go

func TestCreateBooking(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully create a booking", func(t *testing.T) {
		// Clean state
		tbooking.ClearBookingsTable(t, ctx, testPool)
		tavail.ClearAvailabilityTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tcatalog.ClearProductsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

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

		// Setup authenticated users
		_, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@example.com", room.ID, testPool, authCtx.Redis, crypto)
		clientToken, clientID := tsetup.SetupStandardUser(t, ctx, "client@example.com", room.ID, testPool, authCtx.Redis, crypto)

		// Create test availability
		startTime := time.Now().Add(24 * time.Hour).Truncate(time.Hour)
		endTime := startTime.Add(2 * time.Hour)
		availability := tavail.NewTestAvailability(t, partnerID, room.ID, startTime, endTime, 0)
		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

		// Prepare booking request with product ID and slot time using first test product (60-minute massage)
		productID := products[0].ID
		slotStartTime := startTime.Add(30 * time.Minute).Truncate(10 * time.Minute)
		requestBody := map[string]interface{}{
			"availability_id": availability.ID.String(),
			"client_id":       clientID,
			"product_id":      productID.String(),
			"slot_start_time": slotStartTime.Format(time.RFC3339),
			"client_notes":    "Looking forward to this session",
		}
		bodyBytes, err := json.Marshal(requestBody)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", testServerURL+"/bookings", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: clientToken})

		// Execute request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Parse response
		var booking domain.BookingResponse
		err = json.NewDecoder(resp.Body).Decode(&booking)
		require.NoError(t, err)

		// Assert booking details
		assert.NotEqual(t, uuid.Nil, booking.ID)
		assert.Equal(t, availability.ID, booking.AvailabilityID)
		assert.Equal(t, &clientID, booking.ClientID)
		assert.Equal(t, partnerID, booking.PartnerID)
		assert.Equal(t, room.ID, booking.RoomID)
		assert.Equal(t, productID, booking.ProductID)
		assert.Equal(t, slotStartTime.Unix(), booking.SlotStartTime.Unix())
		assert.Equal(t, 0, booking.TotalPriceCents) // Price is 0 (pricing not yet integrated - see TODO in create_booking.go)
		assert.Equal(t, "EUR", booking.Currency)
		assert.Equal(t, domain.BookingStatusConfirmed, booking.Status)
		assert.Equal(t, domain.PaymentStatusPending, booking.PaymentStatus)

		// Verify booking exists in database
		savedBooking, err := bookingRepo.GetByID(ctx, booking.ID)
		require.NoError(t, err)
		assert.Equal(t, booking.ID, savedBooking.ID)
	})

	t.Run("should return 400 for invalid availability", func(t *testing.T) {
		// Clean state
		tbooking.ClearBookingsTable(t, ctx, testPool)
		tavail.ClearAvailabilityTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup standard user for authentication
		userToken := tu.SetupStandardUser(t, ctx, authCtx)

		// Try to create booking with non-existent availability
		nonExistentAvailID := uuid.New()
		productID := uuid.New()
		slotStartTime := time.Now().Add(24 * time.Hour).Truncate(10 * time.Minute)
		requestBody := map[string]interface{}{
			"availability_id": nonExistentAvailID.String(),
			"client_id":       uuid.New(),
			"product_id":      productID.String(),
			"slot_start_time": slotStartTime.Format(time.RFC3339),
		}
		bodyBytes, err := json.Marshal(requestBody)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", testServerURL+"/bookings", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: userToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 415 for missing content-type header", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		userToken := tu.SetupStandardUser(t, ctx, authCtx)

		req, err := http.NewRequest("POST", testServerURL+"/bookings", bytes.NewReader([]byte("{}")))
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
