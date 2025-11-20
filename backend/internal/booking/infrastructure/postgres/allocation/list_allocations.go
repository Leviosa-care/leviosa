package allocationRepository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) List(ctx context.Context, filter ports.RoomAllocationFilter) ([]*domain.RoomAllocationEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, room_id, user_id_encrypted, user_id_hash, allocation_type,
			start_date, end_date, dek_encrypted, key_version,
			is_active, created_at, updated_at
		FROM %s.room_allocations
	`, r.schema)

	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	// Apply filters
	if filter.RoomID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("room_id = $%d", argIndex))
		args = append(args, *filter.RoomID)
		argIndex++
	}

	if filter.UserIDHash != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("user_id_hash = $%d", argIndex))
		args = append(args, *filter.UserIDHash)
		argIndex++
	}

	if filter.AllocationType != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("allocation_type = $%d", argIndex))
		args = append(args, *filter.AllocationType)
		argIndex++
	}

	if filter.IsActive != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *filter.IsActive)
		argIndex++
	}

	if filter.ActiveAt != nil {
		whereConditions = append(whereConditions, fmt.Sprintf(`(
			allocation_type = 'shared'
			OR (
				allocation_type = 'dedicated'
				AND start_date <= $%d
				AND (end_date IS NULL OR end_date >= $%d)
			)
		)`, argIndex, argIndex))
		args = append(args, *filter.ActiveAt)
		argIndex++
	}

	if filter.OverlapsWith != nil {
		whereConditions = append(whereConditions, fmt.Sprintf(`(
			allocation_type = 'shared'
			OR (
				allocation_type = 'dedicated'
				AND start_date < $%d
				AND (end_date IS NULL OR end_date > $%d)
			)
		)`, argIndex+1, argIndex))
		args = append(args, filter.OverlapsWith.Start, filter.OverlapsWith.End)
		argIndex += 2
	}

	// Add WHERE clause if we have conditions
	if len(whereConditions) > 0 {
		query += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Add ordering
	orderBy := "created_at"
	if filter.OrderBy != "" {
		switch filter.OrderBy {
		case "created_at", "start_date", "end_date":
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
		return nil, errs.ClassifyPgError("list room allocations", err)
	}
	defer rows.Close()

	var allocations []*domain.RoomAllocationEncx
	for rows.Next() {
		allocation := &domain.RoomAllocationEncx{}
		err := rows.Scan(
			&allocation.ID,
			&allocation.RoomID,
			&allocation.UserIDEncrypted,
			&allocation.UserIDHash,
			&allocation.AllocationType,
			&allocation.StartDate,
			&allocation.EndDate,
			&allocation.DEKEncrypted,
			&allocation.KeyVersion,
			&allocation.IsActive,
			&allocation.CreatedAt,
			&allocation.UpdatedAt,
		)
		if err != nil {
			return nil, errs.ClassifyPgError("scan room allocation row", err)
		}

		allocations = append(allocations, allocation)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate room allocation rows", err)
	}

	if len(allocations) == 0 {
		return []*domain.RoomAllocationEncx{}, nil
	}

	return allocations, nil
}
