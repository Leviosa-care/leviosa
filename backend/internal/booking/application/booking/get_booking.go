package booking

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (s *BookingService) GetBooking(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	bookingEncx, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("get booking: %w", err)
	}

	booking, err := domain.DecryptBookingEncx(ctx, s.crypto, bookingEncx)
	if err != nil {
		return nil, fmt.Errorf("decrypt booking: %w", err)
	}

	return booking, nil
}

