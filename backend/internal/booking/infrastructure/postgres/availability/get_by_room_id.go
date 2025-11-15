package availabilityRepository

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/google/uuid"
)

func (r *Repository) GetByRoomID(ctx context.Context, roomID uuid.UUID, filter ports.AvailabilityFilter) ([]*domain.AvailabilityEncx, error) {
	// Set room filter
	filter.RoomID = &roomID
	return r.List(ctx, filter)
}
