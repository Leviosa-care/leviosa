package allocationRepository

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetActiveAllocationForPartnerAndRoom(ctx context.Context, partnerID, roomID uuid.UUID, at time.Time) (*domain.RoomAllocation, error) {
	query := fmt.Sprintf(`
		SELECT
			id, room_id, user_id, allocation_type,
			start_date, end_date, is_active, created_at, updated_at
		FROM %s.room_allocations
		WHERE user_id = $1 AND room_id = $2 AND is_active = true
		AND (
			allocation_type = 'shared'
			OR (
				allocation_type = 'dedicated'
				AND start_date <= $3
				AND (end_date IS NULL OR end_date >= $3)
			)
		)
		ORDER BY allocation_type DESC, created_at DESC
		LIMIT 1
	`, r.schema)

	allocation := &domain.RoomAllocation{}
	err := r.pool.QueryRow(ctx, query, partnerID, roomID, at).Scan(
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
		return nil, errs.ClassifyPgError("get active allocation for partner and room", err)
	}

	return allocation, nil
}
