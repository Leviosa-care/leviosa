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

// make test-func TEST_NAME=TestValidateResetSession TEST_PATH=internal/authuser/infrastructure/redis/session/validate_reset_session_test.go

func TestValidateResetSession(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully validate and consume reset token", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "valid_token_hash_123"
		userEmail := "valid@example.com"
		ttl := 15 * time.Minute

		// Store reset session
		key := sessionRepository.FormatResetSessionKey(tokenHash)
		err := testClient.Set(ctx, key, userEmail, ttl).Err()
		require.NoError(t, err)

		// Validate reset session
		retrievedEmail, err := repo.ValidateResetSession(ctx, tokenHash)
		assert.NoError(t, err)
		assert.Equal(t, userEmail, retrievedEmail, "Retrieved email should match stored email")
	})

	t.Run("should verify token is deleted after validation", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "single_use_token"
		userEmail := "singleuse@example.com"
		ttl := 15 * time.Minute

		// Store reset session
		key := sessionRepository.FormatResetSessionKey(tokenHash)
		err := testClient.Set(ctx, key, userEmail, ttl).Err()
		require.NoError(t, err)

		// Verify token exists before validation
		exists, err := testClient.Exists(ctx, key).Result()
		require.NoError(t, err)
		require.Equal(t, int64(1), exists, "Token should exist before validation")

		// Validate reset session
		retrievedEmail, err := repo.ValidateResetSession(ctx, tokenHash)
		require.NoError(t, err)
		require.Equal(t, userEmail, retrievedEmail)

		// Verify token is deleted after validation using raw Redis query
		existsAfter, err := testClient.Exists(ctx, key).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), existsAfter, "Token should be deleted after validation (single-use)")
	})

	t.Run("should fail second validation attempt", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "double_use_attempt_token"
		userEmail := "doubleuse@example.com"
		ttl := 15 * time.Minute

		// Store reset session
		key := sessionRepository.FormatResetSessionKey(tokenHash)
		err := testClient.Set(ctx, key, userEmail, ttl).Err()
		require.NoError(t, err)

		// First validation should succeed
		retrievedEmail, err := repo.ValidateResetSession(ctx, tokenHash)
		require.NoError(t, err)
		require.Equal(t, userEmail, retrievedEmail)

		// Second validation should fail (token consumed)
		_, err = repo.ValidateResetSession(ctx, tokenHash)
		assert.Error(t, err, "Second validation should fail as token is consumed")
	})

	t.Run("should return error for non-existent token", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "nonexistent_token_hash"

		// Try to validate non-existent token
		retrievedEmail, err := repo.ValidateResetSession(ctx, tokenHash)
		assert.Error(t, err, "Should return error for non-existent token")
		assert.Empty(t, retrievedEmail, "Email should be empty for non-existent token")
	})

	t.Run("should return correct email for valid token", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create multiple reset tokens
		tokens := map[string]string{
			"token_1": "user1@example.com",
			"token_2": "user2@example.com",
			"token_3": "user3@example.com",
		}
		ttl := 15 * time.Minute

		// Store all reset sessions
		for tokenHash, email := range tokens {
			key := sessionRepository.FormatResetSessionKey(tokenHash)
			err := testClient.Set(ctx, key, email, ttl).Err()
			require.NoError(t, err)
		}

		// Validate each token and verify correct email is returned
		for tokenHash, expectedEmail := range tokens {
			retrievedEmail, err := repo.ValidateResetSession(ctx, tokenHash)
			assert.NoError(t, err)
			assert.Equal(t, expectedEmail, retrievedEmail, "Email for token %s should match", tokenHash)
		}
	})

	t.Run("should handle special characters in token hash", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		// Base64-encoded token hash with special characters
		tokenHash := "xY9+/zA=456_special"
		userEmail := "special@example.com"
		ttl := 15 * time.Minute

		// Store reset session
		key := sessionRepository.FormatResetSessionKey(tokenHash)
		err := testClient.Set(ctx, key, userEmail, ttl).Err()
		require.NoError(t, err)

		// Validate reset session
		retrievedEmail, err := repo.ValidateResetSession(ctx, tokenHash)
		assert.NoError(t, err)
		assert.Equal(t, userEmail, retrievedEmail, "Should handle special characters in token hash")

		// Verify token is deleted
		exists, err := testClient.Exists(ctx, key).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), exists, "Token with special characters should be deleted")
	})

	t.Run("should handle expired token", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "expired_token"
		userEmail := "expired@example.com"
		ttl := 1 * time.Second

		// Store reset session with very short TTL
		key := sessionRepository.FormatResetSessionKey(tokenHash)
		err := testClient.Set(ctx, key, userEmail, ttl).Err()
		require.NoError(t, err)

		// Wait for token to expire
		time.Sleep(2 * time.Second)

		// Try to validate expired token
		retrievedEmail, err := repo.ValidateResetSession(ctx, tokenHash)
		assert.Error(t, err, "Should return error for expired token")
		assert.Empty(t, retrievedEmail, "Email should be empty for expired token")
	})

	t.Run("should handle Redis connection failure on get", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "connection_failure_token"

		// Close Redis connection to simulate failure
		testClient.Close()

		// Try to validate reset session (should fail on get)
		retrievedEmail, err := repo.ValidateResetSession(ctx, tokenHash)
		assert.Error(t, err)
		assert.Empty(t, retrievedEmail)

		// Reconnect for cleanup
		reconnectRedis()
	})

	t.Run("should handle empty email value", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "empty_email_token"
		userEmail := ""
		ttl := 15 * time.Minute

		// Store reset session with empty email
		key := sessionRepository.FormatResetSessionKey(tokenHash)
		err := testClient.Set(ctx, key, userEmail, ttl).Err()
		require.NoError(t, err)

		// Validate reset session
		retrievedEmail, err := repo.ValidateResetSession(ctx, tokenHash)
		assert.NoError(t, err)
		assert.Equal(t, "", retrievedEmail, "Should handle empty email value")

		// Verify token is deleted
		exists, err := testClient.Exists(ctx, key).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), exists, "Token should be deleted even with empty email")
	})

	t.Run("should handle long email addresses", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "long_email_validation_token"
		userEmail := "very.long.email.with.multiple.subdomains.and.parts@subdomain.example.organization.company.com"
		ttl := 15 * time.Minute

		// Store reset session
		key := sessionRepository.FormatResetSessionKey(tokenHash)
		err := testClient.Set(ctx, key, userEmail, ttl).Err()
		require.NoError(t, err)

		// Validate reset session
		retrievedEmail, err := repo.ValidateResetSession(ctx, tokenHash)
		assert.NoError(t, err)
		assert.Equal(t, userEmail, retrievedEmail, "Should handle long email addresses")
	})

	t.Run("should handle validation near token expiration", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, testClient)

		tokenHash := "near_expiration_token"
		userEmail := "nearexpiry@example.com"
		ttl := 3 * time.Second

		// Store reset session
		key := sessionRepository.FormatResetSessionKey(tokenHash)
		err := testClient.Set(ctx, key, userEmail, ttl).Err()
		require.NoError(t, err)

		// Wait until near expiration
		time.Sleep(2 * time.Second)

		// Validate should still succeed if called before expiration
		retrievedEmail, err := repo.ValidateResetSession(ctx, tokenHash)
		assert.NoError(t, err)
		assert.Equal(t, userEmail, retrievedEmail, "Should validate successfully even near expiration")

		// Verify token is deleted
		exists, err := testClient.Exists(ctx, key).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), exists, "Token should be deleted after validation")
	})
}
