package allocationRepository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) List(ctx context.Context, filter ports.RoomAllocationFilter) ([]*domain.RoomAllocation, error) {
	query := fmt.Sprintf(`
		SELECT
			id, room_id, user_id, allocation_type,
			start_date, end_date, is_active, created_at, updated_at
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

	if filter.UserID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, *filter.UserID)
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

	var allocations []*domain.RoomAllocation
	for rows.Next() {
		allocation := &domain.RoomAllocation{}
		err := rows.Scan(
			&allocation.ID,
			&allocation.RoomID,
			&allocation.UserID,
			&allocation.AllocationType,
			&allocation.StartDate,
			&allocation.EndDate,
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

	return allocations, nil
}

func (r *Repository) CheckConflict(ctx context.Context, roomID, partnerID uuid.UUID, allocationType domain.AllocationType, startDate, endDate *time.Time) (bool, error) {
	// For dedicated allocations, check if there's an overlap with existing dedicated allocations for the same room
	if allocationType == domain.AllocationTypeDedicated {
		query := fmt.Sprintf(`
			SELECT COUNT(*)
			FROM %s.room_allocations
			WHERE room_id = $1
			AND allocation_type = 'dedicated'
			AND is_active = true
			AND start_date < $2
			AND (end_date IS NULL OR end_date > $3)
		`, r.schema)

		var count int
		err := r.pool.QueryRow(ctx, query, roomID, endDate, startDate).Scan(&count)
		if err != nil {
			return false, errs.ClassifyPgError("check room allocation conflict", err)
		}

		return count > 0, nil
	}

	// For shared allocations, check if partner already has allocation for this room
	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s.room_allocations
		WHERE room_id = $1 AND user_id = $2 AND is_active = true
	`, r.schema)

	var count int
	err := r.pool.QueryRow(ctx, query, roomID, partnerID).Scan(&count)
	if err != nil {
		return false, errs.ClassifyPgError("check partner allocation conflict", err)
	}

	return count > 0, nil
}
