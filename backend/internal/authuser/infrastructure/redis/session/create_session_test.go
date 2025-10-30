package sessionRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreateSession TEST_PATH=internal/authuser/infrastructure/redis/session/create_session_test.go

func TestCreateSession(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully create token pair with all mappings", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create test session
		baseSession := td.NewTestSessionEncx(t)
		sessionData := td.EncodeSessionEncx(t, baseSession)
		accessTokenHash := "access_token_hash_123"
		baseSession.AccessTokenHash = accessTokenHash
		refreshTokenHash := "refresh_token_hash_456"
		baseSession.RefreshTokenHash = accessTokenHash
		accessTTL := 1 * time.Hour
		refreshTTL := 24 * time.Hour

		// Create session
		err := repo.CreateSession(ctx, baseSession.ID, accessTokenHash, refreshTokenHash, baseSession.UserIDHash, sessionData, accessTTL, refreshTTL)
		assert.NoError(t, err)

		// Verify session data is stored
		sessionKey := session.FormatSessionKey(baseSession.ID.String())
		storedSessionData, err := testClient.Get(ctx, sessionKey).Result()
		require.NoError(t, err)
		assert.Equal(t, string(sessionData), storedSessionData)

		// Verify session data TTL matches refresh token TTL (longer duration)
		sessionTTL := testClient.TTL(ctx, sessionKey).Val()
		assert.True(t, sessionTTL > 23*time.Hour && sessionTTL <= refreshTTL, "Session TTL should be close to refresh TTL")

		// Verify access token mapping
		accessTokenKey := session.FormatAccessTokenKey(accessTokenHash)
		storedSessionID, err := testClient.Get(ctx, accessTokenKey).Result()
		assert.NoError(t, err)
		assert.Equal(t, baseSession.ID.String(), storedSessionID)

		// Verify access token TTL
		accessTTLStored := testClient.TTL(ctx, accessTokenKey).Val()
		assert.True(t, accessTTLStored > 50*time.Minute && accessTTLStored <= accessTTL, "Access token TTL should be close to specified TTL")

		// Verify refresh token mapping
		refreshTokenKey := session.FormatRefreshTokenKey(refreshTokenHash)
		storedSessionIDFromRefresh, err := testClient.Get(ctx, refreshTokenKey).Result()
		require.NoError(t, err)
		assert.Equal(t, baseSession.ID.String(), storedSessionIDFromRefresh)

		// Verify refresh token TTL
		refreshTTLStored := testClient.TTL(ctx, refreshTokenKey).Val()
		assert.True(t, refreshTTLStored > 23*time.Hour && refreshTTLStored <= refreshTTL, "Refresh token TTL should be close to specified TTL")
	})

	t.Run("should handle duplicate session ID creation", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create first token pair
		baseSession := td.NewTestSessionEncx(t)

		firstAccessToken := "first_access_token"
		baseSession.AccessTokenHash = firstAccessToken
		firstRefreshToken := "first_refresh_token"
		baseSession.RefreshTokenHash = firstRefreshToken

		sessionData := td.EncodeSessionEncx(t, baseSession)

		err := repo.CreateSession(ctx, baseSession.ID, firstAccessToken, firstRefreshToken, baseSession.UserIDHash, sessionData, time.Hour, 24*time.Hour)
		assert.NoError(t, err)

		// Try to create another token pair with the same session ID but different tokens
		// This should overwrite the session data but create new token mappings
		secondAccessToken := "second_access_token"
		baseSession.AccessTokenHash = secondAccessToken
		secondRefreshToken := "second_refresh_token"
		baseSession.RefreshTokenHash = secondRefreshToken

		newSessionData := td.EncodeSessionEncx(t, baseSession)

		err = repo.CreateSession(ctx, baseSession.ID, secondAccessToken, secondRefreshToken, baseSession.UserIDHash, newSessionData, time.Hour, 24*time.Hour)
		assert.NoError(t, err)

		// Both token pairs should work
		retrievedSessionData := td.GetSessionByID(t, ctx, baseSession.ID, testClient)
		assert.NoError(t, err)
		assert.Equal(t, baseSession, retrievedSessionData)

		sessionIDFromFirstAcessToken, err := td.FindSessionByAccessTokenHash(t, ctx, firstAccessToken, testClient)
		assert.NoError(t, err)
		assert.NotNil(t, sessionIDFromFirstAcessToken)

		sessionIDFromFirstRefreshToken, err := td.FindSessionByRefreshTokenHash(t, ctx, firstRefreshToken, testClient)
		assert.NoError(t, err)
		assert.NotNil(t, sessionIDFromFirstRefreshToken)

		// sessionIDFromSecond, err := td.FindSessionByRefreshTokenHash(t, ctx, secondAccessToken, testClient)
		sessionIDFromSecond, err := td.FindSessionByAccessTokenHash(t, ctx, secondAccessToken, testClient)
		assert.NoError(t, err)
		assert.NotNil(t, sessionIDFromSecond, "Second session data should not be nil")

		sessionIDFromSecondRefreshToken, err := td.FindSessionByRefreshTokenHash(t, ctx, secondRefreshToken, testClient)
		assert.NoError(t, err)
		assert.NotNil(t, sessionIDFromSecondRefreshToken)
	})

	t.Run("should rollback access token on refresh token creation failure", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// session := td.CreateTestSessionWithCrypto(t, crypto)
		baseSession := td.NewTestSessionEncx(t)
		accessTokenHash := "access_token_rollback_test"
		baseSession.AccessTokenHash = accessTokenHash
		refreshTokenHash := "refresh_token_rollback_test"
		baseSession.RefreshTokenHash = refreshTokenHash

		sessionData := td.EncodeSessionEncx(t, baseSession)

		// First, create the token pair successfully
		err := repo.CreateSession(ctx, baseSession.ID, accessTokenHash, refreshTokenHash, baseSession.UserIDHash, sessionData, time.Hour, 24*time.Hour)
		assert.NoError(t, err)

		// Close Redis to simulate failure during refresh token creation
		testClient.Close()

		// Try to create another token pair - this should fail
		newAccessToken := "new_access_token"
		baseSession.AccessTokenHash = newAccessToken
		newRefreshToken := "new_refresh_token"
		baseSession.RefreshTokenHash = newRefreshToken

		newSessionData := td.EncodeSessionEncx(t, baseSession)

		err = repo.CreateSession(ctx, baseSession.ID, newAccessToken, newRefreshToken, baseSession.UserIDHash, newSessionData, time.Hour, 24*time.Hour)
		assert.Error(t, err)

		// Reconnect and verify no partial state exists for the failed operation
		reconnectRedis()

		// The failed tokens should not exist
		// _, _, err = repo.FindSessionByAccessTokenHash(ctx, newAccessToken)
		_, err = td.FindSessionByAccessTokenHash(t, ctx, newAccessToken, testClient)
		assert.Error(t, err, "New access token should not exist after rollback")

		// _, _, err = repo.FindSessionByRefreshTokenHash(ctx, newRefreshToken)
		_, err = td.FindSessionByRefreshTokenHash(t, ctx, newRefreshToken, testClient)
		assert.Error(t, err, "New refresh token should not exist after rollback")

		// Original tokens should still work
		// _, _, err = repo.FindSessionByAccessTokenHash(ctx, accessTokenHash)
		_, err = td.FindSessionByAccessTokenHash(t, ctx, accessTokenHash, testClient)
		assert.NoError(t, err, "Original tokens should still work")
	})

	t.Run("should handle zero TTL values", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// session := td.CreateTestSessionWithCrypto(t, crypto)
		baseSession := td.NewTestSessionEncx(t)
		accessTokenHash := "zero_ttl_access_token"
		baseSession.AccessTokenHash = accessTokenHash
		refreshTokenHash := "zero_ttl_refresh_token"
		baseSession.RefreshTokenHash = refreshTokenHash

		sessionData := td.EncodeSessionEncx(t, baseSession)

		// Create token pair with zero TTL (should use Redis default behavior)
		err := repo.CreateSession(ctx, baseSession.ID, accessTokenHash, refreshTokenHash, baseSession.UserIDHash, sessionData, 0, 0)
		assert.NoError(t, err)

		// Tokens should exist and work
		sessionIDFound, err := td.FindSessionByAccessTokenHash(t, ctx, accessTokenHash, testClient)
		assert.NoError(t, err)
		assert.Equal(t, []byte(baseSession.ID.String()), sessionIDFound)

		// Check that keys exist without expiration (TTL = -1 means no expiration)
		accessTokenKey := session.FormatAccessTokenKey(accessTokenHash)
		accessTTL := testClient.TTL(ctx, accessTokenKey).Val()
		assert.Equal(t, time.Duration(-1), accessTTL, "Zero TTL should result in no expiration")
	})

	t.Run("should handle very short TTL values", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// session := td.CreateTestSessionWithCrypto(t, crypto)
		baseSession := td.NewTestSessionEncx(t)
		accessTokenHash := "short_ttl_access_token"
		baseSession.AccessTokenHash = accessTokenHash
		refreshTokenHash := "short_ttl_refresh_token"
		baseSession.RefreshTokenHash = refreshTokenHash

		sessionData := td.EncodeSessionEncx(t, baseSession)
		shortTTL := 50 * time.Millisecond

		// Create token pair with very short TTL
		err := repo.CreateSession(ctx, baseSession.ID, accessTokenHash, refreshTokenHash, baseSession.UserIDHash, sessionData, shortTTL, shortTTL)
		assert.NoError(t, err)

		// Tokens should initially work
		_, err = td.FindSessionByAccessTokenHash(t, ctx, accessTokenHash, testClient)
		assert.NoError(t, err)

		// Wait for expiration
		time.Sleep(60 * time.Millisecond)

		// Tokens should be expired
		_, err = td.FindSessionByAccessTokenHash(t, ctx, accessTokenHash, testClient)
		assert.Error(t, err, "Access token should be expired")

		_, err = td.FindSessionByRefreshTokenHash(t, ctx, refreshTokenHash, testClient)
		assert.Error(t, err, "Refresh token should be expired")
	})

	t.Run("should handle special characters in token hashes", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// session := td.CreateTestSessionWithCrypto(t, crypto)
		baseSession := td.NewTestSessionEncx(t)
		// Use various special characters that might appear in base64/hex encoded hashes
		accessTokenHash := "access-token_hash.with:special@chars+123/456="
		baseSession.AccessTokenHash = accessTokenHash
		refreshTokenHash := "refresh-token_hash.with:special@chars+789/012="
		baseSession.RefreshTokenHash = refreshTokenHash

		sessionData := td.EncodeSessionEncx(t, baseSession)

		// Create token pair with special characters
		err := repo.CreateSession(ctx, baseSession.ID, accessTokenHash, refreshTokenHash, baseSession.UserIDHash, sessionData, time.Hour, 24*time.Hour)
		assert.NoError(t, err)

		// Verify both tokens work
		sessionIDFound, err := td.FindSessionByAccessTokenHash(t, ctx, accessTokenHash, testClient)
		assert.NoError(t, err)
		assert.Equal(t, []byte(baseSession.ID.String()), sessionIDFound)

		sessionIDFound, err = td.FindSessionByRefreshTokenHash(t, ctx, refreshTokenHash, testClient)
		assert.NoError(t, err)
		assert.Equal(t, []byte(baseSession.ID.String()), sessionIDFound)
	})

	t.Run("should handle large session data", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create session with large encrypted data
		baseSession := td.NewTestSessionEncx(t)
		accessTokenHash := "large_data_access_token"
		baseSession.AccessTokenHash = accessTokenHash
		refreshTokenHash := "large_data_refresh_token"
		baseSession.RefreshTokenHash = refreshTokenHash

		sessionData := td.EncodeSessionEncx(t, baseSession)

		// Add extra data to make it larger
		largeSessionData := make([]byte, len(sessionData)+10000)
		copy(largeSessionData, sessionData)
		for i := len(sessionData); i < len(largeSessionData); i++ {
			largeSessionData[i] = byte(i % 256)
		}

		// Create token pair with large session data
		err := repo.CreateSession(ctx, baseSession.ID, accessTokenHash, refreshTokenHash, baseSession.UserIDHash, largeSessionData, time.Hour, 24*time.Hour)
		assert.NoError(t, err)

		// Verify data integrity
		retrievedID, err := td.FindSessionByAccessTokenHash(t, ctx, accessTokenHash, testClient)
		assert.NoError(t, err)
		assert.Equal(t, []byte(baseSession.ID.String()), retrievedID)
	})

	t.Run("should handle Redis connection errors gracefully", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Close Redis connection
		testClient.Close()

		// session := td.CreateTestSessionWithCrypto(t, crypto)
		baseSession := td.NewTestSessionEncx(t)
		accessTokenHash := "connection_error_access_token"
		baseSession.AccessTokenHash = accessTokenHash
		refreshTokenHash := "connection_error_refresh_token"
		baseSession.RefreshTokenHash = refreshTokenHash

		sessionData := td.EncodeSessionEncx(t, baseSession)

		// Try to create token pair with closed connection
		err := repo.CreateSession(ctx, baseSession.ID, accessTokenHash, refreshTokenHash, baseSession.UserIDHash, sessionData, time.Hour, 24*time.Hour)
		assert.Error(t, err)

		// Reconnect for subsequent tests
		reconnectRedis()
	})
}
