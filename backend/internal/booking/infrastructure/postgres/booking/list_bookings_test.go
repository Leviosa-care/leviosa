package bookingRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/booking"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestListBookings TEST_PATH=internal/booking/infrastructure/postgres/booking/list_bookings_test.go

func TestListBookings(t *testing.T) {
	ctx := context.Background()

	t.Run("should list all bookings", func(t *testing.T) {
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

		// List all bookings
		bookings, err := repo.List(ctx, ports.BookingFilter{})
		require.NoError(t, err)
		assert.Len(t, bookings, 3)
	})

	t.Run("should return empty list when no bookings exist", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		bookings, err := repo.List(ctx, ports.BookingFilter{})
		require.NoError(t, err)
		assert.Empty(t, bookings)
	})

	t.Run("should filter by status", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert bookings with different statuses
		confirmed := tb.NewTestBookingEncx(t)
		completed := tb.NewCompletedBookingEncx(t, uuid.New(), uuid.New(), uuid.New(), uuid.New())
		cancelled := tb.NewCancelledBookingEncx(t, uuid.New(), uuid.New(), uuid.New(), uuid.New(), "test")

		err := tb.InsertBookingEncx(t, ctx, testPool, confirmed)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, completed)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, cancelled)
		require.NoError(t, err)

		// Filter for completed bookings
		filter := ports.BookingFilter{
			Status: []domain.BookingStatus{domain.BookingStatusCompleted},
		}
		bookings, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 1)
		assert.Equal(t, domain.BookingStatusCompleted, bookings[0].Status)
	})

	t.Run("should filter by multiple statuses", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert bookings with different statuses
		confirmed := tb.NewTestBookingEncx(t)
		completed := tb.NewCompletedBookingEncx(t, uuid.New(), uuid.New(), uuid.New(), uuid.New())
		cancelled := tb.NewCancelledBookingEncx(t, uuid.New(), uuid.New(), uuid.New(), uuid.New(), "test")

		err := tb.InsertBookingEncx(t, ctx, testPool, confirmed)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, completed)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, cancelled)
		require.NoError(t, err)

		// Filter for confirmed and completed
		filter := ports.BookingFilter{
			Status: []domain.BookingStatus{
				domain.BookingStatusConfirmed,
				domain.BookingStatusCompleted,
			},
		}
		bookings, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 2)
	})

	t.Run("should filter by payment status", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert bookings with different payment statuses
		pending := tb.NewTestBookingEncx(t)
		pending.PaymentStatus = domain.PaymentStatusPending

		paid := tb.NewTestBookingEncx(t)
		paid.PaymentStatus = domain.PaymentStatusPaid

		err := tb.InsertBookingEncx(t, ctx, testPool, pending)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, paid)
		require.NoError(t, err)

		// Filter for paid bookings
		filter := ports.BookingFilter{
			PaymentStatus: []domain.PaymentStatus{domain.PaymentStatusPaid},
		}
		bookings, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 1)
		assert.Equal(t, domain.PaymentStatusPaid, bookings[0].PaymentStatus)
	})

	t.Run("should filter by client ID", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		clientID := uuid.New()
		otherClientID := uuid.New()

		// Insert bookings for different clients
		booking1 := tb.NewTestBookingEncxWithIDs(t, uuid.New(), clientID, uuid.New(), uuid.New())
		booking2 := tb.NewTestBookingEncxWithIDs(t, uuid.New(), clientID, uuid.New(), uuid.New())
		booking3 := tb.NewTestBookingEncxWithIDs(t, uuid.New(), otherClientID, uuid.New(), uuid.New())

		err := tb.InsertBookingEncx(t, ctx, testPool, booking1)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, booking2)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, booking3)
		require.NoError(t, err)

		// Filter by client ID
		filter := ports.BookingFilter{ClientID: &clientID}
		bookings, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 2)
		for _, b := range bookings {
			assert.Equal(t, &clientID, b.ClientID)
		}
	})

	t.Run("should filter by partner ID", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		partnerID := uuid.New()
		otherPartnerID := uuid.New()

		// Insert bookings for different partners
		booking1 := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), partnerID, uuid.New())
		booking2 := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), partnerID, uuid.New())
		booking3 := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), otherPartnerID, uuid.New())

		err := tb.InsertBookingEncx(t, ctx, testPool, booking1)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, booking2)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, booking3)
		require.NoError(t, err)

		// Filter by partner ID
		filter := ports.BookingFilter{PartnerID: &partnerID}
		bookings, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 2)
		for _, b := range bookings {
			assert.Equal(t, partnerID, b.PartnerID)
		}
	})

	t.Run("should filter by created time range", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		now := time.Now()
		yesterday := now.AddDate(0, 0, -1)
		tomorrow := now.AddDate(0, 0, 1)

		// Insert bookings with different creation times
		oldBooking := tb.NewTestBookingEncx(t)
		oldBooking.CreatedAt = yesterday.Add(-2 * time.Hour)

		recentBooking := tb.NewTestBookingEncx(t)
		recentBooking.CreatedAt = now

		err := tb.InsertBookingEncx(t, ctx, testPool, oldBooking)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, recentBooking)
		require.NoError(t, err)

		// Filter for bookings created after yesterday
		filter := ports.BookingFilter{
			CreatedAfter: &yesterday,
		}
		bookings, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 1)

		// Filter for bookings created before tomorrow
		filter = ports.BookingFilter{
			CreatedBefore: &tomorrow,
		}
		bookings, err = repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 2)
	})

	t.Run("should apply pagination with limit", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert 5 bookings
		for i := 0; i < 5; i++ {
			booking := tb.NewTestBookingEncx(t)
			err := tb.InsertBookingEncx(t, ctx, testPool, booking)
			require.NoError(t, err)
		}

		// Get first 2 bookings
		filter := ports.BookingFilter{Limit: 2}
		bookings, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 2)
	})

	t.Run("should apply pagination with offset", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert 5 bookings
		for i := 0; i < 5; i++ {
			booking := tb.NewTestBookingEncx(t)
			err := tb.InsertBookingEncx(t, ctx, testPool, booking)
			require.NoError(t, err)
			time.Sleep(10 * time.Millisecond) // Ensure different timestamps
		}

		// Get bookings with offset
		filter := ports.BookingFilter{Limit: 2, Offset: 2}
		bookings, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 2)
	})

	t.Run("should sort by created_at ascending", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		now := time.Now()

		// Insert bookings with different timestamps
		booking1 := tb.NewTestBookingEncx(t)
		booking1.CreatedAt = now.Add(-2 * time.Hour)

		booking2 := tb.NewTestBookingEncx(t)
		booking2.CreatedAt = now.Add(-1 * time.Hour)

		booking3 := tb.NewTestBookingEncx(t)
		booking3.CreatedAt = now

		err := tb.InsertBookingEncx(t, ctx, testPool, booking1)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, booking2)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, booking3)
		require.NoError(t, err)

		// List with ascending order
		filter := ports.BookingFilter{
			OrderBy:        "created_at",
			OrderDirection: "asc",
		}
		bookings, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 3)

		// Verify order (oldest first)
		assert.True(t, bookings[0].CreatedAt.Before(bookings[1].CreatedAt))
		assert.True(t, bookings[1].CreatedAt.Before(bookings[2].CreatedAt))
	})

	t.Run("should combine multiple filters", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		clientID := uuid.New()

		// Insert multiple bookings
		confirmedPaid := tb.NewTestBookingEncxWithIDs(t, uuid.New(), clientID, uuid.New(), uuid.New())
		confirmedPaid.Status = domain.BookingStatusConfirmed
		confirmedPaid.PaymentStatus = domain.PaymentStatusPaid

		confirmedPending := tb.NewTestBookingEncxWithIDs(t, uuid.New(), clientID, uuid.New(), uuid.New())
		confirmedPending.Status = domain.BookingStatusConfirmed
		confirmedPending.PaymentStatus = domain.PaymentStatusPending

		completedPaid := tb.NewTestBookingEncxWithIDs(t, uuid.New(), clientID, uuid.New(), uuid.New())
		completedPaid.Status = domain.BookingStatusCompleted
		completedPaid.PaymentStatus = domain.PaymentStatusPaid

		err := tb.InsertBookingEncx(t, ctx, testPool, confirmedPaid)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, confirmedPending)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, completedPaid)
		require.NoError(t, err)

		// Filter: confirmed status + paid payment + specific client
		filter := ports.BookingFilter{
			ClientID:      &clientID,
			Status:        []domain.BookingStatus{domain.BookingStatusConfirmed},
			PaymentStatus: []domain.PaymentStatus{domain.PaymentStatusPaid},
		}
		bookings, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 1)
		assert.Equal(t, domain.BookingStatusConfirmed, bookings[0].Status)
		assert.Equal(t, domain.PaymentStatusPaid, bookings[0].PaymentStatus)
	})
}
