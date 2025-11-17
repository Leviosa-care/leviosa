package availabilityRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) Create(ctx context.Context, availabilityEncx *domain.AvailabilityEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.availabilities (
			id, user_id, room_id, start_time, end_time,
			service_type_encrypted, price_cents, max_capacity,
			notes_encrypted, is_recurring, recurrence_pattern_encrypted,
			status, created_at, updated_at, dek_encrypted, key_version, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		availabilityEncx.ID,
		availabilityEncx.UserID,
		availabilityEncx.RoomID,
		availabilityEncx.StartTime,
		availabilityEncx.EndTime,
		availabilityEncx.ServiceTypeEncrypted,
		availabilityEncx.PriceCents,
		availabilityEncx.MaxCapacity,
		availabilityEncx.NotesEncrypted,
		availabilityEncx.IsRecurring,
		availabilityEncx.RecurrencePatternEncrypted,
		availabilityEncx.Status,
		availabilityEncx.CreatedAt,
		availabilityEncx.UpdatedAt,
		availabilityEncx.DEKEncrypted,
		availabilityEncx.KeyVersion,
		availabilityEncx.Metadata,
	)
	if err != nil {
		return errs.ClassifyPgError("create availability", err)
	}

	return nil
}
