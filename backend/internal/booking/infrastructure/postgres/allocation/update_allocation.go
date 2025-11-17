package allocationRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) Update(ctx context.Context, allocation *domain.RoomAllocation) error {
	query := fmt.Sprintf(`
		UPDATE %s.room_allocations SET
			room_id = $2,
			user_id = $3,
			allocation_type = $4,
			start_date = $5,
			end_date = $6,
			is_active = $7,
			updated_at = $8
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		allocation.ID,
		allocation.RoomID,
		allocation.UserID,
		allocation.AllocationType,
		allocation.StartDate,
		allocation.EndDate,
		allocation.IsActive,
		allocation.UpdatedAt,
	)
	if err != nil {
		return errs.ClassifyPgError("update room allocation", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	// Soft delete by marking as inactive
	query := fmt.Sprintf(`
		UPDATE %s.room_allocations
		SET is_active = false, updated_at = NOW()
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return errs.ClassifyPgError("delete room allocation", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrRepositoryNotFound
	}

	return nil
}
