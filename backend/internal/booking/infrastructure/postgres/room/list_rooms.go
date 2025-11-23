package roomRepository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) List(ctx context.Context, filter ports.RoomFilter) ([]*domain.RoomEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, building_id, name_encrypted, name_hash, description_encrypted,
			room_number_encrypted, room_number_hash, capacity, equipment_encrypted,
			operating_start_time, operating_end_time,
			is_active, created_at, updated_at,
			dek_encrypted, key_version, metadata
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

	// Apply name filter using hash
	if filter.NameHash != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("name_hash = $%d", argIndex))
		args = append(args, *filter.NameHash)
		argIndex++
	}

	// Apply room number filter using hash
	if filter.RoomNumberHash != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("room_number_hash = $%d", argIndex))
		args = append(args, *filter.RoomNumberHash)
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
		case "name", "created_at", "capacity":
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

	var roomsEncx []*domain.RoomEncx
	for rows.Next() {
		roomEncx := &domain.RoomEncx{}
		err := rows.Scan(
			&roomEncx.ID,
			&roomEncx.BuildingID,
			&roomEncx.NameEncrypted,
			&roomEncx.NameHash,
			&roomEncx.DescriptionEncrypted,
			&roomEncx.RoomNumberEncrypted,
			&roomEncx.RoomNumberHash,
			&roomEncx.Capacity,
			&roomEncx.EquipmentEncrypted,
			&roomEncx.OperatingStartTime,
			&roomEncx.OperatingEndTime,
			&roomEncx.IsActive,
			&roomEncx.CreatedAt,
			&roomEncx.UpdatedAt,
			&roomEncx.DEKEncrypted,
			&roomEncx.KeyVersion,
			&roomEncx.Metadata,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan room row", err)
		}

		roomsEncx = append(roomsEncx, roomEncx)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate room rows", err)
	}

	return roomsEncx, nil
}
