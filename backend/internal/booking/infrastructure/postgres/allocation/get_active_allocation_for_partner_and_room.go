package allocationRepository

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetActiveAllocationForPartnerAndRoom(ctx context.Context, userIDHash string, roomID uuid.UUID, at time.Time) (*domain.RoomAllocationEncx, error) {
	query := fmt.Sprintf(`
		SELECT
			id, room_id, user_id_encrypted, user_id_hash, allocation_type,
			start_date, end_date, dek_encrypted, key_version,
			is_active, created_at, updated_at
		FROM %s.room_allocations
		WHERE user_id_hash = $1 AND room_id = $2 AND is_active = true
		AND (
			allocation_type = 'shared'
			OR (
				allocation_type = 'dedicated'
				AND start_date <= $3
				AND (end_date IS NULL OR end_date >= $3)
			)
		)
		ORDER BY
			CASE
				WHEN allocation_type = 'dedicated' THEN 1
				WHEN allocation_type = 'shared' THEN 2
			END ASC,
			created_at DESC
		LIMIT 1
	`, r.schema)

	allocation := &domain.RoomAllocationEncx{}
	err := r.pool.QueryRow(ctx, query, userIDHash, roomID, at).Scan(
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
		return nil, errs.ClassifyPgError("get active allocation for partner and room", err)
	}

	return allocation, nil
}
