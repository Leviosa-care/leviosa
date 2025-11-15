package availabilityRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) Create(ctx context.Context, availability *domain.AvailabilityEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.availabilities (
			id, partner_id, room_id, start_time, end_time,
			service_type_encrypted, price_cents, max_capacity,
			notes_encrypted, is_recurring, recurrence_pattern_encrypted,
			status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		availability.ID,
		availability.PartnerID,
		availability.RoomID,
		availability.StartTime,
		availability.EndTime,
		availability.ServiceTypeEncrypted,
		availability.PriceCents,
		availability.MaxCapacity,
		availability.NotesEncrypted,
		availability.IsRecurring,
		availability.RecurrencePatternEncrypted,
		availability.Status,
		availability.CreatedAt,
		availability.UpdatedAt,
	)
	if err != nil {
		return errs.ClassifyPgError("create availability", err)
	}

	return nil
}
