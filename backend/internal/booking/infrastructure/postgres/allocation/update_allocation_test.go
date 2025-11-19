package allocationRepository

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	ta "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/allocation"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestUpdate TEST_PATH=internal/booking/infrastructure/postgres/allocation/update_allocation_test.go

func TestUpdate(t *testing.T) {
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

	t.Run("should successfully update dedicated allocation period", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create initial allocation
		startDate := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 1, 0).Truncate(24 * time.Hour)
		allocation := ta.NewTestDedicatedAllocation(t, roomID, userID, startDate, endDate)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Update period
		newStartDate := time.Now().AddDate(0, 0, 7).Truncate(24 * time.Hour)
		newEndDate := time.Now().AddDate(0, 2, 0).Truncate(24 * time.Hour)
		err := allocation.UpdateDedicatedPeriod(&newStartDate, &newEndDate)
		require.NoError(t, err)

		// Execute update
		err = repo.Update(ctx, allocation)
		require.NoError(t, err)

		// Verify
		updated, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.NotNil(t, updated.StartDate)
		assert.NotNil(t, updated.EndDate)
		assert.WithinDuration(t, newStartDate, *updated.StartDate, time.Second)
		assert.WithinDuration(t, newEndDate, *updated.EndDate, time.Second)
		assert.True(t, updated.UpdatedAt.After(allocation.CreatedAt))
	})

	t.Run("should successfully deactivate allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		allocation := ta.NewTestSharedAllocation(t, roomID, userID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Deactivate
		allocation.Deactivate()

		err := repo.Update(ctx, allocation)
		require.NoError(t, err)

		updated, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.False(t, updated.IsActive)
	})

	t.Run("should successfully activate allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		allocation := ta.NewTestInactiveAllocation(t, roomID, userID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Activate
		allocation.Activate()

		err := repo.Update(ctx, allocation)
		require.NoError(t, err)

		updated, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.True(t, updated.IsActive)
	})

	t.Run("should return ErrRepositoryNotFound when updating non-existent allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		nonExistentAllocation := ta.NewTestSharedAllocation(t, roomID, userID)
		nonExistentAllocation.ID = uuid.New() // Never inserted

		err := repo.Update(ctx, nonExistentAllocation)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("should update updated_at timestamp", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		allocation := ta.NewTestSharedAllocation(t, roomID, userID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Wait a moment to ensure timestamp difference
		time.Sleep(10 * time.Millisecond)

		// Update
		allocation.UpdatedAt = time.Now()
		err := repo.Update(ctx, allocation)
		require.NoError(t, err)

		updated, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.True(t, updated.UpdatedAt.After(updated.CreatedAt))
	})

	t.Run("should update allocation to indefinite end date (NULL)", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create with end date
		startDate := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 1, 0).Truncate(24 * time.Hour)
		allocation := ta.NewTestDedicatedAllocation(t, roomID, userID, startDate, endDate)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Update to remove end date
		allocation.EndDate = nil
		allocation.UpdatedAt = time.Now()

		err := repo.Update(ctx, allocation)
		require.NoError(t, err)

		updated, err := ta.GetAllocationByID(t, ctx, allocation.ID, testPool)
		require.NoError(t, err)
		assert.Nil(t, updated.EndDate)
	})
}
