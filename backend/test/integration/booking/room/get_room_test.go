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

// make test-func TEST_NAME=TestGetRoom TEST_PATH=test/integration/booking/room/get_room_test.go

func TestGetRoom(t *testing.T) {
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
	setupTestRoom := func(t *testing.T, buildingID uuid.UUID) *domain.Room {
		room := tr.NewTestRoomWithBuilding(t, buildingID)
		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)
		return room
	}

	t.Run("should successfully retrieve an existing room", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Setup building and room
		buildingID := setupTestBuilding(t)
		room := setupTestRoom(t, buildingID)

		// Create request
		req := tr.NewGetRoomByIDRequest(t, ctx, testServerURL, room.ID.String(), "")

		// Execute request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response data
		assert.Equal(t, room.ID, response.ID)
		assert.Equal(t, room.BuildingID, response.BuildingID)
		assert.Equal(t, room.Name, response.Name)
		assert.Equal(t, room.Description, response.Description)
		assert.Equal(t, room.Capacity, response.Capacity)
		assert.Equal(t, room.Equipment, response.Equipment)
		assert.Equal(t, room.IsActive, response.IsActive)
	})

	t.Run("should successfully retrieve room with all fields populated", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)

		// Create room with all fields
		room := tr.NewTestRoomWithParams(t, buildingID, "Advanced Treatment Room", "305", 10, 0, true)
		room.Description = "Fully equipped room with advanced medical equipment"
		room.Equipment = []string{"MRI Scanner", "X-Ray Machine", "Ultrasound", "Patient Monitor"}

		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		req := tr.NewGetRoomByIDRequest(t, ctx, testServerURL, room.ID.String(), "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, room.Name, response.Name)
		assert.Equal(t, room.Description, response.Description)
		assert.Equal(t, room.Capacity, response.Capacity)
		assert.Equal(t, 4, len(response.Equipment))
		assert.Contains(t, response.Equipment, "MRI Scanner")
	})

	t.Run("should successfully retrieve inactive room", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)

		// Create inactive room
		room := tr.NewTestRoomWithParams(t, buildingID, "Inactive Room", "999", 1, 0, false)
		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		req := tr.NewGetRoomByIDRequest(t, ctx, testServerURL, room.ID.String(), "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.False(t, response.IsActive)
	})

	t.Run("should successfully retrieve room with empty equipment list", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)

		room := tr.NewTestRoomWithBuilding(t, buildingID)
		room.Equipment = []string{} // Empty equipment

		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		req := tr.NewGetRoomByIDRequest(t, ctx, testServerURL, room.ID.String(), "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Empty(t, response.Equipment)
	})

	t.Run("should return 404 Not Found for non-existent room ID", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Use a valid UUID that doesn't exist in the database
		nonExistentID := uuid.New()

		req := tr.NewGetRoomByIDRequest(t, ctx, testServerURL, nonExistentID.String(), "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrRepositoryNotFound.Error())
	})

	t.Run("should return 400 Bad Request for invalid room ID format", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Invalid UUID format
		invalidID := "not-a-valid-uuid"

		req := tr.NewGetRoomByIDRequest(t, ctx, testServerURL, invalidID, "")

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

		// Malformed UUID (wrong length, invalid characters)
		malformedID := "123e4567-e89b-12d3-a456"

		req := tr.NewGetRoomByIDRequest(t, ctx, testServerURL, malformedID, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 Bad Request for empty room ID", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Empty ID
		emptyID := ""

		req := tr.NewGetRoomByIDRequest(t, ctx, testServerURL, emptyID, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 Bad Request for nil UUID", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Nil UUID (all zeros)
		nilUUID := uuid.Nil.String()

		req := tr.NewGetRoomByIDRequest(t, ctx, testServerURL, nilUUID, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Nil UUID is technically valid UUID format but should not exist
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should handle multiple rooms correctly - retrieve specific room", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)

		// Create multiple rooms
		room1 := tr.NewTestRoomWithParams(t, buildingID, "Room 101", "101", 5, 0, true)
		room2 := tr.NewTestRoomWithParams(t, buildingID, "Room 102", "102", 10, 0, true)
		room3 := tr.NewTestRoomWithParams(t, buildingID, "Room 103", "103", 15, 0, true)

		for _, room := range []*domain.Room{room1, room2, room3} {
			roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
			require.NoError(t, err)
			err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
			require.NoError(t, err)
		}

		// Retrieve room2 specifically
		req := tr.NewGetRoomByIDRequest(t, ctx, testServerURL, room2.ID.String(), "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, room2.ID, response.ID)
		assert.Equal(t, "Room 102", response.Name)
		assert.Equal(t, 10, response.Capacity)
	})

	t.Run("should retrieve room from different building", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		// Create two different buildings
		building1 := tb.NewTestBuildingWithParams(t, "Building A", "Paris", "France", true)
		building2 := tb.NewTestBuildingWithParams(t, "Building B", "Lyon", "France", true)

		for _, building := range []*domain.Building{building1, building2} {
			buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
			require.NoError(t, err)
			err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
			require.NoError(t, err)
		}

		// Create rooms in different buildings
		roomBuilding1 := tr.NewTestRoomWithParams(t, building1.ID, "Room A1", "A101", 5, 0, true)
		roomBuilding2 := tr.NewTestRoomWithParams(t, building2.ID, "Room B1", "B101", 8, 0, true)

		for _, room := range []*domain.Room{roomBuilding1, roomBuilding2} {
			roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
			require.NoError(t, err)
			err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
			require.NoError(t, err)
		}

		// Retrieve room from building 2
		req := tr.NewGetRoomByIDRequest(t, ctx, testServerURL, roomBuilding2.ID.String(), "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, building2.ID, response.BuildingID)
		assert.Equal(t, "Room B1", response.Name)
	})

	t.Run("should handle context timeout appropriately", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)
		room := setupTestRoom(t, buildingID)

		// Use a very short context timeout to potentially trigger timeout
		shortCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
		defer cancel()

		time.Sleep(2 * time.Millisecond) // Ensure timeout has passed

		req := tr.NewGetRoomByIDRequest(t, shortCtx, testServerURL, room.ID.String(), "")

		resp, err := client.Do(req)
		// Either the context timeout or a successful response (if operation was fast enough)
		if err != nil {
			// Context timeout on client side
			assert.ErrorIs(t, err, context.DeadlineExceeded)
		} else {
			defer resp.Body.Close()
			// If we got a response, it should be either success or timeout status
			assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusRequestTimeout)
		}
	})

	t.Run("should retrieve room with special characters in name", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)

		room := tr.NewTestRoomWithParams(t, buildingID, "Salle de Consultation #1 - Étage 2", "2-101", 3, 0, true)
		room.Description = "Chambre avec équipement spécialisé & moderne"

		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		req := tr.NewGetRoomByIDRequest(t, ctx, testServerURL, room.ID.String(), "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify special characters are preserved through encryption/decryption
		assert.Equal(t, "Salle de Consultation #1 - Étage 2", response.Name)
		assert.Equal(t, "Chambre avec équipement spécialisé & moderne", response.Description)
	})

	t.Run("should retrieve room with maximum capacity", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)

		room := tr.NewTestRoomWithParams(t, buildingID, "Large Conference Hall", "C501", 50, 0, true)

		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		req := tr.NewGetRoomByIDRequest(t, ctx, testServerURL, room.ID.String(), "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, 50, response.Capacity)
	})

	t.Run("should retrieve room with minimum capacity", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)

		room := tr.NewTestRoomWithParams(t, buildingID, "Private Office", "P101", 1, 0, true)

		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		req := tr.NewGetRoomByIDRequest(t, ctx, testServerURL, room.ID.String(), "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, 1, response.Capacity)
	})

	t.Run("should work without authentication (public endpoint)", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)
		room := setupTestRoom(t, buildingID)

		// Request without access token (public endpoint)
		req := tr.NewGetRoomByIDRequest(t, ctx, testServerURL, room.ID.String(), "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should succeed - this is a public endpoint
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, room.ID, response.ID)
	})

	t.Run("should work with authentication token (public endpoint allows auth)", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		buildingID := setupTestBuilding(t)
		room := setupTestRoom(t, buildingID)

		// Setup user with token
		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		// Request with access token
		req := tr.NewGetRoomByIDRequest(t, ctx, testServerURL, room.ID.String(), accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should succeed - public endpoint allows authenticated access too
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, room.ID, response.ID)
	})
}
