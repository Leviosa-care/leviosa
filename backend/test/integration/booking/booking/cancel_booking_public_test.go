package booking

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

// make test-func TEST_NAME=TestCancelBookingPublic TEST_PATH=test/integration/booking/booking/cancel_booking_public_test.go

// setupGuestBookingForCancel creates a guest booking with token and returns
// the booking ID, token, and the product ID used.
func setupGuestBookingForCancel(t *testing.T, slotStartTime time.Time) (bookingID uuid.UUID, token string, productID uuid.UUID) {
	t.Helper()
	ctx := context.Background()

	tbooking.ClearBookingsTable(t, ctx, testPool)
	tavail.ClearAvailabilityTable(t, ctx, testPool)
	tr.ClearRoomsTable(t, ctx, testPool)
	tb.ClearBuildingsTable(t, ctx, testPool)
	tcatalog.ClearProductsTable(t, ctx, testPool)

	categoryID := tcatalog.CreateTestCategory(t, ctx, testPool)
	products := tcatalog.CreateDefaultTestProducts(t, ctx, testPool, categoryID)

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

	startTime := slotStartTime.Truncate(time.Hour)
	endTime := startTime.Add(2 * time.Hour)

	_, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@cancel-public-test.com", room.ID, testPool, authCtx.Redis, crypto)

	availability := tavail.NewTestAvailability(t, partnerID, room.ID, startTime, endTime, 5000)
	availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
	require.NoError(t, err)
	tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

	pid := products[0].ID
	slotStart := startTime.Add(30 * time.Minute).Truncate(10 * time.Minute)

	reqBody := map[string]interface{}{
		"availability_id":  availability.ID.String(),
		"product_id":       pid.String(),
		"slot_start_time":  slotStart.Format(time.RFC3339),
		"guest_first_name": "Bob",
		"guest_last_name":  "Durand",
		"guest_email":      "bob@example.com",
		"guest_phone":      "+33698765432",
	}
	bodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodPost, testServerURL+"/bookings", bytes.NewReader(bodyBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode, "guest booking creation must succeed")

	var created domain.BookingResponse
	err = json.NewDecoder(resp.Body).Decode(&created)
	require.NoError(t, err)
	require.NotEmpty(t, created.Token, "booking token must be present")

	return created.ID, created.Token, pid
}

// setupGuestBookingWithCancellationWindow creates a product with a specific
// cancellation window and a booking whose slot is inside that window.
func setupGuestBookingWithCancellationWindow(t *testing.T, cancellationHours int) (bookingID uuid.UUID, token string) {
	t.Helper()
	ctx := context.Background()

	tbooking.ClearBookingsTable(t, ctx, testPool)
	tavail.ClearAvailabilityTable(t, ctx, testPool)
	tr.ClearRoomsTable(t, ctx, testPool)
	tb.ClearBuildingsTable(t, ctx, testPool)
	tcatalog.ClearProductsTable(t, ctx, testPool)

	categoryID := tcatalog.CreateTestCategory(t, ctx, testPool)

	// Create a product with the specified cancellation window
	productID := uuid.New()
	now := time.Now()
	_, err := testPool.Exec(ctx, `
		INSERT INTO catalog.products (id, category_id, name, description, duration, buffer_time, status, availability, cancellation_hours, stripe_product_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		productID, categoryID, "Windowed Massage", "Test", 60, 15, "published", "in-person", cancellationHours,
		"stripe_prod_test_"+productID.String(), now, now,
	)
	require.NoError(t, err)

	// Insert a price
	priceID := uuid.New()
	_, err = testPool.Exec(ctx, `
		INSERT INTO catalog.prices (id, product_id, amount, currency, interval, is_active, stripe_price_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		priceID, productID, 8000, "EUR", "one_time", true, "price_test_windowed", now, now,
	)
	require.NoError(t, err)

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

	// Slot starts soon (inside the cancellation window)
	slotStart := time.Now().Add(1 * time.Hour).Truncate(10 * time.Minute)
	slotEnd := slotStart.Add(2 * time.Hour)

	_, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@cancel-window-test.com", room.ID, testPool, authCtx.Redis, crypto)

	availability := tavail.NewTestAvailability(t, partnerID, room.ID, slotStart, slotEnd, 5000)
	availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
	require.NoError(t, err)
	tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

	reqBody := map[string]interface{}{
		"availability_id":  availability.ID.String(),
		"product_id":       productID.String(),
		"slot_start_time":  slotStart.Format(time.RFC3339),
		"guest_first_name": "Claire",
		"guest_last_name":  "Dupont",
		"guest_email":      "claire@example.com",
		"guest_phone":      "+33611112222",
	}
	bodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodPost, testServerURL+"/bookings", bytes.NewReader(bodyBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var created domain.BookingResponse
	err = json.NewDecoder(resp.Body).Decode(&created)
	require.NoError(t, err)
	require.NotEmpty(t, created.Token)

	return created.ID, created.Token
}

func TestCancelBookingPublic(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("successful token cancel returns updated booking", func(t *testing.T) {
		defer tu.ClearAuthData(t, context.Background(), authCtx)

		// Slot well in the future (default product has 24h cancellation window)
		slotStart := time.Now().Add(48 * time.Hour)
		bookingID, token, _ := setupGuestBookingForCancel(t, slotStart)

		cancelBody := domain.CancelBookingRequest{Reason: "Change of plans"}
		bodyBytes, err := json.Marshal(cancelBody)
		require.NoError(t, err)

		url := fmt.Sprintf("%s/bookings/%s/cancel-public?token=%s", testServerURL, bookingID, token)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result domain.PublicBookingLookupResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Equal(t, bookingID, result.ID)
		assert.Equal(t, domain.BookingStatusCancelled, result.Status)

		// Verify in database
		dbBooking, err := bookingRepo.GetByID(context.Background(), bookingID)
		require.NoError(t, err)
		assert.Equal(t, domain.BookingStatusCancelled, dbBooking.Status)

		// Response must not leak guest contact fields
		raw, _ := json.Marshal(result)
		assert.NotContains(t, string(raw), "bob@example.com")
		assert.NotContains(t, string(raw), "+33698765432")
	})

	t.Run("mismatched token returns 401", func(t *testing.T) {
		defer tu.ClearAuthData(t, context.Background(), authCtx)

		slotStart := time.Now().Add(48 * time.Hour)
		bookingID, _, _ := setupGuestBookingForCancel(t, slotStart)

		// Generate a token for a different booking
		otherID := uuid.New()
		wrongToken := domain.GenerateBookingToken(otherID, time.Now().Add(48*time.Hour), testTokenSecret)

		cancelBody := domain.CancelBookingRequest{Reason: "Trying to cancel someone else's booking"}
		bodyBytes, err := json.Marshal(cancelBody)
		require.NoError(t, err)

		url := fmt.Sprintf("%s/bookings/%s/cancel-public?token=%s", testServerURL, bookingID, wrongToken)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("expired token returns 401", func(t *testing.T) {
		defer tu.ClearAuthData(t, context.Background(), authCtx)

		slotStart := time.Now().Add(48 * time.Hour)
		bookingID, _, _ := setupGuestBookingForCancel(t, slotStart)

		// Generate an expired token (slot end time is well in the past)
		pastSlotEnd := time.Now().Add(-domain.BookingTokenGracePeriod - time.Hour)
		expiredToken := domain.GenerateBookingToken(bookingID, pastSlotEnd, testTokenSecret)

		cancelBody := domain.CancelBookingRequest{Reason: "Late cancel attempt"}
		bodyBytes, err := json.Marshal(cancelBody)
		require.NoError(t, err)

		url := fmt.Sprintf("%s/bookings/%s/cancel-public?token=%s", testServerURL, bookingID, expiredToken)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("cancel inside cancellation window returns 422", func(t *testing.T) {
		defer tu.ClearAuthData(t, context.Background(), authCtx)

		// Product requires 48 hours notice, booking is 1 hour away
		bookingID, token := setupGuestBookingWithCancellationWindow(t, 48)

		cancelBody := domain.CancelBookingRequest{Reason: "Too late to cancel"}
		bodyBytes, err := json.Marshal(cancelBody)
		require.NoError(t, err)

		url := fmt.Sprintf("%s/bookings/%s/cancel-public?token=%s", testServerURL, bookingID, token)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("missing token returns 400", func(t *testing.T) {
		cancelBody := domain.CancelBookingRequest{Reason: "No token"}
		bodyBytes, err := json.Marshal(cancelBody)
		require.NoError(t, err)

		url := fmt.Sprintf("%s/bookings/%s/cancel-public", testServerURL, uuid.New())
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("tampered token returns 401", func(t *testing.T) {
		cancelBody := domain.CancelBookingRequest{Reason: "Tampered"}
		bodyBytes, err := json.Marshal(cancelBody)
		require.NoError(t, err)

		url := fmt.Sprintf("%s/bookings/%s/cancel-public?token=not-a-real-token", testServerURL, uuid.New())
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
