package sessionRepository_test

import (
	"context"
	"testing"
	"time"

	td "github.com/Leviosa-care/authuser/test/helpers"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestRefreshTokenPair make test-unit-session-test

func TestRefreshTokenPair(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully refresh both access and refresh tokens", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create initial token pair
		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)
		oldAccessTokenHash := "old_access_token_hash"
		refreshTokenHash := "refresh_token_hash"
		accessTTL := 1 * time.Hour
		refreshTTL := 24 * time.Hour

		err := repo.CreateSession(ctx, session.ID, oldAccessTokenHash, refreshTokenHash, session.UserIDHash, sessionData, accessTTL, refreshTTL)
		require.NoError(t, err)

		// Refresh with new access and refresh tokens
		newAccessTokenHash := "new_access_token_hash"
		newRefreshTokenHash := "new_refresh_token_hash"
		
		// Update session with new token hashes for test
		session.AccessTokenHash = newAccessTokenHash
		session.RefreshTokenHash = newRefreshTokenHash
		updatedSessionData := td.EncodeSession(t, session)
		
		err = repo.RefreshTokenPair(ctx, refreshTokenHash, newAccessTokenHash, newRefreshTokenHash, session.ID, string(updatedSessionData), accessTTL, refreshTTL)
		require.NoError(t, err)

		// Old access token should be removed
		_, _, err = repo.FindSessionByAccessTokenHash(ctx, oldAccessTokenHash)
		assert.Error(t, err, "Old access token should be removed")

		// New access token should work
		_, sessionDataFound, err := repo.FindSessionByAccessTokenHash(ctx, newAccessTokenHash)
		require.NoError(t, err)
		assert.Equal(t, sessionData, sessionDataFound)

		// Old refresh token should be removed
		_, _, err = repo.FindSessionByRefreshTokenHash(ctx, refreshTokenHash)
		assert.Error(t, err, "Old refresh token should be removed")

		// New refresh token should work
		_, sessionDataFromNewRefresh, err := repo.FindSessionByRefreshTokenHash(ctx, newRefreshTokenHash)
		require.NoError(t, err)
		assert.Equal(t, sessionData, sessionDataFromNewRefresh)
	})

	t.Run("should clean up stale access tokens for the same session", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create initial token pair
		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)
		oldAccessTokenHash1 := "old_access_token_1"
		oldAccessTokenHash2 := "old_access_token_2"
		oldRefreshTokenHash := "old_refresh_token"
		accessTTL := 1 * time.Hour
		refreshTTL := 24 * time.Hour

		// Create multiple access tokens for the same session (simulating multiple devices)
		err := repo.CreateSession(ctx, session.ID, oldAccessTokenHash1, oldRefreshTokenHash, session.UserIDHash, sessionData, accessTTL, refreshTTL)
		require.NoError(t, err)

		// Manually create additional access token for same session
		accessTokenKey2 := auth.FormatAccessTokenKey(oldAccessTokenHash2)
		err = testClient.Set(ctx, accessTokenKey2, session.ID.String(), accessTTL).Err()
		require.NoError(t, err)

		// Verify both old access tokens work initially
		_, _, err = repo.FindSessionByAccessTokenHash(ctx, oldAccessTokenHash1)
		require.NoError(t, err)
		_, _, err = repo.FindSessionByAccessTokenHash(ctx, oldAccessTokenHash2)
		require.NoError(t, err)

		// Refresh tokens
		newAccessTokenHash := "new_access_token"
		newRefreshTokenHash := "new_refresh_token"
		
		// Update session with new token hashes for test
		session.AccessTokenHash = newAccessTokenHash
		session.RefreshTokenHash = newRefreshTokenHash
		updatedSessionData := td.EncodeSession(t, session)
		
		err = repo.RefreshTokenPair(ctx, oldRefreshTokenHash, newAccessTokenHash, newRefreshTokenHash, session.ID, string(updatedSessionData), accessTTL, refreshTTL)
		require.NoError(t, err)

		// All old access tokens should be cleaned up
		_, _, err = repo.FindSessionByAccessTokenHash(ctx, oldAccessTokenHash1)
		assert.Error(t, err, "Old access token 1 should be removed")
		_, _, err = repo.FindSessionByAccessTokenHash(ctx, oldAccessTokenHash2)
		assert.Error(t, err, "Old access token 2 should be removed")

		// Old refresh token should be removed
		_, _, err = repo.FindSessionByRefreshTokenHash(ctx, oldRefreshTokenHash)
		assert.Error(t, err, "Old refresh token should be removed")

		// New tokens should work
		_, sessionDataFromNewAccess, err := repo.FindSessionByAccessTokenHash(ctx, newAccessTokenHash)
		require.NoError(t, err)
		assert.Equal(t, sessionData, sessionDataFromNewAccess)

		_, sessionDataFromNewRefresh, err := repo.FindSessionByRefreshTokenHash(ctx, newRefreshTokenHash)
		require.NoError(t, err)
		assert.Equal(t, sessionData, sessionDataFromNewRefresh)
	})

	t.Run("should rollback on refresh token creation failure", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create initial token pair
		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)
		oldAccessTokenHash := "rollback_test_old_access"
		oldRefreshTokenHash := "rollback_test_old_refresh"

		err := repo.CreateSession(ctx, session.ID, oldAccessTokenHash, oldRefreshTokenHash, session.UserIDHash, sessionData, time.Hour, 24*time.Hour)
		require.NoError(t, err)

		// Verify initial tokens work
		_, _, err = repo.FindSessionByAccessTokenHash(ctx, oldAccessTokenHash)
		require.NoError(t, err)

		// Close Redis to simulate failure
		testClient.Close()

		// Try to refresh tokens - should fail
		newAccessTokenHash := "rollback_test_new_access"
		newRefreshTokenHash := "rollback_test_new_refresh"
		
		// Update session with new token hashes for test
		session.AccessTokenHash = newAccessTokenHash
		session.RefreshTokenHash = newRefreshTokenHash
		updatedSessionData := td.EncodeSession(t, session)
		
		err = repo.RefreshTokenPair(ctx, oldRefreshTokenHash, newAccessTokenHash, newRefreshTokenHash, session.ID, string(updatedSessionData), time.Hour, 24*time.Hour)
		assert.Error(t, err)

		// Reconnect and verify state
		reconnectRedis()

		// Old tokens should still work (no partial updates)
		_, _, err = repo.FindSessionByAccessTokenHash(ctx, oldAccessTokenHash)
		require.NoError(t, err, "Old access token should still work after failed refresh")

		// New tokens should not exist
		_, _, err = repo.FindSessionByAccessTokenHash(ctx, newAccessTokenHash)
		assert.Error(t, err, "New access token should not exist after failed refresh")
	})

	t.Run("should handle TTL updates correctly", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create initial token pair with short TTL
		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)
		oldAccessTokenHash := "ttl_test_old_access"
		refreshTokenHash := "ttl_test_refresh"
		shortTTL := 1 * time.Second
		longTTL := 1 * time.Hour

		err := repo.CreateSession(ctx, session.ID, oldAccessTokenHash, refreshTokenHash, session.UserIDHash, sessionData, shortTTL, longTTL)
		require.NoError(t, err)

		// Refresh with longer TTL
		newAccessTokenHash := "ttl_test_new_access"
		newRefreshTokenHash := "ttl_test_new_refresh"
		newTTL := 2 * time.Hour
		
		// Update session with new token hashes for test
		session.AccessTokenHash = newAccessTokenHash
		session.RefreshTokenHash = newRefreshTokenHash
		updatedSessionData := td.EncodeSession(t, session)
		
		err = repo.RefreshTokenPair(ctx, refreshTokenHash, newAccessTokenHash, newRefreshTokenHash, session.ID, string(updatedSessionData), newTTL, longTTL)
		require.NoError(t, err)

		// Verify new access token has longer TTL
		newAccessTokenKey := auth.FormatAccessTokenKey(newAccessTokenHash)
		actualTTL := testClient.TTL(ctx, newAccessTokenKey).Val()
		assert.True(t, actualTTL > 1*time.Hour, "New access token should have longer TTL")
	})

	t.Run("should handle special characters in token hashes", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create tokens with special characters
		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)
		oldAccessTokenHash := "old-access_token.with:special@chars+123/456="
		oldRefreshTokenHash := "old-refresh_token.with:special@chars+789/012="

		err := repo.CreateSession(ctx, session.ID, oldAccessTokenHash, oldRefreshTokenHash, session.UserIDHash, sessionData, time.Hour, 24*time.Hour)
		require.NoError(t, err)

		// Refresh with new tokens also having special characters
		newAccessTokenHash := "new-access_token.with:special@chars+999/888="
		newRefreshTokenHash := "new-refresh_token.with:special@chars+777/666="

		// Update session with new token hashes for test
		session.AccessTokenHash = newAccessTokenHash
		session.RefreshTokenHash = newRefreshTokenHash
		updatedSessionData := td.EncodeSession(t, session)

		err = repo.RefreshTokenPair(ctx, oldRefreshTokenHash, newAccessTokenHash, newRefreshTokenHash, session.ID, string(updatedSessionData), time.Hour, 24*time.Hour)
		require.NoError(t, err)

		// Verify new tokens work
		_, _, err = repo.FindSessionByAccessTokenHash(ctx, newAccessTokenHash)
		require.NoError(t, err)

		_, _, err = repo.FindSessionByRefreshTokenHash(ctx, newRefreshTokenHash)
		require.NoError(t, err)
	})

	t.Run("should handle Redis connection errors gracefully", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Close Redis connection
		testClient.Close()

		// Try to refresh tokens with closed connection
		sessionID := uuid.New()
		// For this error case, we just need some valid session data
		updatedSessionData := "{\"id\":\"test\",\"access_token_hash\":\"new_access\",\"refresh_token_hash\":\"new_refresh\"}"
		err := repo.RefreshTokenPair(ctx, "old_refresh", "new_access", "new_refresh", sessionID, updatedSessionData, time.Hour, 24*time.Hour)
		assert.Error(t, err)

		// Reconnect for subsequent tests
		reconnectRedis()
	})
}

// TEST=TestInvalidateTokenPair make test-unit-session-test

func TestInvalidateTokenPair(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully invalidate all token pair keys", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create token pair
		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)
		accessTokenHash := "invalidate_test_access_token"
		refreshTokenHash := "invalidate_test_refresh_token"

		err := repo.CreateSession(ctx, session.ID, accessTokenHash, refreshTokenHash, session.UserIDHash, sessionData, time.Hour, 24*time.Hour)
		require.NoError(t, err)

		// Verify tokens work initially
		_, _, err = repo.FindSessionByAccessTokenHash(ctx, accessTokenHash)
		require.NoError(t, err)

		_, _, err = repo.FindSessionByRefreshTokenHash(ctx, refreshTokenHash)
		require.NoError(t, err)

		// Invalidate token pair
		err = repo.InvalidateTokenPair(ctx, accessTokenHash, refreshTokenHash, session.ID)
		require.NoError(t, err)

		// All keys should be removed
		_, _, err = repo.FindSessionByAccessTokenHash(ctx, accessTokenHash)
		assert.Error(t, err, "Access token should be invalidated")

		_, _, err = repo.FindSessionByRefreshTokenHash(ctx, refreshTokenHash)
		assert.Error(t, err, "Refresh token should be invalidated")

		// Session data should also be removed
		sessionKey := auth.FormatSessionKey(session.ID.String())
		_, err = testClient.Get(ctx, sessionKey).Result()
		assert.Error(t, err, "Session data should be invalidated")
	})

	t.Run("should succeed even if some keys don't exist", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Try to invalidate non-existent tokens
		sessionID := uuid.New()
		accessTokenHash := "non_existent_access_token"
		refreshTokenHash := "non_existent_refresh_token"

		// Should not error even if keys don't exist
		err := repo.InvalidateTokenPair(ctx, accessTokenHash, refreshTokenHash, sessionID)
		require.NoError(t, err)
	})

	t.Run("should handle partial token pair states", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create only session data and access token (no refresh token)
		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)
		accessTokenHash := "partial_state_access_token"
		refreshTokenHash := "partial_state_refresh_token"

		// Manually create partial state
		sessionKey := auth.FormatSessionKey(session.ID.String())
		accessTokenKey := auth.FormatAccessTokenKey(accessTokenHash)

		err := testClient.Set(ctx, sessionKey, sessionData, time.Hour).Err()
		require.NoError(t, err)

		err = testClient.Set(ctx, accessTokenKey, session.ID.String(), time.Hour).Err()
		require.NoError(t, err)
		// Note: not creating refresh token mapping

		// Invalidate should work and not fail on missing refresh token
		err = repo.InvalidateTokenPair(ctx, accessTokenHash, refreshTokenHash, session.ID)
		require.NoError(t, err)

		// Verify existing keys are removed
		_, err = testClient.Get(ctx, sessionKey).Result()
		assert.Error(t, err, "Session data should be removed")

		_, err = testClient.Get(ctx, accessTokenKey).Result()
		assert.Error(t, err, "Access token should be removed")
	})

	t.Run("should handle special characters in token hashes", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create token pair with special characters
		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)
		accessTokenHash := "special-access_token.with:chars@123+456/789="
		refreshTokenHash := "special-refresh_token.with:chars@987+654/321="

		err := repo.CreateSession(ctx, session.ID, accessTokenHash, refreshTokenHash, session.UserIDHash, sessionData, time.Hour, 24*time.Hour)
		require.NoError(t, err)

		// Invalidate tokens with special characters
		err = repo.InvalidateTokenPair(ctx, accessTokenHash, refreshTokenHash, session.ID)
		require.NoError(t, err)

		// Verify all keys are removed
		_, _, err = repo.FindSessionByAccessTokenHash(ctx, accessTokenHash)
		assert.Error(t, err, "Access token with special chars should be invalidated")

		_, _, err = repo.FindSessionByRefreshTokenHash(ctx, refreshTokenHash)
		assert.Error(t, err, "Refresh token with special chars should be invalidated")
	})

	t.Run("should not affect other sessions", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create two separate token pairs
		session1 := td.CreateTestSessionWithCrypto(t, crypto)
		session2 := td.CreateTestSessionWithCrypto(t, crypto)

		sessionData1 := td.EncodeSession(t, session1)
		sessionData2 := td.EncodeSession(t, session2)

		// First token pair
		err := repo.CreateSession(ctx, session1.ID, "access_1", "refresh_1", session1.UserIDHash, sessionData1, time.Hour, 24*time.Hour)
		require.NoError(t, err)

		// Second token pair
		err = repo.CreateSession(ctx, session2.ID, "access_2", "refresh_2", session2.UserIDHash, sessionData2, time.Hour, 24*time.Hour)
		require.NoError(t, err)

		// Invalidate only first token pair
		err = repo.InvalidateTokenPair(ctx, "access_1", "refresh_1", session1.ID)
		require.NoError(t, err)

		// First session should be invalidated
		_, _, err = repo.FindSessionByAccessTokenHash(ctx, "access_1")
		assert.Error(t, err)

		// Second session should still work
		_, _, err = repo.FindSessionByAccessTokenHash(ctx, "access_2")
		require.NoError(t, err)

		_, _, err = repo.FindSessionByRefreshTokenHash(ctx, "refresh_2")
		require.NoError(t, err)
	})

	t.Run("should handle Redis connection errors gracefully", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Close Redis connection
		testClient.Close()

		// Try to invalidate tokens with closed connection
		sessionID := uuid.New()
		err := repo.InvalidateTokenPair(ctx, "access_token", "refresh_token", sessionID)
		assert.Error(t, err)

		// Reconnect for subsequent tests
		reconnectRedis()
	})
}

