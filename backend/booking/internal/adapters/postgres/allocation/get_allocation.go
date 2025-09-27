package allocationRepository

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/booking/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.RoomAllocation, error) {
	query := fmt.Sprintf(`
		SELECT
			id, room_id, partner_id, allocation_type,
			start_date, end_date, is_active, created_at, updated_at
		FROM %s.room_allocations
		WHERE id = $1
	`, r.schema)

	allocation := &domain.RoomAllocation{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&allocation.ID,
		&allocation.RoomID,
		&allocation.PartnerID,
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

func (r *Repository) GetByPartnerID(ctx context.Context, partnerID uuid.UUID, activeOnly bool) ([]*domain.RoomAllocation, error) {
	query := fmt.Sprintf(`
		SELECT
			id, room_id, partner_id, allocation_type,
			start_date, end_date, is_active, created_at, updated_at
		FROM %s.room_allocations
		WHERE partner_id = $1
	`, r.schema)

	args := []interface{}{partnerID}
	if activeOnly {
		query += " AND is_active = true"
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, errs.ClassifyPgError("get room allocations by partner id", err)
	}
	defer rows.Close()

	var allocations []*domain.RoomAllocation
	for rows.Next() {
		allocation := &domain.RoomAllocation{}
		err := rows.Scan(
			&allocation.ID,
			&allocation.RoomID,
			&allocation.PartnerID,
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

	return allocations, nil
}

func (r *Repository) GetByRoomID(ctx context.Context, roomID uuid.UUID, activeOnly bool) ([]*domain.RoomAllocation, error) {
	query := fmt.Sprintf(`
		SELECT
			id, room_id, partner_id, allocation_type,
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
			&allocation.PartnerID,
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

	return allocations, nil
}

func (r *Repository) GetActiveAllocationForPartnerAndRoom(ctx context.Context, partnerID, roomID uuid.UUID, at time.Time) (*domain.RoomAllocation, error) {
	query := fmt.Sprintf(`
		SELECT
			id, room_id, partner_id, allocation_type,
			start_date, end_date, is_active, created_at, updated_at
		FROM %s.room_allocations
		WHERE partner_id = $1 AND room_id = $2 AND is_active = true
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
		&allocation.PartnerID,
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