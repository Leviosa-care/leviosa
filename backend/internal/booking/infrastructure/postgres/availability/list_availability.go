package availabilityRepository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) List(ctx context.Context, filter ports.AvailabilityFilter) ([]*domain.AvailabilityEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, user_id, room_id, start_time, end_time,
			service_type_encrypted, price_cents, max_capacity,
			notes_encrypted, is_recurring, recurrence_pattern_encrypted,
			status, created_at, updated_at,
			dek_encrypted, key_version, metadata
		FROM %s.availabilities
	`, r.schema)

	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	// Apply filters
	if filter.UserID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, *filter.UserID)
		argIndex++
	}

	if filter.RoomID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("room_id = $%d", argIndex))
		args = append(args, *filter.RoomID)
		argIndex++
	}

	if len(filter.Status) > 0 {
		placeholders := make([]string, len(filter.Status))
		for i, status := range filter.Status {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, status)
			argIndex++
		}
		whereConditions = append(whereConditions, fmt.Sprintf("status IN (%s)", strings.Join(placeholders, ",")))
	}

	if filter.StartTime != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("start_time >= $%d", argIndex))
		args = append(args, *filter.StartTime)
		argIndex++
	}

	if filter.EndTime != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("end_time <= $%d", argIndex))
		args = append(args, *filter.EndTime)
		argIndex++
	}

	if filter.TimeRange != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("start_time < $%d AND end_time > $%d", argIndex+1, argIndex))
		args = append(args, filter.TimeRange.Start, filter.TimeRange.End)
		argIndex += 2
	}

	if filter.AvailableOnly {
		whereConditions = append(whereConditions, "status = 'available' AND start_time > NOW()")
	}

	// Add WHERE clause if we have conditions
	if len(whereConditions) > 0 {
		query += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Add ordering
	orderBy := "start_time"
	if filter.OrderBy != "" {
		switch filter.OrderBy {
		case "start_time", "end_time", "created_at", "price_cents":
			orderBy = filter.OrderBy
		}
	}

	orderDirection := "ASC"
	if filter.OrderDirection == "desc" {
		orderDirection = "DESC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", orderBy, orderDirection)

	// Add pagination
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, errs.ClassifyPgError("list availabilities", err)
	}
	defer rows.Close()

	var availabilitiesEncx []*domain.AvailabilityEncx
	for rows.Next() {
		availabilityEncx := &domain.AvailabilityEncx{}
		err := rows.Scan(
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
			&availabilityEncx.RecurrencePatternEncrypted,
			&availabilityEncx.Status,
			&availabilityEncx.CreatedAt,
			&availabilityEncx.UpdatedAt,
			&availabilityEncx.DEKEncrypted,
			&availabilityEncx.KeyVersion,
			&availabilityEncx.Metadata,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan availability row", err)
		}

		availabilitiesEncx = append(availabilitiesEncx, availabilityEncx)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate availability rows", err)
	}

	return availabilitiesEncx, nil
}
