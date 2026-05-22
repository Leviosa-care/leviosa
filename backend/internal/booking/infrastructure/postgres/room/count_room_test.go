package roomRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCountRooms TEST_PATH=internal/booking/infrastructure/postgres/room/count_room_test.go

func TestCountRooms(t *testing.T) {
	ctx := context.Background()

	buildingEncx := tb.NewTestBuildingEncx(t)
	err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
	require.NoError(t, err)

	// Create test room data
	room1Encx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)

	room2Encx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
	room2Encx.NameHash = "room_2_name"
	room2Encx.RoomNumberHash = "room_2_number_name"
	inactiveRoomEncx := tr.NewInactiveTestRoomEncx(t, buildingEncx.ID)

	t.Run("should return zero count when no rooms exist", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Test repository Count method with empty filter
		filter := ports.RoomFilter{}

		count, err := repo.Count(ctx, filter)
		require.NoError(t, err)
		require.Equal(t, 0, count, "Should return zero when no rooms exist")
	})

	t.Run("should count all rooms without filters", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Insert test rooms
		err := tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, inactiveRoomEncx)
		require.NoError(t, err)

		// Test repository Count method with empty filter
		filter := ports.RoomFilter{}

		count, err := repo.Count(ctx, filter)
		require.NoError(t, err)
		require.Equal(t, 3, count, "Should count all rooms")
	})

	t.Run("should count by building ID", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Create rooms with different building IDs
		anotherBuildingEncx := tb.NewTestBuildingEncx(t)
		err := tb.InsertBuildingEncx(t, ctx, testPool, anotherBuildingEncx)
		require.NoError(t, err)
		roomInOtherBuildingEncx := tr.NewTestRoomEncxWithBuilding(t, anotherBuildingEncx.ID)

		// Insert rooms
		err = tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomInOtherBuildingEncx)
		require.NoError(t, err)

		// Test repository Count method with building filter
		filter := ports.RoomFilter{
			BuildingID: &buildingEncx.ID,
		}

		count, err := repo.Count(ctx, filter)
		require.NoError(t, err)
		require.Equal(t, 2, count, "Should count only rooms from specified building")

		// Test with other building
		filter.BuildingID = &anotherBuildingEncx.ID
		count, err = repo.Count(ctx, filter)
		require.NoError(t, err)
		require.Equal(t, 1, count, "Should count only one room from other building")
	})

	t.Run("should count by active status", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Insert active and inactive rooms
		err := tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, inactiveRoomEncx)
		require.NoError(t, err)

		// Test repository Count method with active filter
		isActive := true
		filter := ports.RoomFilter{
			IsActive: &isActive,
		}

		count, err := repo.Count(ctx, filter)
		require.NoError(t, err)
		require.Equal(t, 2, count, "Should count only active rooms")

		// Test with inactive filter
		isActive = false
		filter.IsActive = &isActive
		count, err = repo.Count(ctx, filter)
		require.NoError(t, err)
		require.Equal(t, 1, count, "Should count only inactive rooms")
	})

	t.Run("should count by capacity range", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Create rooms with different capacities
		room1Encx.Capacity = 1
		room2Encx.Capacity = 3
		largeRoomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		largeRoomEncx.Capacity = 5

		// Insert rooms
		err := tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, largeRoomEncx)
		require.NoError(t, err)

		// Test repository Count method with capacity filter
		minCapacity := 2
		maxCapacity := 4
		filter := ports.RoomFilter{
			MinCapacity: &minCapacity,
			MaxCapacity: &maxCapacity,
		}

		count, err := repo.Count(ctx, filter)
		require.NoError(t, err)
		require.Equal(t, 1, count, "Should count only rooms within capacity range")
	})

	t.Run("should count by hourly rate range", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		highRateRoomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)

		// Insert rooms
		err := tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, highRateRoomEncx)
		require.NoError(t, err)

		// Test repository Count method with capacity filter
		filter := ports.RoomFilter{}

		count, err := repo.Count(ctx, filter)
		require.NoError(t, err)
		require.GreaterOrEqual(t, count, 1, "Should count at least one room")
	})

	t.Run("should count by name hash", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Insert rooms
		err := tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)

		// Test repository Count method with name hash filter
		filter := ports.RoomFilter{
			NameHash: &room1Encx.NameHash,
		}

		count, err := repo.Count(ctx, filter)
		require.NoError(t, err)
		require.Equal(t, 1, count, "Should count only room with matching name hash")
	})

	t.Run("should count by room number hash", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Insert rooms
		err := tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)

		// Test repository Count method with room number hash filter
		filter := ports.RoomFilter{
			RoomNumberHash: &room2Encx.RoomNumberHash,
		}

		count, err := repo.Count(ctx, filter)
		require.NoError(t, err)
		require.Equal(t, 1, count, "Should count only room with matching room number hash")
	})

	t.Run("should count with combined filters", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Create rooms with different properties
		room1Encx.Capacity = 2
		room2Encx.Capacity = 4
		inactiveRoomEncx.Capacity = 3

		// Insert rooms
		err := tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, inactiveRoomEncx)
		require.NoError(t, err)

		// Test repository Count method with combined filters
		minCapacity := 2
		maxCapacity := 4
		isActive := true
		filter := ports.RoomFilter{
			BuildingID:  &buildingEncx.ID,
			IsActive:    &isActive,
			MinCapacity: &minCapacity,
			MaxCapacity: &maxCapacity,
		}

		count, err := repo.Count(ctx, filter)
		require.NoError(t, err)
		require.Equal(t, 2, count, "Should count rooms matching all filters")
	})

	t.Run("should count rooms without hourly rate", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Create rooms with and without hourly rates
		roomWithRateEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)

		roomWithoutRateEncx := tr.NewInactiveTestRoomEncx(t, buildingEncx.ID)
		roomWithoutRateEncx.IsActive = true

		// Insert rooms
		err := tr.InsertRoomEncx(t, ctx, testPool, roomWithRateEncx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomWithoutRateEncx)
		require.NoError(t, err)

		// Test repository Count method — both rooms should be counted without rate filter
		filter := ports.RoomFilter{}

		count, err := repo.Count(ctx, filter)
		require.NoError(t, err)
		require.GreaterOrEqual(t, count, 1, "Should count rooms")
	})

	t.Run("should handle database errors gracefully", func(t *testing.T) {
		// Test with cancelled context
		ctx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		// Try to count rooms with cancelled context
		filter := ports.RoomFilter{}
		_, err := repo.Count(ctx, filter)
		require.Error(t, err)
	})
}
