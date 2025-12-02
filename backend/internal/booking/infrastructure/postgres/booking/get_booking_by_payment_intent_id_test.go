package bookingRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/booking"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetByPaymentIntentID TEST_PATH=internal/booking/infrastructure/postgres/booking/get_booking_by_payment_intent_id_test.go

func TestGetByPaymentIntentID(t *testing.T) {
	ctx := context.Background()

	t.Run("should retrieve booking by payment intent ID", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		paymentIntentID := "pi_test_1234567890"

		// Insert booking with payment intent
		booking := tb.NewTestBookingEncxWithPaymentIntent(
			t,
			uuid.New(),
			uuid.New(),
			uuid.New(),
			uuid.New(),
			paymentIntentID,
		)
		err := tb.InsertBookingEncx(t, ctx, testPool, booking)
		require.NoError(t, err)

		// Retrieve by payment intent ID
		retrieved, err := repo.GetByPaymentIntentID(ctx, paymentIntentID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, booking.ID, retrieved.ID)
		require.NotNil(t, retrieved.PaymentIntentID)
		assert.Equal(t, paymentIntentID, *retrieved.PaymentIntentID)
	})

	t.Run("should return error when payment intent ID does not exist", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		nonExistentPaymentIntentID := "pi_nonexistent_999"

		retrieved, err := repo.GetByPaymentIntentID(ctx, nonExistentPaymentIntentID)
		assert.Error(t, err)
		assert.Nil(t, retrieved)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("should retrieve booking with all payment fields", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		paymentIntentID := "pi_complete_payment"

		// Insert booking with payment data
		booking := tb.NewTestBookingEncxWithPaymentIntent(
			t,
			uuid.New(),
			uuid.New(),
			uuid.New(),
			uuid.New(),
			paymentIntentID,
		)
		booking.TotalPriceCents = 7500
		booking.Currency = "EUR"

		err := tb.InsertBookingEncx(t, ctx, testPool, booking)
		require.NoError(t, err)

		// Retrieve and verify payment fields
		retrieved, err := repo.GetByPaymentIntentID(ctx, paymentIntentID)
		require.NoError(t, err)
		assert.Equal(t, 7500, retrieved.TotalPriceCents)
		assert.Equal(t, "EUR", retrieved.Currency)
		assert.Equal(t, booking.PaymentStatus, retrieved.PaymentStatus)
	})

	t.Run("should handle Stripe payment intent format", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Realistic Stripe payment intent ID format
		stripePaymentIntentID := "pi_3OMhGz2eZvKYlo2C1234567890"

		booking := tb.NewTestBookingEncxWithPaymentIntent(
			t,
			uuid.New(),
			uuid.New(),
			uuid.New(),
			uuid.New(),
			stripePaymentIntentID,
		)
		err := tb.InsertBookingEncx(t, ctx, testPool, booking)
		require.NoError(t, err)

		// Retrieve
		retrieved, err := repo.GetByPaymentIntentID(ctx, stripePaymentIntentID)
		require.NoError(t, err)
		require.NotNil(t, retrieved.PaymentIntentID)
		assert.Equal(t, stripePaymentIntentID, *retrieved.PaymentIntentID)
	})

	t.Run("should retrieve booking with all fields for payment processing", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		paymentIntentID := "pi_processing_test"

		// Insert booking with complete data
		booking := tb.NewTestBookingEncxWithPaymentIntent(
			t,
			uuid.New(),
			uuid.New(),
			uuid.New(),
			uuid.New(),
			paymentIntentID,
		)
		booking.ClientNotesEncrypted = []byte("encrypted_client_notes")
		booking.PartnerNotesEncrypted = []byte("encrypted_partner_notes")

		err := tb.InsertBookingEncx(t, ctx, testPool, booking)
		require.NoError(t, err)

		// Retrieve and verify all fields
		retrieved, err := repo.GetByPaymentIntentID(ctx, paymentIntentID)
		require.NoError(t, err)
		assert.Equal(t, booking.ID, retrieved.ID)
		assert.Equal(t, booking.AvailabilityID, retrieved.AvailabilityID)
		assert.Equal(t, booking.ClientID, retrieved.ClientID)
		assert.Equal(t, booking.PartnerID, retrieved.PartnerID)
		assert.Equal(t, booking.RoomID, retrieved.RoomID)
		assert.Equal(t, booking.ProductIDEncrypted, retrieved.ProductIDEncrypted)
		assert.Equal(t, booking.SlotStartTimeEncrypted, retrieved.SlotStartTimeEncrypted)
		assert.Equal(t, booking.SlotEndTimeEncrypted, retrieved.SlotEndTimeEncrypted)
		assert.Equal(t, booking.Status, retrieved.Status)
		assert.Equal(t, booking.DEKEncrypted, retrieved.DEKEncrypted)
	})

	t.Run("should distinguish between different payment intents", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		paymentIntent1 := "pi_booking_1"
		paymentIntent2 := "pi_booking_2"

		// Insert two bookings with different payment intents
		booking1 := tb.NewTestBookingEncxWithPaymentIntent(
			t,
			uuid.New(),
			uuid.New(),
			uuid.New(),
			uuid.New(),
			paymentIntent1,
		)
		booking2 := tb.NewTestBookingEncxWithPaymentIntent(
			t,
			uuid.New(),
			uuid.New(),
			uuid.New(),
			uuid.New(),
			paymentIntent2,
		)

		err := tb.InsertBookingEncx(t, ctx, testPool, booking1)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, booking2)
		require.NoError(t, err)

		// Retrieve first booking
		retrieved1, err := repo.GetByPaymentIntentID(ctx, paymentIntent1)
		require.NoError(t, err)
		assert.Equal(t, booking1.ID, retrieved1.ID)

		// Retrieve second booking
		retrieved2, err := repo.GetByPaymentIntentID(ctx, paymentIntent2)
		require.NoError(t, err)
		assert.Equal(t, booking2.ID, retrieved2.ID)

		// Verify they're different
		assert.NotEqual(t, retrieved1.ID, retrieved2.ID)
	})
}
