package booking

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

// CancelBookingByToken cancels a booking using a signed booking token instead of
// a session cookie. It verifies the token, ensures it matches the requested booking
// ID, and then delegates to the standard CancelBooking flow.
func (s *BookingService) CancelBookingByToken(ctx context.Context, bookingID uuid.UUID, token string, reason string) (*domain.Booking, error) {
	if len(s.tokenSecret) == 0 {
		return nil, fmt.Errorf("booking token feature is not configured")
	}

	tokenBookingID, err := domain.VerifyBookingToken(token, s.tokenSecret)
	if err != nil {
		if domain.IsBookingTokenError(err) {
			return nil, errs.NewUnauthorizedErr(err.Error())
		}
		return nil, fmt.Errorf("verify booking token for cancel: %w", err)
	}

	// The token must be scoped to the exact booking being cancelled.
	if tokenBookingID != bookingID {
		return nil, errs.NewUnauthorizedErr("token does not match the requested booking")
	}

	return s.CancelBooking(ctx, bookingID, reason)
}

func (s *BookingService) CancelBooking(ctx context.Context, id uuid.UUID, reason string) (*domain.Booking, error) {
	// Get existing booking
	bookingEncx, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get booking for cancellation: %w", err)
	}

	// Decrypt booking
	booking, err := domain.DecryptBookingEncx(ctx, s.crypto, bookingEncx)
	if err != nil {
		return nil, fmt.Errorf("decrypt booking: %w", err)
	}

	// Check if booking can be cancelled (status-based)
	if !booking.IsCancellable() {
		return nil, errs.NewInvalidValueErr("booking cannot be cancelled: current status is " + string(booking.Status))
	}

	// Enforce cancellation policy based on product's CancellationHours
	product, err := s.productService.GetProductByID(ctx, booking.ProductID.String())
	if err != nil {
		// If product lookup fails, allow cancellation (fail open for customer benefit)
		// but log the issue for investigation
		if !errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, fmt.Errorf("get product for cancellation policy: %w", err)
		}
		// Product not found - proceed without policy enforcement
	} else if product.CancellationHours > 0 {
		// Calculate time remaining until appointment
		now := time.Now()
		timeUntilAppointment := booking.SlotStartTime.Sub(now)
		hoursRemaining := int(timeUntilAppointment.Hours())

		// Check if within cancellation window
		if hoursRemaining < product.CancellationHours {
			return nil, errs.NewCancellationWindowClosedErr(product.CancellationHours, hoursRemaining)
		}
	}

	// Cancel booking
	if err := booking.Cancel(reason); err != nil {
		return nil, fmt.Errorf("cancel booking: %w", err)
	}

	// Mark associated availability as available again (best effort)
	if availabilityEncx, err := s.availabilityRepo.GetByID(ctx, booking.AvailabilityID); err == nil {
		if availability, err := domain.DecryptAvailabilityEncx(ctx, s.crypto, availabilityEncx); err == nil {
			availability.MarkAsAvailable()
			if updatedEncx, err := domain.ProcessAvailabilityEncx(ctx, s.crypto, availability); err == nil {
				_ = s.availabilityRepo.Update(ctx, updatedEncx)
			}
		}
	}

	// Encrypt and persist changes
	bookingEncx, err = domain.ProcessBookingEncx(ctx, s.crypto, booking)
	if err != nil {
		return nil, fmt.Errorf("encrypt booking: %w", err)
	}

	if err := s.bookingRepo.Update(ctx, bookingEncx); err != nil {
		return nil, fmt.Errorf("update cancelled booking: %w", err)
	}

	// Send cancellation notification (best effort)
	if s.notificationService != nil {
		productName := ""
		if product != nil {
			productName = product.Name
		}
		notificationData := s.buildNotificationData(booking, productName)
		if err := s.notificationService.SendBookingCancellation(ctx, notificationData); err != nil {
			// Log error but don't fail the cancellation
			_ = err // TODO: Add proper logging
		}
	}

	return booking, nil
}
