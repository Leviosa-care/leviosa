package booking

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

// HandlePaymentWebhook processes a Stripe payment webhook event and updates the booking status.
//
// This method handles the following webhook event types:
//   - payment_intent.succeeded: Marks the booking payment as paid
//   - payment_intent.payment_failed: Marks the booking payment as failed
//   - payment_intent.canceled: Marks the booking payment as failed and cancels the booking
//   - charge.refunded: Marks the booking payment as refunded
//
// The booking is identified via the metadata attached to the payment intent during creation.
func (s *BookingService) HandlePaymentWebhook(ctx context.Context, event *ports.WebhookEvent) error {
	// Extract booking ID from payment intent metadata
	bookingIDStr, ok := event.Metadata["booking_id"]
	if !ok || bookingIDStr == "" {
		// Payment intent not associated with a booking (might be a test or external payment)
		return nil
	}

	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		slog.ErrorContext(ctx, "invalid booking ID in webhook metadata",
			"stripe_event_id", event.ID,
			"stripe_event_type", event.Type,
			"error", err,
		)
		return errs.NewInvalidValueErr(fmt.Sprintf("invalid booking ID in webhook metadata: %s", bookingIDStr))
	}

	// Fetch the booking
	bookingEncx, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			// Booking might have been deleted, log and continue
			slog.ErrorContext(ctx, "booking not found for webhook",
				"booking_id", bookingID,
				"stripe_event_id", event.ID,
				"stripe_event_type", event.Type,
				"error", err,
			)
			return fmt.Errorf("booking not found for webhook: %w", errs.ErrRepositoryNotFound)
		}
		slog.ErrorContext(ctx, "failed to get booking for webhook",
			"booking_id", bookingID,
			"stripe_event_id", event.ID,
			"stripe_event_type", event.Type,
			"error", err,
		)
		return fmt.Errorf("get booking for webhook: %w", err)
	}

	// Decrypt booking
	booking, err := domain.DecryptBookingEncx(ctx, s.crypto, bookingEncx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to decrypt booking for webhook",
			"booking_id", bookingID,
			"stripe_event_id", event.ID,
			"stripe_event_type", event.Type,
			"error", err,
		)
		return fmt.Errorf("decrypt booking for webhook: %w", err)
	}

	// Verify payment intent ID matches
	if booking.PaymentIntentID == nil || *booking.PaymentIntentID != event.PaymentIntentID {
		return errs.NewInvalidValueErr("payment intent ID mismatch")
	}

	// Process based on event type
	var updated bool
	var sendPaymentConfirmation bool
	var sendPaymentFailed bool

	switch event.Type {
	case ports.WebhookEventPaymentIntentSucceeded:
		if booking.PaymentStatus != domain.PaymentStatusPaid {
			booking.MarkPaymentPaid()
			updated = true
			sendPaymentConfirmation = true
		}

	case ports.WebhookEventPaymentIntentPaymentFailed:
		if booking.PaymentStatus != domain.PaymentStatusFailed {
			booking.MarkPaymentFailed()
			updated = true
			sendPaymentFailed = true
		}

	case ports.WebhookEventPaymentIntentCanceled:
		if booking.PaymentStatus != domain.PaymentStatusFailed {
			booking.MarkPaymentFailed()
			updated = true
			sendPaymentFailed = true
		}
		// Also cancel the booking if payment was canceled
		if booking.IsCancellable() {
			reason := "Payment was canceled"
			if event.FailureMessage != "" {
				reason = fmt.Sprintf("Payment canceled: %s", event.FailureMessage)
			}
			if err := booking.Cancel(reason); err != nil {
				return fmt.Errorf("cancel booking after payment cancellation: %w", err)
			}
		}

	case ports.WebhookEventChargeRefunded:
		if booking.PaymentStatus != domain.PaymentStatusRefunded {
			booking.RefundPayment()
			updated = true
		}

	default:
		// Unknown event type, ignore
		return nil
	}

	// Persist changes if any updates were made
	if updated || booking.Status == domain.BookingStatusCancelled {
		bookingEncx, err = domain.ProcessBookingEncx(ctx, s.crypto, booking)
		if err != nil {
			return fmt.Errorf("encrypt booking after webhook: %w", err)
		}

		if err := s.bookingRepo.Update(ctx, bookingEncx); err != nil {
			return fmt.Errorf("update booking after webhook: %w", err)
		}
	}

	// Send payment notifications (best effort)
	if s.notificationService != nil {
		notificationData := s.buildNotificationData(booking, "")

		if sendPaymentConfirmation {
			if err := s.notificationService.SendPaymentConfirmation(ctx, notificationData); err != nil {
				// Log error but don't fail the webhook processing
				slog.ErrorContext(ctx, "failed to send payment confirmation notification",
					"booking_id", booking.ID,
					"stripe_event_id", event.ID,
					"stripe_event_type", event.Type,
					"error", err,
				)
			}
		}

		if sendPaymentFailed {
			if err := s.notificationService.SendPaymentFailed(ctx, notificationData); err != nil {
				// Log error but don't fail the webhook processing
				slog.ErrorContext(ctx, "failed to send payment failed notification",
					"booking_id", booking.ID,
					"stripe_event_id", event.ID,
					"stripe_event_type", event.Type,
					"error", err,
				)
			}
		}

		// Send cancellation notification if booking was cancelled due to payment failure
		if booking.Status == domain.BookingStatusCancelled {
			if err := s.notificationService.SendBookingCancellation(ctx, notificationData); err != nil {
				slog.ErrorContext(ctx, "failed to send booking cancellation notification",
					"booking_id", booking.ID,
					"stripe_event_id", event.ID,
					"stripe_event_type", event.Type,
					"error", err,
				)
			}
		}
	}

	return nil
}
