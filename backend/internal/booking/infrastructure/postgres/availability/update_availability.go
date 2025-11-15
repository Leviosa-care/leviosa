package availabilityRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) Update(ctx context.Context, availability *domain.Availability) error {
	// Encrypt sensitive fields using ENCX
	availabilityEncx, err := domain.ProcessAvailabilityEncx(ctx, r.crypto, availability)
	if err != nil {
		return fmt.Errorf("encrypt availability data: %w", err)
	}

	query := fmt.Sprintf(`
		UPDATE %s.availabilities SET
			partner_id = $2,
			room_id = $3,
			start_time = $4,
			end_time = $5,
			service_type_encrypted = $6,
			price_cents = $7,
			max_capacity = $8,
			notes_encrypted = $9,
			is_recurring = $10,
			recurrence_pattern_encrypted = $11,
			status = $12,
			updated_at = $13,
			dek_encrypted = $14,
			key_version = $15,
			metadata = $16
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		availability.ID,
		availability.PartnerID,
		availability.RoomID,
		availability.StartTime,
		availability.EndTime,
		availabilityEncx.ServiceTypeEncrypted,
		availability.PriceCents,
		availability.MaxCapacity,
		availabilityEncx.NotesEncrypted,
		availability.IsRecurring,
		availabilityEncx.RecurrencePatternEncrypted,
		availability.Status,
		availability.UpdatedAt,
		availabilityEncx.DEKEncrypted,
		availabilityEncx.KeyVersion,
		availabilityEncx.Metadata,
	)
	if err != nil {
		return errs.ClassifyPgError("update availability", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}