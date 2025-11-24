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
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get booking for cancellation: %w", err)
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
	availability, err := s.availabilityRepo.GetByID(ctx, booking.AvailabilityID)
	if err == nil {
		availability.MarkAsAvailable()
		s.availabilityRepo.Update(ctx, availability) // Best effort, don't fail booking cancellation
	}

	// Persist changes
	if err := s.bookingRepo.Update(ctx, booking); err != nil {
		return nil, fmt.Errorf("update cancelled booking: %w", err)
	}

	return booking, nil
}