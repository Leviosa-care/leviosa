package bookingRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/booking"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarkReminderSent(t *testing.T) {
	ctx := context.Background()

	t.Run("sets reminded_at on an existing booking", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		booking := tb.NewTestBookingEncx(t)
		booking.Status = domain.BookingStatusConfirmed
		assert.Nil(t, booking.RemindedAt, "new booking should have nil reminded_at")

		err := tb.InsertBookingEncx(t, ctx, testPool, booking)
		require.NoError(t, err)

		before := time.Now()
		err = repo.MarkReminderSent(ctx, booking.ID)
		require.NoError(t, err)

		updated, err := tb.GetBookingEncxByID(t, ctx, testPool, booking.ID)
		require.NoError(t, err)
		require.NotNil(t, updated.RemindedAt, "reminded_at should be set")
		assert.True(t, updated.RemindedAt.After(before.Add(-1*time.Second)), "reminded_at should be recent")
	})

	t.Run("returns error for non-existent booking", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		err := repo.MarkReminderSent(ctx, uuid.New())
		assert.Error(t, err, "expected error for non-existent booking")
	})

	t.Run("idempotent — calling twice does not error", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		booking := tb.NewTestBookingEncx(t)
		booking.Status = domain.BookingStatusConfirmed

		err := tb.InsertBookingEncx(t, ctx, testPool, booking)
		require.NoError(t, err)

		err = repo.MarkReminderSent(ctx, booking.ID)
		require.NoError(t, err)

		// Second call should succeed (the row exists, so RowsAffected > 0)
		err = repo.MarkReminderSent(ctx, booking.ID)
		require.NoError(t, err)
	})
}
