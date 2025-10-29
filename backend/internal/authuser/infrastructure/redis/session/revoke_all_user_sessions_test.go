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

// make test-func TEST_NAME=TestRevokeAllUserSessions TEST_PATH=internal/authuser/infrastructure/redis/session/revoke_all_user_sessions_test.go

func TestRevokeAllUserSessions(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully revoke all sessions for user with multiple sessions", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create multiple sessions for the same user
		userIDHash := "test_user_hash"
		sessions := make([]*session.SessionEncx, 3)
		ttl := 1 * time.Hour

		for i := 0; i < 3; i++ {
			baseSession := td.NewTestSessionEncx(t)
			baseSession.UserIDHash = userIDHash
			td.InsertSessionEncx(t, ctx, testClient, baseSession, ttl)
			sessions[i] = baseSession
		}

		// Verify sessions exist before revocation
		userSessionIndexKey := session.FormatUserSessionIndexKey(userIDHash)
		sessionIDs, err := testClient.SMembers(ctx, userSessionIndexKey).Result()
		require.NoError(t, err)
		assert.Equal(t, 3, len(sessionIDs), "Should have 3 sessions before revocation")

		// Revoke all sessions for user
		err = repo.RevokeAllUserSessions(ctx, userIDHash)
		assert.NoError(t, err)

		// Verify all sessions are removed using raw Redis queries
		for _, sess := range sessions {
			sessionKey := session.FormatSessionKey(sess.ID.String())
			exists, err := testClient.Exists(ctx, sessionKey).Result()
			assert.NoError(t, err)
			assert.Equal(t, int64(0), exists, "Session should be removed")

			// Verify token mappings are removed
			accessTokenKey := session.FormatAccessTokenKey(sess.AccessTokenHash)
			tokenExists, err := testClient.Exists(ctx, accessTokenKey).Result()
			assert.NoError(t, err)
			assert.Equal(t, int64(0), tokenExists, "Access token mapping should be removed")

			refreshTokenKey := session.FormatRefreshTokenKey(sess.RefreshTokenHash)
			refreshExists, err := testClient.Exists(ctx, refreshTokenKey).Result()
			assert.NoError(t, err)
			assert.Equal(t, int64(0), refreshExists, "Refresh token mapping should be removed")
		}

		// Verify user session index is removed
		indexExists, err := testClient.Exists(ctx, userSessionIndexKey).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), indexExists, "User session index should be removed")
	})

	t.Run("should successfully revoke single session", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create single session
		userIDHash := "single_user_hash"
		baseSession := td.NewTestSessionEncx(t)
		baseSession.UserIDHash = userIDHash
		ttl := 1 * time.Hour
		td.InsertSessionEncx(t, ctx, testClient, baseSession, ttl)

		// Revoke all sessions for user
		err := repo.RevokeAllUserSessions(ctx, userIDHash)
		assert.NoError(t, err)

		// Verify session is removed
		sessionKey := session.FormatSessionKey(baseSession.ID.String())
		exists, err := testClient.Exists(ctx, sessionKey).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), exists, "Session should be removed")

		// Verify user session index is removed
		userSessionIndexKey := session.FormatUserSessionIndexKey(userIDHash)
		indexExists, err := testClient.Exists(ctx, userSessionIndexKey).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), indexExists, "User session index should be removed")
	})

	t.Run("should handle gracefully when user has no sessions", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Try to revoke sessions for user with no sessions
		userIDHash := "nonexistent_user_hash"
		err := repo.RevokeAllUserSessions(ctx, userIDHash)
		assert.NoError(t, err, "Revoking non-existent sessions should not error")
	})

	t.Run("should ensure revoking one user's sessions doesn't affect other users", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create sessions for two different users
		userIDHash1 := "user_hash_1"
		userIDHash2 := "user_hash_2"
		ttl := 1 * time.Hour

		// User 1 sessions - with unique token hashes
		user1Sessions := make([]*session.SessionEncx, 2)
		for i := 0; i < 2; i++ {
			baseSession := td.NewTestSessionEncx(t)
			baseSession.UserIDHash = userIDHash1
			// Make token hashes unique for each session
			baseSession.AccessTokenHash = "user1_access_token_" + string(rune('a'+i))
			baseSession.RefreshTokenHash = "user1_refresh_token_" + string(rune('a'+i))
			td.InsertSessionEncx(t, ctx, testClient, baseSession, ttl)
			user1Sessions[i] = baseSession
		}

		// User 2 sessions - with unique token hashes
		user2Sessions := make([]*session.SessionEncx, 2)
		for i := 0; i < 2; i++ {
			baseSession := td.NewTestSessionEncx(t)
			baseSession.UserIDHash = userIDHash2
			// Make token hashes unique for each session
			baseSession.AccessTokenHash = "user2_access_token_" + string(rune('a'+i))
			baseSession.RefreshTokenHash = "user2_refresh_token_" + string(rune('a'+i))
			td.InsertSessionEncx(t, ctx, testClient, baseSession, ttl)
			user2Sessions[i] = baseSession
		}

		// Revoke all sessions for user 1
		err := repo.RevokeAllUserSessions(ctx, userIDHash1)
		assert.NoError(t, err)

		// Verify user 1 sessions are removed
		for _, sess := range user1Sessions {
			sessionKey := session.FormatSessionKey(sess.ID.String())
			exists, err := testClient.Exists(ctx, sessionKey).Result()
			assert.NoError(t, err)
			assert.Equal(t, int64(0), exists, "User 1 session should be removed")
		}

		// Verify user 2 sessions still exist
		for _, sess := range user2Sessions {
			sessionKey := session.FormatSessionKey(sess.ID.String())
			exists, err := testClient.Exists(ctx, sessionKey).Result()
			assert.NoError(t, err)
			assert.Equal(t, int64(1), exists, "User 2 session should still exist")

			// Verify user 2 token mappings still exist
			accessTokenKey := session.FormatAccessTokenKey(sess.AccessTokenHash)
			tokenExists, err := testClient.Exists(ctx, accessTokenKey).Result()
			assert.NoError(t, err)
			assert.Equal(t, int64(1), tokenExists, "User 2 access token mapping should still exist")
		}

		// Verify user 1 index is removed but user 2 index still exists
		user1IndexKey := session.FormatUserSessionIndexKey(userIDHash1)
		user1IndexExists, err := testClient.Exists(ctx, user1IndexKey).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), user1IndexExists, "User 1 session index should be removed")

		user2IndexKey := session.FormatUserSessionIndexKey(userIDHash2)
		user2IndexExists, err := testClient.Exists(ctx, user2IndexKey).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), user2IndexExists, "User 2 session index should still exist")
	})

	t.Run("should handle sessions where session data cannot be decoded", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create a valid session
		userIDHash := "corrupted_user_hash"
		baseSession := td.NewTestSessionEncx(t)
		baseSession.UserIDHash = userIDHash
		ttl := 1 * time.Hour
		td.InsertSessionEncx(t, ctx, testClient, baseSession, ttl)

		// Corrupt the session data by directly setting invalid JSON
		sessionKey := session.FormatSessionKey(baseSession.ID.String())
		err := testClient.Set(ctx, sessionKey, "invalid_json_data", ttl).Err()
		require.NoError(t, err)

		// Add session ID to user index manually since we corrupted the session after insertion
		userSessionIndexKey := session.FormatUserSessionIndexKey(userIDHash)
		err = testClient.SAdd(ctx, userSessionIndexKey, baseSession.ID.String()).Err()
		require.NoError(t, err)

		// Revoke all sessions (should succeed even with corrupted data)
		err = repo.RevokeAllUserSessions(ctx, userIDHash)
		assert.NoError(t, err)

		// Verify session is still removed despite corruption
		exists, err := testClient.Exists(ctx, sessionKey).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), exists, "Corrupted session should be removed")

		// Verify user session index is removed
		indexExists, err := testClient.Exists(ctx, userSessionIndexKey).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), indexExists, "User session index should be removed")
	})

	t.Run("should handle Redis connection failure", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		userIDHash := "connection_failure_user"

		// Close Redis connection to simulate failure
		testClient.Close()

		// Try to revoke sessions (should fail)
		err := repo.RevokeAllUserSessions(ctx, userIDHash)
		assert.Error(t, err)

		// Reconnect for cleanup
		reconnectRedis()
	})

	t.Run("should verify all token mappings are removed", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create sessions with unique token hashes
		userIDHash := "token_cleanup_user"
		sessions := make([]*session.SessionEncx, 3)
		ttl := 1 * time.Hour

		for i := 0; i < 3; i++ {
			baseSession := td.NewTestSessionEncx(t)
			baseSession.UserIDHash = userIDHash
			// Ensure unique token hashes
			baseSession.AccessTokenHash = "access_token_" + string(rune('a'+i))
			baseSession.RefreshTokenHash = "refresh_token_" + string(rune('a'+i))
			td.InsertSessionEncx(t, ctx, testClient, baseSession, ttl)
			sessions[i] = baseSession
		}

		// Verify token mappings exist before revocation
		for _, sess := range sessions {
			accessTokenKey := session.FormatAccessTokenKey(sess.AccessTokenHash)
			exists, err := testClient.Exists(ctx, accessTokenKey).Result()
			require.NoError(t, err)
			require.Equal(t, int64(1), exists, "Access token should exist before revocation")

			refreshTokenKey := session.FormatRefreshTokenKey(sess.RefreshTokenHash)
			refreshExists, err := testClient.Exists(ctx, refreshTokenKey).Result()
			require.NoError(t, err)
			require.Equal(t, int64(1), refreshExists, "Refresh token should exist before revocation")
		}

		// Revoke all sessions
		err := repo.RevokeAllUserSessions(ctx, userIDHash)
		assert.NoError(t, err)

		// Verify all token mappings are removed
		for _, sess := range sessions {
			accessTokenKey := session.FormatAccessTokenKey(sess.AccessTokenHash)
			exists, err := testClient.Exists(ctx, accessTokenKey).Result()
			assert.NoError(t, err)
			assert.Equal(t, int64(0), exists, "Access token should be removed")

			refreshTokenKey := session.FormatRefreshTokenKey(sess.RefreshTokenHash)
			refreshExists, err := testClient.Exists(ctx, refreshTokenKey).Result()
			assert.NoError(t, err)
			assert.Equal(t, int64(0), refreshExists, "Refresh token should be removed")
		}
	})
}
