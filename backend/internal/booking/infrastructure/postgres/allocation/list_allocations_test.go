package allocationRepository

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	ta "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/allocation"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestList TEST_PATH=internal/booking/infrastructure/postgres/allocation/list_allocations_test.go

func TestList(t *testing.T) {
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

	t.Run("should list all allocations with empty filter", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		// Create multiple allocations
		allocation1 := ta.NewTestSharedAllocation(t, roomID, uuid.New())
		ta.InsertAllocation(t, ctx, allocation1, testPool)

		allocation2 := ta.NewTestSharedAllocation(t, roomID, uuid.New())
		ta.InsertAllocation(t, ctx, allocation2, testPool)

		allocations, err := repo.List(ctx, ports.RoomAllocationFilter{})
		require.NoError(t, err)
		assert.Len(t, allocations, 2)
	})

	t.Run("should filter by RoomID", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		room1ID := setupTestRoom(t)
		room2ID := setupTestRoom(t)

		allocation1 := ta.NewTestSharedAllocation(t, room1ID, uuid.New())
		ta.InsertAllocation(t, ctx, allocation1, testPool)

		allocation2 := ta.NewTestSharedAllocation(t, room2ID, uuid.New())
		ta.InsertAllocation(t, ctx, allocation2, testPool)

		filter := ports.RoomAllocationFilter{RoomID: &room1ID}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
		assert.Equal(t, room1ID, allocations[0].RoomID)
	})

	t.Run("should filter by UserID", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()

		allocation1 := ta.NewTestSharedAllocation(t, roomID, user1ID)
		ta.InsertAllocation(t, ctx, allocation1, testPool)

		allocation2 := ta.NewTestSharedAllocation(t, roomID, user2ID)
		ta.InsertAllocation(t, ctx, allocation2, testPool)

		filter := ports.RoomAllocationFilter{UserID: &user1ID}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
		assert.Equal(t, user1ID, allocations[0].UserID)
	})

	t.Run("should filter by AllocationType", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create shared
		sharedAllocation := ta.NewTestSharedAllocation(t, roomID, userID)
		ta.InsertAllocation(t, ctx, sharedAllocation, testPool)

		// Create dedicated
		dedicatedAllocation := ta.NewTestActiveDedicatedAllocation(t, roomID, uuid.New())
		ta.InsertAllocation(t, ctx, dedicatedAllocation, testPool)

		allocType := domain.AllocationTypeShared
		filter := ports.RoomAllocationFilter{AllocationType: &allocType}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
		assert.Equal(t, domain.AllocationTypeShared, allocations[0].AllocationType)
	})

	t.Run("should filter by IsActive", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		activeAllocation := ta.NewTestSharedAllocation(t, roomID, uuid.New())
		ta.InsertAllocation(t, ctx, activeAllocation, testPool)

		inactiveAllocation := ta.NewTestInactiveAllocation(t, roomID, uuid.New())
		ta.InsertAllocation(t, ctx, inactiveAllocation, testPool)

		isActive := true
		filter := ports.RoomAllocationFilter{IsActive: &isActive}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
		assert.True(t, allocations[0].IsActive)
	})

	t.Run("should filter by ActiveAt for shared allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		sharedAllocation := ta.NewTestSharedAllocation(t, roomID, uuid.New())
		ta.InsertAllocation(t, ctx, sharedAllocation, testPool)

		// Shared is always active
		activeAt := time.Now()
		filter := ports.RoomAllocationFilter{ActiveAt: &activeAt}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
	})

	t.Run("should filter by ActiveAt for dedicated allocation within period", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		dedicatedAllocation := ta.NewTestActiveDedicatedAllocation(t, roomID, uuid.New())
		ta.InsertAllocation(t, ctx, dedicatedAllocation, testPool)

		// Check at time within period
		activeAt := time.Now()
		filter := ports.RoomAllocationFilter{ActiveAt: &activeAt}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
	})

	t.Run("should not return dedicated allocation outside ActiveAt period", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		pastAllocation := ta.NewTestPastDedicatedAllocation(t, roomID, uuid.New())
		ta.InsertAllocation(t, ctx, pastAllocation, testPool)

		// Check at current time (after period)
		activeAt := time.Now()
		filter := ports.RoomAllocationFilter{ActiveAt: &activeAt}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Empty(t, allocations)
	})

	t.Run("should filter by OverlapsWith for dedicated allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		allocation := ta.NewTestActiveDedicatedAllocation(t, roomID, uuid.New())
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Create range that overlaps with allocation
		overlapRange := ports.TimeRange{
			Start: time.Now().AddDate(0, 0, -1),
			End:   time.Now().AddDate(0, 0, 1),
		}
		filter := ports.RoomAllocationFilter{OverlapsWith: &overlapRange}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
	})

	t.Run("should not return allocation when OverlapsWith has no overlap", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		allocation := ta.NewTestPastDedicatedAllocation(t, roomID, uuid.New())
		ta.InsertAllocation(t, ctx, allocation, testPool)

		// Create range in the future (no overlap with past allocation)
		overlapRange := ports.TimeRange{
			Start: time.Now().AddDate(0, 1, 0),
			End:   time.Now().AddDate(0, 2, 0),
		}
		filter := ports.RoomAllocationFilter{OverlapsWith: &overlapRange}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Empty(t, allocations)
	})

	t.Run("should apply pagination with Limit", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		// Create 5 allocations
		for i := 0; i < 5; i++ {
			allocation := ta.NewTestSharedAllocation(t, roomID, uuid.New())
			ta.InsertAllocation(t, ctx, allocation, testPool)
		}

		filter := ports.RoomAllocationFilter{Limit: 2}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, allocations, 2)
	})

	t.Run("should apply pagination with Offset", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		// Create allocations with staggered times for consistent ordering
		allocations := make([]*domain.RoomAllocation, 5)
		for i := 0; i < 5; i++ {
			allocation := ta.NewTestSharedAllocation(t, roomID, uuid.New())
			allocation.CreatedAt = time.Now().Add(-time.Duration(4-i) * time.Hour)
			ta.InsertAllocation(t, ctx, allocation, testPool)
			allocations[i] = allocation
		}

		filter := ports.RoomAllocationFilter{Offset: 2}
		result, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, result, 3) // 5 total - 2 offset = 3
	})

	t.Run("should order by created_at DESC by default", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		first := ta.NewTestSharedAllocation(t, roomID, uuid.New())
		first.CreatedAt = time.Now().Add(-2 * time.Hour)
		ta.InsertAllocation(t, ctx, first, testPool)

		second := ta.NewTestSharedAllocation(t, roomID, uuid.New())
		second.CreatedAt = time.Now().Add(-1 * time.Hour)
		ta.InsertAllocation(t, ctx, second, testPool)

		third := ta.NewTestSharedAllocation(t, roomID, uuid.New())
		third.CreatedAt = time.Now()
		ta.InsertAllocation(t, ctx, third, testPool)

		allocations, err := repo.List(ctx, ports.RoomAllocationFilter{})
		require.NoError(t, err)
		assert.Equal(t, third.ID, allocations[0].ID)
		assert.Equal(t, second.ID, allocations[1].ID)
		assert.Equal(t, first.ID, allocations[2].ID)
	})

	t.Run("should order by start_date ASC", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		allocation1 := ta.NewTestFutureDedicatedAllocation(t, roomID, uuid.New())
		startDate1 := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		allocation1.StartDate = &startDate1
		ta.InsertAllocation(t, ctx, allocation1, testPool)

		allocation2 := ta.NewTestFutureDedicatedAllocation(t, roomID, uuid.New())
		startDate2 := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		allocation2.StartDate = &startDate2
		ta.InsertAllocation(t, ctx, allocation2, testPool)

		filter := ports.RoomAllocationFilter{
			OrderBy:        "start_date",
			OrderDirection: "asc",
		}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Equal(t, allocation2.ID, allocations[0].ID)
		assert.Equal(t, allocation1.ID, allocations[1].ID)
	})

	t.Run("should combine multiple filters", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		room1ID := setupTestRoom(t)
		room2ID := setupTestRoom(t)
		room3ID := setupTestRoom(t)
		userID := uuid.New()

		// Create allocation that matches all filters
		matchingAllocation := ta.NewTestSharedAllocation(t, room1ID, userID)
		ta.InsertAllocation(t, ctx, matchingAllocation, testPool)

		// Create allocations that don't match
		ta.InsertAllocation(t, ctx, ta.NewTestSharedAllocation(t, room2ID, userID), testPool)
		ta.InsertAllocation(t, ctx, ta.NewTestInactiveAllocation(t, room3ID, userID), testPool)

		allocType := domain.AllocationTypeShared
		isActive := true
		filter := ports.RoomAllocationFilter{
			RoomID:         &room1ID,
			UserID:         &userID,
			AllocationType: &allocType,
			IsActive:       &isActive,
		}

		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
		assert.Equal(t, matchingAllocation.ID, allocations[0].ID)
	})

	t.Run("should return empty slice when no matches", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		nonExistentRoomID := uuid.New()
		filter := ports.RoomAllocationFilter{RoomID: &nonExistentRoomID}

		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Empty(t, allocations)
	})
}
