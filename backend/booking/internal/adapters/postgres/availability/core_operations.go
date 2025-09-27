package availabilityRepository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/booking/internal/ports"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (r *Repository) Update(ctx context.Context, availability *domain.Availability) error {
	// Encrypt sensitive fields
	if err := r.crypto.EncryptStruct(ctx, availability); err != nil {
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
			updated_at = $13
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
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
		availability.UpdatedAt,
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

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.availabilities
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return errs.ClassifyPgError("delete availability", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}

func (r *Repository) List(ctx context.Context, filter ports.AvailabilityFilter) ([]*domain.Availability, error) {
	query := fmt.Sprintf(`
		SELECT
			id, partner_id, room_id, start_time, end_time,
			service_type_encrypted, price_cents, max_capacity,
			notes_encrypted, is_recurring, recurrence_pattern_encrypted,
			status, created_at, updated_at
		FROM %s.availabilities
	`, r.schema)

	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	// Apply filters
	if filter.PartnerID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("partner_id = $%d", argIndex))
		args = append(args, *filter.PartnerID)
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

	var availabilities []*domain.Availability
	for rows.Next() {
		availability := &domain.Availability{}
		err := rows.Scan(
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
			return nil, errs.ClassifyPgError("scan availability row", err)
		}

		// Decrypt sensitive fields
		if err := r.crypto.DecryptStruct(ctx, availability); err != nil {
			return nil, fmt.Errorf("decrypt availability data: %w", err)
		}

		availabilities = append(availabilities, availability)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate availability rows", err)
	}

	return availabilities, nil
}

func (r *Repository) CheckConflict(ctx context.Context, partnerID uuid.UUID, startTime, endTime time.Time, excludeID *uuid.UUID) (bool, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s.availabilities
		WHERE partner_id = $1
		AND start_time < $2
		AND end_time > $3
		AND status IN ('available', 'booked')
	`, r.schema)

	args := []interface{}{partnerID, endTime, startTime}

	if excludeID != nil {
		query += " AND id != $4"
		args = append(args, *excludeID)
	}

	var count int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, errs.ClassifyPgError("check availability conflict", err)
	}

	return count > 0, nil
}

func (r *Repository) GetRecurringAvailabilities(ctx context.Context, until time.Time) ([]*domain.Availability, error) {
	query := fmt.Sprintf(`
		SELECT
			id, partner_id, room_id, start_time, end_time,
			service_type_encrypted, price_cents, max_capacity,
			notes_encrypted, is_recurring, recurrence_pattern_encrypted,
			status, created_at, updated_at
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
		availability := &domain.Availability{}
		err := rows.Scan(
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
			return nil, errs.ClassifyPgError("scan recurring availability row", err)
		}

		// Decrypt sensitive fields
		if err := r.crypto.DecryptStruct(ctx, availability); err != nil {
			return nil, fmt.Errorf("decrypt availability data: %w", err)
		}

		availabilities = append(availabilities, availability)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate recurring availability rows", err)
	}

	return availabilities, nil
}