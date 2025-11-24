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

func (s *BookingService) ProcessPayment(ctx context.Context, id uuid.UUID, paymentIntentID string) (*domain.Booking, error) {
	// Get existing booking
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get booking for payment: %w", err)
	}

	// Verify payment intent status with Stripe
	paymentInfo, err := s.paymentService.RetrievePaymentIntent(ctx, paymentIntentID)
	if err != nil {
		return nil, fmt.Errorf("retrieve payment intent: %w", err)
	}

	// Set payment intent ID if not already set
	if booking.PaymentIntentID == nil || *booking.PaymentIntentID != paymentIntentID {
		booking.SetPaymentIntentID(paymentIntentID)
	}

	// Update payment status based on Stripe payment intent status
	switch paymentInfo.Status {
	case ports.PaymentIntentStatusSucceeded:
		if err := booking.MarkPaymentPaid(); err != nil {
			return nil, fmt.Errorf("mark payment as paid: %w", err)
		}
	case ports.PaymentIntentStatusRequiresPaymentMethod,
		ports.PaymentIntentStatusRequiresConfirmation,
		ports.PaymentIntentStatusRequiresAction,
		ports.PaymentIntentStatusProcessing:
		// Payment is still pending, no action needed
	case ports.PaymentIntentStatusCanceled:
		if err := booking.MarkPaymentFailed(); err != nil {
			return nil, fmt.Errorf("mark payment as failed: %w", err)
		}
	case ports.PaymentIntentStatusPaymentFailed:
		if err := booking.MarkPaymentFailed(); err != nil {
			return nil, fmt.Errorf("mark payment as failed: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown payment status: %s", paymentInfo.Status)
	}

	// Persist changes
	if err := s.bookingRepo.Update(ctx, booking); err != nil {
		return nil, fmt.Errorf("update booking payment: %w", err)
	}

	return booking, nil
}