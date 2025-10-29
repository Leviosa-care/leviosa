package sessionRepository_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestFindSessionByID TEST_PATH=internal/authuser/infrastructure/redis/session/find_session_by_id_test.go

func TestFindSessionByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully find existing session", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create test baseSession directly in Redis
		baseSession := td.NewTestSessionEncx(t)

		ttl := 1 * time.Hour
		td.InsertSessionEncx(t, ctx, testClient, baseSession, ttl)

		// Find session by ID
		sessionData, err := repo.FindSessionByID(ctx, baseSession.ID)
		assert.NoError(t, err)

		var retrievedSessionEncx session.SessionEncx
		err = json.Unmarshal(sessionData, &retrievedSessionEncx)
		assert.NoError(t, err)

		// Verify session data integrity
		assert.Equal(t, baseSession.UserIDHash, retrievedSessionEncx.UserIDHash)
		assert.Equal(t, baseSession.UserIDEncrypted, retrievedSessionEncx.UserIDEncrypted)
		assert.Equal(t, baseSession.RoleEncrypted, retrievedSessionEncx.RoleEncrypted)
		assert.Equal(t, baseSession.StateEncrypted, retrievedSessionEncx.StateEncrypted)
		assert.Equal(t, baseSession.CreatedAtEncrypted, retrievedSessionEncx.CreatedAtEncrypted)
		assert.Equal(t, baseSession.ExpiresAtEncrypted, retrievedSessionEncx.ExpiresAtEncrypted)
		assert.Equal(t, baseSession.AccessTokenHash, retrievedSessionEncx.AccessTokenHash)
		assert.Equal(t, baseSession.RefreshTokenHash, retrievedSessionEncx.RefreshTokenHash)
	})

	t.Run("should return not found error for non-existent session", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		nonExistentID := uuid.New()

		// Try to find non-existent session
		sessionData, err := repo.FindSessionByID(ctx, nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, sessionData)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("should handle Redis connection failure", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		sessionID := uuid.New()

		// Close Redis connection to simulate failure
		testClient.Close()

		// Try to find session (should fail)
		sessionData, err := repo.FindSessionByID(ctx, sessionID)
		assert.Error(t, err)
		assert.Nil(t, sessionData)

		// Reconnect for cleanup
		reconnectRedis()
	})

	t.Run("should find session with empty session data", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		sessionID := uuid.New()
		// sessionKey := "authuser:session:" + sessionID.String()
		sessionKey := session.FormatSessionKey(sessionID.String())

		// Insert empty data directly
		err := testClient.Set(ctx, sessionKey, []byte{}, 1*time.Hour).Err()
		require.NoError(t, err)

		// Find session
		sessionData, err := repo.FindSessionByID(ctx, sessionID)
		require.NoError(t, err)
		assert.Equal(t, []byte{}, sessionData)
	})

	t.Run("should find session near expiration", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)
		baseSession := td.NewTestSessionEncx(t)

		// Create session with short TTL
		ttl := 200 * time.Millisecond
		td.InsertSessionEncx(t, ctx, testClient, baseSession, ttl)

		// Find session immediately (should exist)
		sessionData, err := repo.FindSessionByID(ctx, baseSession.ID)
		require.NoError(t, err)
		assert.NotNil(t, sessionData)

		// Wait for expiration
		time.Sleep(250 * time.Millisecond)

		// Try to find expired session (should not exist)
		expiredSessionData, err := repo.FindSessionByID(ctx, baseSession.ID)
		assert.Error(t, err)
		assert.Nil(t, expiredSessionData)
	})

	t.Run("should find multiple different sessions", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		const count = 3

		// Create multiple sessions
		sessions := make([]*td.SessionTestData, count)
		ttl := 1 * time.Hour

		for i := range count {
			baseSession := td.NewTestSessionEncx(t)
			td.InsertSessionEncx(t, ctx, testClient, baseSession, ttl)

			sessions[i] = &td.SessionTestData{
				SessionID:   baseSession.ID,
				TokenHash:   baseSession.AccessTokenHash,
				Session:     baseSession,
				SessionData: td.EncodeSession(t, baseSession),
			}
		}

		// Find each session and verify data
		for i, sessionData := range sessions {
			retrievedData, err := repo.FindSessionByID(ctx, sessionData.SessionID)
			assert.NoError(t, err, "Should find session %d", i)

			var retrievedSessionEncx session.SessionEncx
			err = json.Unmarshal(retrievedData, &retrievedSessionEncx)
			assert.NoError(t, err)

			assert.Equal(t, sessionData.Session.UserIDEncrypted, retrievedSessionEncx.UserIDEncrypted, "Session %d UserID should match", i)
			assert.Equal(t, sessionData.Session.RoleEncrypted, retrievedSessionEncx.RoleEncrypted, "Session %d Role should match", i)
			assert.Equal(t, sessionData.Session.StateEncrypted, retrievedSessionEncx.StateEncrypted, "Session %d State should match", i)
		}
	})
}
