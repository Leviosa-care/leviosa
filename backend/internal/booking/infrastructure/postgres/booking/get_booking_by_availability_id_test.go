package bookingRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/booking"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetByAvailabilityID TEST_PATH=internal/booking/infrastructure/postgres/booking/get_booking_by_availability_id_test.go

func TestGetByAvailabilityID(t *testing.T) {
	ctx := context.Background()

	t.Run("should retrieve booking for availability", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		availabilityID := uuid.New()

		// Insert booking
		booking := tb.NewTestBookingEncxWithIDs(t, availabilityID, uuid.New(), uuid.New(), uuid.New())
		err := tb.InsertBookingEncx(t, ctx, testPool, booking)
		require.NoError(t, err)

		// Retrieve by availability ID
		retrieved, err := repo.GetByAvailabilityID(ctx, availabilityID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, booking.ID, retrieved.ID)
		assert.Equal(t, availabilityID, retrieved.AvailabilityID)
	})

	t.Run("should return error when no booking exists for availability", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		nonExistentAvailabilityID := uuid.New()

		retrieved, err := repo.GetByAvailabilityID(ctx, nonExistentAvailabilityID)
		assert.Error(t, err)
		assert.Nil(t, retrieved)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("should return most recent booking when multiple exist", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		availabilityID := uuid.New()
		now := time.Now()

		// Insert older booking
		olderBooking := tb.NewTestBookingEncxWithIDs(t, availabilityID, uuid.New(), uuid.New(), uuid.New())
		olderBooking.CreatedAt = now.Add(-2 * time.Hour)
		err := tb.InsertBookingEncx(t, ctx, testPool, olderBooking)
		require.NoError(t, err)

		// Insert more recent booking
		newerBooking := tb.NewTestBookingEncxWithIDs(t, availabilityID, uuid.New(), uuid.New(), uuid.New())
		newerBooking.CreatedAt = now.Add(-1 * time.Hour)
		err = tb.InsertBookingEncx(t, ctx, testPool, newerBooking)
		require.NoError(t, err)

		// Should return the newer booking (created_at DESC, LIMIT 1)
		retrieved, err := repo.GetByAvailabilityID(ctx, availabilityID)
		require.NoError(t, err)
		assert.Equal(t, newerBooking.ID, retrieved.ID)
	})

	t.Run("should retrieve booking with all fields", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		availabilityID := uuid.New()

		// Insert booking with complete data
		booking := tb.NewTestBookingEncxWithPaymentIntent(
			t,
			availabilityID,
			uuid.New(),
			uuid.New(),
			uuid.New(),
			"pi_test_availability",
		)
		booking.ClientNotesEncrypted = []byte("encrypted_client_notes")
		booking.PartnerNotesEncrypted = []byte("encrypted_partner_notes")

		err := tb.InsertBookingEncx(t, ctx, testPool, booking)
		require.NoError(t, err)

		// Retrieve and verify all fields
		retrieved, err := repo.GetByAvailabilityID(ctx, availabilityID)
		require.NoError(t, err)
		assert.Equal(t, booking.ID, retrieved.ID)
		assert.Equal(t, booking.ProductIDEncrypted, retrieved.ProductIDEncrypted)
		assert.Equal(t, booking.SlotStartTimeEncrypted, retrieved.SlotStartTimeEncrypted)
		assert.Equal(t, booking.ClientNotesEncrypted, retrieved.ClientNotesEncrypted)
		assert.Equal(t, booking.PartnerNotesEncrypted, retrieved.PartnerNotesEncrypted)
		require.NotNil(t, retrieved.PaymentIntentID)
		assert.Equal(t, "pi_test_availability", *retrieved.PaymentIntentID)
	})
}
