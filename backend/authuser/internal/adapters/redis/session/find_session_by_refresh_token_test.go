package sessionRepository_test

import (
	"context"
	"testing"
	"time"

	td "github.com/Leviosa-care/authuser/test/helpers"
	sessionRepository "github.com/Leviosa-care/authuser/internal/adapters/redis/session"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestFindSessionByRefreshToken make test-unit-session-test

func TestFindSessionByRefreshToken(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully find session by refresh token", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create test session and tokens
		session := td.CreateTestSessionWithCrypto(t, crypto)
		accessTokenHash := "access_token_hash_123"
		refreshTokenHash := "refresh_token_hash_456"
		accessTTL := 1 * time.Hour
		refreshTTL := 24 * time.Hour

		// Create token pair in Redis (which creates all necessary mappings)
		err := repo.CreateTokenPair(ctx, session.ID, accessTokenHash, refreshTokenHash,
			td.EncodeSession(t, session), accessTTL, refreshTTL)
		require.NoError(t, err)

		// Find session by refresh token
		sessionData, err := repo.FindSessionByRefreshToken(ctx, refreshTokenHash)
		require.NoError(t, err)
		require.NotNil(t, sessionData)

		// Verify session data integrity
		retrievedSession := td.DecodeSessionWithDecryption(t, sessionData, crypto)
		assert.Equal(t, session.UserID, retrievedSession.UserID)
		assert.Equal(t, session.Role, retrievedSession.Role)
		assert.Equal(t, session.State, retrievedSession.State)
		assert.Equal(t, session.TokenHash, retrievedSession.TokenHash)
	})

	t.Run("should return error for non-existent refresh token", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		nonExistentRefreshToken := "non_existent_refresh_token_hash"

		// Try to find session with non-existent refresh token
		sessionData, err := repo.FindSessionByRefreshToken(ctx, nonExistentRefreshToken)
		assert.Error(t, err)
		assert.Nil(t, sessionData)
	})

	t.Run("should return error when refresh token exists but session data is missing", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create only the refresh token mapping without session data
		sessionID := uuid.New()
		refreshTokenHash := "orphaned_refresh_token_hash"
		refreshTokenKey := sessionRepository.FormatRefreshTokenKey(refreshTokenHash)

		// Set refresh token -> session ID mapping
		err := testClient.Set(ctx, refreshTokenKey, sessionID.String(), time.Hour).Err()
		require.NoError(t, err)

		// Try to find session (should fail because session data doesn't exist)
		sessionData, err := repo.FindSessionByRefreshToken(ctx, refreshTokenHash)
		assert.Error(t, err)
		assert.Nil(t, sessionData)
	})

	t.Run("should handle expired refresh token gracefully", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create session with very short-lived refresh token
		session := td.CreateTestSessionWithCrypto(t, crypto)
		accessTokenHash := "access_token_hash"
		refreshTokenHash := "short_lived_refresh_token"
		accessTTL := 1 * time.Hour
		refreshTTL := 10 * time.Millisecond // Very short TTL

		// Create token pair
		err := repo.CreateTokenPair(ctx, session.ID, accessTokenHash, refreshTokenHash,
			td.EncodeSession(t, session), accessTTL, refreshTTL)
		require.NoError(t, err)

		// Wait for refresh token to expire
		time.Sleep(20 * time.Millisecond)

		// Try to find session with expired refresh token
		sessionData, err := repo.FindSessionByRefreshToken(ctx, refreshTokenHash)
		assert.Error(t, err)
		assert.Nil(t, sessionData)
	})

	t.Run("should work with long-lived refresh tokens", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create session with long-lived refresh token
		session := td.CreateTestSessionWithCrypto(t, crypto)
		accessTokenHash := "short_access_token"
		refreshTokenHash := "long_lived_refresh_token"
		accessTTL := 50 * time.Millisecond // Short access token
		refreshTTL := 24 * time.Hour       // Long refresh token

		// Create token pair
		err := repo.CreateTokenPair(ctx, session.ID, accessTokenHash, refreshTokenHash,
			td.EncodeSession(t, session), accessTTL, refreshTTL)
		require.NoError(t, err)

		// Wait for access token to expire but refresh token should still be valid
		time.Sleep(60 * time.Millisecond)

		// Access token should be expired
		_, err = repo.FindSessionByAccessToken(ctx, accessTokenHash)
		assert.Error(t, err, "Access token should be expired")

		// But refresh token should still work
		sessionData, err := repo.FindSessionByRefreshToken(ctx, refreshTokenHash)
		require.NoError(t, err)
		require.NotNil(t, sessionData)

		// Verify session data
		retrievedSession := td.DecodeSessionWithDecryption(t, sessionData, crypto)
		assert.Equal(t, session.UserID, retrievedSession.UserID)
	})

	t.Run("should handle special characters in refresh token hash", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create test session with special characters in token hash
		session := td.CreateTestSessionWithCrypto(t, crypto)
		accessTokenHash := "access_token_hash"
		refreshTokenHash := "refresh-token_hash.with:special@chars+123"
		accessTTL := 1 * time.Hour
		refreshTTL := 24 * time.Hour

		// Create token pair
		err := repo.CreateTokenPair(ctx, session.ID, accessTokenHash, refreshTokenHash,
			td.EncodeSession(t, session), accessTTL, refreshTTL)
		require.NoError(t, err)

		// Find session by refresh token with special characters
		sessionData, err := repo.FindSessionByRefreshToken(ctx, refreshTokenHash)
		require.NoError(t, err)
		require.NotNil(t, sessionData)

		// Verify session data
		retrievedSession := td.DecodeSessionWithDecryption(t, sessionData, crypto)
		assert.Equal(t, session.UserID, retrievedSession.UserID)
	})

	t.Run("should handle Redis connection errors gracefully", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Close Redis connection to simulate connection error
		testClient.Close()

		refreshTokenHash := "test_refresh_token_hash"

		// Try to find session with closed Redis connection
		sessionData, err := repo.FindSessionByRefreshToken(ctx, refreshTokenHash)
		assert.Error(t, err)
		assert.Nil(t, sessionData)

		// Reconnect Redis for subsequent tests
		reconnectRedis()
	})
}