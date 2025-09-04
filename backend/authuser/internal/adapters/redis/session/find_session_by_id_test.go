package sessionRepository_test

import (
	"context"
	"testing"
	"time"

	td "github.com/Leviosa-care/authuser/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestFindSessionByID make test-unit-session-test

func TestFindSessionByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully find existing session", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create test session directly in Redis
		session := td.CreateTestSessionWithCrypto(t, crypto)
		ttl := 1 * time.Hour
		td.InsertSessionDirectly(t, ctx, testClient, session, ttl)

		// Find session by ID
		sessionData, err := repo.FindSessionByID(ctx, session.ID.String())
		require.NoError(t, err)

		// Verify session data integrity
		retrievedSession := td.DecodeSessionWithDecryption(t, sessionData, crypto)
		assert.Equal(t, session.UserID, retrievedSession.UserID)
		assert.Equal(t, session.Role, retrievedSession.Role)
		assert.Equal(t, session.State, retrievedSession.State)
		assert.Equal(t, session.TokenHash, retrievedSession.TokenHash)
	})

	t.Run("should return not found error for non-existent session", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		nonExistentID := uuid.New().String()

		// Try to find non-existent session
		sessionData, err := repo.FindSessionByID(ctx, nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, sessionData)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("should handle Redis connection failure", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		sessionID := uuid.New().String()

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
		sessionKey := "authuser:session:" + sessionID.String()

		// Insert empty data directly
		err := testClient.Set(ctx, sessionKey, []byte{}, 1*time.Hour).Err()
		require.NoError(t, err)

		// Find session
		sessionData, err := repo.FindSessionByID(ctx, sessionID.String())
		require.NoError(t, err)
		assert.Equal(t, []byte{}, sessionData)
	})

	t.Run("should find session near expiration", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create session with short TTL
		session := td.CreateTestSessionWithCrypto(t, crypto)
		ttl := 200 * time.Millisecond
		td.InsertSessionDirectly(t, ctx, testClient, session, ttl)

		// Find session immediately (should exist)
		sessionData, err := repo.FindSessionByID(ctx, session.ID.String())
		require.NoError(t, err)
		assert.NotNil(t, sessionData)

		// Wait for expiration
		time.Sleep(250 * time.Millisecond)

		// Try to find expired session (should not exist)
		expiredSessionData, err := repo.FindSessionByID(ctx, session.ID.String())
		assert.Error(t, err)
		assert.Nil(t, expiredSessionData)
	})

	t.Run("should find multiple different sessions", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create multiple sessions
		sessions := make([]*td.SessionTestData, 3)
		ttl := 1 * time.Hour

		for i := range 3 {
			session := td.CreateTestSessionWithCrypto(t, crypto)
			td.InsertSessionDirectly(t, ctx, testClient, session, ttl)

			sessions[i] = &td.SessionTestData{
				SessionID:   session.ID,
				TokenHash:   session.TokenHash,
				Session:     session,
				SessionData: td.EncodeSession(t, session),
			}
		}

		// Find each session and verify data
		for i, sessionData := range sessions {
			retrievedData, err := repo.FindSessionByID(ctx, sessionData.SessionID.String())
			require.NoError(t, err, "Should find session %d", i)

			retrievedSession := td.DecodeSessionWithDecryption(t, retrievedData, crypto)
			assert.Equal(t, sessionData.Session.UserID, retrievedSession.UserID, "Session %d UserID should match", i)
			assert.Equal(t, sessionData.Session.Role, retrievedSession.Role, "Session %d Role should match", i)
			assert.Equal(t, sessionData.Session.State, retrievedSession.State, "Session %d State should match", i)
		}
	})

	t.Run("should handle malformed session ID gracefully", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		malformedIDs := []string{
			"",
			"not-a-uuid",
			"12345",
			"invalid-session-id-format",
		}

		for _, malformedID := range malformedIDs {
			sessionData, err := repo.FindSessionByID(ctx, malformedID)
			assert.Error(t, err, "Should return error for malformed ID: %s", malformedID)
			assert.Nil(t, sessionData, "Should return nil data for malformed ID: %s", malformedID)
		}
	})
}
