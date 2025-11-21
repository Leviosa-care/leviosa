package allocationRepository

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreate TEST_PATH=internal/booking/infrastructure/postgres/allocation/create_allocation_test.go

func TestCreate(t *testing.T) {
	ctx := context.Background()

	// Setup helper to create building and room dependencies
	setupTestRoom := func(t *testing.T) (uuid.UUID, uuid.UUID) {
		t.Helper()

		// Create building (needed for room foreign key)
		building := tb.NewTestBuildingEncx(t)
		err := tb.InsertBuildingEncx(t, ctx, testPool, building)
		require.NoError(t, err)

		// Create room
		room := tr.NewTestRoomEncxWithBuilding(t, building.ID)
		err = tr.InsertRoomEncx(t, ctx, testPool, room)
		require.NoError(t, err)

		return building.ID, room.ID
	}

	t.Run("should successfully create shared allocation", func(t *testing.T) {
		// Clean tables
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Setup dependencies
		_, roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create encrypted allocation
		allocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, userID)

		// Execute
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		// Verify
		created := GetAllocationEncxByID(t, ctx, testPool, allocationEncx.ID)
		assert.Equal(t, allocationEncx.ID, created.ID)
		assert.Equal(t, roomID, created.RoomID)
		assert.Equal(t, allocationEncx.UserIDHash, created.UserIDHash)
		assert.NotEmpty(t, created.UserIDEncrypted)
		assert.Equal(t, domain.AllocationTypeShared, created.AllocationType)
		assert.Nil(t, created.StartDate)
		assert.Nil(t, created.EndDate)
		assert.True(t, created.IsActive)
	})

	t.Run("should successfully create dedicated allocation with end date", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		_, roomID := setupTestRoom(t)
		userID := uuid.New()

		startDate := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 1, 0).Truncate(24 * time.Hour)

		allocationEncx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, userID, startDate, endDate)

		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		created := GetAllocationEncxByID(t, ctx, testPool, allocationEncx.ID)
		assert.Equal(t, allocationEncx.ID, created.ID)
		assert.Equal(t, domain.AllocationTypeDedicated, created.AllocationType)
		assert.NotNil(t, created.StartDate)
		assert.NotNil(t, created.EndDate)
		assert.WithinDuration(t, startDate, *created.StartDate, time.Second)
		assert.WithinDuration(t, endDate, *created.EndDate, time.Second)
		assert.True(t, created.IsActive)
	})

	t.Run("should successfully create dedicated allocation without end date (indefinite)", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		_, roomID := setupTestRoom(t)
		userID := uuid.New()

		startDate := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)

		allocationEncx := NewTestDedicatedAllocationEncxWithNilEndDate(t, testCrypto, roomID, userID, startDate)

		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		created := GetAllocationEncxByID(t, ctx, testPool, allocationEncx.ID)
		assert.Equal(t, allocationEncx.ID, created.ID)
		assert.Equal(t, domain.AllocationTypeDedicated, created.AllocationType)
		assert.NotNil(t, created.StartDate)
		assert.Nil(t, created.EndDate) // NULL end date
		assert.WithinDuration(t, startDate, *created.StartDate, time.Second)
		assert.True(t, created.IsActive)
	})

	t.Run("should return error when creating allocation with duplicate ID", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		_, roomID := setupTestRoom(t)
		userID := uuid.New()

		allocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, userID)

		// Create first time
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		// Try to create again with same ID
		duplicateAllocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, uuid.New())
		duplicateAllocationEncx.ID = allocationEncx.ID

		err = repo.Create(ctx, duplicateAllocationEncx)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrUniqueViolation)
	})

	t.Run("should return error when room does not exist (foreign key violation)", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		nonExistentRoomID := uuid.New()
		userID := uuid.New()

		allocationEncx := NewTestSharedAllocationEncx(t, testCrypto, nonExistentRoomID, userID)

		err := repo.Create(ctx, allocationEncx)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrForeignKeyViolation)
	})

	t.Run("should successfully create different allocation types for same room and user", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		_, roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create shared allocation
		sharedAllocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, userID)
		err := repo.Create(ctx, sharedAllocationEncx)
		require.NoError(t, err)

		// Create dedicated allocation (different type, so should succeed)
		startDate := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 1, 0).Truncate(24 * time.Hour)
		dedicatedAllocationEncx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, userID, startDate, endDate)
		err = repo.Create(ctx, dedicatedAllocationEncx)
		require.NoError(t, err)

		// Verify both exist using hash
		userIDHash := ComputeUserIDHash(t, ctx, testCrypto, userID)
		count := CountAllocationsByUserIDHash(t, ctx, testPool, userIDHash)
		assert.Equal(t, 2, count)
	})

	t.Run("should successfully create multiple allocations for different users in same room", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		_, roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()

		allocation1Encx := NewTestSharedAllocationEncx(t, testCrypto, roomID, user1ID)
		err := repo.Create(ctx, allocation1Encx)
		require.NoError(t, err)

		allocation2Encx := NewTestSharedAllocationEncx(t, testCrypto, roomID, user2ID)
		err = repo.Create(ctx, allocation2Encx)
		require.NoError(t, err)

		count := CountAllocationsByRoomID(t, ctx, testPool, roomID)
		assert.Equal(t, 2, count)
	})

	t.Run("should preserve created_at and updated_at timestamps", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		_, roomID := setupTestRoom(t)
		userID := uuid.New()

		allocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, userID)
		originalCreatedAt := allocationEncx.CreatedAt
		originalUpdatedAt := allocationEncx.UpdatedAt

		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		created := GetAllocationEncxByID(t, ctx, testPool, allocationEncx.ID)
		assert.WithinDuration(t, originalCreatedAt, created.CreatedAt, time.Second)
		assert.WithinDuration(t, originalUpdatedAt, created.UpdatedAt, time.Second)
	})

	t.Run("should create inactive allocation successfully", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		_, roomID := setupTestRoom(t)
		userID := uuid.New()

		allocationEncx := NewTestInactiveSharedAllocationEncx(t, testCrypto, roomID, userID)

		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		created := GetAllocationEncxByID(t, ctx, testPool, allocationEncx.ID)
		assert.False(t, created.IsActive)
	})
}
