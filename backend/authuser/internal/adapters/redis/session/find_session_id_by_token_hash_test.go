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

// TEST=TestFindSessionIDByTokenHash make test-unit-session-test

func TestFindSessionIDByTokenHash(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully find session ID by token hash", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create test session directly in Redis
		session := td.CreateTestSessionWithCrypto(t, crypto)
		ttl := 1 * time.Hour
		td.InsertSessionDirectly(t, ctx, testClient, session, ttl)

		// Find session ID by token hash
		sessionID, err := repo.FindSessionIDByTokenHash(ctx, session.TokenHash)
		require.NoError(t, err)

		// Verify session ID matches
		assert.Equal(t, session.ID.String(), sessionID)
	})

	t.Run("should return not found error for non-existent token hash", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		nonExistentTokenHash := "non_existent_token_hash"

		// Try to find session ID by non-existent token hash
		sessionID, err := repo.FindSessionIDByTokenHash(ctx, nonExistentTokenHash)
		assert.Error(t, err)
		assert.Empty(t, sessionID)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("should handle Redis connection failure", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "test_token_hash_fail"

		// Close Redis connection to simulate failure
		testClient.Close()

		// Try to find session ID (should fail)
		sessionID, err := repo.FindSessionIDByTokenHash(ctx, tokenHash)
		assert.Error(t, err)
		assert.Empty(t, sessionID)

		// Reconnect for cleanup
		reconnectRedis()
	})

	t.Run("should find session ID with various token hash formats", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		testCases := []struct {
			name      string
			tokenHash string
		}{
			{"simple hash", "simple_hash"},
			{"hash with numbers", "hash123"},
			{"hash with dashes", "hash-with-dashes"},
			{"hash with dots", "hash.with.dots"},
			{"uppercase hash", "UPPERCASE_HASH"},
			{"mixed case", "MixedCase_Hash123"},
			{"long hash", "very_long_token_hash_with_many_characters_for_testing"},
		}

		ttl := 1 * time.Hour
		sessions := make(map[string]uuid.UUID)

		// Create sessions with different token hash formats
		for _, tc := range testCases {
			session := td.CreateTestSessionWithCrypto(t, crypto)
			session.TokenHash = tc.tokenHash
			td.InsertSessionDirectly(t, ctx, testClient, session, ttl)
			sessions[tc.tokenHash] = session.ID
		}

		// Find each session ID by token hash
		for _, tc := range testCases {
			sessionID, err := repo.FindSessionIDByTokenHash(ctx, tc.tokenHash)
			require.NoError(t, err, "Should find session ID for %s", tc.name)

			expectedID := sessions[tc.tokenHash]
			assert.Equal(t, expectedID.String(), sessionID, "Session ID should match for %s", tc.name)
		}
	})

	t.Run("should handle expired token mapping", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create session with short TTL
		session := td.CreateTestSessionWithCrypto(t, crypto)
		ttl := 100 * time.Millisecond
		td.InsertSessionDirectly(t, ctx, testClient, session, ttl)

		// Find session ID immediately (should exist)
		sessionID, err := repo.FindSessionIDByTokenHash(ctx, session.TokenHash)
		require.NoError(t, err)
		assert.Equal(t, session.ID.String(), sessionID)

		// Wait for expiration
		time.Sleep(150 * time.Millisecond)

		// Try to find expired session ID (should not exist)
		expiredSessionID, err := repo.FindSessionIDByTokenHash(ctx, session.TokenHash)
		assert.Error(t, err)
		assert.Empty(t, expiredSessionID)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("should handle empty token hash", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Try to find session ID with empty token hash
		sessionID, err := repo.FindSessionIDByTokenHash(ctx, "")
		assert.Error(t, err)
		assert.Empty(t, sessionID)
	})

	t.Run("should find correct session ID among multiple sessions", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create multiple sessions
		numSessions := 10
		sessions := make(map[string]uuid.UUID)
		ttl := 1 * time.Hour

		for i := 0; i < numSessions; i++ {
			session := td.CreateTestSessionWithCrypto(t, crypto)
			session.TokenHash = "token_hash_" + uuid.New().String()[:8]
			td.InsertSessionDirectly(t, ctx, testClient, session, ttl)
			sessions[session.TokenHash] = session.ID
		}

		// Find each session ID and verify uniqueness
		retrievedSessions := make(map[string]uuid.UUID)

		for tokenHash, expectedID := range sessions {
			retrievedID, err := repo.FindSessionIDByTokenHash(ctx, tokenHash)
			require.NoError(t, err, "Should find session ID for token hash: %s", tokenHash)

			parsedID, err := uuid.Parse(retrievedID)
			require.NoError(t, err, "Retrieved ID should be valid UUID")

			assert.Equal(t, expectedID.String(), retrievedID, "Session ID should match for token hash: %s", tokenHash)

			// Verify this ID hasn't been retrieved for a different token
			for otherToken, otherID := range retrievedSessions {
				assert.NotEqual(t, otherID, parsedID, "Session ID should be unique (found duplicate for tokens %s and %s)", tokenHash, otherToken)
			}

			retrievedSessions[tokenHash] = parsedID
		}

		// Verify we retrieved all sessions
		assert.Equal(t, numSessions, len(retrievedSessions), "Should have retrieved all session IDs")
	})

	t.Run("should handle special characters in token hash", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		specialTokenHashes := []string{
			"token_with_underscores",
			"token-with-hyphens",
			"token.with.dots",
			"token+with+plus",
			"token=with=equals",
			"token%20with%20encoding",
		}

		ttl := 1 * time.Hour
		expectedIDs := make(map[string]uuid.UUID)

		// Create sessions with special character token hashes
		for _, tokenHash := range specialTokenHashes {
			session := td.CreateTestSessionWithCrypto(t, crypto)
			session.TokenHash = tokenHash
			td.InsertSessionDirectly(t, ctx, testClient, session, ttl)
			expectedIDs[tokenHash] = session.ID
		}

		// Find each session ID
		for _, tokenHash := range specialTokenHashes {
			sessionID, err := repo.FindSessionIDByTokenHash(ctx, tokenHash)
			require.NoError(t, err, "Should find session ID for token with special chars: %s", tokenHash)

			expectedID := expectedIDs[tokenHash]
			assert.Equal(t, expectedID.String(), sessionID, "Session ID should match for token: %s", tokenHash)
		}
	})
}
