package allocationRepository

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	ta "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/allocation"
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
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Setup dependencies
		_, roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create allocation
		allocation := ta.NewTestSharedAllocation(t, roomID, userID)

		// Execute
		err := repo.Create(ctx, allocation)
		require.NoError(t, err)

		// Verify
		created, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, allocation.ID, created.ID)
		assert.Equal(t, roomID, created.RoomID)
		assert.Equal(t, userID, created.UserID)
		assert.Equal(t, domain.AllocationTypeShared, created.AllocationType)
		assert.Nil(t, created.StartDate)
		assert.Nil(t, created.EndDate)
		assert.True(t, created.IsActive)
	})

	t.Run("should successfully create dedicated allocation with end date", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		_, roomID := setupTestRoom(t)
		userID := uuid.New()

		startDate := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 1, 0).Truncate(24 * time.Hour)

		allocation := ta.NewTestDedicatedAllocation(t, roomID, userID, startDate, endDate)

		err := repo.Create(ctx, allocation)
		require.NoError(t, err)

		created, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, allocation.ID, created.ID)
		assert.Equal(t, domain.AllocationTypeDedicated, created.AllocationType)
		assert.NotNil(t, created.StartDate)
		assert.NotNil(t, created.EndDate)
		assert.WithinDuration(t, startDate, *created.StartDate, time.Second)
		assert.WithinDuration(t, endDate, *created.EndDate, time.Second)
		assert.True(t, created.IsActive)
	})

	t.Run("should successfully create dedicated allocation without end date (indefinite)", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		_, roomID := setupTestRoom(t)
		userID := uuid.New()

		startDate := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)

		allocation := ta.NewTestDedicatedAllocationIndefinite(t, roomID, userID, startDate)

		err := repo.Create(ctx, allocation)
		require.NoError(t, err)

		created, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.Equal(t, allocation.ID, created.ID)
		assert.Equal(t, domain.AllocationTypeDedicated, created.AllocationType)
		assert.NotNil(t, created.StartDate)
		assert.Nil(t, created.EndDate) // NULL end date
		assert.WithinDuration(t, startDate, *created.StartDate, time.Second)
		assert.True(t, created.IsActive)
	})

	t.Run("should return error when creating allocation with duplicate ID", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		_, roomID := setupTestRoom(t)
		userID := uuid.New()

		allocation := ta.NewTestSharedAllocation(t, roomID, userID)

		// Create first time
		err := repo.Create(ctx, allocation)
		require.NoError(t, err)

		// Try to create again with same ID
		duplicateAllocation := ta.NewTestSharedAllocation(t, roomID, uuid.New())
		duplicateAllocation.ID = allocation.ID

		err = repo.Create(ctx, duplicateAllocation)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrUniqueViolation)
	})

	t.Run("should return error when room does not exist (foreign key violation)", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		nonExistentRoomID := uuid.New()
		userID := uuid.New()

		allocation := ta.NewTestSharedAllocation(t, nonExistentRoomID, userID)

		err := repo.Create(ctx, allocation)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrForeignKeyViolation)
	})

	t.Run("should return error when creating duplicate allocation (same room, user, type)", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		_, roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create first allocation
		allocation1 := ta.NewTestSharedAllocation(t, roomID, userID)
		err := repo.Create(ctx, allocation1)
		require.NoError(t, err)

		// Try to create duplicate (same room, user, type)
		allocation2 := ta.NewTestSharedAllocation(t, roomID, userID)
		err = repo.Create(ctx, allocation2)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrUniqueViolation)
	})

	t.Run("should successfully create different allocation types for same room and user", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		_, roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create shared allocation
		sharedAllocation := ta.NewTestSharedAllocation(t, roomID, userID)
		err := repo.Create(ctx, sharedAllocation)
		require.NoError(t, err)

		// Create dedicated allocation (different type, so should succeed)
		startDate := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 1, 0).Truncate(24 * time.Hour)
		dedicatedAllocation := ta.NewTestDedicatedAllocation(t, roomID, userID, startDate, endDate)
		err = repo.Create(ctx, dedicatedAllocation)
		require.NoError(t, err)

		// Verify both exist
		count := ta.CountAllocationsByUserID(t, ctx, userID, testPool)
		assert.Equal(t, 2, count)
	})

	t.Run("should successfully create multiple allocations for different users in same room", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		_, roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()

		allocation1 := ta.NewTestSharedAllocation(t, roomID, user1ID)
		err := repo.Create(ctx, allocation1)
		require.NoError(t, err)

		allocation2 := ta.NewTestSharedAllocation(t, roomID, user2ID)
		err = repo.Create(ctx, allocation2)
		require.NoError(t, err)

		count := ta.CountAllocationsByRoomID(t, ctx, roomID, testPool)
		assert.Equal(t, 2, count)
	})

	t.Run("should preserve created_at and updated_at timestamps", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		_, roomID := setupTestRoom(t)
		userID := uuid.New()

		allocation := ta.NewTestSharedAllocation(t, roomID, userID)
		originalCreatedAt := allocation.CreatedAt
		originalUpdatedAt := allocation.UpdatedAt

		err := repo.Create(ctx, allocation)
		require.NoError(t, err)

		created, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.WithinDuration(t, originalCreatedAt, created.CreatedAt, time.Second)
		assert.WithinDuration(t, originalUpdatedAt, created.UpdatedAt, time.Second)
	})

	t.Run("should create inactive allocation successfully", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		_, roomID := setupTestRoom(t)
		userID := uuid.New()

		allocation := ta.NewTestInactiveAllocation(t, roomID, userID)

		err := repo.Create(ctx, allocation)
		require.NoError(t, err)

		created, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.False(t, created.IsActive)
	})
}
