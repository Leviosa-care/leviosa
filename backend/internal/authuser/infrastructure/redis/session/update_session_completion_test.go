package sessionRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestUpdateSessionCompletion TEST_PATH=internal/authuser/infrastructure/redis/session/update_session_completion_test.go

func TestUpdateSessionCompletion(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully update session data while preserving TTL", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create initial session
		baseSession := td.NewTestSessionEncx(t)
		initialTTL := 1 * time.Hour
		td.InsertSessionEncx(t, ctx, testClient, baseSession, initialTTL)

		// Get TTL before update using raw Redis query
		sessionKey := session.FormatSessionKey(baseSession.ID.String())
		ttlBefore := testClient.TTL(ctx, sessionKey).Val()
		require.True(t, ttlBefore > 0, "TTL should be positive before update")

		// Update session with new data
		updatedSession := td.NewTestSessionEncx(t)
		updatedSession.ID = baseSession.ID // Keep same ID
		updatedSession.UserIDHash = "updated_user_hash"
		sessionEncoded := td.EncodeSessionEncx(t, updatedSession)

		err := repo.UpdateSessionCompletion(ctx, baseSession.ID, sessionEncoded)
		assert.NoError(t, err)

		// Verify data was updated using raw Redis query
		storedData, err := testClient.Get(ctx, sessionKey).Bytes()
		assert.NoError(t, err)
		decodedSession := td.DecodeSessionEncx(t, storedData)
		assert.Equal(t, "updated_user_hash", decodedSession.UserIDHash, "Session data should be updated")

		// Verify TTL is preserved (allow 1 second margin for execution time)
		ttlAfter := testClient.TTL(ctx, sessionKey).Val()
		timeDiff := ttlBefore - ttlAfter
		assert.True(t, timeDiff < 2*time.Second,
			"TTL should be preserved within 2 seconds, before: %v, after: %v, diff: %v",
			ttlBefore, ttlAfter, timeDiff)
	})

	t.Run("should return ErrRepositoryNotFound for non-existent session", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		nonExistentID := uuid.New()
		sessionEncoded := []byte(`{"id":"test"}`)

		// Verify session doesn't exist before attempting update
		sessionKey := session.FormatSessionKey(nonExistentID.String())
		exists, err := testClient.Exists(ctx, sessionKey).Result()
		require.NoError(t, err)
		require.Equal(t, int64(0), exists, "Session should not exist before update attempt")

		// Check what TTL returns for non-existent key
		// go-redis returns -2 nanoseconds (not seconds) for non-existent keys
		ttl := testClient.TTL(ctx, sessionKey).Val()
		require.Equal(t, time.Duration(-2), ttl, "TTL should be -2 nanoseconds for non-existent key")

		// Try to update non-existent session
		err = repo.UpdateSessionCompletion(ctx, nonExistentID, sessionEncoded)
		assert.Error(t, err, "Should return error for non-existent session")
		if err != nil {
			assert.ErrorIs(t, err, errs.ErrRepositoryNotFound, "Should return ErrRepositoryNotFound")
		}
	})

	t.Run("should update session from pending to active state", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create pending session
		baseSession := td.NewTestSessionEncx(t)
		baseSession.StateEncrypted = []byte("state_pending_encrypted")
		ttl := 1 * time.Hour
		td.InsertSessionEncx(t, ctx, testClient, baseSession, ttl)

		// Update to active state
		updatedSession := td.NewTestSessionEncx(t)
		updatedSession.ID = baseSession.ID
		updatedSession.StateEncrypted = []byte("state_active_encrypted")
		updatedSession.UserIDHash = baseSession.UserIDHash
		updatedSession.AccessTokenHash = baseSession.AccessTokenHash
		updatedSession.RefreshTokenHash = baseSession.RefreshTokenHash
		sessionEncoded := td.EncodeSessionEncx(t, updatedSession)

		err := repo.UpdateSessionCompletion(ctx, baseSession.ID, sessionEncoded)
		assert.NoError(t, err)

		// Verify state was updated using raw Redis query
		sessionKey := session.FormatSessionKey(baseSession.ID.String())
		storedData, err := testClient.Get(ctx, sessionKey).Bytes()
		assert.NoError(t, err)
		decodedSession := td.DecodeSessionEncx(t, storedData)
		assert.Equal(t, []byte("state_active_encrypted"), decodedSession.StateEncrypted, "State should be updated to active")
	})

	t.Run("should handle multiple consecutive updates", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create initial session
		baseSession := td.NewTestSessionEncx(t)
		ttl := 1 * time.Hour
		td.InsertSessionEncx(t, ctx, testClient, baseSession, ttl)

		// Perform multiple updates
		for i := 0; i < 3; i++ {
			updatedSession := td.NewTestSessionEncx(t)
			updatedSession.ID = baseSession.ID
			updatedSession.UserIDHash = "user_hash_" + string(rune('a'+i))
			sessionEncoded := td.EncodeSessionEncx(t, updatedSession)

			err := repo.UpdateSessionCompletion(ctx, baseSession.ID, sessionEncoded)
			require.NoError(t, err, "Update %d should succeed", i+1)

			// Verify update using raw Redis query
			sessionKey := session.FormatSessionKey(baseSession.ID.String())
			storedData, err := testClient.Get(ctx, sessionKey).Bytes()
			require.NoError(t, err)
			decodedSession := td.DecodeSessionEncx(t, storedData)
			expectedHash := "user_hash_" + string(rune('a'+i))
			assert.Equal(t, expectedHash, decodedSession.UserIDHash, "Update %d should persist", i+1)
		}
	})

	t.Run("should verify data integrity after update", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create initial session
		baseSession := td.NewTestSessionEncx(t)
		ttl := 1 * time.Hour
		td.InsertSessionEncx(t, ctx, testClient, baseSession, ttl)

		// Create updated session with all fields populated
		updatedSession := &session.SessionEncx{
			ID:                 baseSession.ID,
			UserIDEncrypted:    []byte("new_user_id_encrypted"),
			UserIDHash:         "new_user_hash",
			RoleEncrypted:      []byte("new_role_encrypted"),
			StateEncrypted:     []byte("new_state_encrypted"),
			CreatedAtEncrypted: []byte("new_created_at_encrypted"),
			ExpiresAtEncrypted: []byte("new_expires_at_encrypted"),
			AccessTokenHash:    "new_access_token_hash",
			RefreshTokenHash:   "new_refresh_token_hash",
		}
		sessionEncoded := td.EncodeSessionEncx(t, updatedSession)

		err := repo.UpdateSessionCompletion(ctx, baseSession.ID, sessionEncoded)
		assert.NoError(t, err)

		// Verify all fields are updated correctly using raw Redis query
		sessionKey := session.FormatSessionKey(baseSession.ID.String())
		storedData, err := testClient.Get(ctx, sessionKey).Bytes()
		assert.NoError(t, err)
		decodedSession := td.DecodeSessionEncx(t, storedData)

		assert.Equal(t, updatedSession.ID, decodedSession.ID)
		assert.Equal(t, updatedSession.UserIDEncrypted, decodedSession.UserIDEncrypted)
		assert.Equal(t, updatedSession.UserIDHash, decodedSession.UserIDHash)
		assert.Equal(t, updatedSession.RoleEncrypted, decodedSession.RoleEncrypted)
		assert.Equal(t, updatedSession.StateEncrypted, decodedSession.StateEncrypted)
		assert.Equal(t, updatedSession.CreatedAtEncrypted, decodedSession.CreatedAtEncrypted)
		assert.Equal(t, updatedSession.ExpiresAtEncrypted, decodedSession.ExpiresAtEncrypted)
		assert.Equal(t, updatedSession.AccessTokenHash, decodedSession.AccessTokenHash)
		assert.Equal(t, updatedSession.RefreshTokenHash, decodedSession.RefreshTokenHash)
	})

	t.Run("should handle large session data updates", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create initial session
		baseSession := td.NewTestSessionEncx(t)
		ttl := 1 * time.Hour
		td.InsertSessionEncx(t, ctx, testClient, baseSession, ttl)

		// Create updated session with large encrypted data
		updatedSession := td.NewTestSessionEncx(t)
		updatedSession.ID = baseSession.ID
		// Simulate large encrypted payload
		largeData := make([]byte, 1024)
		for i := range largeData {
			largeData[i] = byte(i % 256)
		}
		updatedSession.UserIDEncrypted = largeData
		sessionEncoded := td.EncodeSessionEncx(t, updatedSession)

		err := repo.UpdateSessionCompletion(ctx, baseSession.ID, sessionEncoded)
		assert.NoError(t, err)

		// Verify large data is stored correctly using raw Redis query
		sessionKey := session.FormatSessionKey(baseSession.ID.String())
		storedData, err := testClient.Get(ctx, sessionKey).Bytes()
		assert.NoError(t, err)
		decodedSession := td.DecodeSessionEncx(t, storedData)
		assert.Equal(t, largeData, decodedSession.UserIDEncrypted, "Large data should be stored correctly")
	})

	t.Run("should handle empty session data update", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create initial session
		baseSession := td.NewTestSessionEncx(t)
		ttl := 1 * time.Hour
		td.InsertSessionEncx(t, ctx, testClient, baseSession, ttl)

		// Update with minimal data
		emptyData := []byte("{}")

		err := repo.UpdateSessionCompletion(ctx, baseSession.ID, emptyData)
		assert.NoError(t, err)

		// Verify data was updated using raw Redis query
		sessionKey := session.FormatSessionKey(baseSession.ID.String())
		storedData, err := testClient.Get(ctx, sessionKey).Bytes()
		assert.NoError(t, err)
		assert.Equal(t, emptyData, storedData, "Empty data should be stored")
	})

	t.Run("should handle Redis connection failure", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		sessionID := uuid.New()
		sessionEncoded := []byte(`{"id":"test"}`)

		// Close Redis connection to simulate failure
		testClient.Close()

		// Try to update session (should fail)
		err := repo.UpdateSessionCompletion(ctx, sessionID, sessionEncoded)
		assert.Error(t, err)

		// Reconnect for cleanup
		reconnectRedis()
	})

	t.Run("should handle session near expiration", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create session with very short TTL
		baseSession := td.NewTestSessionEncx(t)
		shortTTL := 5 * time.Second
		td.InsertSessionEncx(t, ctx, testClient, baseSession, shortTTL)

		// Wait a bit but not until expiration
		time.Sleep(2 * time.Second)

		// Update session
		updatedSession := td.NewTestSessionEncx(t)
		updatedSession.ID = baseSession.ID
		updatedSession.UserIDHash = "updated_near_expiration"
		sessionEncoded := td.EncodeSessionEncx(t, updatedSession)

		err := repo.UpdateSessionCompletion(ctx, baseSession.ID, sessionEncoded)
		assert.NoError(t, err)

		// Verify data was updated using raw Redis query
		sessionKey := session.FormatSessionKey(baseSession.ID.String())
		storedData, err := testClient.Get(ctx, sessionKey).Bytes()
		assert.NoError(t, err)
		decodedSession := td.DecodeSessionEncx(t, storedData)
		assert.Equal(t, "updated_near_expiration", decodedSession.UserIDHash)

		// Verify TTL is preserved (should be ~3 seconds remaining)
		ttlRemaining := testClient.TTL(ctx, sessionKey).Val()
		assert.True(t, ttlRemaining > 0 && ttlRemaining <= 4*time.Second,
			"TTL should be preserved, got %v", ttlRemaining)
	})
}
