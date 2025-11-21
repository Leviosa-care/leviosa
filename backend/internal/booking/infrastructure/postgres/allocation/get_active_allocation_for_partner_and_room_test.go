package allocationRepository

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
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
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		userIDHash := ComputeUserIDHash(t, ctx, testCrypto, userID)

		allocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, userID)
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		// Check at any time (shared is always active)
		at := time.Now()
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, userIDHash, roomID, at)
		require.NoError(t, err)
		assert.Equal(t, allocationEncx.ID, retrieved.ID)
		assert.Equal(t, domain.AllocationTypeShared, retrieved.AllocationType)
	})

	t.Run("should retrieve dedicated allocation within active period", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		userIDHash := ComputeUserIDHash(t, ctx, testCrypto, userID)

		startDate := time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
		endDate := time.Now().AddDate(0, 0, 7).Truncate(24 * time.Hour)

		allocationEncx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, userID, startDate, endDate)
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		// Check at current time (within period)
		at := time.Now()
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, userIDHash, roomID, at)
		require.NoError(t, err)
		assert.Equal(t, allocationEncx.ID, retrieved.ID)
		assert.Equal(t, domain.AllocationTypeDedicated, retrieved.AllocationType)
	})

	t.Run("should not retrieve dedicated allocation before start_date", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		userIDHash := ComputeUserIDHash(t, ctx, testCrypto, userID)

		allocationEncx := NewTestFutureDedicatedAllocationEncx(t, testCrypto, roomID, userID)
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		// Check at current time (before start_date)
		at := time.Now()
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, userIDHash, roomID, at)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
		assert.Nil(t, retrieved)
	})

	t.Run("should not retrieve dedicated allocation after end_date", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		userIDHash := ComputeUserIDHash(t, ctx, testCrypto, userID)

		allocationEncx := NewTestPastDedicatedAllocationEncx(t, testCrypto, roomID, userID)
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		// Check at current time (after end_date)
		at := time.Now()
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, userIDHash, roomID, at)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
		assert.Nil(t, retrieved)
	})

	t.Run("should retrieve dedicated allocation with NULL end_date (indefinite)", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		userIDHash := ComputeUserIDHash(t, ctx, testCrypto, userID)

		startDate := time.Now().AddDate(0, 0, -30).Truncate(24 * time.Hour)

		allocationEncx := NewTestDedicatedAllocationEncxWithNilEndDate(t, testCrypto, roomID, userID, startDate)
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		// Check at current time and future time
		at := time.Now().AddDate(0, 1, 0) // 1 month in future
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, userIDHash, roomID, at)
		require.NoError(t, err)
		assert.Equal(t, allocationEncx.ID, retrieved.ID)
		assert.Nil(t, retrieved.EndDate)
	})

	t.Run("should prioritize dedicated over shared when both exist", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		userIDHash := ComputeUserIDHash(t, ctx, testCrypto, userID)

		// Create shared allocation
		sharedAllocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, userID)
		err := repo.Create(ctx, sharedAllocationEncx)
		require.NoError(t, err)

		// Create dedicated allocation (currently active)
		dedicatedAllocationEncx := NewTestActiveDedicatedAllocationEncx(t, testCrypto, roomID, userID)
		err = repo.Create(ctx, dedicatedAllocationEncx)
		require.NoError(t, err)

		count := CountActiveAllocations(t, ctx, testPool)
		require.Equal(t, 2, count, "Should have 2 active allocations in the database set up.")

		// Should return dedicated (higher priority)
		at := time.Now()
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, userIDHash, roomID, at)
		require.NoError(t, err)
		assert.Equal(t, dedicatedAllocationEncx.ID, retrieved.ID)
		assert.Equal(t, domain.AllocationTypeDedicated, retrieved.AllocationType)
	})

	t.Run("should not retrieve inactive allocation", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		userIDHash := ComputeUserIDHash(t, ctx, testCrypto, userID)

		allocationEncx := NewTestInactiveSharedAllocationEncx(t, testCrypto, roomID, userID)
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		at := time.Now()
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, userIDHash, roomID, at)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
		assert.Nil(t, retrieved)
	})

	t.Run("should return not found when no allocation exists", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		nonExistentUserID := uuid.New()
		nonExistentUserIDHash := ComputeUserIDHash(t, ctx, testCrypto, nonExistentUserID)
		nonExistentRoomID := uuid.New()

		at := time.Now()
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, nonExistentUserIDHash, nonExistentRoomID, at)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
		assert.Nil(t, retrieved)
	})

	t.Run("should not retrieve allocation for different user", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()
		user2IDHash := ComputeUserIDHash(t, ctx, testCrypto, user2ID)

		allocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, user1ID)
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		at := time.Now()
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, user2IDHash, roomID, at)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
		assert.Nil(t, retrieved)
	})

	t.Run("should not retrieve allocation for different room", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		room1ID := setupTestRoom(t)
		room2ID := setupTestRoom(t)
		userID := uuid.New()
		userIDHash := ComputeUserIDHash(t, ctx, testCrypto, userID)

		allocationEncx := NewTestSharedAllocationEncx(t, testCrypto, room1ID, userID)
		err := repo.Create(ctx, allocationEncx)
		require.NoError(t, err)

		at := time.Now()
		retrieved, err := repo.GetActiveAllocationForPartnerAndRoom(ctx, userIDHash, room2ID, at)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
		assert.Nil(t, retrieved)
	})
}
