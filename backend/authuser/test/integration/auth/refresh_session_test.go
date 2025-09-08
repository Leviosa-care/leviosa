package auth_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/test/helpers"
	"github.com/hengadev/leviosa/core/errs"

	ck "github.com/Leviosa-care/core/auth/cookies"
	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestRefreshSession make test-integration-auth-test

func TestRefreshSession(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully refresh active session", func(t *testing.T) {
		// Clean state
		helpers.ClearSessionsRedis(t, ctx, testClient)

		// Create a valid active session
		activeSession := helpers.CreateTestSessionWithCrypto(t, crypto)
		activeSession.State = session.SessionActive
		activeSession.Role = identity.Standard

		// Re-process after state change
		err := crypto.ProcessStruct(ctx, activeSession)
		require.NoError(t, err)

		// Insert session into Redis
		helpers.InsertSessionDirectly(t, ctx, testClient, activeSession, 24*time.Hour)

		// Make refresh request
		req := helpers.NewRefreshSessionRequest(t, ctx, testServerURL, activeSession.RefreshToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert successful response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		message, status := helpers.ParseRefreshSessionResponse(t, resp)
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
		assert.NotEqual(t, activeSession.AccessTokenHash, accessCookie.Value, "Access token cookie should have a different value than the original")
		assert.NotEqual(t, activeSession.RefreshTokenHash, refreshCookie.Value, "Refresh token cookie should have a different value than the original")

		// Verify session still exists with new tokens
		sessionExists, err := testClient.Exists(ctx, session.FormatSessionKey(activeSession.ID.String())).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), sessionExists, "Session should still exists with new token")

		// Verify old tokens are invalidated and new ones exist
		oldAccessExists, err := testClient.Exists(ctx, session.FormatAccessTokenKey(activeSession.AccessTokenHash)).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(0), oldAccessExists, "Old access token should be invalidated")

		oldRefreshExists, err := testClient.Exists(ctx, session.FormatRefreshTokenKey(activeSession.RefreshTokenHash)).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(0), oldRefreshExists, "Old refresh token should be invalidated")

		newAccessExists, err := testClient.Exists(ctx, session.FormatAccessTokenKey(accessCookie.Value)).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), newAccessExists, "New access token should be set in Redis")

		newRefreshExists, err := testClient.Exists(ctx, session.FormatRefreshTokenKey(refreshCookie.Value)).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), newRefreshExists, "New refresh token should be set in Redis")
	})

	t.Run("should successfully refresh pending session", func(t *testing.T) {
		// Clean state
		helpers.ClearSessionsRedis(t, ctx, testClient)

		// Create a valid pending session
		pendingSession := helpers.CreateTestPendingSessionWithCrypto(t, crypto)

		// Insert session into Redis
		helpers.InsertSessionDirectly(t, ctx, testClient, pendingSession, 24*time.Hour)

		// Make refresh request
		req := helpers.NewRefreshSessionRequest(t, ctx, testServerURL, pendingSession.RefreshToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert successful response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		message, status := helpers.ParseRefreshSessionResponse(t, resp)
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
		helpers.ClearSessionsRedis(t, ctx, testClient)

		// Make request without proper session context (no refresh token cookie)
		req := helpers.NewRefreshSessionRequestWithoutToken(t, ctx, testServerURL)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert error response
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		errorMsg, statusCode := helpers.ParseErrorResponse(t, resp)
		assert.Equal(t, http.StatusUnauthorized, statusCode)
		assert.Contains(t, errorMsg, errs.ErrUnauthorized.Error())
	})

	t.Run("should fail when session not found in Redis", func(t *testing.T) {
		// Clean state
		helpers.ClearSessionsRedis(t, ctx, testClient)

		// Create session but don't insert into Redis
		nonExistentSession := helpers.CreateTestSessionWithCrypto(t, crypto)

		// Make refresh request for non-existent session
		req := helpers.NewRefreshSessionRequest(t, ctx, testServerURL, nonExistentSession.RefreshToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert not found response
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		errorMsg, statusCode := helpers.ParseErrorResponse(t, resp)
		assert.Equal(t, http.StatusUnauthorized, statusCode)
		assert.Contains(t, errorMsg, errs.ErrUnauthorized.Error())
	})

	t.Run("should fail for invalid session state", func(t *testing.T) {
		// Clean state
		helpers.ClearSessionsRedis(t, ctx, testClient)

		// Create session with invalid state (neither active nor pending)
		invalidSession := helpers.CreateTestSessionWithCrypto(t, crypto)
		invalidSession.State = session.SessionState("invalid") // Invalid state for refresh

		// Re-process after state change
		err := crypto.ProcessStruct(ctx, invalidSession)
		require.NoError(t, err)

		// Insert session into Redis
		helpers.InsertSessionDirectly(t, ctx, testClient, invalidSession, 24*time.Hour)

		// Make refresh request
		req := helpers.NewRefreshSessionRequest(t, ctx, testServerURL, invalidSession.RefreshToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert unauthorized response
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		errorMsg, statusCode := helpers.ParseErrorResponse(t, resp)
		assert.Equal(t, http.StatusUnauthorized, statusCode)
		assert.Contains(t, errorMsg, "unauthorized")
	})

	t.Run("should handle expired tokens gracefully", func(t *testing.T) {
		// Clean state
		helpers.ClearSessionsRedis(t, ctx, testClient)

		// Create session with expired timestamps
		expiredSession := helpers.CreateTestSessionWithCrypto(t, crypto)
		expiredSession.ExpiresAt = time.Now().Add(-1 * time.Hour) // Expired 1 hour ago

		// Re-process after timestamp change
		err := crypto.ProcessStruct(ctx, expiredSession)
		require.NoError(t, err)

		// Insert session into Redis with very short TTL
		helpers.InsertSessionDirectly(t, ctx, testClient, expiredSession, 1*time.Millisecond)

		// Wait for Redis to expire the keys
		time.Sleep(10 * time.Millisecond)

		// Make refresh request
		req := helpers.NewRefreshSessionRequest(t, ctx, testServerURL, expiredSession.RefreshToken)
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
		helpers.ClearSessionsRedis(t, ctx, testClient)

		// Create a valid session
		validSession := helpers.CreateTestSessionWithCrypto(t, crypto)
		helpers.InsertSessionDirectly(t, ctx, testClient, validSession, 24*time.Hour)

		// This test would require causing a Redis transaction failure
		// by manipulating the Redis instance or using mocks
		t.Skip("Transaction failure testing requires Redis manipulation techniques")
	})

	t.Run("should preserve session state after successful refresh", func(t *testing.T) {
		// Clean state
		helpers.ClearSessionsRedis(t, ctx, testClient)

		role := identity.Premium
		state := session.SessionActive

		// Create session with specific user ID and role
		userID := uuid.New()
		specificSession := helpers.CreateTestSessionWithUserIDAndCrypto(t, userID, crypto)
		specificSession.Role = role
		specificSession.State = state

		// Re-process after changes
		err := crypto.ProcessStruct(ctx, specificSession)
		require.NoError(t, err)

		helpers.InsertSessionDirectly(t, ctx, testClient, specificSession, 24*time.Hour)

		// Make refresh request
		req := helpers.NewRefreshSessionRequest(t, ctx, testServerURL, specificSession.RefreshToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert successful response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Retrieve session data from Redis to verify preservation
		sessionData, err := testClient.Get(ctx, session.FormatSessionKey(specificSession.ID.String())).Bytes()
		require.NoError(t, err)

		retrievedSession := helpers.DecodeSessionWithDecryption(t, sessionData, crypto)

		// Verify that non-token fields are preserved
		assert.Equal(t, userID, retrievedSession.UserID)
		assert.Equal(t, role, retrievedSession.Role)
		assert.Equal(t, state, retrievedSession.State)

		// Verify that tokens have been updated
		assert.NotEqual(t, specificSession.AccessTokenHash, retrievedSession.AccessTokenHash)
		assert.NotEqual(t, specificSession.RefreshTokenHash, retrievedSession.RefreshTokenHash)
	})

	t.Run("should handle malformed refresh token gracefully", func(t *testing.T) {
		// Clean state
		helpers.ClearSessionsRedis(t, ctx, testClient)

		// Make request with malformed/invalid refresh token
		req := helpers.NewRefreshSessionRequest(t, ctx, testServerURL, "invalid-token-format")
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should fail due to middleware unable to resolve session
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
