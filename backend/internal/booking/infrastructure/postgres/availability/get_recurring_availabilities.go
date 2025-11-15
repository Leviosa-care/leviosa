package availabilityRepository

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) GetRecurringAvailabilities(ctx context.Context, until time.Time) ([]*domain.Availability, error) {
	query := fmt.Sprintf(`
		SELECT
			id, partner_id, room_id, start_time, end_time,
			service_type_encrypted, price_cents, max_capacity,
			notes_encrypted, is_recurring, recurrence_pattern_encrypted,
			status, created_at, updated_at,
			dek_encrypted, key_version, metadata
		FROM %s.availabilities
		WHERE is_recurring = true
		AND status = 'available'
		AND start_time <= $1
		ORDER BY start_time ASC
	`, r.schema)

	rows, err := r.pool.Query(ctx, query, until)
	if err != nil {
		return nil, errs.ClassifyPgError("get recurring availabilities", err)
	}
	defer rows.Close()

	var availabilities []*domain.Availability
	for rows.Next() {
		availabilityEncx := &domain.AvailabilityEncx{}
		err := rows.Scan(
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
			return nil, errs.ClassifyPgError("scan recurring availability row", err)
		}

		// Decrypt sensitive fields using ENCX
		availability, err := domain.DecryptAvailabilityEncx(ctx, r.crypto, availabilityEncx)
		if err != nil {
			return nil, fmt.Errorf("decrypt availability data: %w", err)
		}

		availabilities = append(availabilities, availability)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate recurring availability rows", err)
	}

	return availabilities, nil
}