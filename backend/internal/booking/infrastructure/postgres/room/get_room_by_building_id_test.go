package roomRepository_test

import (
	"context"
	"testing"

	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetRoomByBuildingID TEST_PATH=internal/booking/infrastructure/postgres/room/get_room_by_building_id_test.go

func TestGetRoomByBuildingID(t *testing.T) {
	ctx := context.Background()

	// Create test building data first
	buildingEncx := tb.NewTestBuildingEncx(t)
	err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
	require.NoError(t, err)

	// Create test room data
	room1Encx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
	room2Encx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
	inactiveRoomEncx := tr.NewInactiveTestRoomEncx(t, buildingEncx.ID)

	t.Run("should return empty list when no rooms exist for building", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Test repository GetByBuildingID method
		roomsEncx, err := repo.GetByBuildingID(ctx, buildingEncx.ID, false)
		require.NoError(t, err)
		require.Empty(t, roomsEncx, "Should return empty list when no rooms exist for building")
	})

	t.Run("should return all rooms for building (including inactive)", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Insert test rooms
		err := tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, inactiveRoomEncx)
		require.NoError(t, err)

		// Test repository GetByBuildingID method (activeOnly=false)
		roomsEncx, err := repo.GetByBuildingID(ctx, buildingEncx.ID, false)
		require.NoError(t, err)
		require.Len(t, roomsEncx, 3, "Should return all rooms including inactive")

		// Verify all rooms belong to the specified building
		for _, room := range roomsEncx {
			require.Equal(t, buildingEncx.ID, room.BuildingID)
		}
	})

	t.Run("should return only active rooms for building", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Insert test rooms
		err := tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, inactiveRoomEncx)
		require.NoError(t, err)

		// Test repository GetByBuildingID method (activeOnly=true)
		roomsEncx, err := repo.GetByBuildingID(ctx, buildingEncx.ID, true)
		require.NoError(t, err)
		require.Len(t, roomsEncx, 2, "Should return only active rooms")

		// Verify all returned rooms are active
		for _, room := range roomsEncx {
			require.True(t, room.IsActive)
		}
	})

	t.Run("should return empty list for non-existent building", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Try to get rooms for non-existent building
		nonExistentBuildingID := tr.NewTestRoom(t).BuildingID
		roomsEncx, err := repo.GetByBuildingID(ctx, nonExistentBuildingID, false)
		require.NoError(t, err)
		require.Empty(t, roomsEncx, "Should return empty list for non-existent building")
	})

	t.Run("should return rooms with all fields including hashes", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Insert test rooms with complete data
		err := tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)

		// Test repository GetByBuildingID method
		roomsEncx, err := repo.GetByBuildingID(ctx, buildingEncx.ID, false)
		require.NoError(t, err)
		require.Len(t, roomsEncx, 2, "Should return all rooms")

		// Verify all fields are present
		for _, room := range roomsEncx {
			assert.NotEmpty(t, room.NameHash, "Name hash should be present")
			assert.NotEmpty(t, room.RoomNumberHash, "Room number hash should be present")
			assert.Greater(t, room.Capacity, 0, "Capacity should be positive")
			assert.NotZero(t, room.CreatedAt, "CreatedAt should be set")
			assert.NotZero(t, room.UpdatedAt, "UpdatedAt should be set")
		}
	})

	t.Run("should handle building with no active rooms when activeOnly=true", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Insert only inactive rooms
		err := tr.InsertRoomEncx(t, ctx, testPool, inactiveRoomEncx)
		require.NoError(t, err)

		// Test repository GetByBuildingID method (activeOnly=true)
		roomsEncx, err := repo.GetByBuildingID(ctx, buildingEncx.ID, true)
		require.NoError(t, err)
		require.Empty(t, roomsEncx, "Should return empty list when no active rooms exist")

		// But should return room when activeOnly=false
		roomsEncx, err = repo.GetByBuildingID(ctx, buildingEncx.ID, false)
		require.NoError(t, err)
		require.Len(t, roomsEncx, 1, "Should return inactive room when activeOnly=false")
		require.False(t, roomsEncx[0].IsActive, "Room should be inactive")
	})

	t.Run("should return rooms with and without hourly rates", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Create rooms with different hourly rate scenarios
		rate := 5000
		roomWithRateEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		roomWithRateEncx.HourlyRateCents = &rate

		roomWithoutRateEncx := tr.NewInactiveTestRoomEncx(t, buildingEncx.ID)
		roomWithoutRateEncx.IsActive = true

		// Insert rooms
		err := tr.InsertRoomEncx(t, ctx, testPool, roomWithRateEncx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomWithoutRateEncx)
		require.NoError(t, err)

		// Test repository GetByBuildingID method
		roomsEncx, err := repo.GetByBuildingID(ctx, buildingEncx.ID, false)
		require.NoError(t, err)
		require.Len(t, roomsEncx, 2, "Should return all rooms")

		// Verify hourly rate handling
		var withRateCount, withoutRateCount int
		for _, room := range roomsEncx {
			if room.HourlyRateCents != nil {
				withRateCount++
			} else {
				withoutRateCount++
			}
		}

		require.Equal(t, 1, withRateCount, "Should have one room with hourly rate")
		require.Equal(t, 1, withoutRateCount, "Should have one room without hourly rate")
	})

	t.Run("should handle building with single room", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Insert single room
		err := tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)

		// Test repository GetByBuildingID method
		roomsEncx, err := repo.GetByBuildingID(ctx, buildingEncx.ID, false)
		require.NoError(t, err)
		require.Len(t, roomsEncx, 1, "Should return single room")

		// Verify room details
		singleRoom := roomsEncx[0]
		require.Equal(t, room1Encx.ID, singleRoom.ID)
		require.Equal(t, buildingEncx.ID, singleRoom.BuildingID)
		require.Equal(t, room1Encx.NameHash, singleRoom.NameHash)
		require.Equal(t, room1Encx.RoomNumberHash, singleRoom.RoomNumberHash)
		require.Equal(t, room1Encx.Capacity, singleRoom.Capacity)
	})

	t.Run("should handle database errors gracefully", func(t *testing.T) {
		// Test with cancelled context
		ctx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		// Try to get rooms with cancelled context
		_, err := repo.GetByBuildingID(ctx, buildingEncx.ID, false)
		require.Error(t, err)
	})
}
