package sessionRepository_test

import (
	"context"
	"testing"
	"time"

	td "github.com/Leviosa-care/authuser/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestCreateTokenPair make test-unit-session-test

func TestCreateTokenPair(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully create token pair with all mappings", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create test session
		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)
		accessTokenHash := "access_token_hash_123"
		refreshTokenHash := "refresh_token_hash_456"
		accessTTL := 1 * time.Hour
		refreshTTL := 24 * time.Hour

		// Create token pair
		err := repo.CreateTokenPair(ctx, session.ID, accessTokenHash, refreshTokenHash, sessionData, accessTTL, refreshTTL)
		require.NoError(t, err)

		// Verify session data is stored
		sessionKey := session.FormatSessionKey(session.ID.String())
		storedSessionData, err := testClient.Get(ctx, sessionKey).Result()
		require.NoError(t, err)
		assert.Equal(t, string(sessionData), storedSessionData)

		// Verify session data TTL matches refresh token TTL (longer duration)
		sessionTTL := testClient.TTL(ctx, sessionKey).Val()
		assert.True(t, sessionTTL > 23*time.Hour && sessionTTL <= refreshTTL, "Session TTL should be close to refresh TTL")

		// Verify access token mapping
		accessTokenKey := session.FormatAccessTokenKey(accessTokenHash)
		storedSessionID, err := testClient.Get(ctx, accessTokenKey).Result()
		require.NoError(t, err)
		assert.Equal(t, session.ID.String(), storedSessionID)

		// Verify access token TTL
		accessTTLStored := testClient.TTL(ctx, accessTokenKey).Val()
		assert.True(t, accessTTLStored > 50*time.Minute && accessTTLStored <= accessTTL, "Access token TTL should be close to specified TTL")

		// Verify refresh token mapping
		refreshTokenKey := session.FormatRefreshTokenKey(refreshTokenHash)
		storedSessionIDFromRefresh, err := testClient.Get(ctx, refreshTokenKey).Result()
		require.NoError(t, err)
		assert.Equal(t, session.ID.String(), storedSessionIDFromRefresh)

		// Verify refresh token TTL
		refreshTTLStored := testClient.TTL(ctx, refreshTokenKey).Val()
		assert.True(t, refreshTTLStored > 23*time.Hour && refreshTTLStored <= refreshTTL, "Refresh token TTL should be close to specified TTL")
	})

	t.Run("should handle duplicate session ID creation", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create first token pair
		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)

		firstAccessToken := "first_access_token"
		firstRefreshToken := "first_refresh_token"
		err := repo.CreateTokenPair(ctx, session.ID, firstAccessToken, firstRefreshToken, sessionData, time.Hour, 24*time.Hour)
		require.NoError(t, err)

		// Try to create another token pair with the same session ID but different tokens
		// This should overwrite the session data but create new token mappings
		newSessionData := td.EncodeSession(t, session)
		secondAccessToken := "second_access_token"
		secondRefreshToken := "second_refresh_token"

		err = repo.CreateTokenPair(ctx, session.ID, secondAccessToken, secondRefreshToken, newSessionData, time.Hour, 24*time.Hour)
		require.NoError(t, err)

		// Both token pairs should work
		// TODO: change that to use raw redis query and not use some repository function other than the one tested
		_, sessionDataFromFirst, err := repo.FindSessionByAccessTokenHash(ctx, firstAccessToken)
		require.NoError(t, err)
		assert.NotNil(t, sessionDataFromFirst)

		// TODO: change that to use raw redis query and not use some repository function other than the one tested
		_, sessionDataFromSecond, err := repo.FindSessionByAccessTokenHash(ctx, secondAccessToken)
		require.NoError(t, err)
		assert.NotNil(t, sessionDataFromSecond)
	})

	t.Run("should rollback access token on refresh token creation failure", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)
		accessTokenHash := "access_token_rollback_test"
		refreshTokenHash := "refresh_token_rollback_test"

		// First, create the token pair successfully
		err := repo.CreateTokenPair(ctx, session.ID, accessTokenHash, refreshTokenHash, sessionData, time.Hour, 24*time.Hour)
		require.NoError(t, err)

		// Close Redis to simulate failure during refresh token creation
		testClient.Close()

		// Try to create another token pair - this should fail
		newAccessToken := "new_access_token"
		newRefreshToken := "new_refresh_token"
		err = repo.CreateTokenPair(ctx, session.ID, newAccessToken, newRefreshToken, sessionData, time.Hour, 24*time.Hour)
		assert.Error(t, err)

		// Reconnect and verify no partial state exists for the failed operation
		reconnectRedis()

		// The failed tokens should not exist
		_, _, err = repo.FindSessionByAccessTokenHash(ctx, newAccessToken)
		assert.Error(t, err, "New access token should not exist after rollback")

		_, _, err = repo.FindSessionByRefreshTokenHash(ctx, newRefreshToken)
		assert.Error(t, err, "New refresh token should not exist after rollback")

		// Original tokens should still work
		_, _, err = repo.FindSessionByAccessTokenHash(ctx, accessTokenHash)
		require.NoError(t, err, "Original tokens should still work")
	})

	t.Run("should handle zero TTL values", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)
		accessTokenHash := "zero_ttl_access_token"
		refreshTokenHash := "zero_ttl_refresh_token"

		// Create token pair with zero TTL (should use Redis default behavior)
		err := repo.CreateTokenPair(ctx, session.ID, accessTokenHash, refreshTokenHash, sessionData, 0, 0)
		require.NoError(t, err)

		// Tokens should exist and work
		_, sessionDataFound, err := repo.FindSessionByAccessTokenHash(ctx, accessTokenHash)
		require.NoError(t, err)
		assert.Equal(t, sessionData, sessionDataFound)

		// Check that keys exist without expiration (TTL = -1 means no expiration)
		accessTokenKey := session.FormatAccessTokenKey(accessTokenHash)
		accessTTL := testClient.TTL(ctx, accessTokenKey).Val()
		assert.Equal(t, time.Duration(-1), accessTTL, "Zero TTL should result in no expiration")
	})

	t.Run("should handle very short TTL values", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)
		accessTokenHash := "short_ttl_access_token"
		refreshTokenHash := "short_ttl_refresh_token"
		shortTTL := 50 * time.Millisecond

		// Create token pair with very short TTL
		err := repo.CreateTokenPair(ctx, session.ID, accessTokenHash, refreshTokenHash, sessionData, shortTTL, shortTTL)
		require.NoError(t, err)

		// Tokens should initially work
		_, _, err = repo.FindSessionByAccessTokenHash(ctx, accessTokenHash)
		require.NoError(t, err)

		// Wait for expiration
		time.Sleep(60 * time.Millisecond)

		// Tokens should be expired
		_, _, err = repo.FindSessionByAccessTokenHash(ctx, accessTokenHash)
		assert.Error(t, err, "Access token should be expired")

		_, _, err = repo.FindSessionByRefreshTokenHash(ctx, refreshTokenHash)
		assert.Error(t, err, "Refresh token should be expired")
	})

	t.Run("should handle special characters in token hashes", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)
		// Use various special characters that might appear in base64/hex encoded hashes
		accessTokenHash := "access-token_hash.with:special@chars+123/456="
		refreshTokenHash := "refresh-token_hash.with:special@chars+789/012="

		// Create token pair with special characters
		err := repo.CreateTokenPair(ctx, session.ID, accessTokenHash, refreshTokenHash, sessionData, time.Hour, 24*time.Hour)
		require.NoError(t, err)

		// Verify both tokens work
		_, sessionDataFound, err := repo.FindSessionByAccessTokenHash(ctx, accessTokenHash)
		require.NoError(t, err)
		assert.Equal(t, sessionData, sessionDataFound)

		_, sessionDataFound, err = repo.FindSessionByRefreshTokenHash(ctx, refreshTokenHash)
		require.NoError(t, err)
		assert.Equal(t, sessionData, sessionDataFound)
	})

	t.Run("should handle large session data", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create session with large encrypted data
		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)

		// Add extra data to make it larger
		largeSessionData := make([]byte, len(sessionData)+10000)
		copy(largeSessionData, sessionData)
		for i := len(sessionData); i < len(largeSessionData); i++ {
			largeSessionData[i] = byte(i % 256)
		}

		accessTokenHash := "large_data_access_token"
		refreshTokenHash := "large_data_refresh_token"

		// Create token pair with large session data
		err := repo.CreateTokenPair(ctx, session.ID, accessTokenHash, refreshTokenHash, largeSessionData, time.Hour, 24*time.Hour)
		require.NoError(t, err)

		// Verify data integrity
		_, retrievedData, err := repo.FindSessionByAccessTokenHash(ctx, accessTokenHash)
		require.NoError(t, err)
		assert.Equal(t, largeSessionData, retrievedData)
	})

	t.Run("should handle Redis connection errors gracefully", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Close Redis connection
		testClient.Close()

		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)
		accessTokenHash := "connection_error_access_token"
		refreshTokenHash := "connection_error_refresh_token"

		// Try to create token pair with closed connection
		err := repo.CreateTokenPair(ctx, session.ID, accessTokenHash, refreshTokenHash, sessionData, time.Hour, 24*time.Hour)
		assert.Error(t, err)

		// Reconnect for subsequent tests
		reconnectRedis()
	})
}

