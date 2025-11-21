package allocationRepository

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
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
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		// Create multiple allocations
		allocation1Encx := NewTestSharedAllocationEncx(t, testCrypto, roomID, uuid.New())
		err := repo.Create(ctx, allocation1Encx)
		require.NoError(t, err)

		allocation2Encx := NewTestSharedAllocationEncx(t, testCrypto, roomID, uuid.New())
		err = repo.Create(ctx, allocation2Encx)
		require.NoError(t, err)

		allocations, err := repo.List(ctx, ports.RoomAllocationFilter{})
		require.NoError(t, err)
		assert.Len(t, allocations, 2)
	})

	t.Run("should filter by RoomID", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		room1ID := setupTestRoom(t)
		room2ID := setupTestRoom(t)

		allocation1Encx := NewTestSharedAllocationEncx(t, testCrypto, room1ID, uuid.New())
		err := repo.Create(ctx, allocation1Encx)
		require.NoError(t, err)

		allocation2Encx := NewTestSharedAllocationEncx(t, testCrypto, room2ID, uuid.New())
		err = repo.Create(ctx, allocation2Encx)
		require.NoError(t, err)

		filter := ports.RoomAllocationFilter{RoomID: &room1ID}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
		assert.Equal(t, room1ID, allocations[0].RoomID)
	})

	t.Run("should filter by UserIDHash", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()
		user1IDHash := ComputeUserIDHash(t, ctx, testCrypto, user1ID)

		allocation1Encx := NewTestSharedAllocationEncx(t, testCrypto, roomID, user1ID)
		err := repo.Create(ctx, allocation1Encx)
		require.NoError(t, err)

		allocation2Encx := NewTestSharedAllocationEncx(t, testCrypto, roomID, user2ID)
		err = repo.Create(ctx, allocation2Encx)
		require.NoError(t, err)

		filter := ports.RoomAllocationFilter{UserIDHash: &user1IDHash}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
		assert.Equal(t, user1IDHash, allocations[0].UserIDHash)
	})

	t.Run("should filter by AllocationType", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create shared
		sharedAllocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, userID)
		err := repo.Create(ctx, sharedAllocationEncx)
		require.NoError(t, err)

		// Create dedicated
		dedicatedAllocationEncx := NewTestActiveDedicatedAllocationEncx(t, testCrypto, roomID, uuid.New())
		err = repo.Create(ctx, dedicatedAllocationEncx)
		require.NoError(t, err)

		allocType := domain.AllocationTypeShared
		filter := ports.RoomAllocationFilter{AllocationType: &allocType}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
		assert.Equal(t, domain.AllocationTypeShared, allocations[0].AllocationType)
	})

	t.Run("should filter by IsActive", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		activeAllocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, uuid.New())
		err := repo.Create(ctx, activeAllocationEncx)
		require.NoError(t, err)

		inactiveAllocationEncx := NewTestInactiveSharedAllocationEncx(t, testCrypto, roomID, uuid.New())
		err = repo.Create(ctx, inactiveAllocationEncx)
		require.NoError(t, err)

		isActive := true
		filter := ports.RoomAllocationFilter{IsActive: &isActive}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
		assert.True(t, allocations[0].IsActive)
	})

	t.Run("should filter by ActiveAt for shared allocation", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		sharedAllocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, uuid.New())
		err := repo.Create(ctx, sharedAllocationEncx)
		require.NoError(t, err)

		// Shared is always active
		activeAt := time.Now()
		filter := ports.RoomAllocationFilter{ActiveAt: &activeAt}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
	})

	t.Run("should filter by ActiveAt for dedicated allocation within period", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		dedicatedAllocationEncx := NewTestActiveDedicatedAllocationEncx(t, testCrypto, roomID, uuid.New())
		err := repo.Create(ctx, dedicatedAllocationEncx)
		require.NoError(t, err)

		// Check at time within period
		activeAt := time.Now()
		filter := ports.RoomAllocationFilter{ActiveAt: &activeAt}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
	})

	t.Run("should not return dedicated allocation outside ActiveAt period", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		pastAllocationEncx := NewTestPastDedicatedAllocationEncx(t, testCrypto, roomID, uuid.New())
		err := repo.Create(ctx, pastAllocationEncx)
		require.NoError(t, err)

		// Check at current time (after period)
		activeAt := time.Now()
		filter := ports.RoomAllocationFilter{ActiveAt: &activeAt}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Empty(t, allocations)
	})

	t.Run("should filter by OverlapsWith for dedicated allocation", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		allocationEncx := NewTestActiveDedicatedAllocationEncx(t, testCrypto, roomID, uuid.New())
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

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
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		allocationEncx := NewTestPastDedicatedAllocationEncx(t, testCrypto, roomID, uuid.New())
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

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
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		// Create 5 allocations
		for i := 0; i < 5; i++ {
			allocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, uuid.New())
			err := repo.Create(ctx, allocationEncx)
			require.NoError(t, err)
		}

		filter := ports.RoomAllocationFilter{Limit: 2}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, allocations, 2)
	})

	t.Run("should apply pagination with Offset", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		// Create allocations with staggered times for consistent ordering
		for i := 0; i < 5; i++ {
			allocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, uuid.New())
			allocationEncx.CreatedAt = time.Now().Add(-time.Duration(4-i) * time.Hour)
			err := repo.Create(ctx, allocationEncx)
			require.NoError(t, err)
		}

		filter := ports.RoomAllocationFilter{Offset: 2}
		result, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, result, 3) // 5 total - 2 offset = 3
	})

	t.Run("should order by created_at DESC by default", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		first := NewTestSharedAllocationEncx(t, testCrypto, roomID, uuid.New())
		first.CreatedAt = time.Now().Add(-2 * time.Hour)
		err := repo.Create(ctx, first)
		require.NoError(t, err)

		second := NewTestSharedAllocationEncx(t, testCrypto, roomID, uuid.New())
		second.CreatedAt = time.Now().Add(-1 * time.Hour)
		err = repo.Create(ctx, second)
		require.NoError(t, err)

		third := NewTestSharedAllocationEncx(t, testCrypto, roomID, uuid.New())
		third.CreatedAt = time.Now()
		err = repo.Create(ctx, third)
		require.NoError(t, err)

		allocations, err := repo.List(ctx, ports.RoomAllocationFilter{})
		require.NoError(t, err)
		assert.Equal(t, third.ID, allocations[0].ID)
		assert.Equal(t, second.ID, allocations[1].ID)
		assert.Equal(t, first.ID, allocations[2].ID)
	})

	t.Run("should order by start_date ASC", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		allocation1Encx := NewTestFutureDedicatedAllocationEncx(t, testCrypto, roomID, uuid.New())
		startDate1 := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
		allocation1Encx.StartDate = &startDate1
		err := repo.Create(ctx, allocation1Encx)
		require.NoError(t, err)

		allocation2Encx := NewTestFutureDedicatedAllocationEncx(t, testCrypto, roomID, uuid.New())
		startDate2 := time.Now().AddDate(0, 0, 12).Truncate(24 * time.Hour)
		allocation2Encx.StartDate = &startDate2
		err = repo.Create(ctx, allocation2Encx)
		require.NoError(t, err)

		filter := ports.RoomAllocationFilter{
			OrderBy:        "start_date",
			OrderDirection: "asc",
		}
		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Equal(t, allocation2Encx.ID, allocations[0].ID)
		assert.Equal(t, allocation1Encx.ID, allocations[1].ID)
	})

	t.Run("should combine multiple filters", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		room1ID := setupTestRoom(t)
		room2ID := setupTestRoom(t)
		room3ID := setupTestRoom(t)
		userID := uuid.New()
		userIDHash := ComputeUserIDHash(t, ctx, testCrypto, userID)

		// Create allocation that matches all filters
		matchingAllocationEncx := NewTestSharedAllocationEncx(t, testCrypto, room1ID, userID)
		err := repo.Create(ctx, matchingAllocationEncx)
		require.NoError(t, err)

		// Create allocations that don't match
		allocation2Encx := NewTestSharedAllocationEncx(t, testCrypto, room2ID, userID)
		err = repo.Create(ctx, allocation2Encx)
		require.NoError(t, err)

		allocation3Encx := NewTestInactiveSharedAllocationEncx(t, testCrypto, room3ID, userID)
		err = repo.Create(ctx, allocation3Encx)
		require.NoError(t, err)

		allocType := domain.AllocationTypeShared
		isActive := true
		filter := ports.RoomAllocationFilter{
			RoomID:         &room1ID,
			UserIDHash:     &userIDHash,
			AllocationType: &allocType,
			IsActive:       &isActive,
		}

		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
		assert.Equal(t, matchingAllocationEncx.ID, allocations[0].ID)
	})

	t.Run("should return empty slice when no matches", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		nonExistentRoomID := uuid.New()
		filter := ports.RoomAllocationFilter{RoomID: &nonExistentRoomID}

		allocations, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Empty(t, allocations)
	})
}
