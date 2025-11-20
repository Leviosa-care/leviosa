package allocationRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) Update(ctx context.Context, allocation *domain.RoomAllocationEncx) error {
	query := fmt.Sprintf(`
		UPDATE %s.room_allocations SET
			room_id = $2,
			user_id_encrypted = $3,
			user_id_hash = $4,
			allocation_type = $5,
			start_date = $6,
			end_date = $7,
			dek_encrypted = $8,
			key_version = $9,
			is_active = $10,
			updated_at = $11
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query,
		allocation.ID,
		allocation.RoomID,
		allocation.UserIDEncrypted,
		allocation.UserIDHash,
		allocation.AllocationType,
		allocation.StartDate,
		allocation.EndDate,
		allocation.DEKEncrypted,
		allocation.KeyVersion,
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
