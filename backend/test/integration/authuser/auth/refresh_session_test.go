package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestRefreshSession TEST_PATH=test/integration/authuser/auth/refresh_session_test.go

func TestRefreshSession(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully refresh active session", func(t *testing.T) {
		// Clean state
		th.ClearSessionsRedis(t, ctx, redisClient)

		// Create a valid active session
		// activeSession := th.CreateTestSessionWithCrypto(t, crypto)
		activeSession, err := th.NewTestSession(t, crypto)
		require.NoError(t, err)
		activeSession.State = session.SessionActive
		activeSession.Role = identity.Standard

		// Re-process after state change
		activeSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, activeSession)
		require.NoError(t, err)

		// Insert session into Redis
		th.InsertSessionEncx(t, ctx, redisClient, activeSessionEncx, 24*time.Hour)

		// Make refresh request
		req := th.NewRefreshSessionRequest(t, ctx, testServerURL, activeSession.RefreshToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert successful response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		message, status := th.ParseRefreshSessionResponse(t, resp)
		assert.Equal(t, "Session refreshed successfully", message)
		assert.Equal(t, "success", status)

		// Verify new cookies are set
		cookies := resp.Cookies()
		var accessCookie, refreshCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == ck.AccessTokenCookieName {
				accessCookie = cookie
			}
			if cookie.Name == ck.RefreshTokenCookieName {
				refreshCookie = cookie
			}
		}

		require.NotNil(t, accessCookie, "Access token cookie should be set")
		require.NotNil(t, refreshCookie, "Refresh token cookie should be set")
		assert.NotEmpty(t, accessCookie.Value, "Access token cookie should have a value")
		assert.NotEmpty(t, refreshCookie.Value, "Refresh token cookie should have a value")
		assert.True(t, accessCookie.HttpOnly, "Access token cookie should be HTTP only")
		assert.True(t, refreshCookie.HttpOnly, "Refresh token cookie should be HTTP only")

		// Verify tokens are different from original
		assert.NotEqual(t, activeSessionEncx.AccessTokenHash, accessCookie.Value, "Access token cookie should have a different value than the original")
		assert.NotEqual(t, activeSessionEncx.RefreshTokenHash, refreshCookie.Value, "Refresh token cookie should have a different value than the original")

		// Verify session still exists with new tokens
		sessionExists, err := redisClient.Exists(ctx, session.FormatSessionKey(activeSession.ID.String())).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), sessionExists, "Session should still exists with new token")

		// Verify old tokens are invalidated and new ones exist
		oldAccessExists, err := redisClient.Exists(ctx, session.FormatAccessTokenKey(activeSessionEncx.AccessTokenHash)).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(0), oldAccessExists, "Old access token should be invalidated")

		oldRefreshExists, err := redisClient.Exists(ctx, session.FormatRefreshTokenKey(activeSessionEncx.RefreshTokenHash)).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(0), oldRefreshExists, "Old refresh token should be invalidated")

		accessCookieBytes, err := encx.SerializeValue(accessCookie.Value)
		require.NoError(t, err)
		accessCookieHash := crypto.HashBasic(ctx, accessCookieBytes)
		newAccessExists, err := redisClient.Exists(ctx, session.FormatAccessTokenKey(accessCookieHash)).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), newAccessExists, "New access token should be set in Redis")

		refreshCookieBytes, err := encx.SerializeValue(refreshCookie.Value)
		require.NoError(t, err)
		refreshCookieHash := crypto.HashBasic(ctx, refreshCookieBytes)
		newRefreshExists, err := redisClient.Exists(ctx, session.FormatRefreshTokenKey(refreshCookieHash)).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), newRefreshExists, "New refresh token should be set in Redis")
	})

	t.Run("should successfully refresh pending session", func(t *testing.T) {
		// Clean state
		th.ClearSessionsRedis(t, ctx, redisClient)

		// Create a valid pending session
		// pendingSession := th.CreateTestPendingSessionWithCrypto(t, crypto)
		pendingSession, err := th.NewTestSession(t, crypto)
		require.NoError(t, err)
		pendingSession.State = session.SessionPending

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		// Insert session into Redis
		th.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, 24*time.Hour)

		// Make refresh request
		req := th.NewRefreshSessionRequest(t, ctx, testServerURL, pendingSession.RefreshToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert successful response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		message, status := th.ParseRefreshSessionResponse(t, resp)
		assert.Equal(t, "Session refreshed successfully", message)
		assert.Equal(t, "success", status)

		// Verify new cookies are set
		cookies := resp.Cookies()
		var accessCookie, refreshCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == ck.AccessTokenCookieName {
				accessCookie = cookie
			}
			if cookie.Name == ck.RefreshTokenCookieName {
				refreshCookie = cookie
			}
		}

		require.NotNil(t, accessCookie, "Access token cookie should be set")
		require.NotNil(t, refreshCookie, "Refresh token cookie should be set")
		assert.NotEmpty(t, accessCookie.Value)
		assert.NotEmpty(t, refreshCookie.Value)
	})

	t.Run("should fail with missing session info in context", func(t *testing.T) {
		// Clean state
		th.ClearSessionsRedis(t, ctx, redisClient)

		// Make request without proper session context (no refresh token cookie)
		req := th.NewRefreshSessionRequestWithoutToken(t, ctx, testServerURL)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert error response
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		errorMsg, statusCode := th.ParseErrorResponse(t, resp)
		assert.Equal(t, http.StatusUnauthorized, statusCode)
		assert.Contains(t, errorMsg, errs.ErrUnauthorized.Error())
	})

	t.Run("should fail when session not found in Redis", func(t *testing.T) {
		// Clean state
		th.ClearSessionsRedis(t, ctx, redisClient)

		// Create session but don't insert into Redis
		nonExistentSession, err := th.NewTestSession(t, crypto)
		require.NoError(t, err)

		// Make refresh request for non-existent session
		req := th.NewRefreshSessionRequest(t, ctx, testServerURL, nonExistentSession.RefreshToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert not found response
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		errorMsg, statusCode := th.ParseErrorResponse(t, resp)
		assert.Equal(t, http.StatusUnauthorized, statusCode)
		assert.Contains(t, errorMsg, errs.ErrUnauthorized.Error())
	})

	t.Run("should fail for invalid session state", func(t *testing.T) {
		// Clean state
		th.ClearSessionsRedis(t, ctx, redisClient)

		// Create session with invalid state (neither active nor pending)
		invalidSession, err := th.NewTestSession(t, crypto)
		require.NoError(t, err)
		invalidSession.State = session.SessionState("invalid") // Invalid state for refresh

		// Re-process after state change
		invalidSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, invalidSession)
		require.NoError(t, err)

		// Insert session into Redis
		th.InsertSessionEncx(t, ctx, redisClient, invalidSessionEncx, 24*time.Hour)

		// Make refresh request
		req := th.NewRefreshSessionRequest(t, ctx, testServerURL, invalidSession.RefreshToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert unauthorized response
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		errorMsg, statusCode := th.ParseErrorResponse(t, resp)
		assert.Equal(t, http.StatusUnauthorized, statusCode)
		assert.Contains(t, errorMsg, "unauthorized")
	})

	t.Run("should handle expired tokens gracefully", func(t *testing.T) {
		// Clean state
		th.ClearSessionsRedis(t, ctx, redisClient)

		// Create session with expired timestamps
		expiredSession, err := th.NewTestSession(t, crypto)
		require.NoError(t, err)
		expiredSession.ExpiresAt = time.Now().Add(-1 * time.Hour) // Expired 1 hour ago

		// Re-process after timestamp change
		expiredSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, expiredSession)
		require.NoError(t, err)

		// Insert session into Redis with very short TTL
		th.InsertSessionEncx(t, ctx, redisClient, expiredSessionEncx, 1*time.Millisecond)

		// Wait for Redis to expire the keys
		time.Sleep(10 * time.Millisecond)

		// Make refresh request
		req := th.NewRefreshSessionRequest(t, ctx, testServerURL, expiredSession.RefreshToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should fail because session no longer exists in Redis
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should handle Redis connection failure gracefully", func(t *testing.T) {
		// This test would require mocking Redis client behavior
		// or using a more complex test setup with Redis container manipulation
		// For now, we document the expected behavior:
		// - Redis connection failures should return 503 Service Unavailable
		// - Error should be properly logged as infrastructure connection failure
		t.Skip("Redis connection failure testing requires advanced container manipulation")
	})

	t.Run("should handle database transaction failures", func(t *testing.T) {
		// Clean state
		th.ClearSessionsRedis(t, ctx, redisClient)

		// Create a valid session
		validSession, err := th.NewTestSession(t, crypto)
		require.NoError(t, err)

		validSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, validSession)
		require.NoError(t, err)

		th.InsertSessionEncx(t, ctx, redisClient, validSessionEncx, 24*time.Hour)

		// This test would require causing a Redis transaction failure
		// by manipulating the Redis instance or using mocks
		t.Skip("Transaction failure testing requires Redis manipulation techniques")
	})

	t.Run("should preserve session state after successful refresh", func(t *testing.T) {
		// Clean state
		th.ClearSessionsRedis(t, ctx, redisClient)

		role := identity.Premium
		state := session.SessionActive

		// Create session with specific user ID and role
		userID := uuid.New()
		specificSession, err := th.NewTestSession(t, crypto)
		require.NoError(t, err)
		specificSession.UserID = userID
		specificSession.Role = role
		specificSession.State = state

		// Re-process after changes
		specificSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, specificSession)
		require.NoError(t, err)

		th.InsertSessionEncx(t, ctx, redisClient, specificSessionEncx, 24*time.Hour)

		// Make refresh request
		req := th.NewRefreshSessionRequest(t, ctx, testServerURL, specificSession.RefreshToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert successful response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Retrieve session data from Redis to verify preservation
		sessionData, err := redisClient.Get(ctx, session.FormatSessionKey(specificSession.ID.String())).Bytes()
		require.NoError(t, err)

		var retrievedSessionEncx session.SessionEncx
		err = json.Unmarshal(sessionData, &retrievedSessionEncx)
		require.NoError(t, err, "Failed to unmarshal SessionEncx")

		retrievedSession, err := session.DecryptSessionEncx(context.Background(), crypto, &retrievedSessionEncx)
		require.NoError(t, err, "Failed to decrypt session")

		// Verify that non-token fields are preserved
		assert.Equal(t, userID, retrievedSession.UserID)
		assert.Equal(t, role, retrievedSession.Role)
		assert.Equal(t, state, retrievedSession.State)

		// Verify that tokens have been updated
		assert.NotEqual(t, specificSessionEncx.AccessTokenHash, retrievedSessionEncx.AccessTokenHash)
		assert.NotEqual(t, specificSessionEncx.RefreshTokenHash, retrievedSessionEncx.RefreshTokenHash)
	})

	t.Run("should handle malformed refresh token gracefully", func(t *testing.T) {
		// Clean state
		th.ClearSessionsRedis(t, ctx, redisClient)

		// Make request with malformed/invalid refresh token
		req := th.NewRefreshSessionRequest(t, ctx, testServerURL, "invalid-token-format")
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should fail due to middleware unable to resolve session
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
