package user_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	td "github.com/Leviosa-care/authuser/test/helpers"

	ck "github.com/Leviosa-care/core/auth/cookies"
	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestGetUserByID make test-integration-user-test

func TestGetUserByID(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return 401 without authentication", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		// Create test user first
		testUser := td.NewTestUser("testuser@example.com", "John", "Doe")
		testUser.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, testUser, testPool, crypto)

		// Act - make request without authentication
		req := td.NewGetUserByIDRequestWithoutAuth(t, ctx, testServerURL, testUser.ID)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Should return 401 without authentication")
	})

	t.Run("should return 403 with non-admin role", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create target user to fetch
		targetUser := td.NewTestUser("target@example.com", "Target", "User")
		targetUser.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, targetUser, testPool, crypto)

		// Create standard user (non-admin) requesting the data
		requestingUser := td.NewTestUser("requesting@example.com", "Requesting", "User")
		requestingUser.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, requestingUser, testPool, crypto)

		// Create session for standard user (non-admin)
		now := time.Now()
		sessionID := uuid.New()

		accessToken, err := session.GenerateToken()
		require.NoError(t, err)

		refreshToken, err := session.GenerateToken()
		require.NoError(t, err)

		standardSession := &session.Session{
			ID:           sessionID,
			UserID:       requestingUser.ID,
			Role:         identity.Standard, // Non-admin role
			State:        session.SessionActive,
			CreatedAt:    now,
			ExpiresAt:    now.Add(24 * time.Hour),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		err = crypto.ProcessStruct(ctx, standardSession)
		require.NoError(t, err)

		td.InsertSessionDirectly(t, ctx, testClient, standardSession, time.Hour)

		// Act - request target user with non-admin session
		req := td.NewGetUserByIDRequest(t, ctx, testServerURL, targetUser.ID, standardSession.AccessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusForbidden, resp.StatusCode, "Should return 403 for non-admin user")
	})

	t.Run("should return 400 for invalid UUID format", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create admin user
		adminUser := td.NewTestUser("admin@example.com", "Administrator", "User")
		adminUser.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, adminUser, testPool, crypto)

		// Create admin session
		now := time.Now()
		sessionID := uuid.New()

		accessToken, err := session.GenerateToken()
		require.NoError(t, err)

		refreshToken, err := session.GenerateToken()
		require.NoError(t, err)

		adminSession := &session.Session{
			ID:           sessionID,
			UserID:       adminUser.ID,
			Role:         identity.Administrator,
			State:        session.SessionActive,
			CreatedAt:    now,
			ExpiresAt:    now.Add(24 * time.Hour),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		err = crypto.ProcessStruct(ctx, adminSession)
		require.NoError(t, err)

		td.InsertSessionDirectly(t, ctx, testClient, adminSession, time.Hour)

		// Act - request with invalid UUID format
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			testServerURL+"/admin/users/invalid-uuid-format",
			nil,
		)
		require.NoError(t, err)

		// Add admin session cookie
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: adminSession.AccessToken,
		}
		req.AddCookie(cookie)

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400 for invalid UUID format")
	})

	t.Run("should return 404 for non-existent user", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create admin user
		adminUser := td.NewTestUser("admin@example.com", "Administrator", "User")
		adminUser.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, adminUser, testPool, crypto)

		// Create admin session
		now := time.Now()
		sessionID := uuid.New()

		accessToken, err := session.GenerateToken()
		require.NoError(t, err)

		refreshToken, err := session.GenerateToken()
		require.NoError(t, err)

		adminSession := &session.Session{
			ID:           sessionID,
			UserID:       adminUser.ID,
			Role:         identity.Administrator,
			State:        session.SessionActive,
			CreatedAt:    now,
			ExpiresAt:    now.Add(24 * time.Hour),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		err = crypto.ProcessStruct(ctx, adminSession)
		require.NoError(t, err)

		td.InsertSessionDirectly(t, ctx, testClient, adminSession, time.Hour)

		// Act - request non-existent user
		nonExistentUserID := uuid.New()
		req := td.NewGetUserByIDRequest(t, ctx, testServerURL, nonExistentUserID, adminSession.AccessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Should return 404 for non-existent user")
	})

	t.Run("should successfully return user data with valid admin authentication", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create target user with specific data
		email := "target@example.com"
		city := "Target City"
		telephone := "0612345678"
		firstname := "Target"
		lastname := "User"

		targetUser := td.NewTestUser(email, firstname, lastname)
		targetUser.State = domain.Active
		targetUser.City = city
		targetUser.Telephone = telephone
		td.InsertUserWithEncryption(t, ctx, targetUser, testPool, crypto)

		// Create admin user
		adminUser := td.NewTestUser("admin@example.com", "Administrator", "User")
		adminUser.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, adminUser, testPool, crypto)

		// Create admin session
		now := time.Now()
		sessionID := uuid.New()

		accessToken, err := session.GenerateToken()
		require.NoError(t, err)

		refreshToken, err := session.GenerateToken()
		require.NoError(t, err)

		adminSession := &session.Session{
			ID:           sessionID,
			UserID:       adminUser.ID,
			Role:         identity.Administrator,
			State:        session.SessionActive,
			CreatedAt:    now,
			ExpiresAt:    now.Add(24 * time.Hour),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		err = crypto.ProcessStruct(ctx, adminSession)
		require.NoError(t, err)

		td.InsertSessionDirectly(t, ctx, testClient, adminSession, time.Hour)

		// Act - request target user with admin authentication
		req := td.NewGetUserByIDRequest(t, ctx, testServerURL, targetUser.ID, adminSession.AccessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 for valid admin request")

		user := td.ParseGetUserByIDResponse(t, resp)
		assert.Equal(t, targetUser.ID, user.ID)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, firstname, user.FirstName)
		assert.Equal(t, lastname, user.LastName)
		assert.Equal(t, city, user.City)
		assert.Equal(t, telephone, user.Telephone)
		assert.Equal(t, domain.Active, user.State)
	})

	t.Run("should handle expired admin session", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create target user
		targetUser := td.NewTestUser("target@example.com", "Target", "User")
		targetUser.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, targetUser, testPool, crypto)

		// Create admin user
		adminUser := td.NewTestUser("admin@example.com", "Administrator", "User")
		adminUser.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, adminUser, testPool, crypto)

		// Create expired admin session
		now := time.Now()
		sessionID := uuid.New()

		accessToken, err := session.GenerateToken()
		require.NoError(t, err)

		refreshToken, err := session.GenerateToken()
		require.NoError(t, err)

		expiredAdministratorSession := &session.Session{
			ID:           sessionID,
			UserID:       adminUser.ID,
			Role:         identity.Administrator,
			State:        session.SessionActive,
			CreatedAt:    now,
			ExpiresAt:    now.Add(-1 * time.Hour), // Expired 1 hour ago
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		err = crypto.ProcessStruct(ctx, expiredAdministratorSession)
		require.NoError(t, err)

		td.InsertSessionDirectly(t, ctx, testClient, expiredAdministratorSession, -time.Hour) // Expired

		// Act - request with expired admin session
		req := td.NewGetUserByIDRequest(t, ctx, testServerURL, targetUser.ID, expiredAdministratorSession.AccessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Should return 401 for expired admin session")
	})

	t.Run("should handle malformed access token", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create target user
		targetUser := td.NewTestUser("target@example.com", "Target", "User")
		targetUser.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, targetUser, testPool, crypto)

		// Act - request with malformed access token
		req := td.NewGetUserByIDRequest(t, ctx, testServerURL, targetUser.ID, "malformed-token")
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Should return 401 for malformed access token")
	})
}
