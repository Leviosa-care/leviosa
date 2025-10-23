package availabilityRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/booking/internal/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Availability, error) {
	query := fmt.Sprintf(`
		SELECT
			id, partner_id, room_id, start_time, end_time,
			service_type_encrypted, price_cents, max_capacity,
			notes_encrypted, is_recurring, recurrence_pattern_encrypted,
			status, created_at, updated_at,
			dek_encrypted, key_version, metadata
		FROM %s.availabilities
		WHERE id = $1
	`, r.schema)

	availabilityEncx := &domain.AvailabilityEncx{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&availabilityEncx.ID,
		&availabilityEncx.PartnerID,
		&availabilityEncx.RoomID,
		&availabilityEncx.StartTime,
		&availabilityEncx.EndTime,
		&availabilityEncx.ServiceTypeEncrypted,
		&availabilityEncx.PriceCents,
		&availabilityEncx.MaxCapacity,
		&availabilityEncx.NotesEncrypted,
		&availabilityEncx.IsRecurring,
		&availabilityEncx.RecurrencePatternEncrypted,
		&availabilityEncx.Status,
		&availabilityEncx.CreatedAt,
		&availabilityEncx.UpdatedAt,
		&availabilityEncx.DEKEncrypted,
		&availabilityEncx.KeyVersion,
		&availabilityEncx.Metadata,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get availability by id", err)
	}

	// Decrypt sensitive fields using ENCX
	availability, err := domain.DecryptAvailabilityEncx(ctx, r.crypto, availabilityEncx)
	if err != nil {
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