package booking

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
)

func (s *BookingService) GetUpcomingBookings(ctx context.Context, filter ports.BookingFilter) ([]*domain.Booking, error) {
	bookings, err := s.bookingRepo.GetUpcoming(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get upcoming bookings: %w", err)
	}

	return bookings, nil
}