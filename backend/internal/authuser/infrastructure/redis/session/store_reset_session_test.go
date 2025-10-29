package sessionRepository_test

import (
	"context"
	"testing"
	"time"

	sessionRepository "github.com/Leviosa-care/leviosa/backend/internal/authuser/infrastructure/redis/session"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestStoreResetSession TEST_PATH=internal/authuser/infrastructure/redis/session/store_reset_session_test.go

func TestStoreResetSession(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully store reset session with valid TTL", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "test_token_hash_12345"
		userEmail := "test@example.com"
		ttl := 15 * time.Minute

		// Store reset session
		err := repo.StoreResetSession(ctx, tokenHash, userEmail, ttl)
		assert.NoError(t, err)

		// Verify stored data can be retrieved using raw Redis query
		key := sessionRepository.FormatResetSessionKey(tokenHash)
		storedEmail, err := testClient.Get(ctx, key).Result()
		assert.NoError(t, err)
		assert.Equal(t, userEmail, storedEmail, "Stored email should match")
	})

	t.Run("should verify TTL is set correctly", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "test_token_ttl_hash"
		userEmail := "ttl-test@example.com"
		ttl := 10 * time.Minute

		// Store reset session
		err := repo.StoreResetSession(ctx, tokenHash, userEmail, ttl)
		require.NoError(t, err)

		// Verify TTL is within acceptable range using raw Redis query
		key := sessionRepository.FormatResetSessionKey(tokenHash)
		actualTTL := testClient.TTL(ctx, key).Val()

		// TTL should be slightly less than or equal to set TTL (due to time elapsed)
		lowerBound := 9*time.Minute + 50*time.Second // Allow 10 seconds margin
		upperBound := ttl

		assert.True(t, actualTTL > lowerBound && actualTTL <= upperBound,
			"TTL should be within expected range, got %v", actualTTL)
	})

	t.Run("should overwrite existing reset token", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "overwrite_token_hash"
		firstEmail := "first@example.com"
		secondEmail := "second@example.com"
		ttl := 15 * time.Minute

		// Store first reset session
		err := repo.StoreResetSession(ctx, tokenHash, firstEmail, ttl)
		require.NoError(t, err)

		// Verify first email is stored
		key := sessionRepository.FormatResetSessionKey(tokenHash)
		storedEmail, err := testClient.Get(ctx, key).Result()
		require.NoError(t, err)
		require.Equal(t, firstEmail, storedEmail)

		// Store second reset session with same token hash
		err = repo.StoreResetSession(ctx, tokenHash, secondEmail, ttl)
		assert.NoError(t, err)

		// Verify second email overwrote the first
		storedEmail, err = testClient.Get(ctx, key).Result()
		assert.NoError(t, err)
		assert.Equal(t, secondEmail, storedEmail, "Second email should overwrite first")
	})

	t.Run("should handle special characters in token hash and email", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Base64-encoded token hash with special characters
		tokenHash := "aB3+/cD=123_xyz"
		userEmail := "user+tag@sub.example.com"
		ttl := 15 * time.Minute

		// Store reset session
		err := repo.StoreResetSession(ctx, tokenHash, userEmail, ttl)
		assert.NoError(t, err)

		// Verify stored data using raw Redis query
		key := sessionRepository.FormatResetSessionKey(tokenHash)
		storedEmail, err := testClient.Get(ctx, key).Result()
		assert.NoError(t, err)
		assert.Equal(t, userEmail, storedEmail, "Email with special characters should be stored correctly")
	})

	t.Run("should handle very short TTL", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "short_ttl_token"
		userEmail := "shortttl@example.com"
		ttl := 2 * time.Second

		// Store reset session with short TTL
		err := repo.StoreResetSession(ctx, tokenHash, userEmail, ttl)
		assert.NoError(t, err)

		// Verify data is stored
		key := sessionRepository.FormatResetSessionKey(tokenHash)
		storedEmail, err := testClient.Get(ctx, key).Result()
		assert.NoError(t, err)
		assert.Equal(t, userEmail, storedEmail)

		// Wait for expiration
		time.Sleep(3 * time.Second)

		// Verify key has expired
		exists, err := testClient.Exists(ctx, key).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), exists, "Key should have expired")
	})

	t.Run("should handle empty email string", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "empty_email_token"
		userEmail := ""
		ttl := 15 * time.Minute

		// Store reset session with empty email
		err := repo.StoreResetSession(ctx, tokenHash, userEmail, ttl)
		assert.NoError(t, err)

		// Verify empty string is stored
		key := sessionRepository.FormatResetSessionKey(tokenHash)
		storedEmail, err := testClient.Get(ctx, key).Result()
		assert.NoError(t, err)
		assert.Equal(t, "", storedEmail, "Empty email should be stored")
	})

	t.Run("should handle long email addresses", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "long_email_token"
		// Create a very long but valid email
		userEmail := "very.long.email.address.with.many.parts.and.subdomains@subdomain.example.company.organization.com"
		ttl := 15 * time.Minute

		// Store reset session
		err := repo.StoreResetSession(ctx, tokenHash, userEmail, ttl)
		assert.NoError(t, err)

		// Verify long email is stored correctly
		key := sessionRepository.FormatResetSessionKey(tokenHash)
		storedEmail, err := testClient.Get(ctx, key).Result()
		assert.NoError(t, err)
		assert.Equal(t, userEmail, storedEmail, "Long email should be stored correctly")
	})

	t.Run("should handle Redis connection failure", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "connection_failure_token"
		userEmail := "failure@example.com"
		ttl := 15 * time.Minute

		// Close Redis connection to simulate failure
		testClient.Close()

		// Try to store reset session (should fail)
		err := repo.StoreResetSession(ctx, tokenHash, userEmail, ttl)
		assert.Error(t, err)

		// Reconnect for cleanup
		reconnectRedis()
	})

	t.Run("should handle multiple reset sessions for different users", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Store multiple reset sessions
		sessions := map[string]string{
			"token_hash_1": "user1@example.com",
			"token_hash_2": "user2@example.com",
			"token_hash_3": "user3@example.com",
		}
		ttl := 15 * time.Minute

		// Store all reset sessions
		for tokenHash, email := range sessions {
			err := repo.StoreResetSession(ctx, tokenHash, email, ttl)
			require.NoError(t, err)
		}

		// Verify all sessions are stored correctly using raw Redis queries
		for tokenHash, expectedEmail := range sessions {
			key := sessionRepository.FormatResetSessionKey(tokenHash)
			storedEmail, err := testClient.Get(ctx, key).Result()
			assert.NoError(t, err)
			assert.Equal(t, expectedEmail, storedEmail, "Email for token %s should match", tokenHash)
		}
	})
}
