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

// make test-func TEST_NAME=TestGetActiveAllocationForPartnerAndRoom TEST_PATH=internal/booking/infrastructure/postgres/allocation/get_active_allocation_for_partner_and_room_test.go

func TestGetActiveAllocationForPartnerAndRoom(t *testing.T) {
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

	t.Run("should retrieve shared allocation (always active)", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		allocation := ta.NewTestSharedAllocation(t, roomID, userID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Check at any time (shared is always active)
		at := time.Now()
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, userID, roomID, at)
		require.NoError(t, err)
		assert.Equal(t, allocation.ID, retrieved.ID)
		assert.Equal(t, domain.AllocationTypeShared, retrieved.AllocationType)
	})

	t.Run("should retrieve dedicated allocation within active period", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		startDate := time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 0, 7).Truncate(24 * time.Hour)

		allocation := ta.NewTestDedicatedAllocation(t, roomID, userID, startDate, endDate)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Check at current time (within period)
		at := time.Now()
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, userID, roomID, at)
		require.NoError(t, err)
		assert.Equal(t, allocation.ID, retrieved.ID)
		assert.Equal(t, domain.AllocationTypeDedicated, retrieved.AllocationType)
	})

	t.Run("should not retrieve dedicated allocation before start_date", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		startDate := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 1, 0).Truncate(24 * time.Hour)

		allocation := ta.NewTestFutureDedicatedAllocation(t, roomID, userID)
		allocation.StartDate = &startDate
		allocation.EndDate = &endDate
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Check at current time (before start_date)
		at := time.Now()
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, userID, roomID, at)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
		assert.Nil(t, retrieved)
	})

	t.Run("should not retrieve dedicated allocation after end_date", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		allocation := ta.NewTestPastDedicatedAllocation(t, roomID, userID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Check at current time (after end_date)
		at := time.Now()
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, userID, roomID, at)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
		assert.Nil(t, retrieved)
	})

	t.Run("should retrieve dedicated allocation with NULL end_date (indefinite)", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		startDate := time.Now().AddDate(0, 0, -30).Truncate(24 * time.Hour)

		allocation := ta.NewTestDedicatedAllocationIndefinite(t, roomID, userID, startDate)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Check at current time and future time
		at := time.Now().AddDate(0, 1, 0) // 1 month in future
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, userID, roomID, at)
		require.NoError(t, err)
		assert.Equal(t, allocation.ID, retrieved.ID)
		assert.Nil(t, retrieved.EndDate)
	})

	t.Run("should prioritize dedicated over shared when both exist", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create shared allocation
		sharedAllocation := ta.NewTestSharedAllocation(t, roomID, userID)
		ta.InsertAllocation(t, ctx, sharedAllocation, testPool)

		// Create dedicated allocation (currently active)
		dedicatedAllocation := ta.NewTestActiveDedicatedAllocation(t, roomID, userID)
		ta.InsertAllocation(t, ctx, dedicatedAllocation, testPool)

		count := ta.CountActiveAllocations(t, ctx, testPool)
		require.Equal(t, 2, count, "Should have 2 active allocations in the database set up.")

		// Should return dedicated (higher priority)
		at := time.Now()
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, userID, roomID, at)
		require.NoError(t, err)
		assert.Equal(t, dedicatedAllocation.ID, retrieved.ID)
		assert.Equal(t, domain.AllocationTypeDedicated, retrieved.AllocationType)
	})

	t.Run("should not retrieve inactive allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		allocation := ta.NewTestInactiveAllocation(t, roomID, userID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		at := time.Now()
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, userID, roomID, at)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
		assert.Nil(t, retrieved)
	})

	t.Run("should return not found when no allocation exists", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		nonExistentUserID := uuid.New()
		nonExistentRoomID := uuid.New()

		at := time.Now()
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, nonExistentUserID, nonExistentRoomID, at)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
		assert.Nil(t, retrieved)
	})

	t.Run("should not retrieve allocation for different user", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()

		allocation := ta.NewTestSharedAllocation(t, roomID, user1ID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		at := time.Now()
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, user2ID, roomID, at)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
		assert.Nil(t, retrieved)
	})

	t.Run("should not retrieve allocation for different room", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		room1ID := setupTestRoom(t)
		room2ID := setupTestRoom(t)
		userID := uuid.New()

		allocation := ta.NewTestSharedAllocation(t, room1ID, userID)
		ta.InsertAllocation(t, ctx, allocation, testPool)

		at := time.Now()
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, userID, room2ID, at)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
		assert.Nil(t, retrieved)
	})
}
