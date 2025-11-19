package allocationRepository

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	ta "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/allocation"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestDelete TEST_PATH=internal/booking/infrastructure/postgres/allocation/delete_allocation_test.go

func TestDelete(t *testing.T) {
	ctx := context.Background()

	setupTestRoom := func(t *testing.T) uuid.UUID {
		t.Helper()

		building := tb.NewTestBuildingEncx(t)
		err := tb.InsertBuildingEncx(t, ctx, testPool, building)
		require.NoError(t, err)

		room := tr.NewTestRoomEncxWithBuilding(t, building.ID)
		err = tr.InsertRoomEncx(t, ctx, testPool, room)
		require.NoError(t, err)

		return room.ID
	}

	t.Run("should soft delete allocation (mark as inactive)", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		allocation := ta.NewTestSharedAllocation(t, roomID, userID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Execute soft delete
		err := repo.Delete(ctx, allocation.ID)
		require.NoError(t, err)

		// Verify allocation still exists but is inactive
		deleted, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.False(t, deleted.IsActive)
		assert.True(t, deleted.UpdatedAt.After(allocation.CreatedAt))
	})

	t.Run("should return ErrRepositoryNotFound when deleting non-existent allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)

		nonExistentID := uuid.New()

		err := repo.Delete(ctx, nonExistentID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("should successfully soft delete already inactive allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		allocation := ta.NewTestInactiveAllocation(t, roomID, userID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Soft delete already inactive allocation
		err := repo.Delete(ctx, allocation.ID)
		require.NoError(t, err)

		deleted, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.False(t, deleted.IsActive)
	})

	t.Run("should update updated_at timestamp on soft delete", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		allocation := ta.NewTestSharedAllocation(t, roomID, userID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		originalUpdatedAt := allocation.UpdatedAt

		err := repo.Delete(ctx, allocation.ID)
		require.NoError(t, err)

		deleted, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.True(t, deleted.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("should not affect other allocations when soft deleting", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()

		// Create two allocations
		allocation1 := ta.NewTestSharedAllocation(t, roomID, user1ID)
		ta.InsertAllocation(t, ctx, allocation1, testPool)

		allocation2 := ta.NewTestSharedAllocation(t, roomID, user2ID)
		ta.InsertAllocation(t, ctx, allocation2, testPool)

		// Delete first allocation
		err := repo.Delete(ctx, allocation1.ID)
		require.NoError(t, err)

		// Verify first is inactive
		deleted, err := ta.GetAllocationByID(t, ctx, allocation1.ID, testPool)
		require.NoError(t, err)
		assert.False(t, deleted.IsActive)

		// Verify second is still active
		notDeleted, err := ta.GetAllocationByID(t, ctx, allocation2.ID, testPool)
		require.NoError(t, err)
		assert.True(t, notDeleted.IsActive)
	})
}
