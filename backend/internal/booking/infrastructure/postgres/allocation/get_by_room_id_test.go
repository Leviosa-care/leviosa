package allocationRepository

import (
	"context"
	"testing"
	"time"

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
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()

		// Create active and inactive allocations for the same room
		activeAllocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, user1ID)
		err := repo.Create(ctx, activeAllocationEncx)
		require.NoError(t, err)

		inactiveAllocationEncx := NewTestInactiveSharedAllocationEncx(t, testCrypto, roomID, user2ID)
		err = repo.Create(ctx, inactiveAllocationEncx)
		require.NoError(t, err)

		// Retrieve all
		allocations, err := repo.GetByRoomID(ctx, roomID, false)
		require.NoError(t, err)
		assert.Len(t, allocations, 2)
	})

	t.Run("should retrieve only active allocations for a room (activeOnly=true)", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()

		// Create active and inactive allocations
		activeAllocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, user1ID)
		err := repo.Create(ctx, activeAllocationEncx)
		require.NoError(t, err)

		inactiveAllocationEncx := NewTestInactiveSharedAllocationEncx(t, testCrypto, roomID, user2ID)
		err = repo.Create(ctx, inactiveAllocationEncx)
		require.NoError(t, err)

		// Retrieve only active
		allocations, err := repo.GetByRoomID(ctx, roomID, true)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
		assert.True(t, allocations[0].IsActive)
		assert.Equal(t, activeAllocationEncx.ID, allocations[0].ID)
	})

	t.Run("should return empty slice when room has no allocations", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		setupTestRoom(t) // Create room but no allocations
		nonExistentRoomID := uuid.New()

		allocations, err := repo.GetByRoomID(ctx, nonExistentRoomID, false)
		require.NoError(t, err)
		assert.Empty(t, allocations)
	})

	t.Run("should order allocations by created_at DESC", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		// Create allocations with staggered created_at
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

		allocations, err := repo.GetByRoomID(ctx, roomID, false)
		require.NoError(t, err)
		assert.Len(t, allocations, 3)

		// Should be ordered by created_at DESC (newest first)
		assert.Equal(t, third.ID, allocations[0].ID)
		assert.Equal(t, second.ID, allocations[1].ID)
		assert.Equal(t, first.ID, allocations[2].ID)
	})

	t.Run("should retrieve both shared and dedicated allocations for a room", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()

		// Create shared allocation
		sharedAllocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, user1ID)
		err := repo.Create(ctx, sharedAllocationEncx)
		require.NoError(t, err)

		// Create dedicated allocation
		startDate := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 1, 0).Truncate(24 * time.Hour)
		dedicatedAllocationEncx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, user2ID, startDate, endDate)
		err = repo.Create(ctx, dedicatedAllocationEncx)
		require.NoError(t, err)

		allocations, err := repo.GetByRoomID(ctx, roomID, false)
		require.NoError(t, err)
		assert.Len(t, allocations, 2)
	})

	t.Run("should not retrieve allocations for other rooms", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		room1ID := setupTestRoom(t)
		room2ID := setupTestRoom(t)
		userID := uuid.New()

		// Create allocation for room1
		allocation1Encx := NewTestSharedAllocationEncx(t, testCrypto, room1ID, userID)
		err := repo.Create(ctx, allocation1Encx)
		require.NoError(t, err)

		// Create allocation for room2
		allocation2Encx := NewTestSharedAllocationEncx(t, testCrypto, room2ID, userID)
		err = repo.Create(ctx, allocation2Encx)
		require.NoError(t, err)

		// Retrieve only room1's allocations
		allocations, err := repo.GetByRoomID(ctx, room1ID, false)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
		assert.Equal(t, room1ID, allocations[0].RoomID)
	})

	t.Run("should retrieve multiple users allocated to the same room", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()
		user3ID := uuid.New()

		// Create allocations for multiple users in the same room
		allocation1Encx := NewTestSharedAllocationEncx(t, testCrypto, roomID, user1ID)
		err := repo.Create(ctx, allocation1Encx)
		require.NoError(t, err)

		allocation2Encx := NewTestSharedAllocationEncx(t, testCrypto, roomID, user2ID)
		err = repo.Create(ctx, allocation2Encx)
		require.NoError(t, err)

		allocation3Encx := NewTestSharedAllocationEncx(t, testCrypto, roomID, user3ID)
		err = repo.Create(ctx, allocation3Encx)
		require.NoError(t, err)

		allocations, err := repo.GetByRoomID(ctx, roomID, false)
		require.NoError(t, err)
		assert.Len(t, allocations, 3)
	})
}
