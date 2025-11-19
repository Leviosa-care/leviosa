package allocationHelpers

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// ClearAllocationTable removes all records from the room_allocations table
func ClearAllocationTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE booking.room_allocations RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

// InsertAllocation directly inserts an allocation into the database for testing
func InsertAllocation(t *testing.T, ctx context.Context, allocation *domain.RoomAllocation, pool *pgxpool.Pool) {
	t.Helper()

	query := `
		INSERT INTO booking.room_allocations (
			id, room_id, user_id, allocation_type,
			start_date, end_date, is_active,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := pool.Exec(ctx, query,
		allocation.ID,
		allocation.RoomID,
		allocation.UserID,
		allocation.AllocationType,
		allocation.StartDate,
		allocation.EndDate,
		allocation.IsActive,
		allocation.CreatedAt,
		allocation.UpdatedAt,
	)
	require.NoError(t, err)
}

// GetAllocationByID retrieves an allocation from the database by ID
func GetAllocationByID(t *testing.T, ctx context.Context, id uuid.UUID, pool *pgxpool.Pool) (*domain.RoomAllocation, error) {
	t.Helper()

	query := `
		SELECT id, room_id, user_id, allocation_type,
		       start_date, end_date, is_active,
		       created_at, updated_at
		FROM booking.room_allocations
		WHERE id = $1
	`

	var allocation domain.RoomAllocation
	err := pool.QueryRow(ctx, query, id).Scan(
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
		return nil, err
	}

	return &allocation, nil
}

// DeleteAllocation removes an allocation from the database by ID (hard delete for cleanup)
func DeleteAllocation(t *testing.T, ctx context.Context, id uuid.UUID, pool *pgxpool.Pool) {
	t.Helper()

	query := `DELETE FROM booking.room_allocations WHERE id = $1`
	_, err := pool.Exec(ctx, query, id)
	require.NoError(t, err)
}

// AllocationExistsInTable checks if an allocation exists in the database
func AllocationExistsInTable(t *testing.T, ctx context.Context, id uuid.UUID, pool *pgxpool.Pool) bool {
	t.Helper()

	query := `SELECT EXISTS(SELECT 1 FROM booking.room_allocations WHERE id = $1)`
	var exists bool
	err := pool.QueryRow(ctx, query, id).Scan(&exists)
	require.NoError(t, err)

	return exists
}

// CountAllocationsInTable counts the total number of allocations
func CountAllocationsInTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) int {
	t.Helper()

	query := `SELECT COUNT(*) FROM booking.room_allocations`
	var count int
	err := pool.QueryRow(ctx, query).Scan(&count)
	require.NoError(t, err)

	return count
}

// CountAllocationsByUserID counts allocations for a specific user
func CountAllocationsByUserID(t *testing.T, ctx context.Context, userID uuid.UUID, pool *pgxpool.Pool) int {
	t.Helper()

	query := `SELECT COUNT(*) FROM booking.room_allocations WHERE user_id = $1`
	var count int
	err := pool.QueryRow(ctx, query, userID).Scan(&count)
	require.NoError(t, err)

	return count
}

// CountAllocationsByRoomID counts allocations for a specific room
func CountAllocationsByRoomID(t *testing.T, ctx context.Context, roomID uuid.UUID, pool *pgxpool.Pool) int {
	t.Helper()

	query := `SELECT COUNT(*) FROM booking.room_allocations WHERE room_id = $1`
	var count int
	err := pool.QueryRow(ctx, query, roomID).Scan(&count)
	require.NoError(t, err)

	return count
}

// CountActiveAllocations counts allocations where is_active = true
func CountActiveAllocations(t *testing.T, ctx context.Context, pool *pgxpool.Pool) int {
	t.Helper()

	query := `SELECT COUNT(*) FROM booking.room_allocations WHERE is_active = true`
	var count int
	err := pool.QueryRow(ctx, query).Scan(&count)
	require.NoError(t, err)

	return count
}

// CountInactiveAllocations counts allocations where is_active = false
func CountInactiveAllocations(t *testing.T, ctx context.Context, pool *pgxpool.Pool) int {
	t.Helper()

	query := `SELECT COUNT(*) FROM booking.room_allocations WHERE is_active = false`
	var count int
	err := pool.QueryRow(ctx, query).Scan(&count)
	require.NoError(t, err)

	return count
}

// CountDedicatedAllocations counts dedicated allocations
func CountDedicatedAllocations(t *testing.T, ctx context.Context, pool *pgxpool.Pool) int {
	t.Helper()

	query := `SELECT COUNT(*) FROM booking.room_allocations WHERE allocation_type = 'dedicated'`
	var count int
	err := pool.QueryRow(ctx, query).Scan(&count)
	require.NoError(t, err)

	return count
}

// CountSharedAllocations counts shared allocations
func CountSharedAllocations(t *testing.T, ctx context.Context, pool *pgxpool.Pool) int {
	t.Helper()

	query := `SELECT COUNT(*) FROM booking.room_allocations WHERE allocation_type = 'shared'`
	var count int
	err := pool.QueryRow(ctx, query).Scan(&count)
	require.NoError(t, err)

	return count
}

// GetAllAllocations retrieves all allocations from the database (for verification)
func GetAllAllocations(t *testing.T, ctx context.Context, pool *pgxpool.Pool) []*domain.RoomAllocation {
	t.Helper()

	query := `
		SELECT id, room_id, user_id, allocation_type,
		       start_date, end_date, is_active,
		       created_at, updated_at
		FROM booking.room_allocations
		ORDER BY created_at DESC
	`

	rows, err := pool.Query(ctx, query)
	require.NoError(t, err)
	defer rows.Close()

	var allocations []*domain.RoomAllocation
	for rows.Next() {
		var allocation domain.RoomAllocation
		err := rows.Scan(
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
		require.NoError(t, err)
		allocations = append(allocations, &allocation)
	}

	require.NoError(t, rows.Err())
	return allocations
}
