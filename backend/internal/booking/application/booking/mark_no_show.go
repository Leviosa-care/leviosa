package booking

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *BookingService) MarkNoShow(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	// Get existing booking
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get booking for no-show: %w", err)
	}

	// Mark as no-show
	if err := booking.MarkNoShow(); err != nil {
		return nil, fmt.Errorf("mark booking as no-show: %w", err)
	}

	// Persist changes
	if err := s.bookingRepo.Update(ctx, booking); err != nil {
		return nil, fmt.Errorf("update no-show booking: %w", err)
	}

	return booking, nil
}
