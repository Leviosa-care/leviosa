package room_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetAllRooms TEST_PATH=test/integration/booking/room/get_all_rooms_test.go

func TestGetAllRooms(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 15 * time.Second}

	// Helper function to create a test building and return its ID
	setupTestBuilding := func(t *testing.T) *domain.Building {
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)
		return building
	}

	// Helper function to create and insert a test room
	setupTestRoom := func(t *testing.T, building *domain.Building) *domain.Room {
		room := tr.NewTestRoomWithBuilding(t, building.ID)
		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)
		return room
	}
	_ = setupTestRoom

	t.Run("should successfully get all rooms without authentication", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create building and 3 rooms
		building := setupTestBuilding(t)
		for i := 0; i < 3; i++ {
			setupTestRoom(t, building)
		}

		// Make request without authentication
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, nil, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify we got all 3 rooms
		assert.Len(t, response, 3)
	})

	t.Run("should successfully get all rooms with standard user authentication", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup standard user
		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		// Create building and 2 rooms
		building := setupTestBuilding(t)
		for i := 0; i < 2; i++ {
			setupTestRoom(t, building)
		}

		// Make request with standard user authentication
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, nil, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Len(t, response, 2)
	})

	t.Run("should successfully get all rooms with admin authentication", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create building and room
		building := setupTestBuilding(t)
		setupTestRoom(t, building)

		// Make request with admin authentication
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, nil, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Len(t, response, 1)
	})

	t.Run("should return empty array when no rooms exist", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Make request
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, nil, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Empty(t, response)
	})

	t.Run("should filter rooms by is_active=true", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create building
		building := setupTestBuilding(t)

		// Create 2 active and 1 inactive room
		activeRoom1 := tr.NewTestRoomWithParams(t, building.ID, "Active Room 1", "101", 5, 0, true)
		activeRoom2 := tr.NewTestRoomWithParams(t, building.ID, "Active Room 2", "102", 10, 0, true)
		inactiveRoom := tr.NewTestRoomWithParams(t, building.ID, "Inactive Room", "103", 3, 0, false)

		for _, room := range []*domain.Room{activeRoom1, activeRoom2, inactiveRoom} {
			roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
			require.NoError(t, err)
			err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
			require.NoError(t, err)
		}

		// Make request with is_active=true filter
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"is_active": "true",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should only return 2 active rooms
		assert.Len(t, response, 2)
		for _, room := range response {
			assert.True(t, room.IsActive)
		}
	})

	t.Run("should filter rooms by is_active=false", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create building
		building := setupTestBuilding(t)

		// Create 1 active and 2 inactive rooms
		activeRoom := tr.NewTestRoomWithParams(t, building.ID, "Active Room", "101", 5, 0, true)
		inactiveRoom1 := tr.NewTestRoomWithParams(t, building.ID, "Inactive Room 1", "102", 10, 0, false)
		inactiveRoom2 := tr.NewTestRoomWithParams(t, building.ID, "Inactive Room 2", "103", 3, 0, false)

		for _, room := range []*domain.Room{activeRoom, inactiveRoom1, inactiveRoom2} {
			roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
			require.NoError(t, err)
			err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
			require.NoError(t, err)
		}

		// Make request with is_active=false filter
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"is_active": "false",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should only return 2 inactive rooms
		assert.Len(t, response, 2)
		for _, room := range response {
			assert.False(t, room.IsActive)
		}
	})

	t.Run("should filter rooms by building_id", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create 2 buildings
		building1 := setupTestBuilding(t)
		building2 := setupTestBuilding(t)

		// Create 2 rooms in building1 and 1 room in building2
		for i := 0; i < 2; i++ {
			setupTestRoom(t, building1)
		}
		setupTestRoom(t, building2)

		// Make request filtering by building1 ID
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"building_id": building1.ID.String(),
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should only return 2 rooms from building1
		assert.Len(t, response, 2)
		for _, room := range response {
			assert.Equal(t, building1.ID, room.BuildingID)
		}
	})

	t.Run("should filter rooms by min_capacity", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create building
		building := setupTestBuilding(t)

		// Create rooms with different capacities
		room1 := tr.NewTestRoomWithParams(t, building.ID, "Small Room", "101", 2, 0, true)
		room2 := tr.NewTestRoomWithParams(t, building.ID, "Medium Room", "102", 10, 0, true)
		room3 := tr.NewTestRoomWithParams(t, building.ID, "Large Room", "103", 20, 0, true)

		for _, room := range []*domain.Room{room1, room2, room3} {
			roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
			require.NoError(t, err)
			err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
			require.NoError(t, err)
		}

		// Make request with min_capacity=10 filter
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"min_capacity": "10",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return 2 rooms (capacity >= 10)
		assert.Len(t, response, 2)
		for _, room := range response {
			assert.GreaterOrEqual(t, room.Capacity, 10)
		}
	})

	t.Run("should filter rooms by max_capacity", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create building
		building := setupTestBuilding(t)

		// Create rooms with different capacities
		room1 := tr.NewTestRoomWithParams(t, building.ID, "Small Room", "101", 2, 0, true)
		room2 := tr.NewTestRoomWithParams(t, building.ID, "Medium Room", "102", 10, 0, true)
		room3 := tr.NewTestRoomWithParams(t, building.ID, "Large Room", "103", 20, 0, true)

		for _, room := range []*domain.Room{room1, room2, room3} {
			roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
			require.NoError(t, err)
			err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
			require.NoError(t, err)
		}

		// Make request with max_capacity=10 filter
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"max_capacity": "10",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return 2 rooms (capacity <= 10)
		assert.Len(t, response, 2)
		for _, room := range response {
			assert.LessOrEqual(t, room.Capacity, 10)
		}
	})

	t.Run("should filter rooms by capacity range (min and max)", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create building
		building := setupTestBuilding(t)

		// Create rooms with different capacities
		room1 := tr.NewTestRoomWithParams(t, building.ID, "Tiny Room", "101", 1, 0, true)
		room2 := tr.NewTestRoomWithParams(t, building.ID, "Small Room", "102", 5, 0, true)
		room3 := tr.NewTestRoomWithParams(t, building.ID, "Medium Room", "103", 15, 0, true)
		room4 := tr.NewTestRoomWithParams(t, building.ID, "Large Room", "104", 30, 0, true)

		for _, room := range []*domain.Room{room1, room2, room3, room4} {
			roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
			require.NoError(t, err)
			err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
			require.NoError(t, err)
		}

		// Make request with min_capacity=5 and max_capacity=20 filter
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"min_capacity": "5",
			"max_capacity": "20",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return 2 rooms (capacity between 5 and 20)
		assert.Len(t, response, 2)
		for _, room := range response {
			assert.GreaterOrEqual(t, room.Capacity, 5)
			assert.LessOrEqual(t, room.Capacity, 20)
		}
	})

	t.Run("should filter rooms by name (GDPR-compliant hash search)", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create building
		building := setupTestBuilding(t)

		// Create rooms with specific names
		conferenceRoom := tr.NewTestRoomWithParams(t, building.ID, "Conference Room A", "101", 10, 0, true)
		meetingRoom := tr.NewTestRoomWithParams(t, building.ID, "Meeting Room B", "102", 5, 0, true)
		officeRoom := tr.NewTestRoomWithParams(t, building.ID, "Office Room C", "103", 2, 0, true)

		for _, room := range []*domain.Room{conferenceRoom, meetingRoom, officeRoom} {
			roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
			require.NoError(t, err)
			err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
			require.NoError(t, err)
		}

		// Make request filtering by exact name
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"name": "Conference Room A",
		}, "")
		// req := tr.NewListRoomsRequest(t, ctx, testServerURL, nil, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// bodyBytes, _ := io.ReadAll(resp.Body)
		// fmt.Printf("Status: %d, Body: %s\n", resp.StatusCode, string(bodyBytes))

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return 1 room with matching name
		assert.Len(t, response, 1)
		assert.Equal(t, "Conference Room A", response[0].Name)
	})

	t.Run("should filter rooms by room_number (GDPR-compliant hash search)", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create building
		building := setupTestBuilding(t)

		// Create rooms with specific room numbers
		room1 := tr.NewTestRoomWithParams(t, building.ID, "Room A", "101", 5, 0, true)
		room2 := tr.NewTestRoomWithParams(t, building.ID, "Room B", "201", 10, 0, true)
		room3 := tr.NewTestRoomWithParams(t, building.ID, "Room C", "301", 15, 0, true)

		for _, room := range []*domain.Room{room1, room2, room3} {
			roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
			require.NoError(t, err)
			err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
			require.NoError(t, err)
		}

		// Make request filtering by room number
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"room_number": "201",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return 1 room with matching room number
		assert.Len(t, response, 1)
		// Note: RoomNumber is not in RoomResponse, verify by other fields
		assert.Equal(t, "Room B", response[0].Name)
	})

	t.Run("should combine multiple filters (building_id, is_active, min_capacity)", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create 2 buildings
		building1 := setupTestBuilding(t)
		building2 := setupTestBuilding(t)

		// Create various rooms
		b1Room1Active := tr.NewTestRoomWithParams(t, building1.ID, "B1 Room 1", "101", 10, 0, true)
		b1Room2Inactive := tr.NewTestRoomWithParams(t, building1.ID, "B1 Room 2", "102", 15, 0, false)
		b1Room3ActiveSmall := tr.NewTestRoomWithParams(t, building1.ID, "B1 Room 3", "103", 3, 0, true)
		b2Room1Active := tr.NewTestRoomWithParams(t, building2.ID, "B2 Room 1", "201", 20, 0, true)

		for _, room := range []*domain.Room{b1Room1Active, b1Room2Inactive, b1Room3ActiveSmall, b2Room1Active} {
			roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
			require.NoError(t, err)
			err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
			require.NoError(t, err)
		}

		// Make request with multiple filters
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"building_id":  building1.ID.String(),
			"is_active":    "true",
			"min_capacity": "5",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return 1 room (building1, active, capacity >= 5)
		assert.Len(t, response, 1)
		assert.Equal(t, building1.ID, response[0].BuildingID)
		assert.True(t, response[0].IsActive)
		assert.GreaterOrEqual(t, response[0].Capacity, 5)
	})

	t.Run("should respect limit parameter", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create building
		building := setupTestBuilding(t)

		// Create 5 rooms
		for i := 0; i < 5; i++ {
			setupTestRoom(t, building)
		}

		// Make request with limit=2
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"limit": "2",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should only return 2 rooms
		assert.Len(t, response, 2)
	})

	t.Run("should respect offset parameter", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create building
		building := setupTestBuilding(t)

		// Create 3 rooms
		for i := 0; i < 3; i++ {
			setupTestRoom(t, building)
		}

		// Make request with offset=2
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"offset": "2",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return 1 room (skipped first 2)
		assert.Len(t, response, 1)
	})

	t.Run("should combine limit and offset for pagination", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create building
		building := setupTestBuilding(t)

		// Create 10 rooms
		for i := 0; i < 10; i++ {
			setupTestRoom(t, building)
		}

		// Make request with limit=3 and offset=2
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"limit":  "3",
			"offset": "2",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return 3 rooms (items 3-5)
		assert.Len(t, response, 3)
	})

	t.Run("should return 400 for invalid is_active parameter", func(t *testing.T) {
		// Make request with invalid is_active value
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"is_active": "invalid",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "is_active must be a boolean")
	})

	t.Run("should return 400 for invalid building_id parameter", func(t *testing.T) {
		// Make request with invalid UUID
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"building_id": "not-a-uuid",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "building_id must be a valid UUID")
	})

	t.Run("should return 400 for invalid min_capacity parameter", func(t *testing.T) {
		// Make request with invalid min_capacity
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"min_capacity": "invalid",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "min_capacity must be a non-negative integer")
	})

	t.Run("should return 400 for negative min_capacity", func(t *testing.T) {
		// Make request with negative min_capacity
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"min_capacity": "-5",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "min_capacity must be a non-negative integer")
	})

	t.Run("should return 400 for invalid max_capacity parameter", func(t *testing.T) {
		// Make request with invalid max_capacity
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"max_capacity": "invalid",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "max_capacity must be a non-negative integer")
	})

	t.Run("should return 400 when min_capacity > max_capacity", func(t *testing.T) {
		// Make request with min > max
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"min_capacity": "20",
			"max_capacity": "10",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "min_capacity cannot be greater than max_capacity")
	})

	t.Run("should return 400 for invalid limit parameter", func(t *testing.T) {
		// Make request with invalid limit value
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"limit": "invalid",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "limit must be a positive integer")
	})

	t.Run("should return 400 for limit exceeding maximum", func(t *testing.T) {
		// Make request with limit > 100
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"limit": "101",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "limit must be a positive integer between 1 and 100")
	})

	t.Run("should return 400 for negative offset parameter", func(t *testing.T) {
		// Make request with negative offset
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"offset": "-1",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "offset must be a non-negative integer")
	})

	t.Run("should return 400 for invalid order_by parameter", func(t *testing.T) {
		// Make request with invalid order_by value
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"order_by": "invalid_field",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "order_by must be one of: name, created_at, capacity")
	})

	t.Run("should return 400 for invalid order_direction parameter", func(t *testing.T) {
		// Make request with invalid order_direction value
		req := tr.NewListRoomsRequest(t, ctx, testServerURL, map[string]string{
			"order_direction": "invalid",
		}, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "order_direction must be either 'asc' or 'desc'")
	})

	t.Run("should handle context timeout appropriately", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create building and room
		building := setupTestBuilding(t)
		setupTestRoom(t, building)

		// Use a very short context timeout
		shortCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
		defer cancel()

		time.Sleep(2 * time.Millisecond) // Ensure timeout has passed

		// Make request with short context
		req := tr.NewListRoomsRequest(t, shortCtx, testServerURL, nil, "")

		resp, err := client.Do(req)
		// Either the context timeout or a successful response
		if err != nil {
			// Context timeout on client side
			assert.ErrorIs(t, err, context.DeadlineExceeded)
		} else {
			defer resp.Body.Close()
			// If we got a response, it should be either success or timeout status
			assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusRequestTimeout)
		}
	})
}
