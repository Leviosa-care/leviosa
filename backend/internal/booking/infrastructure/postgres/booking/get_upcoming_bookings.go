package bookingRepository

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
)

func (r *Repository) GetUpcoming(ctx context.Context, filter ports.BookingFilter) ([]*domain.BookingEncx, error) {
	// Force confirmed status and join with availabilities to filter by future start times
	filter.Status = []domain.BookingStatus{domain.BookingStatusConfirmed}

	// This would require a join with availabilities table for more complex filtering
	// For now, return all confirmed bookings and let the service layer filter by time
	return r.List(ctx, filter)
}