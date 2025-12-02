package bookingRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/google/uuid"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/booking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUpcoming(t *testing.T) {
	ctx := context.Background()

	t.Run("should retrieve only confirmed bookings", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert bookings with different statuses
		confirmed := tb.NewTestBookingEncx(t)
		confirmed.Status = domain.BookingStatusConfirmed

		completed := tb.NewCompletedBookingEncx(t, uuid.New(), uuid.New(), uuid.New(), uuid.New())
		cancelled := tb.NewCancelledBookingEncx(t, uuid.New(), uuid.New(), uuid.New(), uuid.New(), "test")
		noShow := tb.NewNoShowBookingEncx(t, uuid.New(), uuid.New(), uuid.New(), uuid.New())

		err := tb.InsertBookingEncx(t, ctx, testPool, confirmed)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, completed)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, cancelled)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, noShow)
		require.NoError(t, err)

		// Get upcoming bookings (confirmed only)
		bookings, err := repo.GetUpcoming(ctx, ports.BookingFilter{})
		require.NoError(t, err)
		assert.Len(t, bookings, 1)
		assert.Equal(t, domain.BookingStatusConfirmed, bookings[0].Status)
	})

	t.Run("should return empty list when no confirmed bookings exist", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert only non-confirmed bookings
		completed := tb.NewCompletedBookingEncx(t, uuid.New(), uuid.New(), uuid.New(), uuid.New())
		cancelled := tb.NewCancelledBookingEncx(t, uuid.New(), uuid.New(), uuid.New(), uuid.New(), "test")

		err := tb.InsertBookingEncx(t, ctx, testPool, completed)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, cancelled)
		require.NoError(t, err)

		bookings, err := repo.GetUpcoming(ctx, ports.BookingFilter{})
		require.NoError(t, err)
		assert.Empty(t, bookings)
	})

	t.Run("should filter upcoming bookings by partner", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		partnerID := uuid.New()
		otherPartnerID := uuid.New()

		// Insert confirmed bookings for different partners
		booking1 := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), partnerID, uuid.New())
		booking1.Status = domain.BookingStatusConfirmed

		booking2 := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), partnerID, uuid.New())
		booking2.Status = domain.BookingStatusConfirmed

		otherBooking := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), otherPartnerID, uuid.New())
		otherBooking.Status = domain.BookingStatusConfirmed

		err := tb.InsertBookingEncx(t, ctx, testPool, booking1)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, booking2)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, otherBooking)
		require.NoError(t, err)

		// Get upcoming bookings for specific partner
		filter := ports.BookingFilter{PartnerID: &partnerID}
		bookings, err := repo.GetUpcoming(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 2)
		for _, b := range bookings {
			assert.Equal(t, partnerID, b.PartnerID)
			assert.Equal(t, domain.BookingStatusConfirmed, b.Status)
		}
	})

	t.Run("should filter upcoming bookings by room", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		roomID := uuid.New()
		otherRoomID := uuid.New()

		// Insert confirmed bookings for different rooms
		booking1 := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), uuid.New(), roomID)
		booking1.Status = domain.BookingStatusConfirmed

		booking2 := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), uuid.New(), otherRoomID)
		booking2.Status = domain.BookingStatusConfirmed

		err := tb.InsertBookingEncx(t, ctx, testPool, booking1)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, booking2)
		require.NoError(t, err)

		// Get upcoming bookings for specific room
		filter := ports.BookingFilter{RoomID: &roomID}
		bookings, err := repo.GetUpcoming(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 1)
		assert.Equal(t, roomID, bookings[0].RoomID)
		assert.Equal(t, domain.BookingStatusConfirmed, bookings[0].Status)
	})

	t.Run("should filter upcoming bookings by availability", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		availabilityID := uuid.New()
		otherAvailabilityID := uuid.New()

		// Insert confirmed bookings for different availabilities
		booking1 := tb.NewTestBookingEncxWithIDs(t, availabilityID, uuid.New(), uuid.New(), uuid.New())
		booking1.Status = domain.BookingStatusConfirmed

		booking2 := tb.NewTestBookingEncxWithIDs(t, otherAvailabilityID, uuid.New(), uuid.New(), uuid.New())
		booking2.Status = domain.BookingStatusConfirmed

		err := tb.InsertBookingEncx(t, ctx, testPool, booking1)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, booking2)
		require.NoError(t, err)

		// Get upcoming bookings for specific availability
		filter := ports.BookingFilter{AvailabilityID: &availabilityID}
		bookings, err := repo.GetUpcoming(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 1)
		assert.Equal(t, availabilityID, bookings[0].AvailabilityID)
	})

	t.Run("should paginate upcoming bookings", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert 5 confirmed bookings
		for i := 0; i < 5; i++ {
			booking := tb.NewTestBookingEncx(t)
			booking.Status = domain.BookingStatusConfirmed
			err := tb.InsertBookingEncx(t, ctx, testPool, booking)
			require.NoError(t, err)
		}

		// Get first 2 upcoming bookings
		filter := ports.BookingFilter{Limit: 2}
		bookings, err := repo.GetUpcoming(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 2)
		for _, b := range bookings {
			assert.Equal(t, domain.BookingStatusConfirmed, b.Status)
		}
	})

	t.Run("should combine filters for upcoming bookings", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		partnerID := uuid.New()
		roomID := uuid.New()

		// Insert various confirmed bookings
		targetBooking := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), partnerID, roomID)
		targetBooking.Status = domain.BookingStatusConfirmed
		targetBooking.PaymentStatus = domain.PaymentStatusPaid

		otherRoomBooking := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), partnerID, uuid.New())
		otherRoomBooking.Status = domain.BookingStatusConfirmed

		pendingPaymentBooking := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), partnerID, roomID)
		pendingPaymentBooking.Status = domain.BookingStatusConfirmed
		pendingPaymentBooking.PaymentStatus = domain.PaymentStatusPending

		err := tb.InsertBookingEncx(t, ctx, testPool, targetBooking)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, otherRoomBooking)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, pendingPaymentBooking)
		require.NoError(t, err)

		// Filter: specific partner + room + paid
		filter := ports.BookingFilter{
			PartnerID:     &partnerID,
			RoomID:        &roomID,
			PaymentStatus: []domain.PaymentStatus{domain.PaymentStatusPaid},
		}
		bookings, err := repo.GetUpcoming(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 1)
		assert.Equal(t, domain.BookingStatusConfirmed, bookings[0].Status)
		assert.Equal(t, domain.PaymentStatusPaid, bookings[0].PaymentStatus)
		assert.Equal(t, partnerID, bookings[0].PartnerID)
		assert.Equal(t, roomID, bookings[0].RoomID)
	})
}
