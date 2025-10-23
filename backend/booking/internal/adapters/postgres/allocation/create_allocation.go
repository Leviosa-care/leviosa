package allocationRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) Create(ctx context.Context, allocation *domain.RoomAllocation) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.room_allocations (
			id, room_id, partner_id, allocation_type,
			start_date, end_date, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		allocation.ID,
		allocation.RoomID,
		allocation.PartnerID,
		allocation.AllocationType,
		allocation.StartDate,
		allocation.EndDate,
		allocation.IsActive,
		allocation.CreatedAt,
		allocation.UpdatedAt,
	)
	if err != nil {
		return errs.ClassifyPgError("create room allocation", err)
	}

	return nil
}