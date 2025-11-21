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

// make test-func TEST_NAME=TestGetByUserIDHash TEST_PATH=internal/booking/infrastructure/postgres/allocation/get_by_user_id_test.go

func TestGetByUserIDHash(t *testing.T) {
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
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID1 := setupTestRoom(t)
		roomID2 := setupTestRoom(t)
		userID := uuid.New()
		userIDHash := ComputeUserIDHash(t, ctx, testCrypto, userID)

		// Create active and inactive allocations
		activeAllocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID1, userID)
		err := repo.Create(ctx, activeAllocationEncx)
		require.NoError(t, err)

		inactiveAllocationEncx := NewTestInactiveSharedAllocationEncx(t, testCrypto, roomID2, userID)
		err = repo.Create(ctx, inactiveAllocationEncx)
		require.NoError(t, err)

		// Retrieve all
		allocations, err := repo.GetByUserIDHash(ctx, userIDHash, false)
		require.NoError(t, err)
		assert.Len(t, allocations, 2)
	})

	t.Run("should retrieve only active allocations for a user (activeOnly=true)", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID1 := setupTestRoom(t)
		roomID2 := setupTestRoom(t)
		userID := uuid.New()
		userIDHash := ComputeUserIDHash(t, ctx, testCrypto, userID)

		// Create active and inactive allocations
		activeAllocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID1, userID)
		err := repo.Create(ctx, activeAllocationEncx)
		require.NoError(t, err)

		inactiveAllocationEncx := NewTestInactiveSharedAllocationEncx(t, testCrypto, roomID2, userID)
		err = repo.Create(ctx, inactiveAllocationEncx)
		require.NoError(t, err)

		// Retrieve only active
		allocations, err := repo.GetByUserIDHash(ctx, userIDHash, true)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
		assert.True(t, allocations[0].IsActive)
		assert.Equal(t, activeAllocationEncx.ID, allocations[0].ID)
	})

	t.Run("should return empty slice when user has no allocations", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		nonExistentUserID := uuid.New()
		nonExistentUserIDHash := ComputeUserIDHash(t, ctx, testCrypto, nonExistentUserID)

		allocations, err := repo.GetByUserIDHash(ctx, nonExistentUserIDHash, false)
		require.NoError(t, err)
		assert.Empty(t, allocations)
	})

	t.Run("should order allocations by created_at DESC", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID1 := setupTestRoom(t)
		roomID2 := setupTestRoom(t)
		roomID3 := setupTestRoom(t)
		userID := uuid.New()
		userIDHash := ComputeUserIDHash(t, ctx, testCrypto, userID)

		// Create allocations with staggered created_at
		first := NewTestSharedAllocationEncx(t, testCrypto, roomID1, userID)
		first.CreatedAt = time.Now().Add(-2 * time.Hour)
		err := repo.Create(ctx, first)
		require.NoError(t, err)

		second := NewTestSharedAllocationEncx(t, testCrypto, roomID2, userID)
		second.CreatedAt = time.Now().Add(-1 * time.Hour)
		err = repo.Create(ctx, second)
		require.NoError(t, err)

		third := NewTestSharedAllocationEncx(t, testCrypto, roomID3, userID)
		third.CreatedAt = time.Now()
		err = repo.Create(ctx, third)
		require.NoError(t, err)

		allocations, err := repo.GetByUserIDHash(ctx, userIDHash, false)
		require.NoError(t, err)
		assert.Len(t, allocations, 3)

		// Should be ordered by created_at DESC (newest first)
		assert.Equal(t, third.ID, allocations[0].ID)
		assert.Equal(t, second.ID, allocations[1].ID)
		assert.Equal(t, first.ID, allocations[2].ID)
	})

	t.Run("should retrieve both shared and dedicated allocations", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		userIDHash := ComputeUserIDHash(t, ctx, testCrypto, userID)

		// Create shared allocation
		sharedAllocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, userID)
		err := repo.Create(ctx, sharedAllocationEncx)
		require.NoError(t, err)

		// Create dedicated allocation
		startDate := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 1, 0).Truncate(24 * time.Hour)
		dedicatedAllocationEncx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, userID, startDate, endDate)
		err = repo.Create(ctx, dedicatedAllocationEncx)
		require.NoError(t, err)

		allocations, err := repo.GetByUserIDHash(ctx, userIDHash, false)
		require.NoError(t, err)
		assert.Len(t, allocations, 2)
	})

	t.Run("should not retrieve allocations for other users", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()
		user1IDHash := ComputeUserIDHash(t, ctx, testCrypto, user1ID)

		// Create allocation for user1
		allocation1Encx := NewTestSharedAllocationEncx(t, testCrypto, roomID, user1ID)
		err := repo.Create(ctx, allocation1Encx)
		require.NoError(t, err)

		// Create allocation for user2
		allocation2Encx := NewTestSharedAllocationEncx(t, testCrypto, roomID, user2ID)
		err = repo.Create(ctx, allocation2Encx)
		require.NoError(t, err)

		// Retrieve only user1's allocations
		allocations, err := repo.GetByUserIDHash(ctx, user1IDHash, false)
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
		assert.Equal(t, user1IDHash, allocations[0].UserIDHash)
	})
}
