package availabilityRepository

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetRoomBookingsForDate retrieves all bookings (availabilities) for a room on a specific date
// Returns availabilities sorted by start time, useful for gap detection
func (r *Repository) GetRoomBookingsForDate(ctx context.Context, roomID uuid.UUID, date time.Time) ([]*domain.AvailabilityEncx, error) {
	query := `
		SELECT
			id, user_id, room_id, start_time, end_time,
			max_capacity, status, price_cents,
			service_type_encrypted, notes_encrypted,
			is_recurring, recurrence_pattern,
			created_at, updated_at,
			dek_encrypted, key_version, metadata
		FROM booking.availabilities
		WHERE room_id = $1
			AND DATE(start_time) = $2
			AND status IN ('available', 'booked')
		ORDER BY start_time ASC
	`

	// Normalize date to midnight for comparison
	normalizedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	rows, err := r.pool.Query(ctx, query, roomID, normalizedDate)
	if err != nil {
		return nil, errs.ClassifyPgError("query room bookings for date", err)
	}
	defer rows.Close()

	availabilities := []*domain.AvailabilityEncx{}

	for rows.Next() {
		availabilityEncx := &domain.AvailabilityEncx{}

		err := rows.Scan(
			&availabilityEncx.ID,
			&availabilityEncx.UserID,
			&availabilityEncx.RoomID,
			&availabilityEncx.StartTime,
			&availabilityEncx.EndTime,
			&availabilityEncx.MaxCapacity,
			&availabilityEncx.Status,
			&availabilityEncx.PriceCents,
			&availabilityEncx.ServiceTypeEncrypted,
			&availabilityEncx.NotesEncrypted,
			&availabilityEncx.IsRecurring,
			&availabilityEncx.RecurrencePattern,
			&availabilityEncx.CreatedAt,
			&availabilityEncx.UpdatedAt,
			&availabilityEncx.DEKEncrypted,
			&availabilityEncx.KeyVersion,
			&availabilityEncx.Metadata,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan availability", err)
		}

		availabilities = append(availabilities, availabilityEncx)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate room bookings for date", err)
	}

	return availabilities, nil
}
