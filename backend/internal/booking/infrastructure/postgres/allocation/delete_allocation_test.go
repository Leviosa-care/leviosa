package allocationRepository

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
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
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		allocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, userID)
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		// Execute soft delete
		err = repo.Delete(ctx, allocationEncx.ID)
		require.NoError(t, err)

		// Verify allocation still exists but is inactive
		deleted := GetAllocationEncxByID(t, ctx, testPool, allocationEncx.ID)
		assert.False(t, deleted.IsActive)
		assert.True(t, deleted.UpdatedAt.After(allocationEncx.CreatedAt))
	})

	t.Run("should return ErrRepositoryNotFound when deleting non-existent allocation", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)

		nonExistentID := uuid.New()

		err := repo.Delete(ctx, nonExistentID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("should successfully soft delete already inactive allocation", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		allocationEncx := NewTestInactiveSharedAllocationEncx(t, testCrypto, roomID, userID)
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		// Soft delete already inactive allocation
		err = repo.Delete(ctx, allocationEncx.ID)
		require.NoError(t, err)

		deleted := GetAllocationEncxByID(t, ctx, testPool, allocationEncx.ID)
		assert.False(t, deleted.IsActive)
	})

	t.Run("should update updated_at timestamp on soft delete", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		allocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, userID)
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		originalUpdatedAt := allocationEncx.UpdatedAt

		err = repo.Delete(ctx, allocationEncx.ID)
		require.NoError(t, err)

		deleted := GetAllocationEncxByID(t, ctx, testPool, allocationEncx.ID)
		assert.True(t, deleted.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("should not affect other allocations when soft deleting", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()

		// Create two allocations
		allocation1Encx := NewTestSharedAllocationEncx(t, testCrypto, roomID, user1ID)
		err := repo.Create(ctx, allocation1Encx)
		require.NoError(t, err)

		allocation2Encx := NewTestSharedAllocationEncx(t, testCrypto, roomID, user2ID)
		err = repo.Create(ctx, allocation2Encx)
		require.NoError(t, err)

		// Delete first allocation
		err = repo.Delete(ctx, allocation1Encx.ID)
		require.NoError(t, err)

		// Verify first is inactive
		deleted := GetAllocationEncxByID(t, ctx, testPool, allocation1Encx.ID)
		assert.False(t, deleted.IsActive)

		// Verify second is still active
		notDeleted := GetAllocationEncxByID(t, ctx, testPool, allocation2Encx.ID)
		assert.True(t, notDeleted.IsActive)
	})
}
