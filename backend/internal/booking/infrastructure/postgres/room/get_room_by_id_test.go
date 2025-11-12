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

// make test-func TEST_NAME=TestGetRoomByID TEST_PATH=internal/booking/infrastructure/postgres/room/get_room_by_id_test.go

func TestGetRoomByID(t *testing.T) {
	ctx := context.Background()

	// Create test building data first
	buildingEncx := tb.NewTestBuildingEncx(t)

	t.Run("should successfully retrieve room by ID", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Insert building first
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Create test room with valid building ID
		roomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)

		// Insert test room directly
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		// Test repository GetByID method
		retrievedRoomEncx, err := repo.GetByID(ctx, roomEncx.ID)
		require.NoError(t, err)

		// Verify all fields match
		assert.Equal(t, roomEncx.ID, retrievedRoomEncx.ID)
		assert.Equal(t, roomEncx.BuildingID, retrievedRoomEncx.BuildingID)
		assert.Equal(t, roomEncx.NameHash, retrievedRoomEncx.NameHash)
		assert.Equal(t, roomEncx.RoomNumberHash, retrievedRoomEncx.RoomNumberHash)
		assert.Equal(t, roomEncx.Capacity, retrievedRoomEncx.Capacity)
		assert.Equal(t, roomEncx.IsActive, retrievedRoomEncx.IsActive)
		assert.WithinDuration(t, roomEncx.CreatedAt, retrievedRoomEncx.CreatedAt, time.Second)
		assert.WithinDuration(t, roomEncx.UpdatedAt, retrievedRoomEncx.UpdatedAt, time.Second)
	})

	t.Run("should return not found error for non-existent room", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Try to get non-existent room
		nonExistentID := uuid.New()
		_, err := repo.GetByID(ctx, nonExistentID)
		require.Error(t, err)

		// Should be a repository not found error (check the error type)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("should handle invalid UUID format", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Test with invalid UUID - this should be caught at a higher level
		// Since the method signature requires uuid.UUID, this test verifies
		// that the repository handles database-level errors correctly
		invalidUUID := uuid.Nil // Use a valid UUID that won't exist

		_, err := repo.GetByID(ctx, invalidUUID)
		require.Error(t, err)
	})

	t.Run("should retrieve room with all hash fields", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Insert building first
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Create room with complete data including hashes
		completeRoomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)

		err = tr.InsertRoomEncx(t, ctx, testPool, completeRoomEncx)
		require.NoError(t, err)

		// Retrieve room
		retrievedRoomEncx, err := repo.GetByID(ctx, completeRoomEncx.ID)
		require.NoError(t, err)

		// Verify hash fields are correctly retrieved
		assert.Equal(t, completeRoomEncx.NameHash, retrievedRoomEncx.NameHash)
		assert.Equal(t, completeRoomEncx.RoomNumberHash, retrievedRoomEncx.RoomNumberHash)
		assert.NotEmpty(t, retrievedRoomEncx.NameHash, "Name hash should not be empty")
		assert.NotEmpty(t, retrievedRoomEncx.RoomNumberHash, "Room number hash should not be empty")
	})

	t.Run("should retrieve room without hourly rate", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Insert building first
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Create room without hourly rate
		noRateRoomEncx := tr.NewInactiveTestRoomEncx(t, buildingEncx.ID)

		err = tr.InsertRoomEncx(t, ctx, testPool, noRateRoomEncx)
		require.NoError(t, err)

		// Retrieve room
		retrievedRoomEncx, err := repo.GetByID(ctx, noRateRoomEncx.ID)
		require.NoError(t, err)

		// Verify hourly rate is null
		assert.Nil(t, retrievedRoomEncx.HourlyRateCents, "Hourly rate should be null")
	})

	t.Run("should retrieve inactive room", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Insert building first
		err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Create inactive room
		inactiveRoomEncx := tr.NewInactiveTestRoomEncx(t, buildingEncx.ID)

		err = tr.InsertRoomEncx(t, ctx, testPool, inactiveRoomEncx)
		require.NoError(t, err)

		// Retrieve room - should work regardless of active status
		retrievedRoomEncx, err := repo.GetByID(ctx, inactiveRoomEncx.ID)
		require.NoError(t, err)

		// Verify room is inactive
		assert.False(t, retrievedRoomEncx.IsActive, "Room should be inactive")
	})
}

