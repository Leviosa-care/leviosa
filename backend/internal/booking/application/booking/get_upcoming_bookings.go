package booking

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
)

func (s *BookingService) GetUpcomingBookings(ctx context.Context, filter ports.BookingFilter) ([]*domain.Booking, error) {
	bookingsEncx, err := s.bookingRepo.GetUpcoming(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get upcoming bookings: %w", err)
	}

	// Decrypt each booking
	bookings := make([]*domain.Booking, 0, len(bookingsEncx))
	for _, bookingEncx := range bookingsEncx {
		booking, err := domain.DecryptBookingEncx(ctx, s.crypto, bookingEncx)
		if err != nil {
			return nil, fmt.Errorf("decrypt booking %s: %w", bookingEncx.ID, err)
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}
