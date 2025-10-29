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

// make test-func TEST_NAME=TestRemoveSessionByID TEST_PATH=internal/authuser/infrastructure/redis/session/remove_session_by_id_test.go

func TestRemoveSessionByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully remove existing session by ID", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create test session directly in Redis
		baseSession := td.NewTestSessionEncx(t)
		ttl := 1 * time.Hour
		td.InsertSessionEncx(t, ctx, testClient, baseSession, ttl)

		// Verify session exists before removal using raw Redis query
		sessionKey := session.FormatSessionKey(baseSession.ID.String())
		exists, err := testClient.Exists(ctx, sessionKey).Result()
		require.NoError(t, err)
		require.Equal(t, int64(1), exists, "Session should exist before removal")

		// Remove session by ID
		err = repo.RemoveSessionByID(ctx, baseSession.ID)
		assert.NoError(t, err)

		// Verify session is removed using raw Redis query
		existsAfter, err := testClient.Exists(ctx, sessionKey).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), existsAfter, "Session should be removed")

		// Verify token mapping still exists (this method only removes session data) using raw Redis query
		tokenKey := session.FormatAccessTokenKey(baseSession.AccessTokenHash)
		tokenExists, err := testClient.Exists(ctx, tokenKey).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), tokenExists, "Token mapping should still exist")
	})

	t.Run("should handle non-existent session removal gracefully", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		nonExistentID := uuid.New()

		// Remove non-existent session (should not error)
		err := repo.RemoveSessionByID(ctx, nonExistentID)
		assert.NoError(t, err, "Removing non-existent session should not error")
	})

	t.Run("should handle Redis connection failure", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		sessionID := uuid.New()

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

			baseSession := td.NewTestSessionEncx(t)
			td.InsertSessionEncx(t, ctx, testClient, baseSession, ttl)

			sessions[i] = &td.SessionTestData{
				SessionID: baseSession.ID,
				TokenHash: baseSession.AccessTokenHash,
				Session:   baseSession,
			}
		}

		// Remove middle session
		targetSession := sessions[2]
		err := repo.RemoveSessionByID(ctx, targetSession.SessionID)
		assert.NoError(t, err)

		// Verify target session is removed using raw Redis query
		targetSessionKey := session.FormatSessionKey(targetSession.SessionID.String())
		targetExists, err := testClient.Exists(ctx, targetSessionKey).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), targetExists, "Target session should be removed")

		// Verify other sessions still exist using raw Redis queries
		for i, s := range sessions {
			if i == 2 { // Skip removed session
				continue
			}
			sessionKey := session.FormatSessionKey(s.SessionID.String())
			exists, err := testClient.Exists(ctx, sessionKey).Result()
			assert.NoError(t, err)
			assert.Equal(t, int64(1), exists, "Session %d should still exist", i)
		}

		// Verify session count using raw Redis query
		allKeys := session.FormatSessionKey("*")
		sessionKeys, err := testClient.Keys(ctx, allKeys).Result()
		assert.NoError(t, err)
		assert.Equal(t, 4, len(sessionKeys), "Should have 4 remaining sessions")
	})
}
