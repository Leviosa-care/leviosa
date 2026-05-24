package bookingRepository_test

import (
	"context"
	"testing"

	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/booking"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreateBooking TEST_PATH=internal/booking/infrastructure/postgres/booking/create_booking_test.go

func TestCreateBooking(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully create a booking", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		bookingEncx := tb.NewTestBookingEncx(t)
		tb.EnsureBookingForeignKeys(t, ctx, testPool, bookingEncx)
		err := repo.Create(ctx, bookingEncx)
		require.NoError(t, err)

		// Verify insertion
		count, err := tb.CountBookings(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, 1, count)

		// Verify booking exists
		exists, err := tb.BookingExists(t, ctx, testPool, bookingEncx.ID)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should handle duplicate booking ID", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		bookingEncx := tb.NewTestBookingEncx(t)
		tb.EnsureBookingForeignKeys(t, ctx, testPool, bookingEncx)
		err := repo.Create(ctx, bookingEncx)
		require.NoError(t, err)

		// Attempt duplicate
		err = repo.Create(ctx, bookingEncx)
		assert.Error(t, err, "Should fail on duplicate ID")
	})

	t.Run("should create booking with all encrypted fields", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		bookingEncx := tb.NewTestBookingEncx(t)
		tb.EnsureBookingForeignKeys(t, ctx, testPool, bookingEncx)

		err := repo.Create(ctx, bookingEncx)
		require.NoError(t, err)

		// Verify all fields persisted
		saved, err := tb.GetBookingEncxByID(t, ctx, testPool, bookingEncx.ID)
		require.NoError(t, err)
		assert.Equal(t, bookingEncx.ID, saved.ID)
		assert.Equal(t, bookingEncx.AvailabilityID, saved.AvailabilityID)
		assert.Equal(t, bookingEncx.ClientID, saved.ClientID)
		assert.Equal(t, bookingEncx.PartnerID, saved.PartnerID)
		assert.Equal(t, bookingEncx.RoomID, saved.RoomID)
		assert.Equal(t, bookingEncx.ProductIDEncrypted, saved.ProductIDEncrypted)
		assert.Equal(t, bookingEncx.SlotStartTimeEncrypted, saved.SlotStartTimeEncrypted)
		assert.Equal(t, bookingEncx.SlotEndTimeEncrypted, saved.SlotEndTimeEncrypted)
		assert.Equal(t, bookingEncx.TotalPriceCents, saved.TotalPriceCents)
		assert.Equal(t, bookingEncx.Currency, saved.Currency)
		assert.Equal(t, bookingEncx.PaymentStatus, saved.PaymentStatus)
		assert.Equal(t, bookingEncx.Status, saved.Status)
	})

	t.Run("should create booking with payment intent", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		paymentIntentID := "pi_test_123456"
		bookingEncx := tb.NewTestBookingEncxWithPaymentIntent(
			t,
			tb.NewTestBookingEncx(t).AvailabilityID,
			*tb.NewTestBookingEncx(t).ClientID,
			tb.NewTestBookingEncx(t).PartnerID,
			tb.NewTestBookingEncx(t).RoomID,
			paymentIntentID,
		)
		tb.EnsureBookingForeignKeys(t, ctx, testPool, bookingEncx)

		err := repo.Create(ctx, bookingEncx)
		require.NoError(t, err)

		// Verify payment intent was saved
		saved, err := tb.GetBookingEncxByID(t, ctx, testPool, bookingEncx.ID)
		require.NoError(t, err)
		require.NotNil(t, saved.PaymentIntentID)
		assert.Equal(t, paymentIntentID, *saved.PaymentIntentID)
	})

	t.Run("should create booking with notes", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		bookingEncx := tb.NewTestBookingEncxWithNotes(
			t,
			tb.NewTestBookingEncx(t).AvailabilityID,
			*tb.NewTestBookingEncx(t).ClientID,
			tb.NewTestBookingEncx(t).PartnerID,
			tb.NewTestBookingEncx(t).RoomID,
			"client note content",
			"partner note content",
		)
		tb.EnsureBookingForeignKeys(t, ctx, testPool, bookingEncx)

		err := repo.Create(ctx, bookingEncx)
		require.NoError(t, err)

		// Verify notes were saved
		saved, err := tb.GetBookingEncxByID(t, ctx, testPool, bookingEncx.ID)
		require.NoError(t, err)
		assert.NotEmpty(t, saved.ClientNotesEncrypted)
		assert.NotEmpty(t, saved.PartnerNotesEncrypted)
	})

	t.Run("should create a guest booking with nil client_id", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		guestBooking := tb.NewGuestBookingEncx(t, "Alice", "Martin", "alice@example.com", "+33612345678")
		tb.EnsureBookingForeignKeys(t, ctx, testPool, guestBooking)

		err := repo.Create(ctx, guestBooking)
		require.NoError(t, err)

		saved, err := tb.GetBookingEncxByID(t, ctx, testPool, guestBooking.ID)
		require.NoError(t, err)
		assert.Nil(t, saved.ClientID, "guest booking should have nil client_id")
		assert.NotEmpty(t, saved.GuestFirstNameEncrypted)
		assert.NotEmpty(t, saved.GuestLastNameEncrypted)
		assert.NotEmpty(t, saved.GuestEmailEncrypted)
		assert.NotEmpty(t, saved.GuestPhoneEncrypted)
	})

	t.Run("should create bookings with different statuses", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Create completed booking
		completedBooking := tb.NewCompletedBookingEncx(
			t,
			tb.NewTestBookingEncx(t).AvailabilityID,
			*tb.NewTestBookingEncx(t).ClientID,
			tb.NewTestBookingEncx(t).PartnerID,
			tb.NewTestBookingEncx(t).RoomID,
		)
		tb.EnsureBookingForeignKeys(t, ctx, testPool, completedBooking)
		err := repo.Create(ctx, completedBooking)
		require.NoError(t, err)

		// Create cancelled booking
		cancelledBooking := tb.NewCancelledBookingEncx(
			t,
			tb.NewTestBookingEncx(t).AvailabilityID,
			*tb.NewTestBookingEncx(t).ClientID,
			tb.NewTestBookingEncx(t).PartnerID,
			tb.NewTestBookingEncx(t).RoomID,
			"client requested",
		)
		tb.EnsureBookingForeignKeys(t, ctx, testPool, cancelledBooking)
		err = repo.Create(ctx, cancelledBooking)
		require.NoError(t, err)

		// Verify both were created
		count, err := tb.CountBookings(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, 2, count)
	})
}
