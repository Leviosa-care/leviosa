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

// make test-func TEST_NAME=TestGetByID TEST_PATH=internal/booking/infrastructure/postgres/booking/get_booking_by_id_test.go

func TestGetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should retrieve existing booking", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert test booking
		bookingEncx := tb.NewTestBookingEncx(t)
		err := tb.InsertBookingEncx(t, ctx, testPool, bookingEncx)
		require.NoError(t, err)

		// Retrieve booking
		retrieved, err := repo.GetByID(ctx, bookingEncx.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		// Verify all fields
		assert.Equal(t, bookingEncx.ID, retrieved.ID)
		assert.Equal(t, bookingEncx.AvailabilityID, retrieved.AvailabilityID)
		assert.Equal(t, bookingEncx.ClientID, retrieved.ClientID)
		assert.Equal(t, bookingEncx.PartnerID, retrieved.PartnerID)
		assert.Equal(t, bookingEncx.RoomID, retrieved.RoomID)
		assert.Equal(t, bookingEncx.ProductIDEncrypted, retrieved.ProductIDEncrypted)
		assert.Equal(t, bookingEncx.SlotStartTimeEncrypted, retrieved.SlotStartTimeEncrypted)
		assert.Equal(t, bookingEncx.SlotEndTimeEncrypted, retrieved.SlotEndTimeEncrypted)
		assert.Equal(t, bookingEncx.TotalPriceCents, retrieved.TotalPriceCents)
		assert.Equal(t, bookingEncx.Currency, retrieved.Currency)
		assert.Equal(t, bookingEncx.PaymentStatus, retrieved.PaymentStatus)
		assert.Equal(t, bookingEncx.Status, retrieved.Status)
		assert.Equal(t, bookingEncx.DEKEncrypted, retrieved.DEKEncrypted)
		assert.Equal(t, bookingEncx.KeyVersion, retrieved.KeyVersion)
	})

	t.Run("should return error for non-existent booking", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		nonExistentID := uuid.New()
		retrieved, err := repo.GetByID(ctx, nonExistentID)

		assert.Error(t, err)
		assert.Nil(t, retrieved)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("should retrieve booking with all optional fields", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert booking with payment intent and notes
		paymentIntentID := "pi_test_full_fields"
		bookingEncx := tb.NewTestBookingEncxWithPaymentIntent(
			t,
			uuid.New(), uuid.New(), uuid.New(), uuid.New(),
			paymentIntentID,
		)
		bookingEncx.ClientNotesEncrypted = []byte("encrypted_client_notes")
		bookingEncx.PartnerNotesEncrypted = []byte("encrypted_partner_notes")

		err := tb.InsertBookingEncx(t, ctx, testPool, bookingEncx)
		require.NoError(t, err)

		// Retrieve and verify
		retrieved, err := repo.GetByID(ctx, bookingEncx.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved.PaymentIntentID)
		assert.Equal(t, paymentIntentID, *retrieved.PaymentIntentID)
		assert.Equal(t, bookingEncx.ClientNotesEncrypted, retrieved.ClientNotesEncrypted)
		assert.Equal(t, bookingEncx.PartnerNotesEncrypted, retrieved.PartnerNotesEncrypted)
	})

	t.Run("should retrieve cancelled booking with cancellation details", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert cancelled booking
		cancelledBooking := tb.NewCancelledBookingEncx(
			t,
			uuid.New(), uuid.New(), uuid.New(), uuid.New(),
			"test cancellation reason",
		)
		err := tb.InsertBookingEncx(t, ctx, testPool, cancelledBooking)
		require.NoError(t, err)

		// Retrieve and verify
		retrieved, err := repo.GetByID(ctx, cancelledBooking.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved.CancelledAt)
		assert.NotEmpty(t, retrieved.CancellationReasonEncrypted)
	})

	t.Run("should retrieve completed booking with completion timestamp", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert completed booking
		completedBooking := tb.NewCompletedBookingEncx(
			t,
			uuid.New(), uuid.New(), uuid.New(), uuid.New(),
		)
		err := tb.InsertBookingEncx(t, ctx, testPool, completedBooking)
		require.NoError(t, err)

		// Retrieve and verify
		retrieved, err := repo.GetByID(ctx, completedBooking.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved.CompletedAt)
	})
}
