package sessionRepository_test

import (
	"testing"

	sessionRepository "github.com/Leviosa-care/leviosa/backend/internal/authuser/infrastructure/redis/session"
	"github.com/stretchr/testify/assert"
)

// TestFormatResetSessionKey tests the FormatResetSessionKey function
func TestFormatResetSessionKey(t *testing.T) {
	t.Run("should_format_reset_session_key_with_token_hash", func(t *testing.T) {
		// Arrange
		tokenHash := "abc123def456ghi789"
		expectedKey := "authuser:reset_session:abc123def456ghi789"

		// Act
		actualKey := sessionRepository.FormatResetSessionKey(tokenHash)

		// Assert
		assert.Equal(t, expectedKey, actualKey)
	})

	t.Run("should_format_reset_session_key_with_custom_string", func(t *testing.T) {
		// Arrange
		tokenHash := "custom-reset-token-123"
		expectedKey := "authuser:reset_session:custom-reset-token-123"

		// Act
		actualKey := sessionRepository.FormatResetSessionKey(tokenHash)

		// Assert
		assert.Equal(t, expectedKey, actualKey)
	})

	t.Run("should_handle_empty_token_hash", func(t *testing.T) {
		// Arrange
		tokenHash := ""
		expectedKey := "authuser:reset_session:"

		// Act
		actualKey := sessionRepository.FormatResetSessionKey(tokenHash)

		// Assert
		assert.Equal(t, expectedKey, actualKey)
	})

	t.Run("should_handle_token_hash_with_special_characters", func(t *testing.T) {
		// Arrange
		tokenHash := "reset-token-with.special@chars+123/456="
		expectedKey := "authuser:reset_session:reset-token-with.special@chars+123/456="

		// Act
		actualKey := sessionRepository.FormatResetSessionKey(tokenHash)

		// Assert
		assert.Equal(t, expectedKey, actualKey)
	})

	t.Run("should_handle_token_hash_with_numbers_and_letters", func(t *testing.T) {
		// Arrange
		tokenHash := "token123abc456def789ghi012"
		expectedKey := "authuser:reset_session:token123abc456def789ghi012"

		// Act
		actualKey := sessionRepository.FormatResetSessionKey(tokenHash)

		// Assert
		assert.Equal(t, expectedKey, actualKey)
	})

	t.Run("should_handle_unicode_characters_in_token_hash", func(t *testing.T) {
		// Arrange
		tokenHash := "unicode-test-ñ-á-é-í-ó-ú"
		expectedKey := "authuser:reset_session:unicode-test-ñ-á-é-í-ó-ú"

		// Act
		actualKey := sessionRepository.FormatResetSessionKey(tokenHash)

		// Assert
		assert.Equal(t, expectedKey, actualKey)
	})

	t.Run("should_handle_long_token_hash", func(t *testing.T) {
		// Arrange
		tokenHash := "very-long-reset-session-token-hash-with-many-characters-1234567890abcdefghijklmnopqrstuvwxyz"
		expectedKey := "authuser:reset_session:very-long-reset-session-token-hash-with-many-characters-1234567890abcdefghijklmnopqrstuvwxyz"

		// Act
		actualKey := sessionRepository.FormatResetSessionKey(tokenHash)

		// Assert
		assert.Equal(t, expectedKey, actualKey)
	})
}

// TestResetSessionKeyFormatConsistency tests the reset session key format consistency
func TestResetSessionKeyFormatConsistency(t *testing.T) {
	t.Run("should_maintain_prefix_constant", func(t *testing.T) {
		// Arrange
		tokenHash := "test-token"
		expectedKey := "authuser:reset_session:test-token"

		// Act
		actualKey := sessionRepository.FormatResetSessionKey(tokenHash)

		// Assert
		assert.Equal(t, expectedKey, actualKey)
	})

	t.Run("should_be_deterministic", func(t *testing.T) {
		// Arrange
		tokenHash := "deterministic-test-123"

		// Act
		firstCall := sessionRepository.FormatResetSessionKey(tokenHash)
		secondCall := sessionRepository.FormatResetSessionKey(tokenHash)

		// Assert
		assert.Equal(t, firstCall, secondCall, "Multiple calls with same input should return identical results")
	})

	t.Run("should_differentiate_different_token_hashes", func(t *testing.T) {
		// Arrange
		tokenHash1 := "token-one-123"
		tokenHash2 := "token-two-456"

		// Act
		key1 := sessionRepository.FormatResetSessionKey(tokenHash1)
		key2 := sessionRepository.FormatResetSessionKey(tokenHash2)

		// Assert
		assert.NotEqual(t, key1, key2, "Different token hashes should produce different keys")
	})
}