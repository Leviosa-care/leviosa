package allocationRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.RoomAllocation, error) {
	query := fmt.Sprintf(`
		SELECT
			id, room_id, user_id, allocation_type,
			start_date, end_date, is_active, created_at, updated_at
		FROM %s.room_allocations
		WHERE id = $1
	`, r.schema)

	allocation := &domain.RoomAllocation{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
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
		return nil, errs.ClassifyPgError("get room allocation by id", err)
	}

	return allocation, nil
}
