package booking

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

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

	// Check if booking can be cancelled
	if !booking.IsCancellable() {
		return nil, fmt.Errorf("booking cannot be cancelled")
	}

	// Cancel booking
	if err := booking.Cancel(reason); err != nil {
		return nil, fmt.Errorf("cancel booking: %w", err)
	}

	// Mark associated availability as available again
	availabilityEncx, err := s.availabilityRepo.GetByID(ctx, booking.AvailabilityID)
	if err == nil {
		availability, err := domain.DecryptAvailabilityEncx(ctx, s.crypto, availabilityEncx)
		if err != nil {

		}
		availability.MarkAsAvailable()
		if err := s.availabilityRepo.Update(ctx, availabilityEncx); err != nil { // Best effort, don't fail booking cancellation

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

	return booking, nil
}

