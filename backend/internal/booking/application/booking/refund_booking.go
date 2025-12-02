package booking

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *BookingService) RefundBooking(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	// Get existing booking
	bookingEncx, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get booking for refund: %w", err)
	}

	// Decrypt booking
	booking, err := domain.DecryptBookingEncx(ctx, s.crypto, bookingEncx)
	if err != nil {
		return nil, fmt.Errorf("decrypt booking: %w", err)
	}

	// Validate booking can be refunded
	if booking.PaymentStatus != domain.PaymentStatusPaid {
		return nil, fmt.Errorf("booking payment must be paid to process refund")
	}

	if booking.PaymentIntentID == nil {
		return nil, fmt.Errorf("booking has no payment intent ID")
	}

	// Process Stripe refund
	refundID, err := s.paymentService.RefundPayment(
		ctx,
		*booking.PaymentIntentID,
		0, // 0 = full refund
		ports.RefundReasonRequestedByCustomer,
	)
	if err != nil {
		return nil, fmt.Errorf("process stripe refund: %w", err)
	}

	// Mark booking as refunded
	if err := booking.RefundPayment(); err != nil {
		return nil, fmt.Errorf("refund booking payment: %w", err)
	}

	// Store refund ID in booking notes for tracking
	notes := fmt.Sprintf("Refund processed: %s", refundID)
	booking.SetPartnerNotes(notes)

	// Encrypt and persist changes
	bookingEncx, err = domain.ProcessBookingEncx(ctx, s.crypto, booking)
	if err != nil {
		return nil, fmt.Errorf("encrypt booking: %w", err)
	}

	if err := s.bookingRepo.Update(ctx, bookingEncx); err != nil {
		return nil, fmt.Errorf("update refunded booking: %w", err)
	}

	return booking, nil
}

