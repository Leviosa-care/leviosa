package sessionRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	td "github.com/Leviosa-care/authuser/test/helpers"
)

// TEST=TestRemoveSessionByToken make test-unit-session-test

func TestRemoveSessionByToken(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully remove token mapping by token hash", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create test session directly in Redis
		session := td.CreateTestSessionWithCrypto(t, crypto)
		ttl := 1 * time.Hour
		td.InsertSessionDirectly(t, ctx, testClient, session, ttl)

		// Verify token mapping exists before removal using raw Redis query
		tokenKey := "authuser:token:" + session.TokenHash
		exists, err := testClient.Exists(ctx, tokenKey).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), exists, "Token mapping should exist before removal")

		// Remove token mapping
		err = repo.RemoveSessionByToken(ctx, session.TokenHash)
		require.NoError(t, err)

		// Verify token mapping is removed using raw Redis query
		existsAfter, err := testClient.Exists(ctx, tokenKey).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(0), existsAfter, "Token mapping should be removed")

		// Verify session data still exists (this method only removes token mapping) using raw Redis query
		sessionKey := "authuser:session:" + session.ID.String()
		sessionExists, err := testClient.Exists(ctx, sessionKey).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), sessionExists, "Session data should still exist")
	})

	t.Run("should handle non-existent token removal gracefully", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		nonExistentTokenHash := "non_existent_token_hash"

		// Remove non-existent token mapping (should not error)
		err := repo.RemoveSessionByToken(ctx, nonExistentTokenHash)
		require.NoError(t, err, "Removing non-existent token mapping should not error")
	})

	t.Run("should handle Redis connection failure", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "test_token_hash_fail"

		// Close Redis connection to simulate failure
		testClient.Close()

		// Try to remove token mapping (should fail)
		err := repo.RemoveSessionByToken(ctx, tokenHash)
		assert.Error(t, err)

		// Reconnect for cleanup
		reconnectRedis()
	})

	t.Run("should remove specific token mapping among multiple mappings", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create multiple sessions
		sessions := make([]*td.SessionTestData, 5)
		ttl := 1 * time.Hour

		for i := 0; i < 5; i++ {
			session := td.CreateTestSessionWithCrypto(t, crypto)
			session.TokenHash = "token_hash_" + uuid.New().String()[:8]
			td.InsertSessionDirectly(t, ctx, testClient, session, ttl)

			sessions[i] = &td.SessionTestData{
				SessionID: session.ID,
				TokenHash: session.TokenHash,
				Session:   session,
			}
		}

		// Remove middle token mapping
		targetSession := sessions[2]
		err := repo.RemoveSessionByToken(ctx, targetSession.TokenHash)
		require.NoError(t, err)

		// Verify target token mapping is removed using raw Redis query
		targetTokenKey := "authuser:token:" + targetSession.TokenHash
		targetExists, err := testClient.Exists(ctx, targetTokenKey).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(0), targetExists, "Target token mapping should be removed")

		// Verify other token mappings still exist using raw Redis queries
		for i, session := range sessions {
			if i == 2 { // Skip removed mapping
				continue
			}
			tokenKey := "authuser:token:" + session.TokenHash
			exists, err := testClient.Exists(ctx, tokenKey).Result()
			require.NoError(t, err)
			assert.Equal(t, int64(1), exists, "Token mapping %d should still exist", i)
		}

		// Verify token mapping count using raw Redis query
		tokenKeys, err := testClient.Keys(ctx, "authuser:token:*").Result()
		require.NoError(t, err)
		assert.Equal(t, 4, len(tokenKeys), "Should have 4 remaining token mappings")

		// Verify all session data still exists using raw Redis query
		sessionKeys, err := testClient.Keys(ctx, "authuser:session:*").Result()
		require.NoError(t, err)
		assert.Equal(t, 5, len(sessionKeys), "All session data should still exist")
	})

	t.Run("should handle empty token hash", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Should not error with empty token hash
		err := repo.RemoveSessionByToken(ctx, "")
		require.NoError(t, err, "Should not error with empty token hash")
	})

	t.Run("should handle various token hash formats", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHashes := []string{
			"simple_hash",
			"hash-with-dashes",
			"hash.with.dots",
			"UPPERCASE_HASH",
			"hash123numbers",
			"very_long_token_hash_with_many_characters",
		}

		ttl := 1 * time.Hour

		// Create token mappings with different formats
		for _, tokenHash := range tokenHashes {
			session := td.CreateTestSessionWithCrypto(t, crypto)
			session.TokenHash = tokenHash
			td.InsertSessionDirectly(t, ctx, testClient, session, ttl)
		}

		// Remove each token mapping
		for _, tokenHash := range tokenHashes {
			err := repo.RemoveSessionByToken(ctx, tokenHash)
			require.NoError(t, err, "Should remove token hash: %s", tokenHash)

			// Verify it's removed using raw Redis query
			tokenKey := "authuser:token:" + tokenHash
			exists, err := testClient.Exists(ctx, tokenKey).Result()
			require.NoError(t, err)
			assert.Equal(t, int64(0), exists, "Token mapping should be removed: %s", tokenHash)
		}

		// Verify all token mappings are removed using raw Redis query
		tokenKeys, err := testClient.Keys(ctx, "authuser:token:*").Result()
		require.NoError(t, err)
		assert.Equal(t, 0, len(tokenKeys), "All token mappings should be removed")
	})
}
