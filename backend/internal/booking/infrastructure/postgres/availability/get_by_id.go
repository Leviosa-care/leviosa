package availabilityRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.AvailabilityEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, user_id, room_id, start_time, end_time,
			service_type_encrypted, price_cents, max_capacity,
			notes_encrypted, is_recurring, recurrence_pattern,
			status, created_at, updated_at,
			dek_encrypted, key_version, metadata
		FROM %s.availabilities
		WHERE id = $1
	`, r.schema)

	availabilityEncx := &domain.AvailabilityEncx{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&availabilityEncx.ID,
		&availabilityEncx.UserID,
		&availabilityEncx.RoomID,
		&availabilityEncx.StartTime,
		&availabilityEncx.EndTime,
		&availabilityEncx.ServiceTypeEncrypted,
		&availabilityEncx.PriceCents,
		&availabilityEncx.MaxCapacity,
		&availabilityEncx.NotesEncrypted,
		&availabilityEncx.IsRecurring,
		&availabilityEncx.RecurrencePattern,
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

	return availabilityEncx, nil
}

