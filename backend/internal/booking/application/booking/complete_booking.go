package booking

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *BookingService) CompleteBooking(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	// Get existing booking
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get booking for completion: %w", err)
	}

	// Complete booking
	if err := booking.Complete(); err != nil {
		return nil, fmt.Errorf("complete booking: %w", err)
	}

	// Persist changes
	if err := s.bookingRepo.Update(ctx, booking); err != nil {
		return nil, fmt.Errorf("update completed booking: %w", err)
	}

	return booking, nil
}

