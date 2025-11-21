package allocationRepository

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// NewTestSharedAllocationEncx creates a pre-encrypted shared allocation for repository tests
func NewTestSharedAllocationEncx(t *testing.T, crypto encx.CryptoService, roomID, userID uuid.UUID) *domain.RoomAllocationEncx {
	t.Helper()

	ctx := context.Background()

	// Create domain entity
	allocation, err := domain.NewSharedAllocation(roomID, userID)
	require.NoError(t, err, "failed to create shared allocation")

	// Encrypt
	allocationEncx, err := domain.ProcessRoomAllocationEncx(ctx, crypto, allocation)
	require.NoError(t, err, "failed to encrypt allocation")

	return allocationEncx
}

// NewTestDedicatedAllocationEncx creates a pre-encrypted dedicated allocation for repository tests
func NewTestDedicatedAllocationEncx(t *testing.T, crypto encx.CryptoService, roomID, userID uuid.UUID, startDate, endDate time.Time) *domain.RoomAllocationEncx {
	t.Helper()

	ctx := context.Background()

	// Create domain entity
	allocation, err := domain.NewDedicatedAllocation(roomID, userID, &startDate, &endDate)
	require.NoError(t, err, "failed to create dedicated allocation")

	// Encrypt
	allocationEncx, err := domain.ProcessRoomAllocationEncx(ctx, crypto, allocation)
	require.NoError(t, err, "failed to encrypt allocation")

	return allocationEncx
}

// NewTestDedicatedAllocationEncxWithNilEndDate creates a pre-encrypted dedicated allocation with nil end date
func NewTestDedicatedAllocationEncxWithNilEndDate(t *testing.T, crypto encx.CryptoService, roomID, userID uuid.UUID, startDate time.Time) *domain.RoomAllocationEncx {
	t.Helper()

	ctx := context.Background()

	// Create domain entity
	allocation, err := domain.NewDedicatedAllocation(roomID, userID, &startDate, nil)
	require.NoError(t, err, "failed to create dedicated allocation with nil end date")

	// Encrypt
	allocationEncx, err := domain.ProcessRoomAllocationEncx(ctx, crypto, allocation)
	require.NoError(t, err, "failed to encrypt allocation")

	return allocationEncx
}

// NewTestInactiveSharedAllocationEncx creates a pre-encrypted inactive shared allocation
func NewTestInactiveSharedAllocationEncx(t *testing.T, crypto encx.CryptoService, roomID, userID uuid.UUID) *domain.RoomAllocationEncx {
	t.Helper()

	ctx := context.Background()

	// Create domain entity
	allocation, err := domain.NewSharedAllocation(roomID, userID)
	require.NoError(t, err, "failed to create shared allocation")

	// Deactivate
	allocation.Deactivate()

	// Encrypt
	allocationEncx, err := domain.ProcessRoomAllocationEncx(ctx, crypto, allocation)
	require.NoError(t, err, "failed to encrypt allocation")

	return allocationEncx
}

// ClearAllocationTable removes all records from room_allocations table
func ClearAllocationTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, "DELETE FROM booking.room_allocations")
	require.NoError(t, err, "failed to clear allocation table")
}

// GetAllocationEncxByID retrieves an encrypted allocation by ID directly from the database
func GetAllocationEncxByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, id uuid.UUID) *domain.RoomAllocationEncx {
	t.Helper()

	query := `
		SELECT id, room_id, user_id_encrypted, user_id_hash, allocation_type,
			start_date, end_date, dek_encrypted, key_version,
			is_active, created_at, updated_at
		FROM booking.room_allocations
		WHERE id = $1
	`

	allocation := &domain.RoomAllocationEncx{}
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
	require.NoError(t, err, "failed to get allocation by ID")

	return allocation
}

// CountAllocationsByUserIDHash counts allocations by user ID hash
func CountAllocationsByUserIDHash(t *testing.T, ctx context.Context, pool *pgxpool.Pool, userIDHash string) int {
	t.Helper()

	var count int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM booking.room_allocations WHERE user_id_hash = $1", userIDHash).Scan(&count)
	require.NoError(t, err, "failed to count allocations by user ID hash")

	return count
}

// CountAllocationsByRoomID counts allocations by room ID
func CountAllocationsByRoomID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, roomID uuid.UUID) int {
	t.Helper()

	var count int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM booking.room_allocations WHERE room_id = $1", roomID).Scan(&count)
	require.NoError(t, err, "failed to count allocations by room ID")

	return count
}

// ComputeUserIDHash computes hash for a user ID (uses same serialization as domain.ProcessRoomAllocationEncx)
func ComputeUserIDHash(t *testing.T, ctx context.Context, crypto encx.CryptoService, userID uuid.UUID) string {
	t.Helper()

	userIDBytes, err := encx.SerializeValue(userID)
	require.NoError(t, err, "failed to serialize user ID")

	return crypto.HashBasic(ctx, userIDBytes)
}

// NewTestFutureDedicatedAllocationEncx creates a pre-encrypted dedicated allocation starting in the future
func NewTestFutureDedicatedAllocationEncx(t *testing.T, crypto encx.CryptoService, roomID, userID uuid.UUID) *domain.RoomAllocationEncx {
	t.Helper()

	ctx := context.Background()

	startDate := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
	endDate := time.Now().AddDate(0, 1, 0).Truncate(24 * time.Hour)

	allocation, err := domain.NewDedicatedAllocation(roomID, userID, &startDate, &endDate)
	require.NoError(t, err, "failed to create future dedicated allocation")

	allocationEncx, err := domain.ProcessRoomAllocationEncx(ctx, crypto, allocation)
	require.NoError(t, err, "failed to encrypt allocation")

	return allocationEncx
}

// NewTestPastDedicatedAllocationEncx creates a pre-encrypted dedicated allocation that has already ended
func NewTestPastDedicatedAllocationEncx(t *testing.T, crypto encx.CryptoService, roomID, userID uuid.UUID) *domain.RoomAllocationEncx {
	t.Helper()

	ctx := context.Background()

	startDate := time.Now().AddDate(0, 0, -30).Truncate(24 * time.Hour)
	endDate := time.Now().AddDate(0, 0, -10).Truncate(24 * time.Hour)

	allocation, err := domain.NewDedicatedAllocation(roomID, userID, &startDate, &endDate)
	require.NoError(t, err, "failed to create past dedicated allocation")

	allocationEncx, err := domain.ProcessRoomAllocationEncx(ctx, crypto, allocation)
	require.NoError(t, err, "failed to encrypt allocation")

	return allocationEncx
}

// NewTestActiveDedicatedAllocationEncx creates a pre-encrypted dedicated allocation that is currently active
func NewTestActiveDedicatedAllocationEncx(t *testing.T, crypto encx.CryptoService, roomID, userID uuid.UUID) *domain.RoomAllocationEncx {
	t.Helper()

	ctx := context.Background()

	startDate := time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
	endDate := time.Now().AddDate(0, 0, 7).Truncate(24 * time.Hour)

	allocation, err := domain.NewDedicatedAllocation(roomID, userID, &startDate, &endDate)
	require.NoError(t, err, "failed to create active dedicated allocation")

	allocationEncx, err := domain.ProcessRoomAllocationEncx(ctx, crypto, allocation)
	require.NoError(t, err, "failed to encrypt allocation")

	return allocationEncx
}

// CountActiveAllocations counts all active allocations in the database
func CountActiveAllocations(t *testing.T, ctx context.Context, pool *pgxpool.Pool) int {
	t.Helper()

	var count int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM booking.room_allocations WHERE is_active = true").Scan(&count)
	require.NoError(t, err, "failed to count active allocations")

	return count
}
