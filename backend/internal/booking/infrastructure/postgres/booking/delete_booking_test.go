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

// make test-func TEST_NAME=TestDeleteBooking TEST_PATH=internal/booking/infrastructure/postgres/booking/delete_booking_test.go

func TestDeleteBooking(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully delete existing booking", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert booking
		bookingEncx := tb.NewTestBookingEncx(t)
		err := tb.InsertBookingEncx(t, ctx, testPool, bookingEncx)
		require.NoError(t, err)

		// Verify booking exists
		exists, err := tb.BookingExists(t, ctx, testPool, bookingEncx.ID)
		require.NoError(t, err)
		assert.True(t, exists)

		// Delete booking
		err = repo.Delete(ctx, bookingEncx.ID)
		require.NoError(t, err)

		// Verify booking no longer exists
		exists, err = tb.BookingExists(t, ctx, testPool, bookingEncx.ID)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("should return error for non-existent booking", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		nonExistentID := uuid.New()
		err := repo.Delete(ctx, nonExistentID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("should perform hard delete (GDPR compliance)", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert booking with encrypted personal data
		bookingEncx := tb.NewTestBookingEncxWithNotes(
			t,
			uuid.New(), uuid.New(), uuid.New(), uuid.New(),
			"personal client notes",
			"personal partner notes",
		)
		err := tb.InsertBookingEncx(t, ctx, testPool, bookingEncx)
		require.NoError(t, err)

		// Delete booking
		err = repo.Delete(ctx, bookingEncx.ID)
		require.NoError(t, err)

		// Verify data is completely removed (not soft deleted)
		count, err := tb.CountBookings(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, 0, count, "Booking should be completely removed, not soft deleted")

		// Verify booking cannot be retrieved
		exists, err := tb.BookingExists(t, ctx, testPool, bookingEncx.ID)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("should delete multiple bookings independently", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert multiple bookings
		booking1 := tb.NewTestBookingEncx(t)
		booking2 := tb.NewTestBookingEncx(t)
		booking3 := tb.NewTestBookingEncx(t)

		err := tb.InsertBookingEncx(t, ctx, testPool, booking1)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, booking2)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, booking3)
		require.NoError(t, err)

		// Verify all exist
		count, err := tb.CountBookings(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, 3, count)

		// Delete first booking
		err = repo.Delete(ctx, booking1.ID)
		require.NoError(t, err)

		// Verify only first booking is deleted
		count, err = tb.CountBookings(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, 2, count)

		exists, err := tb.BookingExists(t, ctx, testPool, booking1.ID)
		require.NoError(t, err)
		assert.False(t, exists)

		exists, err = tb.BookingExists(t, ctx, testPool, booking2.ID)
		require.NoError(t, err)
		assert.True(t, exists)

		exists, err = tb.BookingExists(t, ctx, testPool, booking3.ID)
		require.NoError(t, err)
		assert.True(t, exists)
	})
}
