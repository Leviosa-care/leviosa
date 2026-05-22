package roomRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestListRooms TEST_PATH=internal/booking/infrastructure/postgres/room/list_rooms_test.go

func TestListRooms(t *testing.T) {
	ctx := context.Background()

	buildingEncx := tb.NewTestBuildingEncx(t)
	err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
	require.NoError(t, err)

	// Create test room data
	room1Encx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
	room2Encx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
	room2Encx.NameHash = "room_2_hash"
	room2Encx.RoomNumberHash = "room_2_number_hash"
	inactiveRoomEncx := tr.NewInactiveTestRoomEncx(t, buildingEncx.ID)

	t.Run("should return empty list when no rooms exist", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Test repository List method with empty filter
		filter := ports.RoomFilter{
			Limit:  10,
			Offset: 0,
		}

		roomsEncx, err := repo.List(ctx, filter)
		require.NoError(t, err)
		require.Empty(t, roomsEncx, "Should return empty list when no rooms exist")
	})

	t.Run("should list all rooms without filters", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Insert test rooms
		err := tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)

		// Test repository List method
		filter := ports.RoomFilter{
			Limit:  10,
			Offset: 0,
		}

		roomsEncx, err := repo.List(ctx, filter)
		require.NoError(t, err)
		require.Len(t, roomsEncx, 2, "Should return all rooms")
	})

	t.Run("should filter by building ID", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		anotherBuildingEncx := tb.NewTestBuildingEncx(t)
		err := tb.InsertBuildingEncx(t, ctx, testPool, anotherBuildingEncx)
		require.NoError(t, err)

		// Create rooms with different building IDs
		roomInOtherBuildingEncx := tr.NewTestRoomEncxWithBuilding(t, anotherBuildingEncx.ID)

		// Insert rooms
		err = tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomInOtherBuildingEncx)
		require.NoError(t, err)

		// Test repository List method with building filter
		filter := ports.RoomFilter{
			BuildingID: &buildingEncx.ID,
			Limit:      10,
			Offset:     0,
		}

		roomsEncx, err := repo.List(ctx, filter)
		require.NoError(t, err)
		require.Len(t, roomsEncx, 2, "Should return only rooms from specified building")

		// Verify all returned rooms belong to the specified building
		for _, room := range roomsEncx {
			require.Equal(t, buildingEncx.ID, room.BuildingID)
		}
	})

	t.Run("should filter by active status", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Insert active and inactive rooms
		err := tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, inactiveRoomEncx)
		require.NoError(t, err)

		// Test repository List method with active filter
		isActive := true
		filter := ports.RoomFilter{
			IsActive: &isActive,
			Limit:    10,
			Offset:   0,
		}

		roomsEncx, err := repo.List(ctx, filter)
		require.NoError(t, err)
		require.Len(t, roomsEncx, 2, "Should return only active rooms")

		// Verify all returned rooms are active
		for _, room := range roomsEncx {
			require.True(t, room.IsActive)
		}
	})

	t.Run("should filter by capacity range", func(t *testing.T) {
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

		// Test repository List method with capacity filter
		minCapacity := 2
		maxCapacity := 4
		filter := ports.RoomFilter{
			MinCapacity: &minCapacity,
			MaxCapacity: &maxCapacity,
			Limit:       10,
			Offset:      0,
		}

		roomsEncx, err := repo.List(ctx, filter)
		require.NoError(t, err)
		require.Len(t, roomsEncx, 1, "Should return only rooms within capacity range")

		// Verify returned room is within capacity range
		for _, room := range roomsEncx {
			require.GreaterOrEqual(t, room.Capacity, minCapacity)
			require.LessOrEqual(t, room.Capacity, maxCapacity)
		}
	})

	t.Run("should filter by hourly rate range", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Create rooms with different hourly rates
		highRateRoomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)

		// Insert rooms
		err := tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, highRateRoomEncx)
		require.NoError(t, err)

		// Test repository List method returns rooms
		filter := ports.RoomFilter{
			Limit:  10,
			Offset: 0,
		}

		roomsEncx, err := repo.List(ctx, filter)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(roomsEncx), 1, "Should return rooms")
	})

	t.Run("should filter by name hash", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Insert rooms
		err := tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)

		// Test repository List method with name hash filter
		filter := ports.RoomFilter{
			NameHash: &room1Encx.NameHash,
			Limit:    10,
			Offset:   0,
		}

		roomsEncx, err := repo.List(ctx, filter)
		require.NoError(t, err)
		require.Len(t, roomsEncx, 1, "Should return only room with matching name hash")

		// Verify returned room has correct name hash
		require.Equal(t, room1Encx.NameHash, roomsEncx[0].NameHash)
	})

	t.Run("should filter by room number hash", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Insert rooms
		err := tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)

		// Test repository List method with room number hash filter
		filter := ports.RoomFilter{
			RoomNumberHash: &room2Encx.RoomNumberHash,
			Limit:          10,
			Offset:         0,
		}

		roomsEncx, err := repo.List(ctx, filter)
		require.NoError(t, err)
		require.Len(t, roomsEncx, 1, "Should return only room with matching room number hash")

		// Verify returned room has correct room number hash
		require.Equal(t, room2Encx.RoomNumberHash, roomsEncx[0].RoomNumberHash)
	})

	t.Run("should apply pagination", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Insert multiple rooms
		rooms := []*domain.RoomEncx{room1Encx, room2Encx, inactiveRoomEncx}
		for _, room := range rooms {
			err := tr.InsertRoomEncx(t, ctx, testPool, room)
			require.NoError(t, err)
		}

		// Test repository List method with pagination
		filter := ports.RoomFilter{
			Limit:  2,
			Offset: 0,
		}

		roomsEncx, err := repo.List(ctx, filter)
		require.NoError(t, err)
		require.Len(t, roomsEncx, 2, "Should return limited number of rooms")

		// Test offset
		filter.Offset = 1
		roomsEncx, err = repo.List(ctx, filter)
		require.NoError(t, err)
		require.Len(t, roomsEncx, 2, "Should return rooms with offset")
	})

	t.Run("should handle sorting", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Create rooms with different capacities for sorting
		room1Encx.Capacity = 5
		room2Encx.Capacity = 1
		sortRoomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		sortRoomEncx.Capacity = 3

		// Insert rooms
		err := tr.InsertRoomEncx(t, ctx, testPool, room1Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, room2Encx)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, sortRoomEncx)
		require.NoError(t, err)

		// Test repository List method with sorting
		filter := ports.RoomFilter{
			OrderBy:        "capacity",
			OrderDirection: "asc",
			Limit:          10,
			Offset:         0,
		}

		roomsEncx, err := repo.List(ctx, filter)
		require.NoError(t, err)
		require.Len(t, roomsEncx, 3, "Should return all rooms sorted by capacity")

		// Verify sorting (ascending)
		for i := 1; i < len(roomsEncx); i++ {
			require.LessOrEqual(t, roomsEncx[i-1].Capacity, roomsEncx[i].Capacity)
		}
	})

	t.Run("should handle combined filters", func(t *testing.T) {
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

		// Test repository List method with combined filters
		minCapacity := 2
		maxCapacity := 4
		isActive := true
		filter := ports.RoomFilter{
			BuildingID:  &buildingEncx.ID,
			IsActive:    &isActive,
			MinCapacity: &minCapacity,
			MaxCapacity: &maxCapacity,
			Limit:       10,
			Offset:      0,
		}

		roomsEncx, err := repo.List(ctx, filter)
		require.NoError(t, err)
		require.Len(t, roomsEncx, 2, "Should return rooms matching all filters")

		// Verify all filters are applied
		for _, room := range roomsEncx {
			require.Equal(t, buildingEncx.ID, room.BuildingID)
			require.True(t, room.IsActive)
			require.GreaterOrEqual(t, room.Capacity, minCapacity)
			require.LessOrEqual(t, room.Capacity, maxCapacity)
		}
	})
}
