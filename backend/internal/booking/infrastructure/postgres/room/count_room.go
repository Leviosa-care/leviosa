package roomRepository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

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

	var count int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, errs.ClassifyPgError("count rooms", err)
	}

	return count, nil
}
