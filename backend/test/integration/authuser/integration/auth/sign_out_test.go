package auth_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestSignOut make test-integration-auth-test

func TestSignOut(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully sign out with valid session", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create active user
		user := td.NewTestUser(t, "signout-test@example.com", "Sign", "Out")
		user.State = domain.Active

		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Create active session
		sessionInfo := &session.SessionInfo{
			UserID: user.ID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

		// Make sign-out request
		req := td.NewSignOutRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		var response struct {
			Message string `json:"message"`
			Status  string `json:"status"`
		}
		td.ParseJSONResponse(t, resp, &response)

		// Validate response
		assert.Equal(t, "Successfully signed out", response.Message)
		assert.Equal(t, "signed_out", response.Status)

		// Verify session was removed from Redis
		session := td.GetSessionByID(t, ctx, sessionInfo.ID, redisClient)
		assert.Nil(t, session, "Session should be removed after sign-out")
	})

	t.Run("should fail when no authentication provided", func(t *testing.T) {
		// Make sign-out request without authentication
		req := td.NewSignOutRequestWithoutAuth(t, ctx, testServerURL)
		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should fail with invalid session token", func(t *testing.T) {
		// Make sign-out request with invalid token
		req := td.NewSignOutRequest(t, ctx, testServerURL, "invalid-token")
		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should handle double sign-out gracefully", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create active user
		user := td.NewTestUser(t, "double-signout@example.com", "Double", "SignOut")
		user.State = domain.Active

		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Create active session
		sessionInfo := &session.SessionInfo{
			UserID: user.ID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

		// First sign-out request
		req1 := td.NewSignOutRequest(t, ctx, testServerURL, accessToken)
		resp1, err := client.Do(req1)

		// Assert first sign-out succeeds
		require.NoError(t, err)
		defer resp1.Body.Close()
		assert.Equal(t, http.StatusOK, resp1.StatusCode)

		// Second sign-out request with same token (should fail due to invalid session)
		req2 := td.NewSignOutRequest(t, ctx, testServerURL, accessToken)
		resp2, err := client.Do(req2)

		// Assert second sign-out fails with unauthorized (session no longer exists)
		require.NoError(t, err)
		defer resp2.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp2.StatusCode)
	})

	t.Run("should fail with guest role (below minimum required role)", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create active user
		user := td.NewTestUser(t, "guest-signout@example.com", "Guest", "User")
		user.State = domain.Active

		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Create guest session (below Standard minimum role)
		sessionInfo := &session.SessionInfo{
			UserID: user.ID,
			Role:   identity.Guest,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

		// Make sign-out request with guest role
		req := td.NewSignOutRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert HTTP response (should be forbidden due to insufficient role)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should work with administrator role", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create active user
		user := td.NewTestUser(t, "admin-signout@example.com", "Admin", "User")
		user.State = domain.Active

		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Create administrator session
		sessionInfo := &session.SessionInfo{
			UserID: user.ID,
			Role:   identity.Administrator,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

		// Make sign-out request with administrator role
		req := td.NewSignOutRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify session was removed from Redis
		session := td.GetSessionByID(t, ctx, sessionInfo.ID, redisClient)
		assert.Nil(t, session, "Session should be removed after sign-out")
	})
}
