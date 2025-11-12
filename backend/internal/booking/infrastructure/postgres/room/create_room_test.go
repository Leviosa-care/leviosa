package roomRepository_test

import (
	"context"
	"testing"

	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreateRoom TEST_PATH=internal/booking/infrastructure/postgres/room/create_room_test.go

func TestCreateRoom(t *testing.T) {
	ctx := context.Background()

	// Create test building data first
	buildingEncx := tb.NewTestBuildingEncx(t)

	t.Run("should successfully create a room", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Insert building first
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Create test room with valid building ID
		roomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)

		// Test repository Create method
		err = repo.Create(ctx, roomEncx)
		require.NoError(t, err)

		// Verify the room was inserted by querying directly
		var count int
		err = testPool.QueryRow(ctx,
			"SELECT COUNT(*) FROM booking.rooms WHERE id = $1",
			roomEncx.ID).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count, "Room should be inserted in database")
	})

	t.Run("should handle duplicate room ID", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Insert building first
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Create test room with valid building ID
		roomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)

		// Insert first room
		err = repo.Create(ctx, roomEncx)
		require.NoError(t, err)

		// Try to insert room with same ID
		err = repo.Create(ctx, roomEncx)
		assert.Error(t, err, "Should fail on duplicate ID")
	})

	t.Run("should handle invalid building ID (foreign key constraint)", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create room with non-existent building ID
		invalidRoomEncx := tr.NewTestRoomEncx(t)
		// Use a UUID that doesn't exist in buildings table
		err := repo.Create(ctx, invalidRoomEncx)
		assert.Error(t, err, "Should fail with foreign key constraint violation")
	})

	t.Run("should create room with all fields including hashes", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Insert building first
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Create room with complete data
		completeRoomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)

		err = repo.Create(ctx, completeRoomEncx)
		require.NoError(t, err)

		// Verify all fields were inserted correctly
		savedRoom, err := tr.GetRoomEncxByID(t, ctx, testPool, completeRoomEncx.ID)
		require.NoError(t, err)

		assert.Equal(t, completeRoomEncx.ID, savedRoom.ID)
		assert.Equal(t, completeRoomEncx.BuildingID, savedRoom.BuildingID)
		assert.Equal(t, completeRoomEncx.NameHash, savedRoom.NameHash)
		assert.Equal(t, completeRoomEncx.RoomNumberHash, savedRoom.RoomNumberHash)
		assert.Equal(t, completeRoomEncx.Capacity, savedRoom.Capacity)
		assert.Equal(t, *completeRoomEncx.HourlyRateCents, *savedRoom.HourlyRateCents)
		assert.Equal(t, completeRoomEncx.IsActive, savedRoom.IsActive)
	})

	t.Run("should create room without hourly rate (nullable field)", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Insert building first
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Create room without hourly rate
		noRateRoomEncx := tr.NewInactiveTestRoomEncx(t, buildingEncx.ID)

		err = repo.Create(ctx, noRateRoomEncx)
		require.NoError(t, err)

		// Verify hourly rate is null
		savedRoom, err := tr.GetRoomEncxByID(t, ctx, testPool, noRateRoomEncx.ID)
		require.NoError(t, err)

		assert.Nil(t, savedRoom.HourlyRateCents, "Hourly rate should be null")
	})
}
