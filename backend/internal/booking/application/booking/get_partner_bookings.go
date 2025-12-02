package booking

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/google/uuid"
)

func (s *BookingService) GetPartnerBookings(ctx context.Context, partnerID uuid.UUID, filter ports.BookingFilter) ([]*domain.Booking, error) {
	bookingsEncx, err := s.bookingRepo.GetByPartnerID(ctx, partnerID, filter)
	if err != nil {
		return nil, fmt.Errorf("get partner bookings: %w", err)
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
