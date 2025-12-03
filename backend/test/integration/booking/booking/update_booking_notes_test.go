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
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"
	tavail "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
	tbooking "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/booking"
	tcatalog "github.com/Leviosa-care/leviosa/backend/test/helpers/catalog"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateBookingNotes(t *testing.T) {
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

	t.Run("should successfully update client notes", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup booking
		bookingID, clientToken, _ := setupBooking(t)

		// Update client notes
		clientNotes := "Updated client notes"
		updateRequest := domain.UpdateBookingNotesRequest{
			ClientNotes: clientNotes,
		}
		updateBytes, err := json.Marshal(updateRequest)
		require.NoError(t, err)

		req, err := http.NewRequest("PUT", testServerURL+"/bookings/"+bookingID.String()+"/notes", bytes.NewReader(updateBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: clientToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		var updatedBooking domain.BookingResponse
		err = json.NewDecoder(resp.Body).Decode(&updatedBooking)
		require.NoError(t, err)

		// Assert updated notes
		assert.Equal(t, clientNotes, updatedBooking.ClientNotes)

		// Verify in database
		dbBookingEncx, err := bookingRepo.GetByID(ctx, bookingID)
		require.NoError(t, err)
		dbBooking, err := domain.DecryptBookingEncx(ctx, crypto, dbBookingEncx)
		require.NoError(t, err)
		assert.Equal(t, clientNotes, dbBooking.ClientNotes)
	})

	t.Run("should successfully update partner notes", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup booking
		bookingID, _, partnerToken := setupBooking(t)

		// Update partner notes
		partnerNotes := "Updated partner notes"
		updateRequest := domain.UpdateBookingNotesRequest{
			PartnerNotes: partnerNotes,
		}
		updateBytes, err := json.Marshal(updateRequest)
		require.NoError(t, err)

		req, err := http.NewRequest("PUT", testServerURL+"/bookings/"+bookingID.String()+"/notes", bytes.NewReader(updateBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: partnerToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		var updatedBooking domain.BookingResponse
		err = json.NewDecoder(resp.Body).Decode(&updatedBooking)
		require.NoError(t, err)

		// Assert updated notes
		assert.Equal(t, partnerNotes, updatedBooking.PartnerNotes)
	})

	t.Run("should return 404 for non-existent booking", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup user
		userToken := tu.SetupStandardUser(t, ctx, authCtx)

		// Try to update non-existent booking
		nonExistentID := uuid.New()
		updateRequest := domain.UpdateBookingNotesRequest{
			ClientNotes: "Some notes",
		}
		updateBytes, err := json.Marshal(updateRequest)
		require.NoError(t, err)

		req, err := http.NewRequest("PUT", testServerURL+"/bookings/"+nonExistentID.String()+"/notes", bytes.NewReader(updateBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: userToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert 404
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 415 for missing content-type header", func(t *testing.T) {
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup user
		userToken := tu.SetupStandardUser(t, ctx, authCtx)

		req, err := http.NewRequest("PUT", testServerURL+"/bookings/"+uuid.New().String()+"/notes", bytes.NewReader([]byte("{}")))
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
