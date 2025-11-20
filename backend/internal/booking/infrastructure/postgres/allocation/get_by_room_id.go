package allocationRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetByRoomID(ctx context.Context, roomID uuid.UUID, activeOnly bool) ([]*domain.RoomAllocation, error) {
	query := fmt.Sprintf(`
		SELECT
			id, room_id, user_id, allocation_type,
			start_date, end_date, is_active, created_at, updated_at
		FROM %s.room_allocations
		WHERE room_id = $1
	`, r.schema)

	args := []interface{}{roomID}
	if activeOnly {
		query += " AND is_active = true"
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, errs.ClassifyPgError("get room allocations by room id", err)
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

	if len(allocations) == 0 {
		return []*domain.RoomAllocation{}, nil
	}

	return allocations, nil
}
