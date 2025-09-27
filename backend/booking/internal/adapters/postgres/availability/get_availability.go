package availabilityRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/booking/internal/ports"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Availability, error) {
	query := fmt.Sprintf(`
		SELECT
			id, partner_id, room_id, start_time, end_time,
			service_type_encrypted, price_cents, max_capacity,
			notes_encrypted, is_recurring, recurrence_pattern_encrypted,
			status, created_at, updated_at
		FROM %s.availabilities
		WHERE id = $1
	`, r.schema)

	availability := &domain.Availability{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&availability.ID,
		&availability.PartnerID,
		&availability.RoomID,
		&availability.StartTime,
		&availability.EndTime,
		&availability.ServiceTypeEncrypted,
		&availability.PriceCents,
		&availability.MaxCapacity,
		&availability.NotesEncrypted,
		&availability.IsRecurring,
		&availability.RecurrencePatternEncrypted,
		&availability.Status,
		&availability.CreatedAt,
		&availability.UpdatedAt,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get availability by id", err)
	}

	// Decrypt sensitive fields
	if err := r.crypto.DecryptStruct(ctx, availability); err != nil {
		return nil, fmt.Errorf("decrypt availability data: %w", err)
	}

	return availability, nil
}

func (r *Repository) GetByPartnerID(ctx context.Context, partnerID uuid.UUID, filter ports.AvailabilityFilter) ([]*domain.Availability, error) {
	// Set partner filter
	filter.PartnerID = &partnerID
	return r.List(ctx, filter)
}

func (r *Repository) GetByRoomID(ctx context.Context, roomID uuid.UUID, filter ports.AvailabilityFilter) ([]*domain.Availability, error) {
	// Set room filter
	filter.RoomID = &roomID
	return r.List(ctx, filter)
}

func (r *Repository) GetAvailableSlots(ctx context.Context, filter ports.AvailabilityFilter) ([]*domain.Availability, error) {
	// Force available status filter
	filter.Status = []domain.AvailabilityStatus{domain.AvailabilityStatusAvailable}
	filter.AvailableOnly = true
	return r.List(ctx, filter)
}