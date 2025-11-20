package allocationRepository

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	ta "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/allocation"
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
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create existing dedicated allocation
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
		existingAllocation := ta.NewTestDedicatedAllocation(t, roomID, userID, existingStart, existingEnd)
		ta.InsertAllocation(t, ctx, existingAllocation, testPool)

		// Check for conflict with overlapping period
		newStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, uuid.New(), domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("dedicated: should not detect conflict with non-overlapping dedicated allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create existing dedicated allocation
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		existingAllocation := ta.NewTestDedicatedAllocation(t, roomID, userID, existingStart, existingEnd)
		ta.InsertAllocation(t, ctx, existingAllocation, testPool)

		// Check for conflict with non-overlapping period (starts after existing ends)
		newStart := time.Now().AddDate(0, 0, 11).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, uuid.New(), domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("dedicated: should detect conflict with NULL end_date allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create existing dedicated allocation with NULL end_date (indefinite)
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingAllocation := ta.NewTestDedicatedAllocationIndefinite(t, roomID, userID, existingStart)
		ta.InsertAllocation(t, ctx, existingAllocation, testPool)

		// Try to create new allocation after start_date (should conflict)
		newStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, uuid.New(), domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("dedicated: should not detect conflict with inactive allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create inactive dedicated allocation
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
		existingAllocation := ta.NewTestDedicatedAllocation(t, roomID, userID, existingStart, existingEnd)
		existingAllocation.Deactivate()
		ta.InsertAllocation(t, ctx, existingAllocation, testPool)

		// Check for conflict with overlapping period (should not conflict since existing is inactive)
		newStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, uuid.New(), domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("dedicated: should not detect conflict in different room", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		room1ID := setupTestRoom(t)
		room2ID := setupTestRoom(t)
		userID := uuid.New()

		// Create allocation in room1
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
		existingAllocation := ta.NewTestDedicatedAllocation(t, room1ID, userID, existingStart, existingEnd)
		ta.InsertAllocation(t, ctx, existingAllocation, testPool)

		// Check for conflict in room2 (different room)
		newStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, room2ID, uuid.New(), domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("shared: should detect conflict when partner already has allocation for room", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create existing shared allocation
		existingAllocation := ta.NewTestSharedAllocation(t, roomID, userID)
		ta.InsertAllocation(t, ctx, existingAllocation, testPool)

		// Check for conflict with same partner and room
		hasConflict, err := repo.CheckConflict(ctx, roomID, userID, domain.AllocationTypeShared, nil, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("shared: should not detect conflict when different partner", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()

		// Create allocation for user1
		existingAllocation := ta.NewTestSharedAllocation(t, roomID, user1ID)
		ta.InsertAllocation(t, ctx, existingAllocation, testPool)

		// Check for conflict with user2 (different partner)
		hasConflict, err := repo.CheckConflict(ctx, roomID, user2ID, domain.AllocationTypeShared, nil, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("shared: should not detect conflict in different room", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		room1ID := setupTestRoom(t)
		room2ID := setupTestRoom(t)
		userID := uuid.New()

		// Create allocation in room1
		existingAllocation := ta.NewTestSharedAllocation(t, room1ID, userID)
		ta.InsertAllocation(t, ctx, existingAllocation, testPool)

		// Check for conflict in room2
		hasConflict, err := repo.CheckConflict(ctx, room2ID, userID, domain.AllocationTypeShared, nil, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("shared: should not detect conflict with inactive allocation", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create inactive allocation
		existingAllocation := ta.NewTestInactiveAllocation(t, roomID, userID)
		ta.InsertAllocation(t, ctx, existingAllocation, testPool)

		// Check for conflict (should not conflict since existing is inactive)
		hasConflict, err := repo.CheckConflict(ctx, roomID, userID, domain.AllocationTypeShared, nil, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("should return no conflict when room has no allocations", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)

		newStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, uuid.New(), domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("dedicated: edge case - new allocation starts exactly when existing ends", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create existing allocation ending at specific time
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		existingAllocation := ta.NewTestDedicatedAllocation(t, roomID, userID, existingStart, existingEnd)
		ta.InsertAllocation(t, ctx, existingAllocation, testPool)

		// New allocation starts exactly when existing ends (no overlap)
		newStart := existingEnd
		newEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, uuid.New(), domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("dedicated: edge case - new allocation ends exactly when existing starts", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create existing allocation starting at specific time
		existingStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)
		existingAllocation := ta.NewTestDedicatedAllocation(t, roomID, userID, existingStart, existingEnd)
		ta.InsertAllocation(t, ctx, existingAllocation, testPool)

		// New allocation ends exactly when existing starts (no overlap)
		newStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		newEnd := existingStart

		hasConflict, err := repo.CheckConflict(ctx, roomID, uuid.New(), domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("dedicated: should detect conflict when new allocation completely contains existing", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create existing allocation
		existingStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
		existingAllocation := ta.NewTestDedicatedAllocation(t, roomID, userID, existingStart, existingEnd)
		ta.InsertAllocation(t, ctx, existingAllocation, testPool)

		// New allocation completely contains existing (starts before, ends after)
		newStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, uuid.New(), domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("dedicated: should detect conflict when existing allocation completely contains new", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create existing allocation with wide range
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)
		existingAllocation := ta.NewTestDedicatedAllocation(t, roomID, userID, existingStart, existingEnd)
		ta.InsertAllocation(t, ctx, existingAllocation, testPool)

		// New allocation completely contained within existing
		newStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, uuid.New(), domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("dedicated: should detect conflict when new starts before and ends during existing", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create existing allocation
		existingStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)
		existingAllocation := ta.NewTestDedicatedAllocation(t, roomID, userID, existingStart, existingEnd)
		ta.InsertAllocation(t, ctx, existingAllocation, testPool)

		// New allocation starts before existing but ends during
		newStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, uuid.New(), domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("dedicated: should detect conflict when new starts during and ends after existing", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create existing allocation
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
		existingAllocation := ta.NewTestDedicatedAllocation(t, roomID, userID, existingStart, existingEnd)
		ta.InsertAllocation(t, ctx, existingAllocation, testPool)

		// New allocation starts during existing but ends after
		newStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, uuid.New(), domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("dedicated: should detect conflict with any of multiple overlapping allocations", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		user1ID := uuid.New()
		user2ID := uuid.New()

		// Create first allocation
		allocation1Start := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		allocation1End := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		allocation1 := ta.NewTestDedicatedAllocation(t, roomID, user1ID, allocation1Start, allocation1End)
		ta.InsertAllocation(t, ctx, allocation1, testPool)

		// Create second allocation
		allocation2Start := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
		allocation2End := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)
		allocation2 := ta.NewTestDedicatedAllocation(t, roomID, user2ID, allocation2Start, allocation2End)
		ta.InsertAllocation(t, ctx, allocation2, testPool)

		// Check for conflict with period overlapping second allocation
		newStart := time.Now().AddDate(0, 0, 17).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 25).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, uuid.New(), domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("dedicated: should not detect conflict when new allocation is before all existing allocations", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create existing allocation starting in the future
		existingStart := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 20).Truncate(24 * time.Hour)
		existingAllocation := ta.NewTestDedicatedAllocation(t, roomID, userID, existingStart, existingEnd)
		ta.InsertAllocation(t, ctx, existingAllocation, testPool)

		// New allocation completely before existing
		newStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		newEnd := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, uuid.New(), domain.AllocationTypeDedicated, &newStart, &newEnd, nil)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})

	t.Run("dedicated: new allocation with NULL end_date should conflict with existing", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create existing allocation
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingEnd := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
		existingAllocation := ta.NewTestDedicatedAllocation(t, roomID, userID, existingStart, existingEnd)
		ta.InsertAllocation(t, ctx, existingAllocation, testPool)

		// Try to create new indefinite allocation (NULL end_date) starting during existing
		newStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, uuid.New(), domain.AllocationTypeDedicated, &newStart, nil, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})

	t.Run("dedicated: two indefinite allocations (NULL end_date) should conflict", func(t *testing.T) {
		ta.ClearAllocationTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		roomID := setupTestRoom(t)
		userID := uuid.New()

		// Create existing indefinite allocation
		existingStart := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
		existingAllocation := ta.NewTestDedicatedAllocationIndefinite(t, roomID, userID, existingStart)
		ta.InsertAllocation(t, ctx, existingAllocation, testPool)

		// Try to create new indefinite allocation
		newStart := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)

		hasConflict, err := repo.CheckConflict(ctx, roomID, uuid.New(), domain.AllocationTypeDedicated, &newStart, nil, nil)
		require.NoError(t, err)
		assert.True(t, hasConflict)
	})
}
