package user_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	userEndpoints "github.com/Leviosa-care/leviosa/backend/internal/authuser/interface/user"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME='^TestUpdateUser$$' TEST_PATH=test/integration/authuser/user/update_user_test.go

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return 401 without authentication", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		// Create test user first
		testUser := td.NewTestUser(t, "testuser@example.com", "John", "Doe")
		testUser.State = domain.Active

		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
		require.NoError(t, err)

		updateRequest := domain.UpdateUserRequest{
			FirstName: stringPtr("UpdatedName"),
		}

		// Act - make request without authentication
		req := td.NewUpdateUserRequest(t, ctx, testServerURL, updateRequest)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Should return 401 without authentication")
	})

	t.Run("should successfully update user profile with valid authentication", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		email := "updateuser@example.com"
		originalFirstName := "Original"
		originalLastName := "User"
		originalCity := "OriginalCity"

		// Create test user
		testUser := td.NewTestUser(t, email, originalFirstName, originalLastName)
		testUser.State = domain.Active
		testUser.City = originalCity

		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
		require.NoError(t, err)

		// Create valid session
		now := time.Now()
		sessionID := uuid.New()

		accessToken, err := session.GenerateToken()
		require.NoError(t, err)

		refreshToken, err := session.GenerateToken()
		require.NoError(t, err)

		standardSession := &session.Session{
			ID:           sessionID,
			UserID:       testUser.ID,
			Role:         identity.Standard,
			State:        session.SessionActive,
			CreatedAt:    now,
			ExpiresAt:    now.Add(24 * time.Hour),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		standardSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, standardSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, standardSessionEncx, time.Hour)

		// Prepare update request with new values
		newFirstName := "Updated"
		newCity := "UpdatedCity"
		newAddress1 := "123 New Street"

		updateRequest := domain.UpdateUserRequest{
			FirstName: &newFirstName,
			City:      &newCity,
			Address1:  &newAddress1,
			// Intentionally leave LastName unchanged to test partial updates
		}

		// Create request with authentication
		req := td.NewUpdateUserRequestWithAuth(t, ctx, testServerURL, updateRequest, standardSession.AccessToken)
		resp, err := client.Do(req)

		// Assert response
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		updatedUser := td.ParseUpdateUserResponse(t, resp)
		assert.Equal(t, testUser.ID, updatedUser.ID)
		assert.Equal(t, email, updatedUser.Email)
		assert.Equal(t, newFirstName, updatedUser.FirstName)    // Updated field
		assert.Equal(t, originalLastName, updatedUser.LastName) // Unchanged field
		assert.Equal(t, newCity, updatedUser.City)              // Updated field
		assert.Equal(t, newAddress1, updatedUser.Address1)      // Updated field

		// Verify data was persisted in database
		retrievedUserEncx, err := td.GetUserEnxByID(t, ctx, testUser.ID, testPool)
		assert.NoError(t, err)

		retrievedUser, err := domain.DecryptUserEncx(ctx, crypto, retrievedUserEncx)
		assert.NoError(t, err)

		assert.Equal(t, newFirstName, retrievedUser.FirstName)
		assert.Equal(t, originalLastName, retrievedUser.LastName) // Should remain unchanged
		assert.Equal(t, newCity, retrievedUser.City)
		assert.Equal(t, newAddress1, retrievedUser.Address1)
	})

	t.Run("should handle invalid JSON in request body", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create test user and session
		testUser := td.NewTestUser(t, "jsontest@example.com", "Test", "User")
		testUser.State = domain.Active

		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
		require.NoError(t, err)

		// Create valid standardSession
		standardSession := createTestSession(t, ctx, testUser.ID)

		standardSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, standardSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, standardSessionEncx, time.Hour)

		// Create request with invalid JSON
		req := td.NewInvalidJSONRequest(t, ctx, testServerURL, http.MethodPatch, userEndpoints.UpdateUserEndpoint)
		td.AddAuthCookie(req, standardSession.AccessToken)

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should successfully update only email field", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		originalEmail := "original@example.com"
		newEmail := "newemail@example.com"

		firstname := "Test"
		lastname := "User"

		// Create test user
		testUser := td.NewTestUser(t, originalEmail, firstname, lastname)
		testUser.State = domain.Active

		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
		require.NoError(t, err)

		// Create valid standardSession
		standardSession := createTestSession(t, ctx, testUser.ID)

		standardSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, standardSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, standardSessionEncx, time.Hour)

		// Update only email
		updateRequest := domain.UpdateUserRequest{
			Email: &newEmail,
		}

		req := td.NewUpdateUserRequestWithAuth(t, ctx, testServerURL, updateRequest, standardSession.AccessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		updatedUser := td.ParseUpdateUserResponse(t, resp)
		assert.Equal(t, newEmail, updatedUser.Email)
		assert.Equal(t, firstname, updatedUser.FirstName) // Should remain unchanged
		assert.Equal(t, lastname, updatedUser.LastName)   // Should remain unchanged
	})

	t.Run("should handle birth date updates correctly", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create test user
		testUser := td.NewTestUser(t, "datetest@example.com", "Date", "Tester")
		testUser.State = domain.Active

		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
		require.NoError(t, err)

		// Create valid standardSession
		standardSession := createTestSession(t, ctx, testUser.ID)

		standardSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, standardSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, standardSessionEncx, time.Hour)

		// Update birth date
		newBirthDate := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
		updateRequest := domain.UpdateUserRequest{
			BirthDate: &newBirthDate,
		}

		req := td.NewUpdateUserRequestWithAuth(t, ctx, testServerURL, updateRequest, standardSession.AccessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		updatedUser := td.ParseUpdateUserResponse(t, resp)
		assert.Equal(t, newBirthDate.Format(time.RFC3339), updatedUser.BirthDate.Format(time.RFC3339))
	})

	t.Run("should handle expired session", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create test user
		testUser := td.NewTestUser(t, "expiredtest@example.com", "Expired", "User")
		testUser.State = domain.Active

		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
		require.NoError(t, err)

		// Create expired session
		expiredSession := createTestSession(t, ctx, testUser.ID)
		expiredSession.ExpiresAt = time.Now().Add(-1 * time.Hour) // Expired 1 hour ago

		expiredSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, expiredSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, expiredSessionEncx, -time.Hour) // Expired

		updateRequest := domain.UpdateUserRequest{
			FirstName: stringPtr("ShouldNotUpdate"),
		}

		// Create request with expired session
		req := td.NewUpdateUserRequestWithAuth(t, ctx, testServerURL, updateRequest, expiredSession.AccessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should update multiple fields in single request", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create test user
		testUser := td.NewTestUser(t, "multiupdate@example.com", "Multi", "Update")
		testUser.State = domain.Active

		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
		require.NoError(t, err)

		// Create valid session
		standardSession := createTestSession(t, ctx, testUser.ID)

		standardSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, standardSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, standardSessionEncx, time.Hour)

		// Update multiple fields
		updateRequest := domain.UpdateUserRequest{
			FirstName:  stringPtr("NewFirst"),
			LastName:   stringPtr("NewLast"),
			City:       stringPtr("NewCity"),
			PostalCode: stringPtr("12345"),
			Address1:   stringPtr("456 New Avenue"),
			Gender:     stringPtr("other"),
		}

		req := td.NewUpdateUserRequestWithAuth(t, ctx, testServerURL, updateRequest, standardSession.AccessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		updatedUser := td.ParseUpdateUserResponse(t, resp)
		assert.Equal(t, "NewFirst", updatedUser.FirstName)
		assert.Equal(t, "NewLast", updatedUser.LastName)
		assert.Equal(t, "NewCity", updatedUser.City)
		assert.Equal(t, "12345", updatedUser.PostalCode)
		assert.Equal(t, "456 New Avenue", updatedUser.Address1)
		assert.Equal(t, "other", updatedUser.Gender)
	})
}

// Helper function to create a pointer to string
func stringPtr(s string) *string {
	return &s
}

// Helper function to create a test session
func createTestSession(t *testing.T, ctx context.Context, userID uuid.UUID) *session.Session {
	t.Helper()

	now := time.Now()
	sessionID := uuid.New()

	accessToken, err := session.GenerateToken()
	require.NoError(t, err)

	refreshToken, err := session.GenerateToken()
	require.NoError(t, err)

	return &session.Session{
		ID:           sessionID,
		UserID:       userID,
		Role:         identity.Standard,
		State:        session.SessionActive,
		CreatedAt:    now,
		ExpiresAt:    now.Add(24 * time.Hour),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
