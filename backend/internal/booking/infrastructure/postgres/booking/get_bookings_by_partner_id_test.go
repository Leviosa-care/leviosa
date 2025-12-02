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

// make test-func TEST_NAME=TestGetByPartnerID TEST_PATH=internal/booking/infrastructure/postgres/booking/get_bookings_by_partner_id_test.go

func TestGetByPartnerID(t *testing.T) {
	ctx := context.Background()

	t.Run("should retrieve all bookings for a partner", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		partnerID := uuid.New()
		otherPartnerID := uuid.New()

		// Insert bookings for the partner
		booking1 := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), partnerID, uuid.New())
		booking2 := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), partnerID, uuid.New())

		// Insert booking for different partner
		otherBooking := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), otherPartnerID, uuid.New())

		err := tb.InsertBookingEncx(t, ctx, testPool, booking1)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, booking2)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, otherBooking)
		require.NoError(t, err)

		// Get bookings for partner
		bookings, err := repo.GetByPartnerID(ctx, partnerID, ports.BookingFilter{})
		require.NoError(t, err)
		assert.Len(t, bookings, 2)

		// Verify all bookings belong to the partner
		for _, b := range bookings {
			assert.Equal(t, partnerID, b.PartnerID)
		}
	})

	t.Run("should return empty list for partner with no bookings", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		nonExistentPartnerID := uuid.New()

		bookings, err := repo.GetByPartnerID(ctx, nonExistentPartnerID, ports.BookingFilter{})
		require.NoError(t, err)
		assert.Empty(t, bookings)
	})

	t.Run("should filter partner bookings by status", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		partnerID := uuid.New()

		// Insert bookings with different statuses
		confirmed := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), partnerID, uuid.New())
		confirmed.Status = domain.BookingStatusConfirmed

		completed := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), partnerID, uuid.New())
		completed.Status = domain.BookingStatusCompleted
		completed.PaymentStatus = domain.PaymentStatusPaid

		noShow := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), partnerID, uuid.New())
		noShow.Status = domain.BookingStatusNoShow

		err := tb.InsertBookingEncx(t, ctx, testPool, confirmed)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, completed)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, noShow)
		require.NoError(t, err)

		// Get only completed bookings
		filter := ports.BookingFilter{
			Status: []domain.BookingStatus{domain.BookingStatusCompleted},
		}
		bookings, err := repo.GetByPartnerID(ctx, partnerID, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 1)
		assert.Equal(t, domain.BookingStatusCompleted, bookings[0].Status)
	})

	t.Run("should filter partner bookings by room", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		partnerID := uuid.New()
		roomID := uuid.New()
		otherRoomID := uuid.New()

		// Insert bookings for different rooms
		booking1 := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), partnerID, roomID)
		booking2 := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), partnerID, roomID)
		booking3 := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), partnerID, otherRoomID)

		err := tb.InsertBookingEncx(t, ctx, testPool, booking1)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, booking2)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, booking3)
		require.NoError(t, err)

		// Get bookings for specific room
		filter := ports.BookingFilter{RoomID: &roomID}
		bookings, err := repo.GetByPartnerID(ctx, partnerID, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 2)
		for _, b := range bookings {
			assert.Equal(t, roomID, b.RoomID)
		}
	})

	t.Run("should paginate partner bookings", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		partnerID := uuid.New()

		// Insert 5 bookings for the partner
		for i := 0; i < 5; i++ {
			booking := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), partnerID, uuid.New())
			err := tb.InsertBookingEncx(t, ctx, testPool, booking)
			require.NoError(t, err)
		}

		// Get first 3 bookings
		filter := ports.BookingFilter{Limit: 3}
		bookings, err := repo.GetByPartnerID(ctx, partnerID, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 3)
	})

	t.Run("should combine multiple filters for partner bookings", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		partnerID := uuid.New()
		roomID := uuid.New()

		// Insert various bookings
		confirmedPaid := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), partnerID, roomID)
		confirmedPaid.Status = domain.BookingStatusConfirmed
		confirmedPaid.PaymentStatus = domain.PaymentStatusPaid

		confirmedPending := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), partnerID, roomID)
		confirmedPending.Status = domain.BookingStatusConfirmed
		confirmedPending.PaymentStatus = domain.PaymentStatusPending

		completedInRoom := tb.NewTestBookingEncxWithIDs(t, uuid.New(), uuid.New(), partnerID, roomID)
		completedInRoom.Status = domain.BookingStatusCompleted

		err := tb.InsertBookingEncx(t, ctx, testPool, confirmedPaid)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, confirmedPending)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, completedInRoom)
		require.NoError(t, err)

		// Filter: confirmed + paid + specific room
		filter := ports.BookingFilter{
			RoomID:        &roomID,
			Status:        []domain.BookingStatus{domain.BookingStatusConfirmed},
			PaymentStatus: []domain.PaymentStatus{domain.PaymentStatusPaid},
		}
		bookings, err := repo.GetByPartnerID(ctx, partnerID, filter)
		require.NoError(t, err)
		assert.Len(t, bookings, 1)
		assert.Equal(t, domain.BookingStatusConfirmed, bookings[0].Status)
		assert.Equal(t, domain.PaymentStatusPaid, bookings[0].PaymentStatus)
		assert.Equal(t, roomID, bookings[0].RoomID)
	})
}
