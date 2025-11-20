package allocationRepository

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (r *Repository) GetByUserIDHash(ctx context.Context, userIDHash string, activeOnly bool) ([]*domain.RoomAllocationEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, room_id, user_id_encrypted, user_id_hash, allocation_type,
			start_date, end_date, dek_encrypted, key_version,
			is_active, created_at, updated_at
		FROM %s.room_allocations
		WHERE user_id_hash = $1
	`, r.schema)

	args := []interface{}{userIDHash}
	if activeOnly {
		query += " AND is_active = true"
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, errs.ClassifyPgError("get room allocations by user id hash", err)
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
		return nil, errs.ClassifyPgError("iterate allocations rows", err)
	}

	if len(allocations) == 0 {
		return []*domain.RoomAllocationEncx{}, nil
	}

	return allocations, nil
}
