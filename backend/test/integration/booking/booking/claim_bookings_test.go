package booking

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/services"
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

// make test-func TEST_NAME=TestClaimBookings TEST_PATH=test/integration/booking/booking/claim_bookings_test.go

func TestClaimBookings(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should link guest bookings to client when email matches", func(t *testing.T) {
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

		// Setup partner with allocation
		_, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@example.com", room.ID, testPool, authCtx.Redis, crypto)

		// Create availability
		availability := tavail.NewTestAvailability(t, room.ID, partnerID,
			time.Now().Add(24*time.Hour),
			time.Now().Add(25*time.Hour),
			5000,
		)
		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

		// Create guest bookings via the service (with real encryption)
		guestEmail := "guest@example.com"
		booking1, err := service.CreateBooking(ctx, availability.ID, nil, products[0].ID, time.Now().Add(24*time.Hour).Truncate(10*time.Minute), "", "Jean", "Dupont", guestEmail, "0612345678")
		require.NoError(t, err)
		assert.True(t, booking1.IsGuestBooking())
		assert.Nil(t, booking1.ClientID)

		// Create a second guest booking with the same email
		availability2 := tavail.NewTestAvailability(t, room.ID, partnerID,
			time.Now().Add(48*time.Hour),
			time.Now().Add(49*time.Hour),
			5000,
		)
		availabilityEncx2, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability2)
		require.NoError(t, err)
		tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx2, testPool)

		booking2, err := service.CreateBooking(ctx, availability2.ID, nil, products[0].ID, time.Now().Add(48*time.Hour).Truncate(10*time.Minute), "", "Jean", "Dupont", guestEmail, "0612345678")
		require.NoError(t, err)
		assert.True(t, booking2.IsGuestBooking())

		// Create a client user
		_, clientID := tsetup.SetupStandardUser(t, ctx, "newclient@example.com", room.ID, testPool, authCtx.Redis, crypto)

		// Call claim endpoint
		claimReq := testClaimRequest{ClientID: clientID.String(), Email: guestEmail}
		body, _ := json.Marshal(claimReq)
		req, err := http.NewRequest("POST", testServerURL+"/bookings/claim", bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(services.ServiceNameHeader, services.AuthUser)
		req.Header.Set(services.ServiceKeyHeader, testClaimServiceKey)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var claimResp testClaimResponse
		err = json.NewDecoder(resp.Body).Decode(&claimResp)
		require.NoError(t, err)
		assert.Equal(t, 2, claimResp.Claimed)

		// Verify bookings now belong to the client
		clientBookings, err := service.GetClientBookings(ctx, clientID, ports.BookingFilter{})
		require.NoError(t, err)
		assert.Len(t, clientBookings, 2)

		bookingIDs := map[uuid.UUID]bool{}
		for _, b := range clientBookings {
			assert.Equal(t, &clientID, b.ClientID)
			bookingIDs[b.ID] = true
		}
		assert.True(t, bookingIDs[booking1.ID])
		assert.True(t, bookingIDs[booking2.ID])
	})

	t.Run("should be a no-op when no guest bookings match the email", func(t *testing.T) {
		// Clean state
		tbooking.ClearBookingsTable(t, ctx, testPool)
		tavail.ClearAvailabilityTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tcatalog.ClearProductsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup minimal infrastructure
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

		_, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner2@example.com", room.ID, testPool, authCtx.Redis, crypto)

		availability := tavail.NewTestAvailability(t, room.ID, partnerID,
			time.Now().Add(24*time.Hour),
			time.Now().Add(25*time.Hour),
			5000,
		)
		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

		// Create a guest booking with a DIFFERENT email
		_, err = service.CreateBooking(ctx, availability.ID, nil, products[0].ID, time.Now().Add(24*time.Hour).Truncate(10*time.Minute), "", "Other", "Person", "other@example.com", "0600000000")
		require.NoError(t, err)

		// Create a client user
		_, clientID := tsetup.SetupStandardUser(t, ctx, "noclaim@example.com", room.ID, testPool, authCtx.Redis, crypto)

		// Call claim with an email that doesn't match any guest booking
		claimReq := testClaimRequest{ClientID: clientID.String(), Email: "nonexistent@example.com"}
		body, _ := json.Marshal(claimReq)
		req, err := http.NewRequest("POST", testServerURL+"/bookings/claim", bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(services.ServiceNameHeader, services.AuthUser)
		req.Header.Set(services.ServiceKeyHeader, testClaimServiceKey)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var claimResp testClaimResponse
		err = json.NewDecoder(resp.Body).Decode(&claimResp)
		require.NoError(t, err)
		assert.Equal(t, 0, claimResp.Claimed)

		// Client should have zero bookings
		clientBookings, err := service.GetClientBookings(ctx, clientID, ports.BookingFilter{})
		require.NoError(t, err)
		assert.Len(t, clientBookings, 0)
	})

	t.Run("should skip already-owned bookings", func(t *testing.T) {
		// Clean state
		tbooking.ClearBookingsTable(t, ctx, testPool)
		tavail.ClearAvailabilityTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tcatalog.ClearProductsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup minimal infrastructure
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

		_, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner3@example.com", room.ID, testPool, authCtx.Redis, crypto)

		// Create two availabilities
		availability := tavail.NewTestAvailability(t, room.ID, partnerID,
			time.Now().Add(24*time.Hour),
			time.Now().Add(25*time.Hour),
			5000,
		)
		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

		availability2 := tavail.NewTestAvailability(t, room.ID, partnerID,
			time.Now().Add(48*time.Hour),
			time.Now().Add(49*time.Hour),
			5000,
		)
		availabilityEncx2, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability2)
		require.NoError(t, err)
		tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx2, testPool)

		// Create two users
		_, existingID := tsetup.SetupStandardUser(t, ctx, "existing@example.com", room.ID, testPool, authCtx.Redis, crypto)
		_, newID := tsetup.SetupStandardUser(t, ctx, "newuser@example.com", room.ID, testPool, authCtx.Redis, crypto)

		// Create a booking already owned by existing user
		existingBooking, err := service.CreateBooking(ctx, availability.ID, &existingID, products[0].ID, time.Now().Add(24*time.Hour).Truncate(10*time.Minute), "", "", "", "", "")
		require.NoError(t, err)
		assert.NotNil(t, existingBooking.ClientID)

		// Create a guest booking with an email
		guestEmail := "guestshared@example.com"
		guestBooking, err := service.CreateBooking(ctx, availability2.ID, nil, products[0].ID, time.Now().Add(48*time.Hour).Truncate(10*time.Minute), "", "Guest", "User", guestEmail, "")
		require.NoError(t, err)
		assert.Nil(t, guestBooking.ClientID)

		// Claim for new user — should only get the guest booking, not the one owned by existingID
		claimReq := testClaimRequest{ClientID: newID.String(), Email: guestEmail}
		body, _ := json.Marshal(claimReq)
		req, err := http.NewRequest("POST", testServerURL+"/bookings/claim", bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(services.ServiceNameHeader, services.AuthUser)
		req.Header.Set(services.ServiceKeyHeader, testClaimServiceKey)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var claimResp testClaimResponse
		err = json.NewDecoder(resp.Body).Decode(&claimResp)
		require.NoError(t, err)
		assert.Equal(t, 1, claimResp.Claimed)

		// Existing user still has their booking
		existingBookings, err := service.GetClientBookings(ctx, existingID, ports.BookingFilter{})
		require.NoError(t, err)
		assert.Len(t, existingBookings, 1)
		assert.Equal(t, existingBooking.ID, existingBookings[0].ID)

		// New user has only the guest booking
		newBookings, err := service.GetClientBookings(ctx, newID, ports.BookingFilter{})
		require.NoError(t, err)
		assert.Len(t, newBookings, 1)
		assert.Equal(t, guestBooking.ID, newBookings[0].ID)
	})

	t.Run("should return 400 for missing client_id", func(t *testing.T) {
		claimReq := testClaimRequest{Email: "test@example.com"}
		body, _ := json.Marshal(claimReq)
		req, err := http.NewRequest("POST", testServerURL+"/bookings/claim", bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(services.ServiceNameHeader, services.AuthUser)
		req.Header.Set(services.ServiceKeyHeader, testClaimServiceKey)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for missing email", func(t *testing.T) {
		claimReq := testClaimRequest{ClientID: uuid.New().String()}
		body, _ := json.Marshal(claimReq)
		req, err := http.NewRequest("POST", testServerURL+"/bookings/claim", bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(services.ServiceNameHeader, services.AuthUser)
		req.Header.Set(services.ServiceKeyHeader, testClaimServiceKey)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 401 for missing service auth", func(t *testing.T) {
		claimReq := testClaimRequest{ClientID: uuid.New().String(), Email: "test@example.com"}
		body, _ := json.Marshal(claimReq)
		req, err := http.NewRequest("POST", testServerURL+"/bookings/claim", bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		// intentionally no service headers

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

// testClaimRequest is the request DTO for the claim endpoint test.
type testClaimRequest struct {
	ClientID string `json:"client_id"`
	Email    string `json:"email"`
}

// testClaimResponse is the response DTO for the claim endpoint test.
type testClaimResponse struct {
	Claimed int `json:"claimed"`
}
