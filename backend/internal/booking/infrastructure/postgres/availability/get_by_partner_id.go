package availabilityRepository

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/google/uuid"
)

func (r *Repository) GetByPartnerID(ctx context.Context, partnerID uuid.UUID, filter ports.AvailabilityFilter) ([]*domain.AvailabilityEncx, error) {
	// Set partner filter
	filter.UserID = &partnerID
	return r.List(ctx, filter)
}
