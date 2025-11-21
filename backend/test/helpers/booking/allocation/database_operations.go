package allocationHelpers

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// ClearAllocationTable removes all records from the room_allocations table
func ClearAllocationTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE booking.room_allocations RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

// InsertAllocationEncx directly inserts an encrypted allocation into the database for testing
func InsertAllocationEncx(t *testing.T, ctx context.Context, allocation *domain.RoomAllocationEncx, pool *pgxpool.Pool) {
	t.Helper()

	query := `
		INSERT INTO booking.room_allocations (
			id, room_id, user_id_encrypted, user_id_hash, allocation_type,
			start_date, end_date, dek_encrypted, key_version,
			is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := pool.Exec(ctx, query,
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
	require.NoError(t, err)
}

// GetAllocationEncxByID retrieves an encrypted allocation from the database by ID
func GetAllocationEncxByID(t *testing.T, ctx context.Context, id uuid.UUID, pool *pgxpool.Pool) (*domain.RoomAllocationEncx, error) {
	t.Helper()

	query := `
		SELECT id, room_id, user_id_encrypted, user_id_hash, allocation_type,
		       start_date, end_date, dek_encrypted, key_version,
		       is_active, created_at, updated_at
		FROM booking.room_allocations
		WHERE id = $1
	`

	var allocation domain.RoomAllocationEncx
	err := pool.QueryRow(ctx, query, id).Scan(
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

// CountAllocationsByUserIDHash counts allocations for a specific user by their hash
func CountAllocationsByUserIDHash(t *testing.T, ctx context.Context, userIDHash string, pool *pgxpool.Pool) int {
	t.Helper()

	query := `SELECT COUNT(*) FROM booking.room_allocations WHERE user_id_hash = $1`
	var count int
	err := pool.QueryRow(ctx, query, userIDHash).Scan(&count)
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

// GetAllAllocationsEncx retrieves all encrypted allocations from the database (for verification)
func GetAllAllocationsEncx(t *testing.T, ctx context.Context, pool *pgxpool.Pool) []*domain.RoomAllocationEncx {
	t.Helper()

	query := `
		SELECT id, room_id, user_id_encrypted, user_id_hash, allocation_type,
		       start_date, end_date, dek_encrypted, key_version,
		       is_active, created_at, updated_at
		FROM booking.room_allocations
		ORDER BY created_at DESC
	`

	rows, err := pool.Query(ctx, query)
	require.NoError(t, err)
	defer rows.Close()

	var allocations []*domain.RoomAllocationEncx
	for rows.Next() {
		var allocation domain.RoomAllocationEncx
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
		require.NoError(t, err)
		allocations = append(allocations, &allocation)
	}

	require.NoError(t, rows.Err())
	return allocations
}

// ComputeUserIDHash computes hash for a user ID using the crypto service
func ComputeUserIDHash(t *testing.T, ctx context.Context, crypto encx.CryptoService, userID uuid.UUID) string {
	t.Helper()

	userIDBytes, err := userID.MarshalBinary()
	require.NoError(t, err, "failed to serialize user ID")

	return crypto.HashBasic(ctx, userIDBytes)
}

// InsertAllocation is a convenience wrapper that encrypts a RoomAllocation and inserts it.
// This is for integration tests that work with normal domain objects.
func InsertAllocation(t *testing.T, ctx context.Context, allocation *domain.RoomAllocation, pool *pgxpool.Pool, crypto encx.CryptoService) {
	t.Helper()

	allocationEncx, err := domain.ProcessRoomAllocationEncx(ctx, crypto, allocation)
	require.NoError(t, err, "failed to encrypt allocation")

	InsertAllocationEncx(t, ctx, allocationEncx, pool)
}

// GetAllocationByID is a convenience wrapper that retrieves and decrypts an allocation.
// This is for integration tests that need to verify decrypted data.
func GetAllocationByID(t *testing.T, ctx context.Context, id uuid.UUID, pool *pgxpool.Pool, crypto encx.CryptoService) (*domain.RoomAllocation, error) {
	t.Helper()

	allocationEncx, err := GetAllocationEncxByID(t, ctx, id, pool)
	if err != nil {
		return nil, err
	}

	allocation, err := domain.DecryptRoomAllocationEncx(ctx, crypto, allocationEncx)
	if err != nil {
		return nil, err
	}

	return allocation, nil
}
