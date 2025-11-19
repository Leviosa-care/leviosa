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

// make test-func TEST_NAME=TestGetByUserID TEST_PATH=internal/booking/infrastructure/postgres/allocation/get_by_user_id_test.go

func TestGetByUserID(t *testing.T) {
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

	t.Run("should retrieve all allocations for a user (activeOnly=false)", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID1 := setupTestRoom(t)
		roomID2 := setupTestRoom(t)
		userID := uuid.New()

		// Create active and inactive allocations
		activeAllocation := ta.NewTestSharedAllocation(t, roomID1, userID)
		ta.InsertAllocation(t, ctx, activeAllocation, testPool)

		inactiveAllocation := ta.NewTestInactiveAllocation(t, roomID2, userID)
		ta.InsertAllocation(t, ctx, inactiveAllocation, testPool)

		// Retrieve all
		allocations, err := repo.GetByUserID(ctx, userID, false)
		require.NoError(t, err)
		assert.Len(t, allocations, 2)
	})

	t.Run("should retrieve only active allocations for a user (activeOnly=true)", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID1 := setupTestRoom(t)
		roomID2 := setupTestRoom(t)
		userID := uuid.New()

		// Create active and inactive allocations
		activeAllocation := ta.NewTestSharedAllocation(t, roomID1, userID)
		ta.InsertAllocation(t, ctx, activeAllocation, testPool)

		inactiveAllocation := ta.NewTestInactiveAllocation(t, roomID2, userID)
		ta.InsertAllocation(t, ctx, inactiveAllocation, testPool)

		// Retrieve only active
		allocations, err := repo.GetByUserID(ctx, userID, true)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
		assert.True(t, allocations[0].IsActive)
		assert.Equal(t, activeAllocation.ID, allocations[0].ID)
	})

	t.Run("should return empty slice when user has no allocations", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		nonExistentUserID := uuid.New()

		allocations, err := repo.GetByUserID(ctx, nonExistentUserID, false)
		require.NoError(t, err)
		assert.Empty(t, allocations)
	})

	t.Run("should order allocations by created_at DESC", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID1 := setupTestRoom(t)
		roomID2 := setupTestRoom(t)
		roomID3 := setupTestRoom(t)
		userID := uuid.New()

		// Create allocations with staggered created_at
		first := ta.NewTestSharedAllocation(t, roomID1, userID)
		first.CreatedAt = time.Now().Add(-2 * time.Hour)
		ta.InsertAllocation(t, ctx, first, testPool)

		second := ta.NewTestSharedAllocation(t, roomID2, userID)
		second.CreatedAt = time.Now().Add(-1 * time.Hour)
		ta.InsertAllocation(t, ctx, second, testPool)

		third := ta.NewTestSharedAllocation(t, roomID3, userID)
		third.CreatedAt = time.Now()
		ta.InsertAllocation(t, ctx, third, testPool)

		allocations, err := repo.GetByUserID(ctx, userID, false)
		require.NoError(t, err)
		assert.Len(t, allocations, 3)

		// Should be ordered by created_at DESC (newest first)
		assert.Equal(t, third.ID, allocations[0].ID)
		assert.Equal(t, second.ID, allocations[1].ID)
		assert.Equal(t, first.ID, allocations[2].ID)
	})

	t.Run("should retrieve both shared and dedicated allocations", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create shared allocation
		sharedAllocation := ta.NewTestSharedAllocation(t, roomID, userID)
		ta.InsertAllocation(t, ctx, sharedAllocation, testPool)

		// Create dedicated allocation
		startDate := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 1, 0).Truncate(24 * time.Hour)
		dedicatedAllocation := ta.NewTestDedicatedAllocation(t, roomID, userID, startDate, endDate)
		ta.InsertAllocation(t, ctx, dedicatedAllocation, testPool)

		allocations, err := repo.GetByUserID(ctx, userID, false)
		require.NoError(t, err)
		assert.Len(t, allocations, 2)
	})

	t.Run("should not retrieve allocations for other users", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()

		// Create allocation for user1
		allocation1 := ta.NewTestSharedAllocation(t, roomID, user1ID)
		ta.InsertAllocation(t, ctx, allocation1, testPool)

		// Create allocation for user2
		allocation2 := ta.NewTestSharedAllocation(t, roomID, user2ID)
		ta.InsertAllocation(t, ctx, allocation2, testPool)

		// Retrieve only user1's allocations
		allocations, err := repo.GetByUserID(ctx, user1ID, false)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
		assert.Equal(t, user1ID, allocations[0].UserID)
	})
}
