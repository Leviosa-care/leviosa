package bookingRepository

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/google/uuid"
)

func (r *Repository) GetByClientID(ctx context.Context, clientID uuid.UUID, filter ports.BookingFilter) ([]*domain.BookingEncx, error) {
	filter.ClientID = &clientID
	return r.List(ctx, filter)
}