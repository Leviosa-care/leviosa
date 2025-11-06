package user_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetUserByID TEST_PATH=test/integration/authuser/user/get_user_by_id_test.go

func TestGetUserByID(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	setTargetUser := func() uuid.UUID {
		targetUser := td.NewTestUser(t, "target@example.com", "Target", "User")
		targetUser.State = domain.Active

		targetUserEncx, err := domain.ProcessUserEncx(ctx, crypto, targetUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, targetUserEncx, testPool)
		require.NoError(t, err)

		return targetUser.ID
	}

	t.Run("should return 401 without authentication", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		// Create test user first
		userID := setTargetUser()

		// Act - make request without authentication
		req := td.NewGetUserByIDRequestWithoutAuth(t, ctx, testServerURL, userID)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Should return 401 without authentication")
	})

	t.Run("should return 403 with non-admin role", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create target user to fetch
		userID := setTargetUser()

		// Create standard user (non-admin) requesting the data
		requestingUser := td.NewTestUser(t, "requesting@example.com", "Requesting", "User")
		requestingUser.State = domain.Active
		requestingUserEncx, err := domain.ProcessUserEncx(ctx, crypto, requestingUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, requestingUserEncx, testPool)
		require.NoError(t, err)

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

		standardSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, standardSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, standardSessionEncx, time.Hour)

		// Act - request target user with non-admin session
		req := td.NewGetUserByIDRequest(t, ctx, testServerURL, userID, standardSession.AccessToken)

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusForbidden, resp.StatusCode, "Should return 403 for non-admin user")
	})

	t.Run("should return 400 for invalid UUID format", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create admin user
		adminUser := td.NewTestUser(t, "admin@example.com", "Administrator", "User")
		adminUser.State = domain.Active
		adminUserEncx, err := domain.ProcessUserEncx(ctx, crypto, adminUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, adminUserEncx, testPool)
		require.NoError(t, err)

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

		adminSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, adminSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, adminSessionEncx, time.Hour)

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
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create admin user
		adminUser := td.NewTestUser(t, "admin@example.com", "Administrator", "User")
		adminUser.State = domain.Active
		adminUserEncx, err := domain.ProcessUserEncx(ctx, crypto, adminUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, adminUserEncx, testPool)
		require.NoError(t, err)

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

		adminSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, adminSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, adminSessionEncx, time.Hour)

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
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create target user with specific data
		email := "target@example.com"
		city := "Target City"
		telephone := "0612345678"
		firstname := "Target"
		lastname := "User"

		targetUser := td.NewTestUser(t, email, firstname, lastname)
		targetUser.State = domain.Active
		targetUser.City = city
		targetUser.Telephone = telephone
		targetUserEncx, err := domain.ProcessUserEncx(ctx, crypto, targetUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, targetUserEncx, testPool)
		require.NoError(t, err)

		// Create admin user
		adminUser := td.NewTestUser(t, "admin@example.com", "Administrator", "User")
		adminUser.State = domain.Active
		adminUser.Role = identity.Administrator.String()
		adminUserEncx, err := domain.ProcessUserEncx(ctx, crypto, adminUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, adminUserEncx, testPool)
		require.NoError(t, err)

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

		adminSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, adminSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, adminSessionEncx, time.Hour)

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
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create target user
		userID := setTargetUser()

		// Create admin user
		adminUser := td.NewTestUser(t, "admin@example.com", "Administrator", "User")
		adminUser.State = domain.Active
		adminUserEncx, err := domain.ProcessUserEncx(ctx, crypto, adminUser)
		require.NoError(t, err)
		td.InsertUserEncx(t, ctx, adminUserEncx, testPool)
		require.NoError(t, err)

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

		expiredAdministratorSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, expiredAdministratorSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, expiredAdministratorSessionEncx, -time.Hour) // Expired

		// Act - request with expired admin session
		req := td.NewGetUserByIDRequest(t, ctx, testServerURL, userID, expiredAdministratorSession.AccessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Should return 401 for expired admin session")
	})

	t.Run("should handle malformed access token", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		userID := setTargetUser()

		// Act - request with malformed access token
		req := td.NewGetUserByIDRequest(t, ctx, testServerURL, userID, "malformed-token")
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Should return 401 for malformed access token")
	})
}
