package allocationRepository

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCheckConflict TEST_PATH=internal/booking/infrastructure/postgres/allocation/check_conflict_test.go

func TestCheckConflict(t *testing.T) {
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

	t.Run("dedicated: should detect conflict with overlapping dedicated allocation", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		newUserIDHash := ComputeUserIDHash(t, ctx, testCrypto, uuid.New())

		// Create existing dedicated allocation
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
		existingAllocationEncx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, userID, existingStart, existingEnd)
		err := repo.Create(ctx, existingAllocationEncx)
		require.NoError(t, err)

		// Check for conflict with overlapping period
		newStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, newUserIDHash, domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("dedicated: should not detect conflict with non-overlapping dedicated allocation", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		newUserIDHash := ComputeUserIDHash(t, ctx, testCrypto, uuid.New())

		// Create existing dedicated allocation
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		existingAllocationEncx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, userID, existingStart, existingEnd)
		err := repo.Create(ctx, existingAllocationEncx)
		require.NoError(t, err)

		// Check for conflict with non-overlapping period (starts after existing ends)
		newStart := time.Now().AddDate(0, 0, 11).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, newUserIDHash, domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("dedicated: should detect conflict with NULL end_date allocation", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		newUserIDHash := ComputeUserIDHash(t, ctx, testCrypto, uuid.New())

		// Create existing dedicated allocation with NULL end_date (indefinite)
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingAllocationEncx := NewTestDedicatedAllocationEncxWithNilEndDate(t, testCrypto, roomID, userID, existingStart)
		err := repo.Create(ctx, existingAllocationEncx)
		require.NoError(t, err)

		// Try to create new allocation after start_date (should conflict)
		newStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, newUserIDHash, domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("dedicated: should not detect conflict with inactive allocation", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		newUserIDHash := ComputeUserIDHash(t, ctx, testCrypto, uuid.New())

		// Create inactive dedicated allocation
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
		existingAllocationEncx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, userID, existingStart, existingEnd)
		// Deactivate by setting IsActive to false
		existingAllocationEncx.IsActive = false
		err := repo.Create(ctx, existingAllocationEncx)
		require.NoError(t, err)

		// Check for conflict with overlapping period (should not conflict since existing is inactive)
		newStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, newUserIDHash, domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("dedicated: should not detect conflict in different room", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		room1ID := setupTestRoom(t)
		room2ID := setupTestRoom(t)
		userID := uuid.New()
		newUserIDHash := ComputeUserIDHash(t, ctx, testCrypto, uuid.New())

		// Create allocation in room1
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
		existingAllocationEncx := NewTestDedicatedAllocationEncx(t, testCrypto, room1ID, userID, existingStart, existingEnd)
		err := repo.Create(ctx, existingAllocationEncx)
		require.NoError(t, err)

		// Check for conflict in room2 (different room)
		newStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, room2ID, newUserIDHash, domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("shared: should detect conflict when partner already has allocation for room", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		userIDHash := ComputeUserIDHash(t, ctx, testCrypto, userID)

		// Create existing shared allocation
		existingAllocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, userID)
		err := repo.Create(ctx, existingAllocationEncx)
		require.NoError(t, err)

		// Check for conflict with same partner and room
		hasConflict, err := repo.CheckConflict(ctx, roomID, userIDHash, domain.AllocationTypeShared, nil, nil, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("shared: should not detect conflict when different partner", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2IDHash := ComputeUserIDHash(t, ctx, testCrypto, uuid.New())

		// Create allocation for user1
		existingAllocationEncx := NewTestSharedAllocationEncx(t, testCrypto, roomID, user1ID)
		err := repo.Create(ctx, existingAllocationEncx)
		require.NoError(t, err)

		// Check for conflict with user2 (different partner)
		hasConflict, err := repo.CheckConflict(ctx, roomID, user2IDHash, domain.AllocationTypeShared, nil, nil, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("shared: should not detect conflict in different room", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		room1ID := setupTestRoom(t)
		room2ID := setupTestRoom(t)
		userID := uuid.New()
		userIDHash := ComputeUserIDHash(t, ctx, testCrypto, userID)

		// Create allocation in room1
		existingAllocationEncx := NewTestSharedAllocationEncx(t, testCrypto, room1ID, userID)
		err := repo.Create(ctx, existingAllocationEncx)
		require.NoError(t, err)

		// Check for conflict in room2
		hasConflict, err := repo.CheckConflict(ctx, room2ID, userIDHash, domain.AllocationTypeShared, nil, nil, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("shared: should not detect conflict with inactive allocation", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		userIDHash := ComputeUserIDHash(t, ctx, testCrypto, userID)

		// Create inactive allocation
		existingAllocationEncx := NewTestInactiveSharedAllocationEncx(t, testCrypto, roomID, userID)
		err := repo.Create(ctx, existingAllocationEncx)
		require.NoError(t, err)

		// Check for conflict (should not conflict since existing is inactive)
		hasConflict, err := repo.CheckConflict(ctx, roomID, userIDHash, domain.AllocationTypeShared, nil, nil, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("should return no conflict when room has no allocations", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		newUserIDHash := ComputeUserIDHash(t, ctx, testCrypto, uuid.New())

		newStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, newUserIDHash, domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("dedicated: edge case - new allocation starts exactly when existing ends", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		newUserIDHash := ComputeUserIDHash(t, ctx, testCrypto, uuid.New())

		// Create existing allocation ending at specific time
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		existingAllocationEncx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, userID, existingStart, existingEnd)
		err := repo.Create(ctx, existingAllocationEncx)
		require.NoError(t, err)

		// New allocation starts exactly when existing ends (no overlap)
		newStart := existingEnd
		newEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, newUserIDHash, domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("dedicated: edge case - new allocation ends exactly when existing starts", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		newUserIDHash := ComputeUserIDHash(t, ctx, testCrypto, uuid.New())

		// Create existing allocation starting at specific time
		existingStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)
		existingAllocationEncx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, userID, existingStart, existingEnd)
		err := repo.Create(ctx, existingAllocationEncx)
		require.NoError(t, err)

		// New allocation ends exactly when existing starts (no overlap)
		newStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		newEnd := existingStart

		hasConflict, err := repo.CheckConflict(ctx, roomID, newUserIDHash, domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("dedicated: should detect conflict when new allocation completely contains existing", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		newUserIDHash := ComputeUserIDHash(t, ctx, testCrypto, uuid.New())

		// Create existing allocation
		existingStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
		existingAllocationEncx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, userID, existingStart, existingEnd)
		err := repo.Create(ctx, existingAllocationEncx)
		require.NoError(t, err)

		// New allocation completely contains existing (starts before, ends after)
		newStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, newUserIDHash, domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("dedicated: should detect conflict when existing allocation completely contains new", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		newUserIDHash := ComputeUserIDHash(t, ctx, testCrypto, uuid.New())

		// Create existing allocation with wide range
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)
		existingAllocationEncx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, userID, existingStart, existingEnd)
		err := repo.Create(ctx, existingAllocationEncx)
		require.NoError(t, err)

		// New allocation completely contained within existing
		newStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, newUserIDHash, domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("dedicated: should detect conflict when new starts before and ends during existing", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		newUserIDHash := ComputeUserIDHash(t, ctx, testCrypto, uuid.New())

		// Create existing allocation
		existingStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)
		existingAllocationEncx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, userID, existingStart, existingEnd)
		err := repo.Create(ctx, existingAllocationEncx)
		require.NoError(t, err)

		// New allocation starts before existing but ends during
		newStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, newUserIDHash, domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("dedicated: should detect conflict when new starts during and ends after existing", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		newUserIDHash := ComputeUserIDHash(t, ctx, testCrypto, uuid.New())

		// Create existing allocation
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
		existingAllocationEncx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, userID, existingStart, existingEnd)
		err := repo.Create(ctx, existingAllocationEncx)
		require.NoError(t, err)

		// New allocation starts during existing but ends after
		newStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, newUserIDHash, domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("dedicated: should detect conflict with any of multiple overlapping allocations", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()
		newUserIDHash := ComputeUserIDHash(t, ctx, testCrypto, uuid.New())

		// Create first allocation
		allocation1Start := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		allocation1End := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		allocation1Encx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, user1ID, allocation1Start, allocation1End)
		err := repo.Create(ctx, allocation1Encx)
		require.NoError(t, err)

		// Create second allocation
		allocation2Start := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
		allocation2End := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)
		allocation2Encx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, user2ID, allocation2Start, allocation2End)
		err = repo.Create(ctx, allocation2Encx)
		require.NoError(t, err)

		// Check for conflict with period overlapping second allocation
		newStart := time.Now().AddDate(0, 0, 17).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 25).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, newUserIDHash, domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("dedicated: should not detect conflict when new allocation is before all existing allocations", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		newUserIDHash := ComputeUserIDHash(t, ctx, testCrypto, uuid.New())

		// Create existing allocation starting in the future
		existingStart := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)
		existingAllocationEncx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, userID, existingStart, existingEnd)
		err := repo.Create(ctx, existingAllocationEncx)
		require.NoError(t, err)

		// New allocation completely before existing
		newStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, newUserIDHash, domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("dedicated: new allocation with NULL end_date should conflict with existing", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		newUserIDHash := ComputeUserIDHash(t, ctx, testCrypto, uuid.New())

		// Create existing allocation
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
		existingAllocationEncx := NewTestDedicatedAllocationEncx(t, testCrypto, roomID, userID, existingStart, existingEnd)
		err := repo.Create(ctx, existingAllocationEncx)
		require.NoError(t, err)

		// Try to create new indefinite allocation (NULL end_date) starting during existing
		newStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, newUserIDHash, domain.AllocationTypeDedicated, &newStart, nil, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("dedicated: two indefinite allocations (NULL end_date) should conflict", func(t *testing.T) {
		ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()
		newUserIDHash := ComputeUserIDHash(t, ctx, testCrypto, uuid.New())

		// Create existing indefinite allocation
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingAllocationEncx := NewTestDedicatedAllocationEncxWithNilEndDate(t, testCrypto, roomID, userID, existingStart)
		err := repo.Create(ctx, existingAllocationEncx)
		require.NoError(t, err)

		// Try to create new indefinite allocation
		newStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, newUserIDHash, domain.AllocationTypeDedicated, &newStart, nil, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})
}
