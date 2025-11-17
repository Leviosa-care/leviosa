package availabilityRepository

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
)

func (r *Repository) GetAvailableSlots(ctx context.Context, filter ports.AvailabilityFilter) ([]*domain.AvailabilityEncx, error) {
	// Force available status filter
	filter.Status = []domain.AvailabilityStatus{domain.AvailabilityStatusAvailable}
	filter.AvailableOnly = true
	return r.List(ctx, filter)
}
