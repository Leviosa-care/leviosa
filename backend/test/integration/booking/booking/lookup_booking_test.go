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

// make test-func TEST_NAME=TestLookupBooking TEST_PATH=test/integration/booking/booking/lookup_booking_test.go

const (
	guestEmail = "alice@example.com"
	guestPhone = "+33612345678"
)

// setupGuestBookingForLookup creates all infra and returns a guest booking with token.
func setupGuestBookingForLookup(t *testing.T) (bookingID uuid.UUID, token string) {
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

	startTime := time.Now().Add(24 * time.Hour).Truncate(time.Hour)
	endTime := startTime.Add(2 * time.Hour)

	_, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@lookup-test.com", room.ID, testPool, authCtx.Redis, crypto)

	availability := tavail.NewTestAvailability(t, partnerID, room.ID, startTime, endTime, 5000)
	availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
	require.NoError(t, err)
	tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

	productID := products[0].ID
	slotStartTime := startTime.Add(30 * time.Minute).Truncate(10 * time.Minute)

	reqBody := map[string]interface{}{
		"availability_id":  availability.ID.String(),
		"product_id":       productID.String(),
		"slot_start_time":  slotStartTime.Format(time.RFC3339),
		"guest_first_name": "Alice",
		"guest_last_name":  "Martin",
		"guest_email":      guestEmail,
		"guest_phone":      guestPhone,
	}
	bodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, testServerURL+"/bookings", bytes.NewReader(bodyBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode, "guest booking creation must succeed")

	var created domain.BookingResponse
	err = json.NewDecoder(resp.Body).Decode(&created)
	require.NoError(t, err)
	require.NotEmpty(t, created.Token, "booking token must be present when WithTokenSecret is configured")

	return created.ID, created.Token
}

func TestLookupBooking(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("token path - valid token returns booking details", func(t *testing.T) {
		defer tu.ClearAuthData(t, context.Background(), authCtx)

		bookingID, token := setupGuestBookingForLookup(t)

		req, err := http.NewRequest(http.MethodGet,
			testServerURL+"/bookings/lookup?token="+token, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result domain.PublicBookingLookupResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Equal(t, bookingID, result.ID)
		assert.Equal(t, domain.BookingStatusConfirmed, result.Status)
		// Guest contact fields must never appear in the public response
		raw, _ := json.Marshal(result)
		assert.NotContains(t, string(raw), guestEmail)
		assert.NotContains(t, string(raw), guestPhone)
	})

	t.Run("token path - expired token returns 401", func(t *testing.T) {
		defer tu.ClearAuthData(t, context.Background(), authCtx)

		// Generate a token whose slot end time is in the past (already expired)
		pastSlotEnd := time.Now().Add(-domain.BookingTokenGracePeriod - time.Hour)
		expiredToken := domain.GenerateBookingToken(uuid.New(), pastSlotEnd, testTokenSecret)

		req, err := http.NewRequest(http.MethodGet,
			testServerURL+"/bookings/lookup?token="+expiredToken, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("token path - tampered token returns 401", func(t *testing.T) {
		defer tu.ClearAuthData(t, context.Background(), authCtx)

		req, err := http.NewRequest(http.MethodGet,
			testServerURL+"/bookings/lookup?token=not-a-valid-token", nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("manual fallback - valid ref + email returns booking details", func(t *testing.T) {
		defer tu.ClearAuthData(t, context.Background(), authCtx)

		bookingID, _ := setupGuestBookingForLookup(t)

		req, err := http.NewRequest(http.MethodGet,
			fmt.Sprintf("%s/bookings/lookup?ref=%s&email=%s", testServerURL, bookingID, guestEmail), nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result domain.PublicBookingLookupResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Equal(t, bookingID, result.ID)
	})

	t.Run("manual fallback - valid ref + phone returns booking details", func(t *testing.T) {
		defer tu.ClearAuthData(t, context.Background(), authCtx)

		bookingID, _ := setupGuestBookingForLookup(t)

		req, err := http.NewRequest(http.MethodGet,
			fmt.Sprintf("%s/bookings/lookup?ref=%s&phone=%s", testServerURL, bookingID, guestPhone), nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result domain.PublicBookingLookupResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Equal(t, bookingID, result.ID)
	})

	t.Run("manual fallback - wrong booking reference returns 401", func(t *testing.T) {
		defer tu.ClearAuthData(t, context.Background(), authCtx)

		wrongRef := uuid.New()
		req, err := http.NewRequest(http.MethodGet,
			fmt.Sprintf("%s/bookings/lookup?ref=%s&email=%s", testServerURL, wrongRef, guestEmail), nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("manual fallback - mismatched contact returns 401", func(t *testing.T) {
		defer tu.ClearAuthData(t, context.Background(), authCtx)

		bookingID, _ := setupGuestBookingForLookup(t)

		req, err := http.NewRequest(http.MethodGet,
			fmt.Sprintf("%s/bookings/lookup?ref=%s&email=wrong@example.com", testServerURL, bookingID), nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("missing params returns 400", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, testServerURL+"/bookings/lookup", nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
