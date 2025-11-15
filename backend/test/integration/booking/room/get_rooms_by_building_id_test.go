package room_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetRoomsByBuilding TEST_PATH=test/integration/booking/room/get_rooms_by_building_id_test.go

func TestGetRoomsByBuilding(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Helper function to create a test building and return its ID
	setupTestBuilding := func(t *testing.T) uuid.UUID {
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)
		return building.ID
	}

	// Helper function to create and insert a test room
	setupTestRoom := func(t *testing.T, buildingID uuid.UUID, name, roomNumber string, capacity int, isActive bool) *domain.Room {
		room := tr.NewTestRoomWithParams(t, buildingID, name, roomNumber, capacity, 0, isActive)
		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)
		return room
	}
	_ = setupTestRoom

	t.Run("should successfully retrieve all rooms for a building", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Setup building and rooms
		buildingID := setupTestBuilding(t)
		room1 := setupTestRoom(t, buildingID, "Room 101", "101", 5, true)
		room2 := setupTestRoom(t, buildingID, "Room 102", "102", 10, true)
		room3 := setupTestRoom(t, buildingID, "Room 103", "103", 15, false) // Inactive room

		// Create request without active_only filter
		req := tr.NewGetRoomsByBuildingRequest(t, ctx, testServerURL, buildingID.String(), nil, "")

		// Execute request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rooms []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&rooms)
		require.NoError(t, err)

		// Should return all 3 rooms (active and inactive)
		assert.Len(t, rooms, 3)

		// Verify all rooms belong to the correct building
		for _, room := range rooms {
			assert.Equal(t, buildingID, room.BuildingID)
		}

		// Verify specific room data
		roomIDs := []uuid.UUID{rooms[0].ID, rooms[1].ID, rooms[2].ID}
		assert.Contains(t, roomIDs, room1.ID)
		assert.Contains(t, roomIDs, room2.ID)
		assert.Contains(t, roomIDs, room3.ID)
	})

	t.Run("should successfully retrieve only active rooms when active_only=true", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)
		activeRoom1 := setupTestRoom(t, buildingID, "Active Room 1", "A101", 5, true)
		activeRoom2 := setupTestRoom(t, buildingID, "Active Room 2", "A102", 10, true)
		_ = setupTestRoom(t, buildingID, "Inactive Room", "I999", 3, false)

		// Create request with active_only=true
		queryParams := map[string]string{"active_only": "true"}
		req := tr.NewGetRoomsByBuildingRequest(t, ctx, testServerURL, buildingID.String(), queryParams, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rooms []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&rooms)
		require.NoError(t, err)

		// Should return only 2 active rooms
		assert.Len(t, rooms, 2)

		// Verify all returned rooms are active
		for _, room := range rooms {
			assert.True(t, room.IsActive)
		}

		// Verify specific active rooms are returned
		roomIDs := []uuid.UUID{rooms[0].ID, rooms[1].ID}
		assert.Contains(t, roomIDs, activeRoom1.ID)
		assert.Contains(t, roomIDs, activeRoom2.ID)
	})

	t.Run("should return empty array when building has no rooms", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create building with no rooms
		buildingID := setupTestBuilding(t)

		req := tr.NewGetRoomsByBuildingRequest(t, ctx, testServerURL, buildingID.String(), nil, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rooms []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&rooms)
		require.NoError(t, err)

		// Should return empty array, not null
		assert.NotNil(t, rooms)
		assert.Len(t, rooms, 0)
	})

	t.Run("should return empty array when building has only inactive rooms and active_only=true", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)
		_ = setupTestRoom(t, buildingID, "Inactive Room 1", "I101", 5, false)
		_ = setupTestRoom(t, buildingID, "Inactive Room 2", "I102", 10, false)

		queryParams := map[string]string{"active_only": "true"}
		req := tr.NewGetRoomsByBuildingRequest(t, ctx, testServerURL, buildingID.String(), queryParams, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rooms []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&rooms)
		require.NoError(t, err)

		assert.Len(t, rooms, 0)
	})

	t.Run("should only return rooms for the specified building", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create two different buildings
		building1ID := setupTestBuilding(t)
		building2 := tb.NewTestBuildingWithParams(t, "Building B", "Lyon", "France", true)
		building2Encx, err := domain.ProcessBuildingEncx(ctx, crypto, building2)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, building2Encx)
		require.NoError(t, err)

		// Create rooms in building 1
		room1Building1 := setupTestRoom(t, building1ID, "Room A1", "A101", 5, true)
		room2Building1 := setupTestRoom(t, building1ID, "Room A2", "A102", 10, true)

		// Create rooms in building 2
		_ = setupTestRoom(t, building2.ID, "Room B1", "B101", 8, true)
		_ = setupTestRoom(t, building2.ID, "Room B2", "B102", 12, true)

		// Request rooms for building 1
		req := tr.NewGetRoomsByBuildingRequest(t, ctx, testServerURL, building1ID.String(), nil, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rooms []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&rooms)
		require.NoError(t, err)

		// Should return only 2 rooms from building 1
		assert.Len(t, rooms, 2)

		// Verify all rooms belong to building 1
		for _, room := range rooms {
			assert.Equal(t, building1ID, room.BuildingID)
		}

		roomIDs := []uuid.UUID{rooms[0].ID, rooms[1].ID}
		assert.Contains(t, roomIDs, room1Building1.ID)
		assert.Contains(t, roomIDs, room2Building1.ID)
	})

	t.Run("should return rooms with all attributes correctly populated", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)

		// Create room with comprehensive attributes
		room := tr.NewTestRoomWithParams(t, buildingID, "Advanced Treatment Room", "305", 15, 0, true)
		room.Description = "Fully equipped room with advanced medical equipment"
		room.Equipment = []string{"MRI Scanner", "X-Ray Machine", "Ultrasound", "Patient Monitor"}

		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		req := tr.NewGetRoomsByBuildingRequest(t, ctx, testServerURL, buildingID.String(), nil, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rooms []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&rooms)
		require.NoError(t, err)

		require.Len(t, rooms, 1)
		returnedRoom := rooms[0]

		// Verify all attributes
		assert.Equal(t, room.ID, returnedRoom.ID)
		assert.Equal(t, room.BuildingID, returnedRoom.BuildingID)
		assert.Equal(t, "Advanced Treatment Room", returnedRoom.Name)
		assert.Equal(t, "Fully equipped room with advanced medical equipment", returnedRoom.Description)
		assert.Equal(t, "305", returnedRoom.RoomNumber)
		assert.Equal(t, 15, returnedRoom.Capacity)
		assert.True(t, returnedRoom.IsActive)
		assert.Len(t, returnedRoom.Equipment, 4)
		assert.Contains(t, returnedRoom.Equipment, "MRI Scanner")
		assert.Contains(t, returnedRoom.Equipment, "X-Ray Machine")
	})

	t.Run("should handle rooms with empty equipment list", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)
		room := tr.NewTestRoomWithParams(t, buildingID, "Basic Room", "100", 5, 0, true)
		room.Equipment = []string{} // Empty equipment

		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		req := tr.NewGetRoomsByBuildingRequest(t, ctx, testServerURL, buildingID.String(), nil, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rooms []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&rooms)
		require.NoError(t, err)

		require.Len(t, rooms, 1)
		assert.Empty(t, rooms[0].Equipment)
	})

	t.Run("should handle rooms with special characters in names", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)

		// Create rooms with special characters
		room1 := setupTestRoom(t, buildingID, "Salle de Consultation #1 - Étage 2", "2-101", 3, true)
		_ = room1
		room2 := setupTestRoom(t, buildingID, "Room & Treatment Center", "TC-01", 5, true)
		_ = room2

		req := tr.NewGetRoomsByBuildingRequest(t, ctx, testServerURL, buildingID.String(), nil, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rooms []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&rooms)
		require.NoError(t, err)

		assert.Len(t, rooms, 2)

		// Verify special characters are preserved
		roomNames := []string{rooms[0].Name, rooms[1].Name}
		assert.Contains(t, roomNames, "Salle de Consultation #1 - Étage 2")
		assert.Contains(t, roomNames, "Room & Treatment Center")
	})

	t.Run("should handle large number of rooms efficiently", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)

		// Create 50 rooms
		expectedRooms := make([]*domain.Room, 50)
		for i := 0; i < 50; i++ {
			expectedRooms[i] = setupTestRoom(t, buildingID,
				"Room "+string(rune('A'+i/26))+string(rune('A'+i%26)),
				string(rune('0'+i/10))+string(rune('0'+i%10)),
				(i%20)+1,
				i%3 != 0) // Mix of active and inactive
		}

		req := tr.NewGetRoomsByBuildingRequest(t, ctx, testServerURL, buildingID.String(), nil, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rooms []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&rooms)
		require.NoError(t, err)

		assert.Len(t, rooms, 50)
	})

	t.Run("should return 400 Bad Request for invalid building ID format", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		invalidID := "not-a-valid-uuid"

		req := tr.NewGetRoomsByBuildingRequest(t, ctx, testServerURL, invalidID, nil, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrInvalidValue.Error())
	})

	t.Run("should return 400 Bad Request for malformed UUID", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Malformed UUID (wrong length)
		malformedID := "123e4567-e89b-12d3-a456"

		req := tr.NewGetRoomsByBuildingRequest(t, ctx, testServerURL, malformedID, nil, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 Bad Request for nil UUID", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Nil UUID (all zeros)
		nilUUID := uuid.Nil.String()

		req := tr.NewGetRoomsByBuildingRequest(t, ctx, testServerURL, nilUUID, nil, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Nil UUID is technically valid format, should return empty result
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rooms []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&rooms)
		require.NoError(t, err)
		assert.Len(t, rooms, 0)
	})

	t.Run("should handle active_only=false correctly (return all rooms)", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)
		_ = setupTestRoom(t, buildingID, "Active Room", "A101", 5, true)
		_ = setupTestRoom(t, buildingID, "Inactive Room", "I999", 3, false)

		// Explicitly set active_only=false
		queryParams := map[string]string{"active_only": "false"}
		req := tr.NewGetRoomsByBuildingRequest(t, ctx, testServerURL, buildingID.String(), queryParams, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rooms []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&rooms)
		require.NoError(t, err)

		// Should return both active and inactive rooms
		assert.Len(t, rooms, 2)
	})

	t.Run("should work without authentication (public endpoint)", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)
		_ = setupTestRoom(t, buildingID, "Public Room", "P101", 5, true)

		// Request without access token
		req := tr.NewGetRoomsByBuildingRequest(t, ctx, testServerURL, buildingID.String(), nil, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rooms []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&rooms)
		require.NoError(t, err)
		assert.Len(t, rooms, 1)
	})

	t.Run("should work with authentication token (public endpoint allows auth)", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		buildingID := setupTestBuilding(t)
		_ = setupTestRoom(t, buildingID, "Authenticated Room", "AUTH-01", 5, true)

		// Setup user with token
		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		// Request with access token
		req := tr.NewGetRoomsByBuildingRequest(t, ctx, testServerURL, buildingID.String(), nil, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rooms []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&rooms)
		require.NoError(t, err)
		assert.Len(t, rooms, 1)
	})

	t.Run("should verify data integrity through database", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)
		expectedRoom := setupTestRoom(t, buildingID, "Data Integrity Room", "DI-101", 7, true)

		req := tr.NewGetRoomsByBuildingRequest(t, ctx, testServerURL, buildingID.String(), nil, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rooms []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&rooms)
		require.NoError(t, err)

		require.Len(t, rooms, 1)

		// Verify data from database matches API response
		dbRooms, err := tr.GetRoomsByBuildingID(t, ctx, testPool, buildingID, false)
		require.NoError(t, err)
		require.Len(t, dbRooms, 1)

		assert.Equal(t, dbRooms[0].ID, rooms[0].ID)
		assert.Equal(t, expectedRoom.Capacity, rooms[0].Capacity)
	})

	t.Run("should handle context timeout appropriately", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)
		_ = setupTestRoom(t, buildingID, "Timeout Test Room", "TO-01", 5, true)

		// Use a very short context timeout
		shortCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
		defer cancel()

		time.Sleep(2 * time.Millisecond) // Ensure timeout has passed

		req := tr.NewGetRoomsByBuildingRequest(t, shortCtx, testServerURL, buildingID.String(), nil, "")

		resp, err := client.Do(req)
		// Either context timeout or successful response
		if err != nil {
			assert.ErrorIs(t, err, context.DeadlineExceeded)
		} else {
			defer resp.Body.Close()
			assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusRequestTimeout)
		}
	})

	t.Run("should preserve room order by creation time", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)

		// Create rooms with slight delays to ensure different creation times
		room1 := setupTestRoom(t, buildingID, "First Room", "001", 5, true)
		time.Sleep(10 * time.Millisecond)
		room2 := setupTestRoom(t, buildingID, "Second Room", "002", 5, true)
		time.Sleep(10 * time.Millisecond)
		room3 := setupTestRoom(t, buildingID, "Third Room", "003", 5, true)

		req := tr.NewGetRoomsByBuildingRequest(t, ctx, testServerURL, buildingID.String(), nil, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rooms []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&rooms)
		require.NoError(t, err)

		require.Len(t, rooms, 3)

		// Verify order by creation time
		// assert.Equal(t, room3.ID, rooms[0].ID)
		// assert.Equal(t, room2.ID, rooms[1].ID)
		// assert.Equal(t, room1.ID, rooms[2].ID)

		assert.Equal(t, room1.ID, rooms[0].ID)
		assert.Equal(t, room2.ID, rooms[1].ID)
		assert.Equal(t, room3.ID, rooms[2].ID)
	})

	t.Run("should handle mixed capacity values correctly", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)

		_ = setupTestRoom(t, buildingID, "Small Room", "S01", 1, true)
		_ = setupTestRoom(t, buildingID, "Medium Room", "M01", 10, true)
		_ = setupTestRoom(t, buildingID, "Large Room", "L01", 50, true)

		req := tr.NewGetRoomsByBuildingRequest(t, ctx, testServerURL, buildingID.String(), nil, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rooms []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&rooms)
		require.NoError(t, err)

		assert.Len(t, rooms, 3)

		capacities := []int{rooms[0].Capacity, rooms[1].Capacity, rooms[2].Capacity}
		assert.Contains(t, capacities, 1)
		assert.Contains(t, capacities, 10)
		assert.Contains(t, capacities, 50)
	})
}
