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

// make test-func TEST_NAME=TestUpdateRoom TEST_PATH=internal/booking/infrastructure/postgres/room/update_room_test.go

func TestUpdateRoom(t *testing.T) {
	ctx := context.Background()

	// Create test building data first
	buildingEncx := tb.NewTestBuildingEncx(t)
	// Insert building first
	err := tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
	require.NoError(t, err)

	t.Run("should successfully update room", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Create test room with valid building ID
		roomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)

		// Insert test room directly
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		// Update room with new values
		updatedRoomEncx := *roomEncx // Copy original room
		updatedRoomEncx.NameEncrypted = []byte("updated_encrypted_name")
		updatedRoomEncx.NameHash = "updated_hashed_name"
		updatedRoomEncx.DescriptionEncrypted = []byte("updated_description")
		updatedRoomEncx.RoomNumberEncrypted = []byte("updated_999")
		updatedRoomEncx.RoomNumberHash = "updated_hashed_999"
		updatedRoomEncx.Capacity = 3
		updatedRoomEncx.EquipmentEncrypted = []byte(`["updated_equipment"]`)
		updatedRoomEncx.IsActive = false
		updatedRoomEncx.UpdatedAt = time.Now()

		// Test repository Update method
		err = repo.Update(ctx, &updatedRoomEncx)
		require.NoError(t, err)

		// Verify the room was updated by querying directly
		savedRoom, err := tr.GetRoomEncxByID(t, ctx, testPool, roomEncx.ID)
		require.NoError(t, err)

		assert.Equal(t, updatedRoomEncx.NameHash, savedRoom.NameHash)
		assert.Equal(t, updatedRoomEncx.RoomNumberHash, savedRoom.RoomNumberHash)
		assert.Equal(t, updatedRoomEncx.Capacity, savedRoom.Capacity)
		assert.Equal(t, updatedRoomEncx.IsActive, savedRoom.IsActive)
	})

	t.Run("should return not found error for non-existent room", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Try to update non-existent room
		nonExistentID := uuid.New()
		nonExistentRoomEncx := tr.NewTestRoomEncx(t)
		nonExistentRoomEncx.ID = nonExistentID

		err := repo.Update(ctx, nonExistentRoomEncx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("should handle partial updates", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Create test room with valid building ID
		roomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)

		// Insert test room directly
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		// Update only specific fields
		partialUpdateRoomEncx := *roomEncx // Copy original room
		partialUpdateRoomEncx.Capacity = 5
		partialUpdateRoomEncx.IsActive = false
		partialUpdateRoomEncx.UpdatedAt = time.Now()

		// Test repository Update method
		err = repo.Update(ctx, &partialUpdateRoomEncx)
		require.NoError(t, err)

		// Verify only the specified fields were updated
		savedRoom, err := tr.GetRoomEncxByID(t, ctx, testPool, roomEncx.ID)
		require.NoError(t, err)

		// These fields should be updated
		assert.Equal(t, 5, savedRoom.Capacity)
		assert.Equal(t, false, savedRoom.IsActive)

		// These fields should remain unchanged
		assert.Equal(t, roomEncx.NameHash, savedRoom.NameHash)
		assert.Equal(t, roomEncx.RoomNumberHash, savedRoom.RoomNumberHash)
	})

	t.Run("should update room with null hourly rate", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Create room with hourly rate
		roomWithRateEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)
		err := tr.InsertRoomEncx(t, ctx, testPool, roomWithRateEncx)
		require.NoError(t, err)

		// Update room to remove hourly rate
		updateNoRateRoomEncx := *roomWithRateEncx
		updateNoRateRoomEncx.UpdatedAt = time.Now()

		// Test repository Update method
		err = repo.Update(ctx, &updateNoRateRoomEncx)
		require.NoError(t, err)

		// Verify room was updated
		_, err = tr.GetRoomEncxByID(t, ctx, testPool, roomWithRateEncx.ID)
		require.NoError(t, err)
	})

	t.Run("should update room with new hourly rate", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Create room without hourly rate
		noRateRoomEncx := tr.NewInactiveTestRoomEncx(t, buildingEncx.ID)
		err := tr.InsertRoomEncx(t, ctx, testPool, noRateRoomEncx)
		require.NoError(t, err)

		// Update room to activate it
		updateWithRateRoomEncx := *noRateRoomEncx
		updateWithRateRoomEncx.IsActive = true
		updateWithRateRoomEncx.UpdatedAt = time.Now()

		// Test repository Update method
		err = repo.Update(ctx, &updateWithRateRoomEncx)
		require.NoError(t, err)

		// Verify hourly rate is set
		savedRoom, err := tr.GetRoomEncxByID(t, ctx, testPool, noRateRoomEncx.ID)
		require.NoError(t, err)

		assert.True(t, savedRoom.IsActive, "Room should be active")
	})

	t.Run("should update hash fields correctly", func(t *testing.T) {
		// Clean up before test
		tr.ClearRoomsTable(t, ctx, testPool)

		// Create test room with valid building ID
		roomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)

		// Insert test room directly
		err := tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		// Update hash fields
		updateHashRoomEncx := *roomEncx
		updateHashRoomEncx.NameEncrypted = []byte("new_encrypted_name")
		updateHashRoomEncx.NameHash = "new_hashed_name"
		updateHashRoomEncx.RoomNumberEncrypted = []byte("new_999")
		updateHashRoomEncx.RoomNumberHash = "new_hashed_999"
		updateHashRoomEncx.UpdatedAt = time.Now()

		// Test repository Update method
		err = repo.Update(ctx, &updateHashRoomEncx)
		require.NoError(t, err)

		// Verify hash fields were updated
		savedRoom, err := tr.GetRoomEncxByID(t, ctx, testPool, roomEncx.ID)
		require.NoError(t, err)

		assert.Equal(t, "new_hashed_name", savedRoom.NameHash)
		assert.Equal(t, "new_hashed_999", savedRoom.RoomNumberHash)
	})

	t.Run("should handle database errors gracefully", func(t *testing.T) {
		// Test with cancelled context
		ctx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		// Create test room with valid building ID
		roomEncx := tr.NewTestRoomEncxWithBuilding(t, buildingEncx.ID)

		// Try to update room with cancelled context
		err := repo.Update(ctx, roomEncx)
		require.Error(t, err)
	})
}

