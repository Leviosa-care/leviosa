package sessionRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestInvalidateTokenPair TEST_PATH=internal/authuser/infrastructure/redis/session/invalidate_token_pair_test.go

func TestInvalidateTokenPair(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully invalidate all token pair keys", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create token pair
		baseSession := td.NewTestSessionEncx(t)
		accessTokenHash := "invalidate_test_access_token"
		baseSession.AccessTokenHash = accessTokenHash
		refreshTokenHash := "invalidate_test_refresh_token"
		baseSession.RefreshTokenHash = refreshTokenHash

		td.InsertSessionEncx(t, ctx, testClient, baseSession, 1*time.Hour)

		// Verify tokens work initially
		_, err := td.FindSessionByAccessTokenHash(t, ctx, accessTokenHash, testClient)
		require.NoError(t, err)

		_, err = td.FindSessionByRefreshTokenHash(t, ctx, refreshTokenHash, testClient)
		require.NoError(t, err)

		// Invalidate token pair
		err = repo.InvalidateTokenPair(ctx, accessTokenHash, refreshTokenHash, baseSession.ID)
		assert.NoError(t, err)

		// All keys should be removed
		_, err = td.FindSessionByAccessTokenHash(t, ctx, accessTokenHash, testClient)
		assert.Error(t, err, "Access token should be invalidated")

		_, err = td.FindSessionByRefreshTokenHash(t, ctx, refreshTokenHash, testClient)
		assert.Error(t, err, "Refresh token should be invalidated")

		// Session data should also be removed
		sessionKey := session.FormatSessionKey(baseSession.ID.String())
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
		assert.NoError(t, err)
	})

	t.Run("should handle partial token pair states", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create only baseSession data and access token (no refresh token)
		baseSession := td.NewTestSessionEncx(t)
		sessionData := td.EncodeSession(t, baseSession)
		accessTokenHash := "partial_state_access_token"
		refreshTokenHash := "partial_state_refresh_token"

		// Manually create partial state
		sessionKey := session.FormatSessionKey(baseSession.ID.String())
		accessTokenKey := session.FormatAccessTokenKey(accessTokenHash)

		err := testClient.Set(ctx, sessionKey, sessionData, time.Hour).Err()
		require.NoError(t, err)

		err = testClient.Set(ctx, accessTokenKey, baseSession.ID.String(), time.Hour).Err()
		require.NoError(t, err)
		// Note: not creating refresh token mapping

		// Invalidate should work and not fail on missing refresh token
		err = repo.InvalidateTokenPair(ctx, accessTokenHash, refreshTokenHash, baseSession.ID)
		assert.NoError(t, err)

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
		baseSession := td.NewTestSessionEncx(t)

		td.InsertSessionEncx(t, ctx, testClient, baseSession, time.Hour)

		// Invalidate tokens with special characters
		err := repo.InvalidateTokenPair(ctx, baseSession.AccessTokenHash, baseSession.RefreshTokenHash, baseSession.ID)
		assert.NoError(t, err)

		// Verify all keys are removed
		_, err = td.FindSessionByAccessTokenHash(t, ctx, baseSession.AccessTokenHash, testClient)
		assert.Error(t, err, "Access token with special chars should be invalidated")

		_, err = td.FindSessionByRefreshTokenHash(t, ctx, baseSession.RefreshTokenHash, testClient)
		assert.Error(t, err, "Refresh token with special chars should be invalidated")
	})

	t.Run("should not affect other sessions", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create two separate token pairs
		session1 := td.NewTestSessionEncx(t)
		session2 := td.NewTestSessionEncx(t)
		session2.ID = uuid.New()
		session2.AccessTokenHash = "session_2_access_token_hash"
		session2.RefreshTokenHash = "session_2_refresh_token_hash"

		// First token pair
		td.InsertSessionEncx(t, ctx, testClient, session1, time.Hour)
		// Second token pair
		td.InsertSessionEncx(t, ctx, testClient, session2, time.Hour)

		// Invalidate only first token pair
		err := repo.InvalidateTokenPair(ctx, session1.AccessTokenHash, session1.RefreshTokenHash, session1.ID)
		assert.NoError(t, err)

		// First session should be invalidated
		_, err = td.FindSessionByAccessTokenHash(t, ctx, session1.AccessTokenHash, testClient)
		assert.Error(t, err)

		// Second session should still work
		_, err = td.FindSessionByAccessTokenHash(t, ctx, session2.AccessTokenHash, testClient)
		assert.NoError(t, err)

		_, err = td.FindSessionByRefreshTokenHash(t, ctx, session2.RefreshTokenHash, testClient)
		assert.NoError(t, err)
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
