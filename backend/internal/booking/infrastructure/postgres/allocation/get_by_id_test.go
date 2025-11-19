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

// make test-func TEST_NAME=TestGetByID TEST_PATH=internal/booking/infrastructure/postgres/allocation/get_by_id_test.go

func TestGetByID(t *testing.T) {
	ctx := context.Background()

	// Setup helper to create building and room dependencies
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

	t.Run("should successfully retrieve shared allocation by ID", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create allocation
		allocation := ta.NewTestSharedAllocation(t, roomID, userID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Retrieve by ID
		retrieved, err := repo.GetByID(ctx, allocation.ID)
		require.NoError(t, err)
		assert.Equal(t, allocation.ID, retrieved.ID)
		assert.Equal(t, roomID, retrieved.RoomID)
		assert.Equal(t, userID, retrieved.UserID)
		assert.Equal(t, domain.AllocationTypeShared, retrieved.AllocationType)
		assert.Nil(t, retrieved.StartDate)
		assert.Nil(t, retrieved.EndDate)
		assert.True(t, retrieved.IsActive)
	})

	t.Run("should successfully retrieve dedicated allocation with end date", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		startDate := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 1, 0).Truncate(24 * time.Hour)

		allocation := ta.NewTestDedicatedAllocation(t, roomID, userID, startDate, endDate)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		retrieved, err := repo.GetByID(ctx, allocation.ID)
		require.NoError(t, err)
		assert.Equal(t, allocation.ID, retrieved.ID)
		assert.Equal(t, domain.AllocationTypeDedicated, retrieved.AllocationType)
		assert.NotNil(t, retrieved.StartDate)
		assert.NotNil(t, retrieved.EndDate)
		assert.WithinDuration(t, startDate, *retrieved.StartDate, time.Second)
		assert.WithinDuration(t, endDate, *retrieved.EndDate, time.Second)
	})

	t.Run("should successfully retrieve dedicated allocation without end date", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		startDate := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)

		allocation := ta.NewTestDedicatedAllocationIndefinite(t, roomID, userID, startDate)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		retrieved, err := repo.GetByID(ctx, allocation.ID)
		require.NoError(t, err)
		assert.Equal(t, allocation.ID, retrieved.ID)
		assert.Equal(t, domain.AllocationTypeDedicated, retrieved.AllocationType)
		assert.NotNil(t, retrieved.StartDate)
		assert.Nil(t, retrieved.EndDate) // NULL end date
		assert.WithinDuration(t, startDate, *retrieved.StartDate, time.Second)
	})

	t.Run("should return ErrRepositoryNotFound when allocation does not exist", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		nonExistentID := uuid.New()

		retrieved, err := repo.GetByID(ctx, nonExistentID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
		assert.Nil(t, retrieved)
	})

	t.Run("should retrieve inactive allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		allocation := ta.NewTestInactiveAllocation(t, roomID, userID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		retrieved, err := repo.GetByID(ctx, allocation.ID)
		require.NoError(t, err)
		assert.Equal(t, allocation.ID, retrieved.ID)
		assert.False(t, retrieved.IsActive)
	})

	t.Run("should retrieve allocation with all timestamps", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		allocation := ta.NewTestSharedAllocation(t, roomID, userID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		retrieved, err := repo.GetByID(ctx, allocation.ID)
		require.NoError(t, err)
		assert.NotZero(t, retrieved.CreatedAt)
		assert.NotZero(t, retrieved.UpdatedAt)
		assert.WithinDuration(t, allocation.CreatedAt, retrieved.CreatedAt, time.Second)
		assert.WithinDuration(t, allocation.UpdatedAt, retrieved.UpdatedAt, time.Second)
	})
}
