package sessionRepository_test

import (
	"context"
	"testing"
	"time"

	sessionRepository "github.com/Leviosa-care/authuser/internal/adapters/redis/session"
	td "github.com/Leviosa-care/authuser/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestFindSessionByAccessTokenHash make test-unit-session-test

func TestFindSessionByAccessTokenHash(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully find session by access token", func(t *testing.T) {
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

		// Find session by access token
		_, sessionData, err := repo.FindSessionByAccessTokenHash(ctx, accessTokenHash)
		require.NoError(t, err)
		require.NotNil(t, sessionData)

		// Verify session data integrity
		retrievedSession := td.DecodeSessionWithDecryption(t, sessionData, crypto)
		assert.Equal(t, session.UserID, retrievedSession.UserID)
		assert.Equal(t, session.Role, retrievedSession.Role)
		assert.Equal(t, session.State, retrievedSession.State)
		assert.Equal(t, session.AccessTokenHash, retrievedSession.AccessTokenHash)
	})

	t.Run("should return error for non-existent access token", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		nonExistentAccessToken := "non_existent_access_token_hash"

		// Try to find session with non-existent access token
		_, sessionData, err := repo.FindSessionByAccessTokenHash(ctx, nonExistentAccessToken)
		assert.Error(t, err)
		assert.Nil(t, sessionData)
	})

	t.Run("should return error when access token exists but session data is missing", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create only the access token mapping without session data
		sessionID := uuid.New()
		accessTokenHash := "orphaned_access_token_hash"
		accessTokenKey := sessionRepository.FormatAccessTokenKey(accessTokenHash)

		// Set access token -> session ID mapping
		err := testClient.Set(ctx, accessTokenKey, sessionID.String(), time.Hour).Err()
		require.NoError(t, err)

		// Try to find session (should fail because session data doesn't exist)
		_, sessionData, err := repo.FindSessionByAccessTokenHash(ctx, accessTokenHash)
		assert.Error(t, err)
		assert.Nil(t, sessionData)
	})

	t.Run("should handle expired access token gracefully", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create session with very short-lived access token
		session := td.CreateTestSessionWithCrypto(t, crypto)
		accessTokenHash := "short_lived_access_token"
		refreshTokenHash := "refresh_token_hash"
		accessTTL := 10 * time.Millisecond // Very short TTL
		refreshTTL := 24 * time.Hour

		// Create token pair
		err := repo.CreateTokenPair(ctx, session.ID, accessTokenHash, refreshTokenHash,
			td.EncodeSession(t, session), accessTTL, refreshTTL)
		require.NoError(t, err)

		// Wait for access token to expire
		time.Sleep(20 * time.Millisecond)

		// Try to find session with expired access token
		_, sessionData, err := repo.FindSessionByAccessTokenHash(ctx, accessTokenHash)
		assert.Error(t, err)
		assert.Nil(t, sessionData)
	})

	t.Run("should handle special characters in access token hash", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create test session with special characters in token hash
		session := td.CreateTestSessionWithCrypto(t, crypto)
		accessTokenHash := "access-token_hash.with:special@chars"
		refreshTokenHash := "refresh_token_hash"
		accessTTL := 1 * time.Hour
		refreshTTL := 24 * time.Hour

		// Create token pair
		err := repo.CreateTokenPair(ctx, session.ID, accessTokenHash, refreshTokenHash,
			td.EncodeSession(t, session), accessTTL, refreshTTL)
		require.NoError(t, err)

		// Find session by access token with special characters
		_, sessionData, err := repo.FindSessionByAccessTokenHash(ctx, accessTokenHash)
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

		accessTokenHash := "test_access_token_hash"

		// Try to find session with closed Redis connection
		_, sessionData, err := repo.FindSessionByAccessTokenHash(ctx, accessTokenHash)
		assert.Error(t, err)
		assert.Nil(t, sessionData)

		// Reconnect Redis for subsequent tests
		reconnectRedis()
	})
}
