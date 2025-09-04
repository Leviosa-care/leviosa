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

// TEST=TestFindSessionByTokenHash make test-unit-session-test

func TestFindSessionByTokenHash(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully find session by token hash", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create test session directly in Redis
		session := td.CreateTestSessionWithCrypto(t, crypto)
		ttl := 1 * time.Hour
		td.InsertSessionDirectly(t, ctx, testClient, session, ttl)

		// Find session by token hash
		sessionData, err := repo.FindSessionByTokenHash(ctx, session.TokenHash)
		require.NoError(t, err)

		// Verify session data integrity
		retrievedSession := td.DecodeSessionWithDecryption(t, sessionData, crypto)
		assert.Equal(t, session.UserID, retrievedSession.UserID)
		assert.Equal(t, session.TokenHash, retrievedSession.TokenHash)
		assert.Equal(t, session.Role, retrievedSession.Role)
		assert.Equal(t, session.State, retrievedSession.State)
	})

	t.Run("should return not found error for non-existent token hash", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		nonExistentTokenHash := "non_existent_token_hash"

		// Try to find session by non-existent token hash
		sessionData, err := repo.FindSessionByTokenHash(ctx, nonExistentTokenHash)
		assert.Error(t, err)
		assert.Nil(t, sessionData)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("should return not found when token mapping exists but session data is missing", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "orphaned_token_hash"
		sessionID := uuid.New()

		// Create only token mapping without session data
		tokenKey := "authuser:token:" + tokenHash
		err := testClient.Set(ctx, tokenKey, sessionID.String(), 1*time.Hour).Err()
		require.NoError(t, err)

		// Try to find session (token exists but session data doesn't)
		sessionData, err := repo.FindSessionByTokenHash(ctx, tokenHash)
		assert.Error(t, err)
		assert.Nil(t, sessionData)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("should handle Redis connection failure", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "test_token_hash_fail"

		// Close Redis connection to simulate failure
		testClient.Close()

		// Try to find session (should fail)
		sessionData, err := repo.FindSessionByTokenHash(ctx, tokenHash)
		assert.Error(t, err)
		assert.Nil(t, sessionData)

		// Reconnect for cleanup
		reconnectRedis()
	})

	t.Run("should find session with different token hash formats", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHashes := []string{
			"simple_hash",
			"hash_with_numbers_123",
			"hash-with-dashes",
			"hash.with.dots",
			"HASH_WITH_UPPERCASE",
			"very_long_token_hash_with_many_characters_to_test_length_limits",
		}

		sessions := make([]*td.SessionTestData, len(tokenHashes))
		ttl := 1 * time.Hour

		// Create sessions with different token hash formats
		for i, tokenHash := range tokenHashes {
			session := td.CreateTestSessionWithCrypto(t, crypto)
			session.TokenHash = tokenHash
			td.InsertSessionDirectly(t, ctx, testClient, session, ttl)

			sessions[i] = &td.SessionTestData{
				SessionID:   session.ID,
				TokenHash:   tokenHash,
				Session:     session,
				SessionData: td.EncodeSession(t, session),
			}
		}

		// Find each session by token hash
		for i, sessionData := range sessions {
			retrievedData, err := repo.FindSessionByTokenHash(ctx, sessionData.TokenHash)
			require.NoError(t, err, "Should find session %d with token hash: %s", i, sessionData.TokenHash)

			retrievedSession := td.DecodeSessionWithDecryption(t, retrievedData, crypto)
			assert.Equal(t, sessionData.Session.TokenHash, retrievedSession.TokenHash, "Session %d should match", i)
			assert.Equal(t, sessionData.Session.UserID, retrievedSession.UserID, "Session %d should match", i)
			assert.Equal(t, sessionData.Session.Role, retrievedSession.Role, "Session %d should match", i)
			assert.Equal(t, sessionData.Session.State, retrievedSession.State, "Session %d should match", i)
		}
	})

	t.Run("should handle expired token mapping", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create session with short TTL
		session := td.CreateTestSessionWithCrypto(t, crypto)
		ttl := 100 * time.Millisecond
		td.InsertSessionDirectly(t, ctx, testClient, session, ttl)

		// Find session immediately (should exist)
		sessionData, err := repo.FindSessionByTokenHash(ctx, session.TokenHash)
		require.NoError(t, err)
		assert.NotNil(t, sessionData)

		// Wait for expiration
		time.Sleep(150 * time.Millisecond)

		// Try to find expired session (should not exist)
		expiredSessionData, err := repo.FindSessionByTokenHash(ctx, session.TokenHash)
		assert.Error(t, err)
		assert.Nil(t, expiredSessionData)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("should handle empty token hash", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Try to find session with empty token hash
		sessionData, err := repo.FindSessionByTokenHash(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, sessionData)
	})

	t.Run("should find correct session among multiple sessions", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create multiple sessions
		sessions := make([]*td.SessionTestData, 5)
		ttl := 1 * time.Hour

		for i := range 5 {
			session := td.CreateTestSessionWithCrypto(t, crypto)
			session.TokenHash = "token_hash_" + uuid.New().String()[:8]
			td.InsertSessionDirectly(t, ctx, testClient, session, ttl)

			sessions[i] = &td.SessionTestData{
				SessionID: session.ID,
				TokenHash: session.TokenHash,
				Session:   session,
			}
		}

		// Find a specific session (middle one)
		targetSession := sessions[2]
		retrievedData, err := repo.FindSessionByTokenHash(ctx, targetSession.TokenHash)
		require.NoError(t, err)

		retrievedSession := td.DecodeSessionWithDecryption(t, retrievedData, crypto)
		assert.Equal(t, targetSession.Session.UserID, retrievedSession.UserID)
		assert.Equal(t, targetSession.Session.Role, retrievedSession.Role, "Session %d should match")
		assert.Equal(t, targetSession.Session.State, retrievedSession.State, "Session %d should match")
		assert.Equal(t, targetSession.TokenHash, retrievedSession.TokenHash)

		// Verify we didn't get a different session
		for i, session := range sessions {
			if i == 2 { // Skip the target session
				continue
			}
			assert.NotEqual(t, session.Session.UserID, retrievedSession.ID, "Should not return session %d", i)
			assert.NotEqual(t, session.Session.Role, retrievedSession.ID, "Should not return session %d", i)
			assert.NotEqual(t, session.Session.State, retrievedSession.ID, "Should not return session %d", i)
			assert.NotEqual(t, session.Session.TokenHash, retrievedSession.ID, "Should not return session %d", i)
		}
	})
}
