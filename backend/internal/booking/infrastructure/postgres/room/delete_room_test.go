package roomRepository_test

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

// make test-func TEST_NAME=TestDeleteRoom TEST_PATH=internal/booking/infrastructure/postgres/room/delete_room_test.go

func TestDeleteRoom(t *testing.T) {
	ctx := context.Background()

	// Create test building data first
	buildingEncx := tb.NewTestBuildingEncx(t)
	err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
	require.NoError(t, err)

	// Create test room data
	roomEncx := tr.NewTestRoomEncx(t)
	roomEncx.BuildingID = buildingEncx.ID

	t.Run("should successfully soft delete room", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Insert test room directly
		err := tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		// Verify room is initially active
		savedRoom, err := tr.GetRoomEncxByID(t, ctx, testPool, roomEncx.ID)
		require.NoError(t, err)
		require.True(t, savedRoom.IsActive, "Room should be initially active")

		// Test repository Delete method
		err = repo.Delete(ctx, roomEncx.ID)
		require.NoError(t, err)

		// Verify room is now inactive (soft delete)
		deletedRoom, err := tr.GetRoomEncxByID(t, ctx, testPool, roomEncx.ID)
		require.NoError(t, err)
		require.False(t, deletedRoom.IsActive, "Room should be inactive after delete")

		// Verify other fields remain unchanged
		assert.Equal(t, roomEncx.ID, deletedRoom.ID)
		assert.Equal(t, roomEncx.BuildingID, deletedRoom.BuildingID)
		assert.Equal(t, roomEncx.NameHash, deletedRoom.NameHash)
		assert.Equal(t, roomEncx.RoomNumberHash, deletedRoom.RoomNumberHash)
		assert.Equal(t, roomEncx.Capacity, deletedRoom.Capacity)
	})

	t.Run("should return not found error for non-existent room", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Try to delete non-existent room
		nonExistentID := uuid.New()
		err := repo.Delete(ctx, nonExistentID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("should handle delete of already inactive room", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Create and insert inactive room
		inactiveRoomEncx := tr.NewInactiveTestRoomEncx(t, roomEncx.BuildingID)
		err := tr.InsertRoomEncx(t, ctx, testPool, inactiveRoomEncx)
		require.NoError(t, err)

		// Verify room is initially inactive
		savedRoom, err := tr.GetRoomEncxByID(t, ctx, testPool, inactiveRoomEncx.ID)
		require.NoError(t, err)
		require.False(t, savedRoom.IsActive, "Room should be initially inactive")

		// Test repository Delete method (should still succeed)
		err = repo.Delete(ctx, inactiveRoomEncx.ID)
		require.NoError(t, err)

		// Verify room is still inactive
		deletedRoom, err := tr.GetRoomEncxByID(t, ctx, testPool, inactiveRoomEncx.ID)
		require.NoError(t, err)
		require.False(t, deletedRoom.IsActive, "Room should remain inactive")
	})

	t.Run("should preserve all other fields during soft delete", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Create room with complete data
		completeRoomEncx := tr.NewTestRoomEncxWithBuilding(t, roomEncx.BuildingID)
		err := tr.InsertRoomEncx(t, ctx, testPool, completeRoomEncx)
		require.NoError(t, err)

		// Test repository Delete method
		err = repo.Delete(ctx, completeRoomEncx.ID)
		require.NoError(t, err)

		// Verify all fields except is_active are preserved
		deletedRoom, err := tr.GetRoomEncxByID(t, ctx, testPool, completeRoomEncx.ID)
		require.NoError(t, err)

		assert.Equal(t, completeRoomEncx.ID, deletedRoom.ID)
		assert.Equal(t, completeRoomEncx.BuildingID, deletedRoom.BuildingID)
		assert.Equal(t, completeRoomEncx.NameHash, deletedRoom.NameHash)
		assert.Equal(t, completeRoomEncx.RoomNumberHash, deletedRoom.RoomNumberHash)
		assert.Equal(t, completeRoomEncx.Capacity, deletedRoom.Capacity)
		assert.Equal(t, completeRoomEncx.EquipmentEncrypted, deletedRoom.EquipmentEncrypted)
		assert.Equal(t, completeRoomEncx.HourlyRateCents, deletedRoom.HourlyRateCents)
		assert.False(t, deletedRoom.IsActive, "Room should be inactive")
		assert.WithinDuration(t, completeRoomEncx.CreatedAt, deletedRoom.CreatedAt, time.Second)
		assert.WithinDuration(t, completeRoomEncx.UpdatedAt, deletedRoom.UpdatedAt, time.Second, "UpdatedAt should be updated")
	})

	t.Run("should handle room without hourly rate deletion", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Create room without hourly rate
		noRateRoomEncx := tr.NewInactiveTestRoomEncx(t, roomEncx.BuildingID)
		noRateRoomEncx.IsActive = true // Make it active for deletion test
		err := tr.InsertRoomEncx(t, ctx, testPool, noRateRoomEncx)
		require.NoError(t, err)

		// Test repository Delete method
		err = repo.Delete(ctx, noRateRoomEncx.ID)
		require.NoError(t, err)

		// Verify room is inactive but other fields preserved
		deletedRoom, err := tr.GetRoomEncxByID(t, ctx, testPool, noRateRoomEncx.ID)
		require.NoError(t, err)

		assert.False(t, deletedRoom.IsActive, "Room should be inactive")
		assert.Nil(t, deletedRoom.HourlyRateCents, "Hourly rate should remain null")
		assert.Equal(t, noRateRoomEncx.NameHash, deletedRoom.NameHash)
	})

	t.Run("should handle multiple deletions", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Create multiple rooms
		room1Encx := tr.NewTestRoomEncxWithBuilding(t, roomEncx.BuildingID)
		room2Encx := tr.NewTestRoomEncxWithBuilding(t, roomEncx.BuildingID)

		err := tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)

		// Delete both rooms
		err = repo.Delete(ctx, room1Encx.ID)
		require.NoError(t, err)
		err = repo.Delete(ctx, room2Encx.ID)
		require.NoError(t, err)

		// Verify both rooms are inactive
		deletedRoom1, err := tr.GetRoomEncxByID(t, ctx, testPool, room1Encx.ID)
		require.NoError(t, err)
		assert.False(t, deletedRoom1.IsActive)

		deletedRoom2, err := tr.GetRoomEncxByID(t, ctx, testPool, room2Encx.ID)
		require.NoError(t, err)
		assert.False(t, deletedRoom2.IsActive)
	})

	t.Run("should handle database errors gracefully", func(t *testing.T) {
		// Test with cancelled context
		ctx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		// Try to delete room with cancelled context
		err := repo.Delete(ctx, roomEncx.ID)
		require.Error(t, err)
	})

	t.Run("should verify soft delete doesn't remove record from database", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Insert test room
		err := tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		// Count rooms before deletion
		var countBefore int
		err = testPool.QueryRow(ctx, "SELECT COUNT(*) FROM booking.rooms WHERE id = $1", roomEncx.ID).Scan(&countBefore)
		require.NoError(t, err)
		require.Equal(t, 1, countBefore)

		// Delete room (soft delete)
		err = repo.Delete(ctx, roomEncx.ID)
		require.NoError(t, err)

		// Count rooms after deletion - should still exist
		var countAfter int
		err = testPool.QueryRow(ctx, "SELECT COUNT(*) FROM booking.rooms WHERE id = $1", roomEncx.ID).Scan(&countAfter)
		require.NoError(t, err)
		require.Equal(t, 1, countAfter, "Room record should still exist after soft delete")

		// But should be inactive
		var activeCount int
		err = testPool.QueryRow(ctx, "SELECT COUNT(*) FROM booking.rooms WHERE id = $1 AND is_active = true", roomEncx.ID).Scan(&activeCount)
		require.NoError(t, err)
		require.Equal(t, 0, activeCount, "Room should be inactive")
	})
}

