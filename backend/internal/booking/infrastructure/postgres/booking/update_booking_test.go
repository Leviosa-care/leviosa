package bookingRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/booking"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestUpdateBooking TEST_PATH=internal/booking/infrastructure/postgres/booking/update_booking_test.go

func TestUpdateBooking(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully update existing booking", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert original booking
		bookingEncx := tb.NewTestBookingEncx(t)
		err := tb.InsertBookingEncx(t, ctx, testPool, bookingEncx)
		require.NoError(t, err)

		// Modify booking
		bookingEncx.TotalPriceCents = 7500 // Change price
		bookingEncx.ClientNotesEncrypted = []byte("encrypted_updated_notes")
		bookingEncx.UpdatedAt = time.Now()

		// Update
		err = repo.Update(ctx, bookingEncx)
		require.NoError(t, err)

		// Verify changes
		updated, err := tb.GetBookingEncxByID(t, ctx, testPool, bookingEncx.ID)
		require.NoError(t, err)
		assert.Equal(t, 7500, updated.TotalPriceCents)
		assert.Equal(t, bookingEncx.ClientNotesEncrypted, updated.ClientNotesEncrypted)
	})

	t.Run("should return error for non-existent booking", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		nonExistentBooking := tb.NewTestBookingEncx(t)
		err := repo.Update(ctx, nonExistentBooking)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("should update booking status to completed", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert confirmed booking
		bookingEncx := tb.NewTestBookingEncx(t)
		err := tb.InsertBookingEncx(t, ctx, testPool, bookingEncx)
		require.NoError(t, err)

		// Update to completed
		now := time.Now()
		bookingEncx.Status = domain.BookingStatusCompleted
		bookingEncx.CompletedAt = &now
		bookingEncx.PaymentStatus = domain.PaymentStatusPaid
		bookingEncx.UpdatedAt = now

		err = repo.Update(ctx, bookingEncx)
		require.NoError(t, err)

		// Verify
		updated, err := tb.GetBookingEncxByID(t, ctx, testPool, bookingEncx.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.BookingStatusCompleted, updated.Status)
		require.NotNil(t, updated.CompletedAt)
		assert.Equal(t, domain.PaymentStatusPaid, updated.PaymentStatus)
	})

	t.Run("should update booking status to cancelled", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert confirmed booking
		bookingEncx := tb.NewTestBookingEncx(t)
		err := tb.InsertBookingEncx(t, ctx, testPool, bookingEncx)
		require.NoError(t, err)

		// Update to cancelled
		now := time.Now()
		bookingEncx.Status = domain.BookingStatusCancelled
		bookingEncx.CancelledAt = &now
		bookingEncx.CancellationReasonEncrypted = []byte("encrypted_cancellation_reason")
		bookingEncx.UpdatedAt = now

		err = repo.Update(ctx, bookingEncx)
		require.NoError(t, err)

		// Verify
		updated, err := tb.GetBookingEncxByID(t, ctx, testPool, bookingEncx.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.BookingStatusCancelled, updated.Status)
		require.NotNil(t, updated.CancelledAt)
		assert.NotEmpty(t, updated.CancellationReasonEncrypted)
	})

	t.Run("should update payment status", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert booking with pending payment
		bookingEncx := tb.NewTestBookingEncx(t)
		bookingEncx.PaymentStatus = domain.PaymentStatusPending
		err := tb.InsertBookingEncx(t, ctx, testPool, bookingEncx)
		require.NoError(t, err)

		// Update to paid
		bookingEncx.PaymentStatus = domain.PaymentStatusPaid
		bookingEncx.UpdatedAt = time.Now()

		err = repo.Update(ctx, bookingEncx)
		require.NoError(t, err)

		// Verify
		updated, err := tb.GetBookingEncxByID(t, ctx, testPool, bookingEncx.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.PaymentStatusPaid, updated.PaymentStatus)
	})

	t.Run("should update partner notes", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert booking
		bookingEncx := tb.NewTestBookingEncx(t)
		err := tb.InsertBookingEncx(t, ctx, testPool, bookingEncx)
		require.NoError(t, err)

		// Update partner notes
		bookingEncx.PartnerNotesEncrypted = []byte("encrypted_new_partner_notes")
		bookingEncx.UpdatedAt = time.Now()

		err = repo.Update(ctx, bookingEncx)
		require.NoError(t, err)

		// Verify
		updated, err := tb.GetBookingEncxByID(t, ctx, testPool, bookingEncx.ID)
		require.NoError(t, err)
		assert.Equal(t, bookingEncx.PartnerNotesEncrypted, updated.PartnerNotesEncrypted)
	})

	t.Run("should update payment intent ID", func(t *testing.T) {
		tb.ClearBookingsTable(t, ctx, testPool)

		// Insert booking without payment intent
		bookingEncx := tb.NewTestBookingEncx(t)
		bookingEncx.PaymentIntentID = nil
		err := tb.InsertBookingEncx(t, ctx, testPool, bookingEncx)
		require.NoError(t, err)

		// Update with payment intent
		newPaymentIntentID := "pi_updated_" + uuid.New().String()[:8]
		bookingEncx.PaymentIntentID = &newPaymentIntentID
		bookingEncx.UpdatedAt = time.Now()

		err = repo.Update(ctx, bookingEncx)
		require.NoError(t, err)

		// Verify
		updated, err := tb.GetBookingEncxByID(t, ctx, testPool, bookingEncx.ID)
		require.NoError(t, err)
		require.NotNil(t, updated.PaymentIntentID)
		assert.Equal(t, newPaymentIntentID, *updated.PaymentIntentID)
	})
}
