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

// make test-func TEST_NAME=TestRefreshTokenPair TEST_PATH=internal/authuser/infrastructure/redis/session/refresh_token_pair_test.go

func TestRefreshTokenPair(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully refresh both access and refresh tokens", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create initial token pair
		baseSession := td.NewTestSessionEncx(t)
		oldAccessTokenHash := "old_access_token_hash"
		baseSession.AccessTokenHash = oldAccessTokenHash
		refreshTokenHash := "refresh_token_hash"
		baseSession.RefreshTokenHash = refreshTokenHash

		accessTTL := 1 * time.Hour
		refreshTTL := 24 * time.Hour

		td.InsertSessionEncx(t, ctx, testClient, baseSession, accessTTL)

		// Refresh with new access and refresh tokens
		newAccessTokenHash := "new_access_token_hash"
		newRefreshTokenHash := "new_refresh_token_hash"

		// Update session with new token hashes for test
		baseSession.AccessTokenHash = newAccessTokenHash
		baseSession.RefreshTokenHash = newRefreshTokenHash
		updatedSessionData := td.EncodeSession(t, baseSession)

		err := repo.RefreshTokenPair(ctx, refreshTokenHash, newAccessTokenHash, newRefreshTokenHash, baseSession.ID, updatedSessionData, accessTTL, refreshTTL)
		assert.NoError(t, err)

		// Old access token should be removed
		_, err = td.FindSessionByAccessTokenHash(t, ctx, oldAccessTokenHash, testClient)
		assert.Error(t, err, "Old access token should be removed")

		// New access token should work
		sessionIDFound, err := td.FindSessionByAccessTokenHash(t, ctx, newAccessTokenHash, testClient)
		assert.NoError(t, err)
		assert.Equal(t, []byte(baseSession.ID.String()), sessionIDFound)

		// Old refresh token should be removed
		_, err = td.FindSessionByRefreshTokenHash(t, ctx, refreshTokenHash, testClient)
		assert.Error(t, err, "Old refresh token should be removed")

		// New refresh token should work
		sessionIDFound, err = td.FindSessionByRefreshTokenHash(t, ctx, newRefreshTokenHash, testClient)
		assert.NoError(t, err)
		assert.Equal(t, []byte(baseSession.ID.String()), sessionIDFound)
	})

	t.Run("should clean up stale access tokens for the same session", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create initial token pair
		baseSession := td.NewTestSessionEncx(t)
		oldAccessTokenHash1 := "old_access_token_1"
		oldAccessTokenHash2 := "old_access_token_2"
		baseSession.AccessTokenHash = oldAccessTokenHash1
		oldRefreshTokenHash := "old_refresh_token"
		baseSession.RefreshTokenHash = oldRefreshTokenHash

		accessTTL := 1 * time.Hour
		refreshTTL := 24 * time.Hour

		// Create multiple access tokens for the same session (simulating multiple devices)
		td.InsertSessionEncx(t, ctx, testClient, baseSession, accessTTL)

		// Manually create additional access token for same session
		accessTokenKey2 := session.FormatAccessTokenKey(oldAccessTokenHash2)
		err := testClient.Set(ctx, accessTokenKey2, baseSession.ID.String(), accessTTL).Err()
		require.NoError(t, err)

		// Verify both old access tokens work initially
		_, err = td.FindSessionByAccessTokenHash(t, ctx, oldAccessTokenHash1, testClient)
		require.NoError(t, err)
		_, err = td.FindSessionByAccessTokenHash(t, ctx, oldAccessTokenHash2, testClient)
		require.NoError(t, err)

		// Refresh tokens
		newAccessTokenHash := "new_access_token"
		newRefreshTokenHash := "new_refresh_token"

		// Update session with new token hashes for test
		baseSession.AccessTokenHash = newAccessTokenHash
		baseSession.RefreshTokenHash = newRefreshTokenHash
		updatedSessionData := td.EncodeSession(t, baseSession)

		err = repo.RefreshTokenPair(ctx, oldRefreshTokenHash, newAccessTokenHash, newRefreshTokenHash, baseSession.ID, updatedSessionData, accessTTL, refreshTTL)
		assert.NoError(t, err)

		// All old access tokens should be cleaned up
		_, err = td.FindSessionByAccessTokenHash(t, ctx, oldAccessTokenHash1, testClient)
		assert.Error(t, err, "Old access token 1 should be removed")
		_, err = td.FindSessionByAccessTokenHash(t, ctx, oldAccessTokenHash2, testClient)
		assert.Error(t, err, "Old access token 2 should be removed")

		// Old refresh token should be removed
		_, err = td.FindSessionByRefreshTokenHash(t, ctx, oldRefreshTokenHash, testClient)
		assert.Error(t, err, "Old refresh token should be removed")

		// New tokens should work
		sessionDataFromNewAccess, err := td.FindSessionByAccessTokenHash(t, ctx, newAccessTokenHash, testClient)
		assert.NoError(t, err)
		assert.Equal(t, []byte(baseSession.ID.String()), sessionDataFromNewAccess)

		sessionDataFromNewRefresh, err := td.FindSessionByRefreshTokenHash(t, ctx, newRefreshTokenHash, testClient)
		assert.NoError(t, err)
		assert.Equal(t, []byte(baseSession.ID.String()), sessionDataFromNewRefresh)
	})

	t.Run("should rollback on refresh token creation failure", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create initial token pair
		// session := td.CreateTestSessionWithCrypto(t, crypto)
		baseSession := td.NewTestSessionEncx(t)
		oldAccessTokenHash := "rollback_test_old_access"
		baseSession.AccessTokenHash = oldAccessTokenHash
		oldRefreshTokenHash := "rollback_test_old_refresh"
		baseSession.RefreshTokenHash = oldRefreshTokenHash

		td.InsertSessionEncx(t, ctx, testClient, baseSession, 1*time.Hour)
		// Verify initial tokens work
		_, err := td.FindSessionByAccessTokenHash(t, ctx, oldAccessTokenHash, testClient)
		require.NoError(t, err)

		// Close Redis to simulate failure
		testClient.Close()

		// Try to refresh tokens - should fail
		newAccessTokenHash := "rollback_test_new_access"
		newRefreshTokenHash := "rollback_test_new_refresh"

		// Update session with new token hashes for test
		baseSession.AccessTokenHash = newAccessTokenHash
		baseSession.RefreshTokenHash = newRefreshTokenHash
		updatedSessionData := td.EncodeSession(t, baseSession)

		err = repo.RefreshTokenPair(ctx, oldRefreshTokenHash, newAccessTokenHash, newRefreshTokenHash, baseSession.ID, updatedSessionData, time.Hour, 24*time.Hour)
		assert.Error(t, err)

		// Reconnect and verify state
		reconnectRedis()

		// Old tokens should still work (no partial updates)
		_, err = td.FindSessionByAccessTokenHash(t, ctx, oldAccessTokenHash, testClient)
		assert.NoError(t, err, "Old access token should still work after failed refresh")

		// New tokens should not exist
		_, err = td.FindSessionByAccessTokenHash(t, ctx, newAccessTokenHash, testClient)
		assert.Error(t, err, "New access token should not exist after failed refresh")
	})

	t.Run("should handle TTL updates correctly", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create initial token pair with short TTL
		baseSession := td.NewTestSessionEncx(t)
		oldAccessTokenHash := "ttl_test_old_access"
		baseSession.AccessTokenHash = oldAccessTokenHash
		refreshTokenHash := "ttl_test_refresh"
		baseSession.RefreshTokenHash = refreshTokenHash

		shortTTL := 1 * time.Second
		longTTL := 1 * time.Hour

		td.InsertSessionEncx(t, ctx, testClient, baseSession, shortTTL)

		// Refresh with longer TTL
		newAccessTokenHash := "ttl_test_new_access"
		newRefreshTokenHash := "ttl_test_new_refresh"
		newTTL := 2 * time.Hour

		// Update session with new token hashes for test
		baseSession.AccessTokenHash = newAccessTokenHash
		baseSession.RefreshTokenHash = newRefreshTokenHash
		updatedSessionData := td.EncodeSession(t, baseSession)

		err := repo.RefreshTokenPair(ctx, refreshTokenHash, newAccessTokenHash, newRefreshTokenHash, baseSession.ID, updatedSessionData, newTTL, longTTL)
		assert.NoError(t, err)

		// Verify new access token has longer TTL
		newAccessTokenKey := session.FormatAccessTokenKey(newAccessTokenHash)
		actualTTL := testClient.TTL(ctx, newAccessTokenKey).Val()
		assert.True(t, actualTTL > 1*time.Hour, "New access token should have longer TTL")
	})

	t.Run("should handle special characters in token hashes", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create tokens with special characters
		baseSession := td.NewTestSessionEncx(t)
		oldAccessTokenHash := "old-access_token.with:special@chars+123/456="
		baseSession.AccessTokenHash = oldAccessTokenHash
		oldRefreshTokenHash := "old-refresh_token.with:special@chars+789/012="
		baseSession.RefreshTokenHash = oldRefreshTokenHash

		td.InsertSessionEncx(t, ctx, testClient, baseSession, 1*time.Hour)

		// Refresh with new tokens also having special characters
		newAccessTokenHash := "new-access_token.with:special@chars+999/888="
		newRefreshTokenHash := "new-refresh_token.with:special@chars+777/666="

		// Update session with new token hashes for test
		baseSession.AccessTokenHash = newAccessTokenHash
		baseSession.RefreshTokenHash = newRefreshTokenHash
		updatedSessionData := td.EncodeSession(t, baseSession)

		err := repo.RefreshTokenPair(ctx, oldRefreshTokenHash, newAccessTokenHash, newRefreshTokenHash, baseSession.ID, updatedSessionData, time.Hour, 24*time.Hour)
		assert.NoError(t, err)

		// Verify new tokens work
		_, err = td.FindSessionByAccessTokenHash(t, ctx, newAccessTokenHash, testClient)
		assert.NoError(t, err)

		_, err = td.FindSessionByRefreshTokenHash(t, ctx, newRefreshTokenHash, testClient)
		assert.NoError(t, err)
	})

	t.Run("should handle Redis connection errors gracefully", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Close Redis connection
		testClient.Close()

		// Try to refresh tokens with closed connection
		sessionID := uuid.New()
		// For this error case, we just need some valid session data
		updatedSessionData := []byte("{\"id\":\"test\",\"access_token_hash\":\"new_access\",\"refresh_token_hash\":\"new_refresh\"}")
		err := repo.RefreshTokenPair(ctx, "old_refresh", "new_access", "new_refresh", sessionID, updatedSessionData, time.Hour, 24*time.Hour)
		assert.Error(t, err)

		// Reconnect for subsequent tests
		reconnectRedis()
	})
}
