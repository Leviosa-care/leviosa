package allocationRepository

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

// CheckConflict checks if a new allocation would conflict with existing active allocations
// For dedicated allocations: checks for time period overlaps in the same room
// For shared allocations: checks if the user already has an allocation for the room
// excludeID allows excluding a specific allocation from the conflict check (useful for updates)
func (r *Repository) CheckConflict(
	ctx context.Context,
	roomID uuid.UUID,
	userIDHash string,
	allocationType domain.AllocationType,
	startDate, endDate *time.Time,
	excludeID *uuid.UUID,
) (bool, error) {
	if allocationType == domain.AllocationTypeDedicated {
		return r.checkDedicatedConflict(ctx, roomID, startDate, endDate, excludeID)
	}
	return r.checkSharedConflict(ctx, roomID, userIDHash, excludeID)
}

// checkDedicatedConflict checks for time period overlaps with existing dedicated allocations
func (r *Repository) checkDedicatedConflict(
	ctx context.Context,
	roomID uuid.UUID,
	startDate, endDate *time.Time,
	excludeID *uuid.UUID,
) (bool, error) {
	if startDate == nil {
		return false, fmt.Errorf("start_date is required for dedicated allocation conflict check: %w", errs.ErrInvalidValue)
	}

	// Query to detect overlapping time periods
	// Two periods [A_start, A_end] and [B_start, B_end] overlap if:
	// A_start < B_end AND B_start < A_end
	// When end_date is NULL (indefinite), we treat it as extending infinitely
	// excludeID allows excluding a specific allocation (useful when updating)
	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM %s.room_allocations
			WHERE room_id = $1
			AND is_active = true
			AND allocation_type = 'dedicated'
			AND ($4::uuid IS NULL OR id != $4)
			AND (
				-- Case 1: Existing has end_date, check overlap
				(end_date IS NOT NULL AND start_date < $3 AND $2 < end_date)
				OR
				-- Case 2: Existing has NULL end_date (indefinite), conflicts if new starts after existing
				(end_date IS NULL AND $2 >= start_date)
				OR
				-- Case 3: New allocation has NULL end_date, conflicts if new start is before existing end
				($3 IS NULL AND (end_date IS NULL OR start_date < end_date))
			)
		)
	`, r.schema)

	var hasConflict bool
	err := r.pool.QueryRow(ctx, query, roomID, startDate, endDate, excludeID).Scan(&hasConflict)
	if err != nil {
		return false, errs.ClassifyPgError("check dedicated allocation conflict", err)
	}

	return hasConflict, nil
}

// checkSharedConflict checks if the user already has an active allocation for the room
func (r *Repository) checkSharedConflict(
	ctx context.Context,
	roomID uuid.UUID,
	userIDHash string,
	excludeID *uuid.UUID,
) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM %s.room_allocations
			WHERE room_id = $1
			AND user_id_hash = $2
			AND is_active = true
			AND allocation_type = 'shared'
			AND ($3::uuid IS NULL OR id != $3)
		)
	`, r.schema)

	var hasConflict bool
	err := r.pool.QueryRow(ctx, query, roomID, userIDHash, excludeID).Scan(&hasConflict)
	if err != nil {
		return false, errs.ClassifyPgError("check shared allocation conflict", err)
	}

	return hasConflict, nil
}
