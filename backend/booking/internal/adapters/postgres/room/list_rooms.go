package roomRepository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/booking/internal/ports"
	"github.com/Leviosa-care/core/errs"
)

func (r *Repository) List(ctx context.Context, filter ports.RoomFilter) ([]*domain.Room, error) {
	query := fmt.Sprintf(`
		SELECT
			id, building_id, name_encrypted, description_encrypted,
			room_number_encrypted, capacity, equipment_encrypted,
			hourly_rate_cents, is_active, created_at, updated_at
		FROM %s.rooms
	`, r.schema)

	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	// Apply filters
	if filter.BuildingID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("building_id = $%d", argIndex))
		args = append(args, *filter.BuildingID)
		argIndex++
	}

	if filter.IsActive != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *filter.IsActive)
		argIndex++
	}

	if filter.MinCapacity != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("capacity >= $%d", argIndex))
		args = append(args, *filter.MinCapacity)
		argIndex++
	}

	if filter.MaxCapacity != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("capacity <= $%d", argIndex))
		args = append(args, *filter.MaxCapacity)
		argIndex++
	}

	if filter.MinHourlyRate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("hourly_rate_cents >= $%d", argIndex))
		args = append(args, *filter.MinHourlyRate)
		argIndex++
	}

	if filter.MaxHourlyRate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("hourly_rate_cents <= $%d", argIndex))
		args = append(args, *filter.MaxHourlyRate)
		argIndex++
	}

	// Add WHERE clause if we have conditions
	if len(whereConditions) > 0 {
		query += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Add ordering
	orderBy := "created_at"
	if filter.OrderBy != "" {
		switch filter.OrderBy {
		case "name", "created_at", "capacity", "hourly_rate_cents":
			orderBy = filter.OrderBy
		}
	}

	orderDirection := "DESC"
	if filter.OrderDirection == "asc" {
		orderDirection = "ASC"
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
		return nil, errs.ClassifyPgError("list rooms", err)
	}
	defer rows.Close()

	var rooms []*domain.Room
	for rows.Next() {
		room := &domain.Room{}
		err := rows.Scan(
			&room.ID,
			&room.BuildingID,
			&room.NameEncrypted,
			&room.DescriptionEncrypted,
			&room.RoomNumberEncrypted,
			&room.Capacity,
			&room.EquipmentEncrypted,
			&room.HourlyRateCents,
			&room.IsActive,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan room row", err)
		}

		// Decrypt sensitive fields
		if err := r.crypto.DecryptStruct(ctx, room); err != nil {
			return nil, fmt.Errorf("decrypt room data: %w", err)
		}

		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate room rows", err)
	}

	return rooms, nil
}

func (r *Repository) Count(ctx context.Context, filter ports.RoomFilter) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s.rooms", r.schema)

	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	// Apply filters
	if filter.BuildingID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("building_id = $%d", argIndex))
		args = append(args, *filter.BuildingID)
		argIndex++
	}

	if filter.IsActive != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *filter.IsActive)
		argIndex++
	}

	if filter.MinCapacity != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("capacity >= $%d", argIndex))
		args = append(args, *filter.MinCapacity)
		argIndex++
	}

	if filter.MaxCapacity != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("capacity <= $%d", argIndex))
		args = append(args, *filter.MaxCapacity)
		argIndex++
	}

	if filter.MinHourlyRate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("hourly_rate_cents >= $%d", argIndex))
		args = append(args, *filter.MinHourlyRate)
		argIndex++
	}

	if filter.MaxHourlyRate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("hourly_rate_cents <= $%d", argIndex))
		args = append(args, *filter.MaxHourlyRate)
		argIndex++
	}

	// Add WHERE clause if we have conditions
	if len(whereConditions) > 0 {
		query += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	var count int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, errs.ClassifyPgError("count rooms", err)
	}

	return count, nil
}