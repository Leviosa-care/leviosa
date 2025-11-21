package allocationRepository

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
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
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create initial allocation
		startDate := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 1, 0).Truncate(24 * time.Hour)
		allocationEncx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, userID, startDate, endDate)
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		// Update period
		newStartDate := time.Now().AddDate(0, 0, 7).Truncate(24 * time.Hour)
		newEndDate := time.Now().AddDate(0, 2, 0).Truncate(24 * time.Hour)
		allocationEncx.StartDate = &newStartDate
		allocationEncx.EndDate = &newEndDate
		allocationEncx.UpdatedAt = time.Now()

		// Execute update
		err = repo.Update(ctx, allocationEncx)
		require.NoError(t, err)

		// Verify
		updated := GetAllocationEncxByID(t, ctx, testPool, allocationEncx.ID)
		assert.NotNil(t, updated.StartDate)
		assert.NotNil(t, updated.EndDate)
		assert.WithinDuration(t, newStartDate, *updated.StartDate, time.Second)
		assert.WithinDuration(t, newEndDate, *updated.EndDate, time.Second)
		assert.True(t, updated.UpdatedAt.After(allocationEncx.CreatedAt))
	})

	t.Run("should successfully deactivate allocation", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		allocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, userID)
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		// Deactivate
		allocationEncx.IsActive = false
		allocationEncx.UpdatedAt = time.Now()

		err = repo.Update(ctx, allocationEncx)
		require.NoError(t, err)

		updated := GetAllocationEncxByID(t, ctx, testPool, allocationEncx.ID)
		assert.False(t, updated.IsActive)
	})

	t.Run("should successfully activate allocation", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		allocationEncx := NewTestInactiveSharedAllocationEncx(t, testCrypto, roomID, userID)
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		// Activate
		allocationEncx.IsActive = true
		allocationEncx.UpdatedAt = time.Now()

		err = repo.Update(ctx, allocationEncx)
		require.NoError(t, err)

		updated := GetAllocationEncxByID(t, ctx, testPool, allocationEncx.ID)
		assert.True(t, updated.IsActive)
	})

	t.Run("should return ErrRepositoryNotFound when updating non-existent allocation", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		nonExistentAllocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, userID)
		nonExistentAllocationEncx.ID = uuid.New() // Never inserted

		err := repo.Update(ctx, nonExistentAllocationEncx)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("should update updated_at timestamp", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		allocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, userID)
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		// Wait a moment to ensure timestamp difference
		time.Sleep(10 * time.Millisecond)

		// Update
		allocationEncx.UpdatedAt = time.Now()
		err = repo.Update(ctx, allocationEncx)
		require.NoError(t, err)

		updated := GetAllocationEncxByID(t, ctx, testPool, allocationEncx.ID)
		assert.True(t, updated.UpdatedAt.After(updated.CreatedAt))
	})

	t.Run("should update allocation to indefinite end date (NULL)", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create with end date
		startDate := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 1, 0).Truncate(24 * time.Hour)
		allocationEncx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, userID, startDate, endDate)
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		// Update to remove end date
		allocationEncx.EndDate = nil
		allocationEncx.UpdatedAt = time.Now()

		err = repo.Update(ctx, allocationEncx)
		require.NoError(t, err)

		updated := GetAllocationEncxByID(t, ctx, testPool, allocationEncx.ID)
		assert.Nil(t, updated.EndDate)
	})
}
