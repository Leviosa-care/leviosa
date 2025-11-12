package room_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	tb "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
	tr "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreateRoom TEST_PATH=test/integration/booking/room/create_room_test.go

func TestCreateRoom(t *testing.T) {
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

	t.Run("should successfully create a room with valid admin token", func(t *testing.T) {
		// Clean test data
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and get access token
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create test building
		buildingID := setupTestBuilding(t)

		// Prepare request
		request := domain.CreateRoomRequest{
			BuildingID:  buildingID,
			Name:        "Conference Room A",
			Description: "Large conference room with projector",
			RoomNumber:  "101",
			Capacity:    20,
			Equipment:   []string{"Projector", "Whiteboard", "Video Conferencing"},
			IsActive:    true,
		}

		req := tr.NewCreateRoomRequest(t, ctx, testServerURL, request, accessToken)

		// Execute request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response data
		assert.NotEqual(t, uuid.Nil, response.ID)
		assert.Equal(t, buildingID, response.BuildingID)
		assert.Equal(t, request.Name, response.Name)
		assert.Equal(t, request.Description, response.Description)
		assert.Equal(t, request.Capacity, response.Capacity)
		assert.Equal(t, request.Equipment, response.Equipment)
		assert.True(t, response.IsActive)

		// Verify room exists in database
		roomEncx, err := tr.GetRoomEncxByID(t, ctx, testPool, response.ID)
		require.NoError(t, err)

		room, err := domain.DecryptRoomEncx(ctx, crypto, roomEncx)
		require.NoError(t, err)

		assert.Equal(t, buildingID, room.BuildingID)
		assert.Equal(t, request.Name, room.Name)
		assert.Equal(t, request.Description, room.Description)
		assert.Equal(t, request.RoomNumber, room.RoomNumber)
		assert.Equal(t, request.Capacity, room.Capacity)
		assert.Equal(t, request.Equipment, room.Equipment)
		assert.True(t, room.IsActive)
	})

	t.Run("should successfully create room with minimal required fields", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		buildingID := setupTestBuilding(t)

		request := domain.CreateRoomRequest{
			BuildingID: buildingID,
			Name:       "Meeting Room",
			Capacity:   5,
			IsActive:   true,
		}

		req := tr.NewCreateRoomRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotEqual(t, uuid.Nil, response.ID)
		assert.Equal(t, buildingID, response.BuildingID)
		assert.Equal(t, request.Name, response.Name)
		assert.Empty(t, response.Description)
		assert.Equal(t, request.Capacity, response.Capacity)
	})

	t.Run("should successfully create room with equipment list", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		buildingID := setupTestBuilding(t)

		equipment := []string{"Desk", "Chair", "Computer", "Monitor", "Phone"}
		request := domain.CreateRoomRequest{
			BuildingID: buildingID,
			Name:       "Office Room",
			Capacity:   1,
			Equipment:  equipment,
			IsActive:   true,
		}

		req := tr.NewCreateRoomRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response domain.RoomResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, equipment, response.Equipment)

		// Verify in database
		roomEncx, err := tr.GetRoomEncxByID(t, ctx, testPool, response.ID)
		require.NoError(t, err)

		room, err := domain.DecryptRoomEncx(ctx, crypto, roomEncx)
		require.NoError(t, err)

		assert.Equal(t, equipment, room.Equipment)
	})

	t.Run("should return 400 Bad Request for empty name", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		buildingID := setupTestBuilding(t)

		request := domain.CreateRoomRequest{
			BuildingID: buildingID,
			Name:       "",
			Capacity:   10,
			IsActive:   true,
		}

		req := tr.NewCreateRoomRequest(t, ctx, testServerURL, request, accessToken)

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

	t.Run("should return 400 Bad Request for nil building ID", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.CreateRoomRequest{
			BuildingID: uuid.Nil,
			Name:       "Test Room",
			Capacity:   10,
			IsActive:   true,
		}

		req := tr.NewCreateRoomRequest(t, ctx, testServerURL, request, accessToken)

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

	t.Run("should return 400 Bad Request for non-existent building", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.CreateRoomRequest{
			BuildingID: uuid.New(), // Non-existent building
			Name:       "Test Room",
			Capacity:   10,
			IsActive:   true,
		}

		req := tr.NewCreateRoomRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 Bad Request for inactive building", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create inactive building
		building := tb.NewTestBuildingWithParams(t, "Inactive Building", "Paris", "France", false)
		buildingEncx, err := domain.ProcessBuildingEncx(ctx, crypto, building)
		require.NoError(t, err)
		err = tb.InsertBuildingEncx(t, ctx, testPool, buildingEncx)
		require.NoError(t, err)

		request := domain.CreateRoomRequest{
			BuildingID: building.ID,
			Name:       "Test Room",
			Capacity:   10,
			IsActive:   true,
		}

		req := tr.NewCreateRoomRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "inactive building")
	})

	t.Run("should return 400 Bad Request for invalid capacity (zero)", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		buildingID := setupTestBuilding(t)

		request := domain.CreateRoomRequest{
			BuildingID: buildingID,
			Name:       "Test Room",
			Capacity:   0,
			IsActive:   true,
		}

		req := tr.NewCreateRoomRequest(t, ctx, testServerURL, request, accessToken)

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

	t.Run("should return 400 Bad Request for invalid capacity (negative)", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		buildingID := setupTestBuilding(t)

		request := domain.CreateRoomRequest{
			BuildingID: buildingID,
			Name:       "Test Room",
			Capacity:   -5,
			IsActive:   true,
		}

		req := tr.NewCreateRoomRequest(t, ctx, testServerURL, request, accessToken)

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

	t.Run("should return 400 Bad Request for capacity exceeding maximum", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		buildingID := setupTestBuilding(t)

		request := domain.CreateRoomRequest{
			BuildingID: buildingID,
			Name:       "Test Room",
			Capacity:   51, // Max is 50
			IsActive:   true,
		}

		req := tr.NewCreateRoomRequest(t, ctx, testServerURL, request, accessToken)

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

	t.Run("should return 400 Bad Request for name exceeding maximum length", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		buildingID := setupTestBuilding(t)

		longName := string(make([]byte, 256)) // Max is 255
		for i := range longName {
			longName = string([]byte(longName)[:i]) + "a"
		}

		request := domain.CreateRoomRequest{
			BuildingID: buildingID,
			Name:       longName,
			Capacity:   10,
			IsActive:   true,
		}

		req := tr.NewCreateRoomRequest(t, ctx, testServerURL, request, accessToken)

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

	t.Run("should return 400 Bad Request for invalid JSON body", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create malformed JSON
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/rooms", bytes.NewBuffer([]byte("{invalid json")))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		if accessToken != "" {
			cookie := &http.Cookie{
				Name:  ck.AccessTokenCookieName,
				Value: accessToken,
			}
			req.AddCookie(cookie)
		}

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

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)

		request := domain.CreateRoomRequest{
			BuildingID: buildingID,
			Name:       "Test Room",
			Capacity:   10,
			IsActive:   true,
		}

		req := tr.NewCreateRoomRequest(t, ctx, testServerURL, request, "") // Empty token

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
		buildingID := setupTestBuilding(t)

		request := domain.CreateRoomRequest{
			BuildingID: buildingID,
			Name:       "Test Room",
			Capacity:   10,
			IsActive:   true,
		}

		req := tr.NewCreateRoomRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role (standard user)", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create standard user (not admin)
		accessToken := tu.SetupStandardUser(t, ctx, authCtx)
		buildingID := setupTestBuilding(t)

		request := domain.CreateRoomRequest{
			BuildingID: buildingID,
			Name:       "Test Room",
			Capacity:   10,
			IsActive:   true,
		}

		req := tr.NewCreateRoomRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role (partner)", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create partner user (not admin)
		accessToken := tu.SetupUserWithRole(t, ctx, identity.Partner, authCtx)
		buildingID := setupTestBuilding(t)

		request := domain.CreateRoomRequest{
			BuildingID: buildingID,
			Name:       "Test Room",
			Capacity:   10,
			IsActive:   true,
		}

		req := tr.NewCreateRoomRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)

		buildingID := setupTestBuilding(t)

		request := domain.CreateRoomRequest{
			BuildingID: buildingID,
			Name:       "Test Room",
			Capacity:   10,
			IsActive:   true,
		}

		req := tr.NewCreateRoomRequest(t, ctx, testServerURL, request, "invalid-token-12345")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should handle context timeout appropriately", func(t *testing.T) {
		tr.ClearRoomsTable(t, ctx, testPool)
		tb.ClearBuildingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		buildingID := setupTestBuilding(t)

		request := domain.CreateRoomRequest{
			BuildingID: buildingID,
			Name:       "Test Room",
			Capacity:   10,
			IsActive:   true,
		}

		// Use a very short context timeout to potentially trigger timeout
		shortCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
		defer cancel()

		time.Sleep(2 * time.Millisecond) // Ensure timeout has passed

		req := tr.NewCreateRoomRequest(t, shortCtx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		// Either the context timeout or a successful response (if operation was fast enough)
		if err != nil {
			// Context timeout on client side
			assert.ErrorIs(t, err, context.DeadlineExceeded)
		} else {
			defer resp.Body.Close()
			// If we got a response, it should be either success or timeout status
			assert.True(t, resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusRequestTimeout)
		}
	})
}
