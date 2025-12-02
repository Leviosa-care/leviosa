package bookingRepository_test

import (
	"context"
	"testing"

	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/booking"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetBookingsByAvailability TEST_PATH=internal/booking/infrastructure/postgres/booking/get_bookings_by_availability_test.go

func TestGetBookingsByAvailability(t *testing.T) {
	ctx := context.Background()

	t.Run("should retrieve all bookings for availability", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		availabilityID := uuid.New()
		otherAvailabilityID := uuid.New()

		// Insert bookings for the availability
		booking1 := tb.NewTestBookingEncxWithIDs(t, availabilityID, uuid.New(), uuid.New(), uuid.New())
		booking2 := tb.NewTestBookingEncxWithIDs(t, availabilityID, uuid.New(), uuid.New(), uuid.New())

		// Insert booking for different availability
		otherBooking := tb.NewTestBookingEncxWithIDs(t, otherAvailabilityID, uuid.New(), uuid.New(), uuid.New())

		err := tb.InsertBookingEncx(t, ctx, testPool, booking1)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, booking2)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, otherBooking)
		require.NoError(t, err)

		// Get bookings for availability
		bookings, err := repo.GetBookingsByAvailability(ctx, availabilityID)
		require.NoError(t, err)
		assert.Len(t, bookings, 2)

		// Verify all bookings belong to the availability
		for _, b := range bookings {
			assert.Equal(t, availabilityID, b.AvailabilityID)
		}
	})

	t.Run("should return empty slice when no bookings exist", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		nonExistentAvailabilityID := uuid.New()

		bookings, err := repo.GetBookingsByAvailability(ctx, nonExistentAvailabilityID)
		require.NoError(t, err)
		assert.NotNil(t, bookings)
		assert.Empty(t, bookings)
	})

	t.Run("should sort bookings by slot start time", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		availabilityID := uuid.New()

		// Note: In real implementation, slot_start_time is encrypted
		// The sorting happens on the encrypted field
		// This test just verifies the query returns multiple bookings in order

		// Insert bookings (encrypted times will be compared lexicographically)
		booking1 := tb.NewTestBookingEncxWithIDs(t, availabilityID, uuid.New(), uuid.New(), uuid.New())
		booking1.SlotStartTimeEncrypted = []byte("encrypted_2024_01_01_10_00")

		booking2 := tb.NewTestBookingEncxWithIDs(t, availabilityID, uuid.New(), uuid.New(), uuid.New())
		booking2.SlotStartTimeEncrypted = []byte("encrypted_2024_01_01_11_00")

		booking3 := tb.NewTestBookingEncxWithIDs(t, availabilityID, uuid.New(), uuid.New(), uuid.New())
		booking3.SlotStartTimeEncrypted = []byte("encrypted_2024_01_01_09_00")

		err := tb.InsertBookingEncx(t, ctx, testPool, booking1)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, booking2)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, booking3)
		require.NoError(t, err)

		// Get bookings (should be ordered by slot_start_time_encrypted ASC)
		bookings, err := repo.GetBookingsByAvailability(ctx, availabilityID)
		require.NoError(t, err)
		assert.Len(t, bookings, 3)

		// Verify order (lexicographic on encrypted field)
		assert.Equal(t, booking3.ID, bookings[0].ID) // 09:00
		assert.Equal(t, booking1.ID, bookings[1].ID) // 10:00
		assert.Equal(t, booking2.ID, bookings[2].ID) // 11:00
	})

	t.Run("should retrieve bookings with all fields for overlap detection", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		availabilityID := uuid.New()

		// Insert booking with complete slot information
		booking := tb.NewTestBookingEncxWithIDs(t, availabilityID, uuid.New(), uuid.New(), uuid.New())
		booking.ProductIDEncrypted = []byte("encrypted_product_id")
		booking.SlotStartTimeEncrypted = []byte("encrypted_2024_01_01_10_00")
		booking.SlotEndTimeEncrypted = []byte("encrypted_2024_01_01_11_00")

		err := tb.InsertBookingEncx(t, ctx, testPool, booking)
		require.NoError(t, err)

		// Retrieve for overlap detection
		bookings, err := repo.GetBookingsByAvailability(ctx, availabilityID)
		require.NoError(t, err)
		require.Len(t, bookings, 1)

		// Verify slot fields are present (used for overlap detection)
		assert.NotEmpty(t, bookings[0].ProductIDEncrypted)
		assert.NotEmpty(t, bookings[0].SlotStartTimeEncrypted)
		assert.NotEmpty(t, bookings[0].SlotEndTimeEncrypted)
	})

	t.Run("should handle multiple bookings in slot-based system", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		availabilityID := uuid.New()

		// Insert 5 bookings for the same availability (slot-based system)
		for i := 0; i < 5; i++ {
			booking := tb.NewTestBookingEncxWithIDs(t, availabilityID, uuid.New(), uuid.New(), uuid.New())
			err := tb.InsertBookingEncx(t, ctx, testPool, booking)
			require.NoError(t, err)
		}

		// Get all bookings
		bookings, err := repo.GetBookingsByAvailability(ctx, availabilityID)
		require.NoError(t, err)
		assert.Len(t, bookings, 5)
	})

	t.Run("should include cancelled and completed bookings", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		availabilityID := uuid.New()

		// Insert bookings with different statuses
		confirmed := tb.NewTestBookingEncxWithIDs(t, availabilityID, uuid.New(), uuid.New(), uuid.New())
		completed := tb.NewCompletedBookingEncx(t, availabilityID, uuid.New(), uuid.New(), uuid.New())
		cancelled := tb.NewCancelledBookingEncx(t, availabilityID, uuid.New(), uuid.New(), uuid.New(), "test")

		err := tb.InsertBookingEncx(t, ctx, testPool, confirmed)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, completed)
		require.NoError(t, err)
		err = tb.InsertBookingEncx(t, ctx, testPool, cancelled)
		require.NoError(t, err)

		// Get all bookings (no status filter)
		bookings, err := repo.GetBookingsByAvailability(ctx, availabilityID)
		require.NoError(t, err)
		assert.Len(t, bookings, 3, "Should include all bookings regardless of status")
	})
}
