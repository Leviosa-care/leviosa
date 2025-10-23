package sessionRepository_test

import (
	"context"
	"testing"
	"time"

	td "github.com/Leviosa-care/leviosa/backend/test/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestRemoveSessionByID make test-unit-session-test

func TestRemoveSessionByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully remove existing session by ID", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create test session directly in Redis
		session := td.CreateTestSessionWithCrypto(t, crypto)
		ttl := 1 * time.Hour
		td.InsertSessionDirectly(t, ctx, testClient, session, ttl)

		// Verify session exists before removal using raw Redis query
		sessionKey := "authuser:session:" + session.ID.String()
		exists, err := testClient.Exists(ctx, sessionKey).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), exists, "Session should exist before removal")

		// Remove session by ID
		err = repo.RemoveSessionByID(ctx, session.ID.String())
		require.NoError(t, err)

		// Verify session is removed using raw Redis query
		existsAfter, err := testClient.Exists(ctx, sessionKey).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(0), existsAfter, "Session should be removed")

		// Verify token mapping still exists (this method only removes session data) using raw Redis query
		tokenKey := "authuser:token:" + session.TokenHash
		tokenExists, err := testClient.Exists(ctx, tokenKey).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), tokenExists, "Token mapping should still exist")
	})

	t.Run("should handle non-existent session removal gracefully", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		nonExistentID := uuid.New().String()

		// Remove non-existent session (should not error)
		err := repo.RemoveSessionByID(ctx, nonExistentID)
		require.NoError(t, err, "Removing non-existent session should not error")
	})

	t.Run("should handle Redis connection failure", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		sessionID := uuid.New().String()

		// Close Redis connection to simulate failure
		testClient.Close()

		// Try to remove session (should fail)
		err := repo.RemoveSessionByID(ctx, sessionID)
		assert.Error(t, err)

		// Reconnect for cleanup
		reconnectRedis()
	})

	t.Run("should remove specific session among multiple sessions", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create multiple sessions
		sessions := make([]*td.SessionTestData, 5)
		ttl := 1 * time.Hour

		for i := 0; i < 5; i++ {
			session := td.CreateTestSessionWithCrypto(t, crypto)
			td.InsertSessionDirectly(t, ctx, testClient, session, ttl)

			sessions[i] = &td.SessionTestData{
				SessionID: session.ID,
				TokenHash: session.TokenHash,
				Session:   session,
			}
		}

		// Remove middle session
		targetSession := sessions[2]
		err := repo.RemoveSessionByID(ctx, targetSession.SessionID.String())
		require.NoError(t, err)

		// Verify target session is removed using raw Redis query
		targetSessionKey := "authuser:session:" + targetSession.SessionID.String()
		targetExists, err := testClient.Exists(ctx, targetSessionKey).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(0), targetExists, "Target session should be removed")

		// Verify other sessions still exist using raw Redis queries
		for i, session := range sessions {
			if i == 2 { // Skip removed session
				continue
			}
			sessionKey := "authuser:session:" + session.SessionID.String()
			exists, err := testClient.Exists(ctx, sessionKey).Result()
			require.NoError(t, err)
			assert.Equal(t, int64(1), exists, "Session %d should still exist", i)
		}

		// Verify session count using raw Redis query
		sessionKeys, err := testClient.Keys(ctx, "authuser:session:*").Result()
		require.NoError(t, err)
		assert.Equal(t, 4, len(sessionKeys), "Should have 4 remaining sessions")
	})

	t.Run("should handle malformed session ID", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		malformedIDs := []string{
			"",
			"not-a-uuid",
			"12345",
			"invalid-session-id-format",
		}

		for _, malformedID := range malformedIDs {
			// Should not error even with malformed ID
			err := repo.RemoveSessionByID(ctx, malformedID)
			require.NoError(t, err, "Should not error with malformed ID: %s", malformedID)
		}
	})
}
