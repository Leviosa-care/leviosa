package booking

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/google/uuid"
)

func (s *BookingService) GetClientBookings(ctx context.Context, clientID uuid.UUID, filter ports.BookingFilter) ([]*domain.Booking, error) {
	bookings, err := s.bookingRepo.GetByClientID(ctx, clientID, filter)
	if err != nil {
		return nil, fmt.Errorf("get client bookings: %w", err)
	}

	return bookings, nil
}
