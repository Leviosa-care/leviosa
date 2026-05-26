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

func TestBookingConfirmationNotification(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should send booking confirmation notification after booking is created", func(t *testing.T) {
		// Clean state
		tbooking.ClearBookingsTable(t, ctx, testPool)
		tavail.ClearAvailabilityTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tcatalog.ClearProductsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Reset spy
		notificationSpy.BookingConfirmations = nil
		notificationSpy.PaymentConfirmations = nil

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
		_, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner-notify@example.com", room.ID, testPool, authCtx.Redis, crypto)
		clientToken, clientID := tsetup.SetupStandardUser(t, ctx, "client-notify@example.com", room.ID, testPool, authCtx.Redis, crypto)

		// Create test availability
		startTime := time.Now().Add(24 * time.Hour).Truncate(time.Hour)
		endTime := startTime.Add(2 * time.Hour)
		availability := tavail.NewTestAvailability(t, partnerID, room.ID, startTime, endTime, 0)
		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

		// Prepare booking request
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

		// Parse response to get booking ID
		var booking domain.BookingResponse
		err = json.NewDecoder(resp.Body).Decode(&booking)
		require.NoError(t, err)

		// Assert notification was sent
		require.Len(t, notificationSpy.BookingConfirmations, 1, "expected exactly one booking confirmation notification")

		notification := notificationSpy.BookingConfirmations[0]
		assert.Equal(t, booking.ID, notification.BookingID, "notification should reference the created booking")
		assert.Equal(t, clientID, notification.ClientID, "notification should reference the client")
		assert.Equal(t, partnerID, notification.PartnerID, "notification should reference the partner")
		assert.Equal(t, productID, notification.ProductID, "notification should reference the product")
		assert.Equal(t, room.ID, notification.RoomID, "notification should reference the room")
		assert.Equal(t, slotStartTime.Unix(), notification.SlotStartTime.Unix(), "notification should have correct slot start time")
	})
}

func TestBookingCancellationNotification(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should send cancellation notifications to both client and partner when booking is cancelled", func(t *testing.T) {
		// Clean state
		tbooking.ClearBookingsTable(t, ctx, testPool)
		tavail.ClearAvailabilityTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tcatalog.ClearProductsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Reset spy
		notificationSpy.BookingConfirmations = nil
		notificationSpy.BookingCancellations = nil
		notificationSpy.PaymentConfirmations = nil

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
		_, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner-cancel@example.com", room.ID, testPool, authCtx.Redis, crypto)
		clientToken, clientID := tsetup.SetupStandardUser(t, ctx, "client-cancel@example.com", room.ID, testPool, authCtx.Redis, crypto)

		// Create test availability (far enough in the future to pass cancellation policy)
		startTime := time.Now().Add(72 * time.Hour).Truncate(time.Hour)
		endTime := startTime.Add(2 * time.Hour)
		availability := tavail.NewTestAvailability(t, partnerID, room.ID, startTime, endTime, 0)
		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

		// Create booking
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

		req, err := http.NewRequest("POST", testServerURL+"/bookings", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: clientToken})

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var createdBooking domain.BookingResponse
		err = json.NewDecoder(resp.Body).Decode(&createdBooking)
		require.NoError(t, err)

		// Assert booking confirmation was sent on creation
		require.Len(t, notificationSpy.BookingConfirmations, 1, "expected booking confirmation after creation")

		// Reset spy to isolate cancellation notifications
		notificationSpy.BookingConfirmations = nil
		notificationSpy.BookingCancellations = nil

		// Cancel the booking
		cancelRequest := domain.CancelBookingRequest{
			Reason: "Client needs to reschedule",
		}
		cancelBytes, err := json.Marshal(cancelRequest)
		require.NoError(t, err)

		cancelReq, err := http.NewRequest("POST", testServerURL+"/bookings/"+createdBooking.ID.String()+"/cancel", bytes.NewReader(cancelBytes))
		require.NoError(t, err)
		cancelReq.Header.Set("Content-Type", "application/json")
		cancelReq.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: clientToken})

		cancelResp, err := client.Do(cancelReq)
		require.NoError(t, err)
		defer cancelResp.Body.Close()
		assert.Equal(t, http.StatusOK, cancelResp.StatusCode)

		// Assert cancellation notification was sent
		require.Len(t, notificationSpy.BookingCancellations, 1, "expected exactly one booking cancellation notification call")

		notification := notificationSpy.BookingCancellations[0]
		assert.Equal(t, createdBooking.ID, notification.BookingID, "notification should reference the cancelled booking")
		assert.Equal(t, clientID, notification.ClientID, "notification should reference the client")
		assert.Equal(t, partnerID, notification.PartnerID, "notification should reference the partner")
		assert.Equal(t, productID, notification.ProductID, "notification should reference the product")
		assert.Equal(t, room.ID, notification.RoomID, "notification should reference the room")
		assert.Equal(t, "Client needs to reschedule", notification.CancellationReason, "notification should include the cancellation reason")
		assert.NotNil(t, notification.CancelledAt, "notification should include cancelled_at timestamp")
	})
}

func TestPaymentConfirmationNotification(t *testing.T) {
	ctx := context.Background()
	httpClient := &http.Client{Timeout: 10 * time.Second}

	t.Run("should send payment confirmation notification after payment webhook", func(t *testing.T) {
		// Clean state
		tbooking.ClearBookingsTable(t, ctx, testPool)
		tavail.ClearAvailabilityTable(t, ctx, testPool)
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		tcatalog.ClearProductsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Reset spy
		notificationSpy.BookingConfirmations = nil
		notificationSpy.PaymentConfirmations = nil
		notificationSpy.PaymentFailed = nil

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
		_, partnerID := tsetup.SetupAuthenticatedPartnerWithAllocation(t, ctx, "partner-payment@example.com", room.ID, testPool, authCtx.Redis, crypto)
		clientToken, clientID := tsetup.SetupStandardUser(t, ctx, "client-payment@example.com", room.ID, testPool, authCtx.Redis, crypto)

		// Create test availability
		startTime := time.Now().Add(24 * time.Hour).Truncate(time.Hour)
		endTime := startTime.Add(2 * time.Hour)
		availability := tavail.NewTestAvailability(t, partnerID, room.ID, startTime, endTime, 0)
		availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, crypto, availability)
		require.NoError(t, err)
		tavail.InsertAvailabilityEncx(t, ctx, availabilityEncx, testPool)

		// Prepare booking request
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

		req, err := http.NewRequest("POST", testServerURL+"/bookings", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "leviosa_access_token", Value: clientToken})

		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var booking domain.BookingResponse
		err = json.NewDecoder(resp.Body).Decode(&booking)
		require.NoError(t, err)

		// Simulate payment success via webhook
		paymentIntentID := "pi_test_" + uuid.New().String()[:8]
		webhookEvent := &ports.WebhookEvent{
			ID:              "evt_test_" + uuid.New().String(),
			Type:            ports.WebhookEventPaymentIntentSucceeded,
			PaymentIntentID: paymentIntentID,
			Status:          "succeeded",
			Amount:          booking.TotalPriceCents,
			Currency:        booking.Currency,
			Metadata: map[string]string{
				"booking_id": booking.ID.String(),
			},
		}

		// First set the payment intent ID on the booking via mock
		mockPayment, ok := paymentService.(*MockPaymentService)
		require.True(t, ok, "payment service should be MockPaymentService")
		mockPayment.paymentIntents[paymentIntentID] = &ports.PaymentIntentInfo{
			ID:           paymentIntentID,
			Status:       ports.PaymentIntentStatusSucceeded,
			Amount:       booking.TotalPriceCents,
			Currency:     booking.Currency,
			ClientSecret: "cs_test",
			Metadata: map[string]string{
				"booking_id": booking.ID.String(),
			},
		}

		// Update booking with payment intent ID
		bookingEncx, err := bookingRepo.GetByID(ctx, booking.ID)
		require.NoError(t, err)
		decryptedBooking, err := domain.DecryptBookingEncx(ctx, crypto, bookingEncx)
		require.NoError(t, err)
		decryptedBooking.SetPaymentIntentID(paymentIntentID)
		updatedEncx, err := domain.ProcessBookingEncx(ctx, crypto, decryptedBooking)
		require.NoError(t, err)
		err = bookingRepo.Update(ctx, updatedEncx)
		require.NoError(t, err)

		// Process webhook
		err = service.HandlePaymentWebhook(ctx, webhookEvent)
		require.NoError(t, err)

		// Assert booking confirmation was sent on creation
		require.Len(t, notificationSpy.BookingConfirmations, 1, "expected booking confirmation after creation")

		// Assert payment confirmation was sent
		require.Len(t, notificationSpy.PaymentConfirmations, 1, "expected payment confirmation after payment success")

		paymentNotif := notificationSpy.PaymentConfirmations[0]
		assert.Equal(t, booking.ID, paymentNotif.BookingID, "payment notification should reference the booking")
		assert.Equal(t, clientID, paymentNotif.ClientID, "payment notification should reference the client")
	})
}
