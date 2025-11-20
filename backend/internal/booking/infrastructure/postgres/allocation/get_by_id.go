package allocationRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.RoomAllocationEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, room_id, user_id_encrypted, user_id_hash, allocation_type,
			start_date, end_date, dek_encrypted, key_version,
			is_active, created_at, updated_at
		FROM %s.room_allocations
		WHERE id = $1
	`, r.schema)

	allocation := &domain.RoomAllocationEncx{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
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
		return nil, errs.ClassifyPgError("get room allocation by id", err)
	}

	return allocation, nil
}
