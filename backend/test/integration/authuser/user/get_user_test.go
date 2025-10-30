package user_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetUser TEST_PATH=test/integration/authuser/user/get_user_test.go

func TestGetUser(t *testing.T) {
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

		// Act - make request without authentication
		req := td.NewGetUserRequest(t, ctx, testServerURL)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Should return 401 without authentication")
	})

	t.Run("should return user data with valid authentication", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		email := "authenticateduser@example.com"
		city := "Test City"
		telephone := "0612345678"
		firstname := "Jane"
		lastname := "Smith"

		// Create test user with specific data
		testUser := td.NewTestUser(t, email, firstname, lastname)
		testUser.State = domain.Active
		testUser.City = city
		testUser.Telephone = telephone

		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
		require.NoError(t, err)

		// Create a valid session for the user
		now := time.Now()
		sessionID := uuid.New()

		// Generate valid base64 tokens for testing
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

		// Create request with authentication cookie
		req := td.NewGetUserRequestWithAuth(t, ctx, testServerURL, standardSession.AccessToken)
		resp, err := client.Do(req)

		// Assert response
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return 200 with user data since we have proper auth middleware
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		user := td.ParseGetUserResponse(t, resp)
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, firstname, user.FirstName)
		assert.Equal(t, lastname, user.LastName)
		assert.Equal(t, city, user.City)
		assert.Equal(t, telephone, user.Telephone)
	})

	t.Run("should handle invalid session token format", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create request with malformed session token
		req := td.NewGetUserRequestWithAuth(t, ctx, testServerURL, "invalid-token-format")
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode,
			"Should return 401 for malformed session token")
	})

	t.Run("should handle expired session", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create test user
		testUser := td.NewTestUser(t, "expireduser@example.com", "Expired", "User")
		testUser.State = domain.Active
		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
		require.NoError(t, err)

		// Create expired session (negative duration to ensure it's expired)
		now := time.Now()
		sessionID := uuid.New()

		// Generate valid base64 tokens for testing
		accessToken, err := session.GenerateToken()
		require.NoError(t, err)

		refreshToken, err := session.GenerateToken()
		require.NoError(t, err)

		expiredSession := &session.Session{
			ID:           sessionID,
			UserID:       testUser.ID,
			Role:         identity.Standard,
			State:        session.SessionActive,
			CreatedAt:    now,
			ExpiresAt:    now.Add(-1 * time.Hour), // Expired 1 hour ago
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		expiredSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, expiredSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, expiredSessionEncx, -time.Hour) // Expired 1 hour ago

		// Create request with expired session
		req := td.NewGetUserRequestWithAuth(t, ctx, testServerURL, expiredSession.AccessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode,
			"Should return 401 for expired session")
	})
}
