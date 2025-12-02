package bookingRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/booking"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetByClientID TEST_PATH=internal/booking/infrastructure/postgres/booking/get_bookings_by_client_id_test.go

func TestGetByClientID(t *testing.T) {
	ctx := context.Background()

	t.Run("should retrieve all bookings for a client", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		clientID := uuid.New()
		otherClientID := uuid.New()

		// Insert bookings for the client
		booking1 := tb.NewTestBookingEncxWithIDs(t, uuid.New(), clientID, uuid.New(), uuid.New())
		booking2 := tb.NewTestBookingEncxWithIDs(t, uuid.New(), clientID, uuid.New(), uuid.New())

		// Insert booking for different client
		otherBooking := tb.NewTestBookingEncxWithIDs(t, uuid.New(), otherClientID, uuid.New(), uuid.New())

		err := tb.InsertBookingEncx(t, ctx, testPool, booking1)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, booking2)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, otherBooking)
		require.NoError(t, err)

		// Get bookings for client
		bookings, err := repo.GetByClientID(ctx, clientID, ports.BookingFilter{})
		require.NoError(t, err)
		assert.Len(t, bookings, 2)

		// Verify all bookings belong to the client
		for _, b := range bookings {
			assert.Equal(t, clientID, b.ClientID)
		}
	})

	t.Run("should return empty list for client with no bookings", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		nonExistentClientID := uuid.New()

		bookings, err := repo.GetByClientID(ctx, nonExistentClientID, ports.BookingFilter{})
		require.NoError(t, err)
		assert.Empty(t, bookings)
	})

	t.Run("should filter client bookings by status", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		clientID := uuid.New()

		// Insert bookings with different statuses
		confirmed := tb.NewTestBookingEncxWithIDs(t, uuid.New(), clientID, uuid.New(), uuid.New())
		confirmed.Status = domain.BookingStatusConfirmed

		completed := tb.NewTestBookingEncxWithIDs(t, uuid.New(), clientID, uuid.New(), uuid.New())
		completed.Status = domain.BookingStatusCompleted
		completed.PaymentStatus = domain.PaymentStatusPaid

		cancelled := tb.NewTestBookingEncxWithIDs(t, uuid.New(), clientID, uuid.New(), uuid.New())
		cancelled.Status = domain.BookingStatusCancelled

		err := tb.InsertBookingEncx(t, ctx, testPool, confirmed)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, completed)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, cancelled)
		require.NoError(t, err)

		// Get only confirmed bookings
		filter := ports.BookingFilter{
			Status: []domain.BookingStatus{domain.BookingStatusConfirmed},
		}
		bookings, err := repo.GetByClientID(ctx, clientID, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 1)
		assert.Equal(t, domain.BookingStatusConfirmed, bookings[0].Status)
	})

	t.Run("should filter client bookings by payment status", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		clientID := uuid.New()

		// Insert bookings with different payment statuses
		pending := tb.NewTestBookingEncxWithIDs(t, uuid.New(), clientID, uuid.New(), uuid.New())
		pending.PaymentStatus = domain.PaymentStatusPending

		paid := tb.NewTestBookingEncxWithIDs(t, uuid.New(), clientID, uuid.New(), uuid.New())
		paid.PaymentStatus = domain.PaymentStatusPaid

		err := tb.InsertBookingEncx(t, ctx, testPool, pending)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, paid)
		require.NoError(t, err)

		// Get only paid bookings
		filter := ports.BookingFilter{
			PaymentStatus: []domain.PaymentStatus{domain.PaymentStatusPaid},
		}
		bookings, err := repo.GetByClientID(ctx, clientID, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 1)
		assert.Equal(t, domain.PaymentStatusPaid, bookings[0].PaymentStatus)
	})

	t.Run("should paginate client bookings", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		clientID := uuid.New()

		// Insert 5 bookings for the client
		for i := 0; i < 5; i++ {
			booking := tb.NewTestBookingEncxWithIDs(t, uuid.New(), clientID, uuid.New(), uuid.New())
			err := tb.InsertBookingEncx(t, ctx, testPool, booking)
			require.NoError(t, err)
		}

		// Get first 2 bookings
		filter := ports.BookingFilter{Limit: 2}
		bookings, err := repo.GetByClientID(ctx, clientID, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 2)

		// Get next 2 bookings
		filter = ports.BookingFilter{Limit: 2, Offset: 2}
		bookings, err = repo.GetByClientID(ctx, clientID, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 2)
	})

	t.Run("should sort client bookings", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		clientID := uuid.New()

		// Insert bookings with different prices
		cheap := tb.NewTestBookingEncxWithIDs(t, uuid.New(), clientID, uuid.New(), uuid.New())
		cheap.TotalPriceCents = 1000

		expensive := tb.NewTestBookingEncxWithIDs(t, uuid.New(), clientID, uuid.New(), uuid.New())
		expensive.TotalPriceCents = 5000

		err := tb.InsertBookingEncx(t, ctx, testPool, cheap)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, expensive)
		require.NoError(t, err)

		// Sort by price ascending
		filter := ports.BookingFilter{
			OrderBy:        "total_price_cents",
			OrderDirection: "asc",
		}
		bookings, err := repo.GetByClientID(ctx, clientID, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 2)
		assert.Equal(t, 1000, bookings[0].TotalPriceCents)
		assert.Equal(t, 5000, bookings[1].TotalPriceCents)
	})
}
