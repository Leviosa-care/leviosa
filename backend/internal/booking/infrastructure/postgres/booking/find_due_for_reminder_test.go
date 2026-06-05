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

func TestFindBookingsDueForReminder(t *testing.T) {
	ctx := context.Background()

	t.Run("returns only confirmed bookings with reminded_at IS NULL", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Confirmed, not reminded → should be returned
		confirmed := tb.NewTestBookingEncx(t)
		confirmed.Status = domain.BookingStatusConfirmed

		// Confirmed, already reminded → should NOT be returned
		reminded := tb.NewTestBookingEncx(t)
		reminded.Status = domain.BookingStatusConfirmed
		now := time.Now()
		reminded.RemindedAt = &now

		err := tb.InsertBookingEncx(t, ctx, testPool, confirmed)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, reminded)
		require.NoError(t, err)

		bookings, err := repo.FindBookingsDueForReminder(ctx)
		require.NoError(t, err)

		ids := make(map[uuid.UUID]bool)
		for _, b := range bookings {
			ids[b.ID] = true
		}

		assert.True(t, ids[confirmed.ID], "confirmed unreminded booking should be returned")
		assert.False(t, ids[reminded.ID], "already-reminded booking should NOT be returned")
	})

	t.Run("excludes cancelled bookings", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		cancelled := tb.NewCancelledBookingEncx(t, uuid.New(), uuid.New(), uuid.New(), uuid.New(), "test")
		err := tb.InsertBookingEncx(t, ctx, testPool, cancelled)
		require.NoError(t, err)

		bookings, err := repo.FindBookingsDueForReminder(ctx)
		require.NoError(t, err)

		for _, b := range bookings {
			assert.NotEqual(t, cancelled.ID, b.ID, "cancelled booking should be excluded")
		}
	})

	t.Run("excludes completed bookings", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		completed := tb.NewCompletedBookingEncx(t, uuid.New(), uuid.New(), uuid.New(), uuid.New())
		err := tb.InsertBookingEncx(t, ctx, testPool, completed)
		require.NoError(t, err)

		bookings, err := repo.FindBookingsDueForReminder(ctx)
		require.NoError(t, err)

		for _, b := range bookings {
			assert.NotEqual(t, completed.ID, b.ID, "completed booking should be excluded")
		}
	})

	t.Run("excludes no_show bookings", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		noShow := tb.NewNoShowBookingEncx(t, uuid.New(), uuid.New(), uuid.New(), uuid.New())
		err := tb.InsertBookingEncx(t, ctx, testPool, noShow)
		require.NoError(t, err)

		bookings, err := repo.FindBookingsDueForReminder(ctx)
		require.NoError(t, err)

		for _, b := range bookings {
			assert.NotEqual(t, noShow.ID, b.ID, "no_show booking should be excluded")
		}
	})

	t.Run("returns empty when all bookings are reminded", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		now := time.Now()
		b1 := tb.NewTestBookingEncx(t)
		b1.Status = domain.BookingStatusConfirmed
		b1.RemindedAt = &now

		b2 := tb.NewTestBookingEncx(t)
		b2.Status = domain.BookingStatusConfirmed
		b2.RemindedAt = &now

		err := tb.InsertBookingEncx(t, ctx, testPool, b1)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, b2)
		require.NoError(t, err)

		bookings, err := repo.FindBookingsDueForReminder(ctx)
		require.NoError(t, err)
		assert.Empty(t, bookings)
	})

	t.Run("returns multiple confirmed unreminded bookings", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		b1 := tb.NewTestBookingEncx(t)
		b1.Status = domain.BookingStatusConfirmed

		b2 := tb.NewTestBookingEncx(t)
		b2.Status = domain.BookingStatusConfirmed

		err := tb.InsertBookingEncx(t, ctx, testPool, b1)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, b2)
		require.NoError(t, err)

		bookings, err := repo.FindBookingsDueForReminder(ctx)
		require.NoError(t, err)

		ids := make(map[uuid.UUID]bool)
		for _, b := range bookings {
			ids[b.ID] = true
		}
		assert.True(t, ids[b1.ID])
		assert.True(t, ids[b2.ID])
	})
}
