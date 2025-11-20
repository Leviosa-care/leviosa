package allocationRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) Create(ctx context.Context, allocation *domain.RoomAllocationEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.room_allocations (
			id, room_id, user_id_encrypted, user_id_hash, allocation_type,
			start_date, end_date, dek_encrypted, key_version,
			is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
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
		allocation.CreatedAt,
		allocation.UpdatedAt,
	)
	if err != nil {
		return errs.ClassifyPgError("create room allocation", err)
	}

	return nil
}
