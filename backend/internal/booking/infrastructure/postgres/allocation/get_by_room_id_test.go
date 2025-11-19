package allocationRepository

import (
	"context"
	"testing"
	"time"

	ta "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/allocation"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetByRoomID TEST_PATH=internal/booking/infrastructure/postgres/allocation/get_by_room_id_test.go

func TestGetByRoomID(t *testing.T) {
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

	t.Run("should retrieve all allocations for a room (activeOnly=false)", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()

		// Create active and inactive allocations for the same room
		activeAllocation := ta.NewTestSharedAllocation(t, roomID, user1ID)
		ta.InsertAllocation(t, ctx, activeAllocation, testPool)

		inactiveAllocation := ta.NewTestInactiveAllocation(t, roomID, user2ID)
		ta.InsertAllocation(t, ctx, inactiveAllocation, testPool)

		// Retrieve all
		allocations, err := repo.GetByRoomID(ctx, roomID, false)
		require.NoError(t, err)
		assert.Len(t, allocations, 2)
	})

	t.Run("should retrieve only active allocations for a room (activeOnly=true)", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()

		// Create active and inactive allocations
		activeAllocation := ta.NewTestSharedAllocation(t, roomID, user1ID)
		ta.InsertAllocation(t, ctx, activeAllocation, testPool)

		inactiveAllocation := ta.NewTestInactiveAllocation(t, roomID, user2ID)
		ta.InsertAllocation(t, ctx, inactiveAllocation, testPool)

		// Retrieve only active
		allocations, err := repo.GetByRoomID(ctx, roomID, true)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
		assert.True(t, allocations[0].IsActive)
		assert.Equal(t, activeAllocation.ID, allocations[0].ID)
	})

	t.Run("should return empty slice when room has no allocations", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		setupTestRoom(t) // Create room but no allocations
		nonExistentRoomID := uuid.New()

		allocations, err := repo.GetByRoomID(ctx, nonExistentRoomID, false)
		require.NoError(t, err)
		assert.Empty(t, allocations)
	})

	t.Run("should order allocations by created_at DESC", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		// Create allocations with staggered created_at
		first := ta.NewTestSharedAllocation(t, roomID, uuid.New())
		first.CreatedAt = time.Now().Add(-2 * time.Hour)
		ta.InsertAllocation(t, ctx, first, testPool)

		second := ta.NewTestSharedAllocation(t, roomID, uuid.New())
		second.CreatedAt = time.Now().Add(-1 * time.Hour)
		ta.InsertAllocation(t, ctx, second, testPool)

		third := ta.NewTestSharedAllocation(t, roomID, uuid.New())
		third.CreatedAt = time.Now()
		ta.InsertAllocation(t, ctx, third, testPool)

		allocations, err := repo.GetByRoomID(ctx, roomID, false)
		require.NoError(t, err)
		assert.Len(t, allocations, 3)

		// Should be ordered by created_at DESC (newest first)
		assert.Equal(t, third.ID, allocations[0].ID)
		assert.Equal(t, second.ID, allocations[1].ID)
		assert.Equal(t, first.ID, allocations[2].ID)
	})

	t.Run("should retrieve both shared and dedicated allocations for a room", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()

		// Create shared allocation
		sharedAllocation := ta.NewTestSharedAllocation(t, roomID, user1ID)
		ta.InsertAllocation(t, ctx, sharedAllocation, testPool)

		// Create dedicated allocation
		startDate := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 1, 0).Truncate(24 * time.Hour)
		dedicatedAllocation := ta.NewTestDedicatedAllocation(t, roomID, user2ID, startDate, endDate)
		ta.InsertAllocation(t, ctx, dedicatedAllocation, testPool)

		allocations, err := repo.GetByRoomID(ctx, roomID, false)
		require.NoError(t, err)
		assert.Len(t, allocations, 2)
	})

	t.Run("should not retrieve allocations for other rooms", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		room1ID := setupTestRoom(t)
		room2ID := setupTestRoom(t)
		userID := uuid.New()

		// Create allocation for room1
		allocation1 := ta.NewTestSharedAllocation(t, room1ID, userID)
		ta.InsertAllocation(t, ctx, allocation1, testPool)

		// Create allocation for room2
		allocation2 := ta.NewTestSharedAllocation(t, room2ID, userID)
		ta.InsertAllocation(t, ctx, allocation2, testPool)

		// Retrieve only room1's allocations
		allocations, err := repo.GetByRoomID(ctx, room1ID, false)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
		assert.Equal(t, room1ID, allocations[0].RoomID)
	})

	t.Run("should retrieve multiple users allocated to the same room", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()
		user3ID := uuid.New()

		// Create allocations for multiple users in the same room
		allocation1 := ta.NewTestSharedAllocation(t, roomID, user1ID)
		ta.InsertAllocation(t, ctx, allocation1, testPool)

		allocation2 := ta.NewTestSharedAllocation(t, roomID, user2ID)
		ta.InsertAllocation(t, ctx, allocation2, testPool)

		allocation3 := ta.NewTestSharedAllocation(t, roomID, user3ID)
		ta.InsertAllocation(t, ctx, allocation3, testPool)

		allocations, err := repo.GetByRoomID(ctx, roomID, false)
		require.NoError(t, err)
		assert.Len(t, allocations, 3)
	})
}
