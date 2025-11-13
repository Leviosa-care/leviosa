package room_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestUpdateRoom TEST_PATH=test/integration/booking/room/update_room_test.go

func TestUpdateRoom(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Helper function to setup test data (building + room)
	setupTestRoomWithBuilding := func(t *testing.T) (*domain.Room, uuid.UUID) {
		// Create building
		building := tb.NewTestBuilding(t)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		// Create room
		room := tr.NewTestRoomWithBuilding(t, building.ID)
		roomEncx, err := domain.ProcessRoomEncx(ctx, crypto, room)
		require.NoError(t, err)
		err = tr.InsertRoomEncx(t, ctx, testPool, roomEncx)
		require.NoError(t, err)

		return room, building.ID
	}

	t.Run("SuccessfulUpdates", func(t *testing.T) {
		t.Run("should successfully update all room fields", func(t *testing.T) {
			// Clean test data
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			// Setup admin user
			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			// Setup test data
			room, _ := setupTestRoomWithBuilding(t)

			// Prepare update request with all fields
			newName := "Updated Treatment Room"
			newDescription := "Updated description for modern treatments"
			newRoomNumber := "302"
			newCapacity := 5
			newEquipment := []string{"new equipment 1", "new equipment 2"}
			newIsActive := false

			updateRequest := domain.UpdateRoomRequest{
				ID:          room.ID,
				Name:        &newName,
				Description: &newDescription,
				RoomNumber:  &newRoomNumber,
				Capacity:    &newCapacity,
				Equipment:   &newEquipment,
				IsActive:    &newIsActive,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, accessToken)

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
			assert.Equal(t, newName, response.Name)
			assert.Equal(t, newDescription, response.Description)
			assert.Equal(t, newRoomNumber, response.RoomNumber)
			assert.Equal(t, newCapacity, response.Capacity)
			assert.Equal(t, newEquipment, response.Equipment)
			assert.Equal(t, newIsActive, response.IsActive)

			// Verify database state
			roomEncx, err := tr.GetRoomEncxByID(t, ctx, testPool, room.ID)
			require.NoError(t, err)

			updatedRoom, err := domain.DecryptRoomEncx(ctx, crypto, roomEncx)
			require.NoError(t, err)

			assert.Equal(t, newName, updatedRoom.Name)
			assert.Equal(t, newDescription, updatedRoom.Description)
			assert.Equal(t, newRoomNumber, updatedRoom.RoomNumber)
			assert.Equal(t, newCapacity, updatedRoom.Capacity)
			assert.Equal(t, newEquipment, updatedRoom.Equipment)
			assert.Equal(t, newIsActive, updatedRoom.IsActive)
		})

		t.Run("should successfully update only room name", func(t *testing.T) {
			// Clean test data
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)
			room, _ := setupTestRoomWithBuilding(t)

			// Update only name
			newName := "Only Name Updated"
			updateRequest := domain.UpdateRoomRequest{
				ID:   room.ID,
				Name: &newName,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, accessToken)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response domain.RoomResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			// Verify only name changed
			assert.Equal(t, newName, response.Name)
			assert.Equal(t, room.Description, response.Description)
			assert.Equal(t, room.RoomNumber, response.RoomNumber)
			assert.Equal(t, room.Capacity, response.Capacity)
		})

		t.Run("should successfully update only room description", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)
			room, _ := setupTestRoomWithBuilding(t)

			newDescription := "Only description updated"
			updateRequest := domain.UpdateRoomRequest{
				ID:          room.ID,
				Description: &newDescription,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, accessToken)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response domain.RoomResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.Equal(t, newDescription, response.Description)
			assert.Equal(t, room.Name, response.Name)
		})

		t.Run("should successfully update only room number", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)
			room, _ := setupTestRoomWithBuilding(t)

			newRoomNumber := "999"
			updateRequest := domain.UpdateRoomRequest{
				ID:         room.ID,
				RoomNumber: &newRoomNumber,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, accessToken)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response domain.RoomResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.Equal(t, newRoomNumber, response.RoomNumber)
			assert.Equal(t, room.Name, response.Name)
		})

		t.Run("should successfully update only capacity", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)
			room, _ := setupTestRoomWithBuilding(t)

			newCapacity := 10
			updateRequest := domain.UpdateRoomRequest{
				ID:       room.ID,
				Capacity: &newCapacity,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, accessToken)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response domain.RoomResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.Equal(t, newCapacity, response.Capacity)
			assert.Equal(t, room.Name, response.Name)
		})

		t.Run("should successfully update only equipment", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)
			room, _ := setupTestRoomWithBuilding(t)

			newEquipment := []string{"brand new equipment"}
			updateRequest := domain.UpdateRoomRequest{
				ID:        room.ID,
				Equipment: &newEquipment,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, accessToken)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response domain.RoomResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.Equal(t, newEquipment, response.Equipment)
			assert.Equal(t, room.Name, response.Name)
		})

		t.Run("should successfully update only is_active", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)
			room, _ := setupTestRoomWithBuilding(t)

			newIsActive := false
			updateRequest := domain.UpdateRoomRequest{
				ID:       room.ID,
				IsActive: &newIsActive,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, accessToken)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response domain.RoomResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.Equal(t, newIsActive, response.IsActive)
			assert.Equal(t, room.Name, response.Name)
		})

		t.Run("should successfully update multiple fields but not all", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)
			room, _ := setupTestRoomWithBuilding(t)

			// Update name, capacity, and is_active only
			newName := "Partially Updated Room"
			newCapacity := 7
			newIsActive := false

			updateRequest := domain.UpdateRoomRequest{
				ID:       room.ID,
				Name:     &newName,
				Capacity: &newCapacity,
				IsActive: &newIsActive,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, accessToken)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response domain.RoomResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			// Verify updated fields
			assert.Equal(t, newName, response.Name)
			assert.Equal(t, newCapacity, response.Capacity)
			assert.Equal(t, newIsActive, response.IsActive)

			// Verify unchanged fields
			assert.Equal(t, room.Description, response.Description)
			assert.Equal(t, room.RoomNumber, response.RoomNumber)
			assert.Equal(t, room.Equipment, response.Equipment)
		})
	})

	t.Run("ValidationErrors", func(t *testing.T) {
		t.Run("should return 400 for invalid room ID format in URL", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			newName := "Test"
			updateRequest := domain.UpdateRoomRequest{
				ID:   uuid.New(), // This is not used because URL parsing happens first
				Name: &newName,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, "invalid-uuid", updateRequest, accessToken)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("should return 400 for empty name", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)
			room, _ := setupTestRoomWithBuilding(t)

			emptyName := ""
			updateRequest := domain.UpdateRoomRequest{
				ID:   room.ID,
				Name: &emptyName,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, accessToken)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("should return 400 for name exceeding maximum length", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)
			room, _ := setupTestRoomWithBuilding(t)

			// Create a name with 256 characters (exceeds 255 limit)
			longName := string(make([]byte, 256))
			for i := range longName {
				longName = longName[:i] + "a"
			}

			updateRequest := domain.UpdateRoomRequest{
				ID:   room.ID,
				Name: &longName,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, accessToken)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("should return 400 for description exceeding maximum length", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)
			room, _ := setupTestRoomWithBuilding(t)

			// Create a description with 1001 characters (exceeds 1000 limit)
			longDescription := string(make([]byte, 1001))
			for i := range longDescription {
				longDescription = longDescription[:i] + "a"
			}

			updateRequest := domain.UpdateRoomRequest{
				ID:          room.ID,
				Description: &longDescription,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, accessToken)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("should return 400 for room number exceeding maximum length", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)
			room, _ := setupTestRoomWithBuilding(t)

			// Create a room number with 51 characters (exceeds 50 limit)
			longRoomNumber := string(make([]byte, 51))
			for i := range longRoomNumber {
				longRoomNumber = longRoomNumber[:i] + "1"
			}

			updateRequest := domain.UpdateRoomRequest{
				ID:         room.ID,
				RoomNumber: &longRoomNumber,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, accessToken)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("should return 400 for negative capacity", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)
			room, _ := setupTestRoomWithBuilding(t)

			negativeCapacity := -1
			updateRequest := domain.UpdateRoomRequest{
				ID:       room.ID,
				Capacity: &negativeCapacity,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, accessToken)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("should return 400 for zero capacity", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)
			room, _ := setupTestRoomWithBuilding(t)

			zeroCapacity := 0
			updateRequest := domain.UpdateRoomRequest{
				ID:       room.ID,
				Capacity: &zeroCapacity,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, accessToken)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("should return 400 for invalid building ID format", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)
			room, _ := setupTestRoomWithBuilding(t)

			invalidBuildingID := "not-a-valid-uuid"
			updateRequest := domain.UpdateRoomRequest{
				ID:         room.ID,
				BuildingID: &invalidBuildingID,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, accessToken)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})

	t.Run("NotFoundErrors", func(t *testing.T) {
		t.Run("should return 404 for non-existent room", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)

			nonExistentID := uuid.New()
			newName := "Test"
			updateRequest := domain.UpdateRoomRequest{
				ID:   nonExistentID,
				Name: &newName,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, nonExistentID.String(), updateRequest, accessToken)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	})

	t.Run("AuthorizationErrors", func(t *testing.T) {
		t.Run("should return 401 when access token is missing", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)

			room, _ := setupTestRoomWithBuilding(t)

			newName := "Test"
			updateRequest := domain.UpdateRoomRequest{
				ID:   room.ID,
				Name: &newName,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, "")
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should return 401 when token is invalid", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)

			room, _ := setupTestRoomWithBuilding(t)

			newName := "Test"
			updateRequest := domain.UpdateRoomRequest{
				ID:   room.ID,
				Name: &newName,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, "invalid-token-12345")
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should return 401 when session is expired", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			// Create expired admin session
			accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

			room, _ := setupTestRoomWithBuilding(t)

			newName := "Test"
			updateRequest := domain.UpdateRoomRequest{
				ID:   room.ID,
				Name: &newName,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, accessToken)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			// Create standard user (not admin)
			accessToken := tu.SetupStandardUser(t, ctx, authCtx)

			room, _ := setupTestRoomWithBuilding(t)

			newName := "Test"
			updateRequest := domain.UpdateRoomRequest{
				ID:   room.ID,
				Name: &newName,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, accessToken)
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})
	})

	t.Run("HTTPErrors", func(t *testing.T) {
		t.Run("should return 415 for non-JSON content type", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)
			room, _ := setupTestRoomWithBuilding(t)

			newName := "Test"
			updateRequest := domain.UpdateRoomRequest{
				ID:   room.ID,
				Name: &newName,
			}

			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), updateRequest, accessToken)
			req.Header.Set("Content-Type", "text/plain")

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
		})

		t.Run("should return 400 for invalid JSON body", func(t *testing.T) {
			tr.ClearRoomsTable(t, ctx, testPool)
			tb.ClearBuildingsTable(t, ctx, testPool)
			defer tu.ClearAuthData(t, ctx, authCtx)

			accessToken := tu.SetupAdminUser(t, ctx, authCtx)
			room, _ := setupTestRoomWithBuilding(t)

			// Create request with invalid JSON
			req := tr.NewUpdateRoomRequest(t, ctx, testServerURL, room.ID.String(), "invalid json", accessToken)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
}
