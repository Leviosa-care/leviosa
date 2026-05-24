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

// make test-func TEST_NAME=TestGetClientBookings TEST_PATH=test/integration/booking/booking/get_client_bookings_test.go

func TestGetClientBookings(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get all bookings for a client", func(t *testing.T) {
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
		_, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner@test.com", room.ID, testPool, authCtx.Redis, crypto)
		clientToken, clientID := tsetup.SetupStandardUser(t, ctx, "client@test.com", room.ID, testPool, authCtx.Redis, crypto)

		// Create multiple bookings for the same client
		numBookings := 3
		for i := 0; i < numBookings; i++ {
			startTime := time.Now().Add(time.Duration(24+i) * time.Hour).Truncate(time.Hour)
			endTime := startTime.Add(2 * time.Hour)
			availability := tavail.NewTestAvailability(t, partnerID, room.ID, startTime, endTime, 5000+i*1000)
			availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
			require.NoError(t, err)
			tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)
			require.NoError(t, err)

			// Use different products for each booking to vary the tests
			productID := products[i%len(products)].ID
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
			respCreate.Body.Close()
			require.Equal(t, http.StatusCreated, respCreate.StatusCode)
		}

		// Get client bookings
		req, err := http.NewRequest("GET", testServerURL+"/clients/"+clientID.String()+"/bookings", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: clientToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		var bookings []domain.BookingResponse
		err = json.NewDecoder(resp.Body).Decode(&bookings)
		require.NoError(t, err)

		// Assert bookings count
		assert.Len(t, bookings, numBookings)

		// Verify all bookings belong to the client
		for _, booking := range bookings {
			assert.Equal(t, &clientID, booking.ClientID)
		}
	})

	t.Run("should return empty array for client with no bookings", func(t *testing.T) {
		// Clean state
		tbooking.ClearBookingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup user with no bookings
		userToken := tu.SetupStandardUser(t, ctx, authCtx)

		// Get a user ID for the test
		userID := uuid.New()

		// Get client bookings
		req, err := http.NewRequest("GET", testServerURL+"/clients/"+userID.String()+"/bookings", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: userToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		var bookings []domain.BookingResponse
		err = json.NewDecoder(resp.Body).Decode(&bookings)
		require.NoError(t, err)

		// Assert empty array
		assert.Empty(t, bookings)
	})

	t.Run("should return 400 for invalid client ID", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup user for authentication
		userToken := tu.SetupStandardUser(t, ctx, authCtx)

		// Try to get bookings with invalid UUID
		req, err := http.NewRequest("GET", testServerURL+"/clients/invalid-uuid/bookings", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: userToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert 400
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
