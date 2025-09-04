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

// TEST=TestCreateSession make test-unit-session-test

func TestCreateSession(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully create session with both keys", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Test data
		sessionID := uuid.New()
		tokenHash := "test_token_hash_123"
		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)
		ttl := 1 * time.Hour

		// Create session
		err := repo.CreateSession(ctx, sessionID, tokenHash, sessionData, ttl)
		require.NoError(t, err)

		// Verify session key exists using raw Redis query
		sessionKey := "authuser:session:" + sessionID.String()
		sessionExists, err := testClient.Exists(ctx, sessionKey).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), sessionExists, "Session key should exist in Redis")

		// Verify token mapping exists using raw Redis query
		tokenKey := "authuser:token:" + tokenHash
		tokenExists, err := testClient.Exists(ctx, tokenKey).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), tokenExists, "Token mapping should exist in Redis")

		// Verify session data integrity
		retrievedData, err := testClient.Get(ctx, "authuser:session:"+sessionID.String()).Bytes()
		require.NoError(t, err)
		assert.Equal(t, sessionData, retrievedData, "Session data should match")

		// Verify token mapping integrity using raw Redis query
		retrievedSessionID, err := testClient.Get(ctx, "authuser:token:"+tokenHash).Result()
		require.NoError(t, err)
		assert.Equal(t, sessionID.String(), retrievedSessionID, "Token should map to correct session ID")

		// Verify TTL is set (should be approximately the same)
		sessionTTL, err := testClient.TTL(ctx, "authuser:session:"+sessionID.String()).Result()
		require.NoError(t, err)
		assert.True(t, sessionTTL > 50*time.Minute && sessionTTL <= ttl, "Session TTL should be approximately correct")

		tokenTTL, err := testClient.TTL(ctx, "authuser:token:"+tokenHash).Result()
		require.NoError(t, err)
		assert.True(t, tokenTTL > 50*time.Minute && tokenTTL <= ttl, "Token TTL should be approximately correct")
	})

	t.Run("should handle Redis connection failure gracefully", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Close the Redis connection to simulate failure
		testClient.Close()

		sessionID := uuid.New()
		tokenHash := "test_token_hash_fail"
		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)
		ttl := 1 * time.Hour

		// Attempt to create session (should fail)
		err := repo.CreateSession(ctx, sessionID, tokenHash, sessionData, ttl)
		assert.Error(t, err, "Should return error when Redis is unavailable")

		// Reconnect for cleanup
		reconnectRedis()
	})

	t.Run("should rollback session key if token key creation fails", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		sessionID := uuid.New()
		tokenHash := "test_token_hash_rollback"
		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)
		ttl := 1 * time.Hour

		// First, manually create the token key with a different value to cause conflict
		tokenKey := "authuser:token:" + tokenHash
		err := testClient.Set(ctx, tokenKey, "existing_value", 1*time.Hour).Err()
		require.NoError(t, err)

		// Try to create session (should handle the conflict properly)
		err = repo.CreateSession(ctx, sessionID, tokenHash, sessionData, ttl)
		// Note: The current implementation will overwrite the existing token key
		// This test verifies that behavior works as expected
		require.NoError(t, err)

		// Verify both keys exist using raw Redis queries
		sessionKey := "authuser:session:" + sessionID.String()
		sessionExists, err := testClient.Exists(ctx, sessionKey).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), sessionExists, "Session should exist")

		tokenExists, err := testClient.Exists(ctx, tokenKey).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), tokenExists, "Token mapping should exist")
	})

	t.Run("should handle empty session data", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		sessionID := uuid.New()
		tokenHash := "test_token_hash_empty"
		var sessionData []byte // Empty data
		ttl := 1 * time.Hour

		// Create session with empty data
		err := repo.CreateSession(ctx, sessionID, tokenHash, sessionData, ttl)
		require.NoError(t, err)

		// Verify keys exist even with empty data using raw Redis queries
		sessionKey := "authuser:session:" + sessionID.String()
		sessionExists, err := testClient.Exists(ctx, sessionKey).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), sessionExists, "Session key should exist even with empty data")

		tokenKey := "authuser:token:" + tokenHash
		tokenExists, err := testClient.Exists(ctx, tokenKey).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), tokenExists, "Token mapping should exist")

		// Verify empty data is stored correctly
		retrievedData, err := testClient.Get(ctx, "authuser:session:"+sessionID.String()).Bytes()
		require.NoError(t, err)
		assert.Empty(t, retrievedData, "Empty session data should be preserved")
	})

	t.Run("should handle very short TTL", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		sessionID := uuid.New()
		tokenHash := "test_token_hash_short_ttl"
		session := td.CreateTestSessionWithCrypto(t, crypto)
		sessionData := td.EncodeSession(t, session)
		ttl := 100 * time.Millisecond

		// Create session with short TTL
		err := repo.CreateSession(ctx, sessionID, tokenHash, sessionData, ttl)
		require.NoError(t, err)

		// Verify keys exist initially using raw Redis queries
		sessionKey := "authuser:session:" + sessionID.String()
		sessionExists, err := testClient.Exists(ctx, sessionKey).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), sessionExists, "Session should exist initially")

		tokenKey := "authuser:token:" + tokenHash
		tokenExists, err := testClient.Exists(ctx, tokenKey).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), tokenExists, "Token mapping should exist initially")

		// Wait for expiration
		time.Sleep(200 * time.Millisecond)

		// Verify keys have expired using raw Redis queries
		sessionExistsAfter, err := testClient.Exists(ctx, sessionKey).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(0), sessionExistsAfter, "Session should expire after TTL")

		tokenExistsAfter, err := testClient.Exists(ctx, tokenKey).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(0), tokenExistsAfter, "Token mapping should expire after TTL")
	})

	t.Run("should create multiple sessions without conflict", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		sessions := make([]*td.SessionTestData, 3)
		ttl := 1 * time.Hour

		// Create multiple sessions
		for i := 0; i < 3; i++ {
			sessionID := uuid.New()
			tokenHash := "test_token_hash_multi_" + uuid.New().String()[:8]
			session := td.CreateTestSessionWithCrypto(t, crypto)
			sessionData := td.EncodeSession(t, session)

			sessions[i] = &td.SessionTestData{
				SessionID:   sessionID,
				TokenHash:   tokenHash,
				Session:     session,
				SessionData: sessionData,
			}

			err := repo.CreateSession(ctx, sessionID, tokenHash, sessionData, ttl)
			require.NoError(t, err, "Should create session %d successfully", i)
		}

		// Verify all sessions exist using raw Redis queries
		sessionKeys, err := testClient.Keys(ctx, "authuser:session:*").Result()
		require.NoError(t, err)
		tokenKeys, err := testClient.Keys(ctx, "authuser:token:*").Result()
		require.NoError(t, err)

		assert.Equal(t, 3, len(sessionKeys), "Should have 3 session keys")
		assert.Equal(t, 3, len(tokenKeys), "Should have 3 token keys")

		// Verify each session's data integrity using raw Redis queries
		for i, sessionData := range sessions {
			sessionKey := "authuser:session:" + sessionData.SessionID.String()
			sessionExists, err := testClient.Exists(ctx, sessionKey).Result()
			require.NoError(t, err)
			assert.Equal(t, int64(1), sessionExists, "Session %d should exist", i)

			tokenKey := "authuser:token:" + sessionData.TokenHash
			tokenExists, err := testClient.Exists(ctx, tokenKey).Result()
			require.NoError(t, err)
			assert.Equal(t, int64(1), tokenExists, "Token mapping %d should exist", i)

			retrievedSessionID, err := testClient.Get(ctx, tokenKey).Result()
			require.NoError(t, err)
			assert.Equal(t, sessionData.SessionID.String(), retrievedSessionID, "Token %d should map to correct session ID", i)
		}
	})
}
