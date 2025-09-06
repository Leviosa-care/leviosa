package sessionRepository_test

import (
	"strings"
	"testing"

	"github.com/Leviosa-care/core/auth/session"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TEST=TestFormatKeys make test-unit-session-test

func TestFormatSessionKey(t *testing.T) {
	t.Run("should format session key with UUID", func(t *testing.T) {
		sessionID := uuid.New().String()
		expected := "authuser:session:" + sessionID

		result := session.FormatSessionKey(sessionID)

		assert.Equal(t, expected, result)
		assert.Contains(t, result, "authuser:session:")
		assert.Contains(t, result, sessionID)
	})

	t.Run("should format session key with custom string", func(t *testing.T) {
		sessionID := "custom-session-id-123"
		expected := "authuser:session:custom-session-id-123"

		result := session.FormatSessionKey(sessionID)

		assert.Equal(t, expected, result)
	})

	t.Run("should handle empty session ID", func(t *testing.T) {
		sessionID := ""
		expected := "authuser:session:"

		result := session.FormatSessionKey(sessionID)

		assert.Equal(t, expected, result)
	})

	t.Run("should handle session ID with special characters", func(t *testing.T) {
		sessionID := "session-id_with.special@chars+123/456="
		expected := "authuser:session:session-id_with.special@chars+123/456="

		result := session.FormatSessionKey(sessionID)

		assert.Equal(t, expected, result)
	})
}

func TestFormatTokenKey(t *testing.T) {
	t.Run("should format token key with hash", func(t *testing.T) {
		tokenHash := "abc123def456ghi789"
		expected := "authuser:token:abc123def456ghi789"

		result := session.FormatAccessTokenKey(tokenHash)

		assert.Equal(t, expected, result)
		assert.Contains(t, result, "authuser:token:")
		assert.Contains(t, result, tokenHash)
	})

	t.Run("should handle token hash with special characters", func(t *testing.T) {
		tokenHash := "token-hash_with.special@chars+123/456="
		expected := "authuser:token:token-hash_with.special@chars+123/456="

		result := session.FormatAccessTokenKey(tokenHash)

		assert.Equal(t, expected, result)
	})

	t.Run("should handle empty token hash", func(t *testing.T) {
		tokenHash := ""
		expected := "authuser:token:"

		result := session.FormatAccessTokenKey(tokenHash)

		assert.Equal(t, expected, result)
	})
}

func TestFormatAccessTokenKey(t *testing.T) {
	t.Run("should format access token key with hash", func(t *testing.T) {
		accessTokenHash := "access-token-hash-123"
		expected := "authuser:access:access-token-hash-123"

		result := session.FormatAccessTokenKey(accessTokenHash)

		assert.Equal(t, expected, result)
		assert.Contains(t, result, "authuser:access:")
		assert.Contains(t, result, accessTokenHash)
	})

	t.Run("should handle base64-like access token hash", func(t *testing.T) {
		accessTokenHash := "dGVzdC1hY2Nlc3MtdG9rZW4="
		expected := "authuser:access:dGVzdC1hY2Nlc3MtdG9rZW4="

		result := session.FormatAccessTokenKey(accessTokenHash)

		assert.Equal(t, expected, result)
	})

	t.Run("should handle access token hash with special characters", func(t *testing.T) {
		accessTokenHash := "access-token_hash.with:special@chars+123/456="
		expected := "authuser:access:access-token_hash.with:special@chars+123/456="

		result := session.FormatAccessTokenKey(accessTokenHash)

		assert.Equal(t, expected, result)
	})

	t.Run("should handle empty access token hash", func(t *testing.T) {
		accessTokenHash := ""
		expected := "authuser:access:"

		result := session.FormatAccessTokenKey(accessTokenHash)

		assert.Equal(t, expected, result)
	})

	t.Run("should handle very long access token hash", func(t *testing.T) {
		// Simulate a long JWT-like token hash
		accessTokenHash := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
		expected := "authuser:access:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

		result := session.FormatAccessTokenKey(accessTokenHash)

		assert.Equal(t, expected, result)
	})
}

func TestFormatRefreshTokenKey(t *testing.T) {
	t.Run("should format refresh token key with hash", func(t *testing.T) {
		refreshTokenHash := "refresh-token-hash-456"
		expected := "authuser:refresh:refresh-token-hash-456"

		result := session.FormatRefreshTokenKey(refreshTokenHash)

		assert.Equal(t, expected, result)
		assert.Contains(t, result, "authuser:refresh:")
		assert.Contains(t, result, refreshTokenHash)
	})

	t.Run("should handle base64-like refresh token hash", func(t *testing.T) {
		refreshTokenHash := "dGVzdC1yZWZyZXNoLXRva2Vu"
		expected := "authuser:refresh:dGVzdC1yZWZyZXNoLXRva2Vu"

		result := session.FormatRefreshTokenKey(refreshTokenHash)

		assert.Equal(t, expected, result)
	})

	t.Run("should handle refresh token hash with special characters", func(t *testing.T) {
		refreshTokenHash := "refresh-token_hash.with:special@chars+789/012="
		expected := "authuser:refresh:refresh-token_hash.with:special@chars+789/012="

		result := session.FormatRefreshTokenKey(refreshTokenHash)

		assert.Equal(t, expected, result)
	})

	t.Run("should handle empty refresh token hash", func(t *testing.T) {
		refreshTokenHash := ""
		expected := "authuser:refresh:"

		result := session.FormatRefreshTokenKey(refreshTokenHash)

		assert.Equal(t, expected, result)
	})

	t.Run("should handle very long refresh token hash", func(t *testing.T) {
		// Simulate a long refresh token hash
		refreshTokenHash := "very-long-refresh-token-hash-with-many-characters-that-could-be-used-in-production-environments-to-ensure-uniqueness-and-security-abc123def456ghi789"
		expected := "authuser:refresh:very-long-refresh-token-hash-with-many-characters-that-could-be-used-in-production-environments-to-ensure-uniqueness-and-security-abc123def456ghi789"

		result := session.FormatRefreshTokenKey(refreshTokenHash)

		assert.Equal(t, expected, result)
	})
}

func TestKeyFormatConsistency(t *testing.T) {
	t.Run("should have different prefixes for different key types", func(t *testing.T) {
		testValue := "test-value-123"

		sessionKey := session.FormatSessionKey(testValue)
		tokenKey := session.FormatAccessTokenKey(testValue)
		accessTokenKey := session.FormatAccessTokenKey(testValue)
		refreshTokenKey := session.FormatRefreshTokenKey(testValue)

		// All should contain the test value
		assert.Contains(t, sessionKey, testValue)
		assert.Contains(t, tokenKey, testValue)
		assert.Contains(t, accessTokenKey, testValue)
		assert.Contains(t, refreshTokenKey, testValue)

		// All should start with authuser prefix
		assert.True(t, strings.HasPrefix(sessionKey, "authuser:"))
		assert.True(t, strings.HasPrefix(tokenKey, "authuser:"))
		assert.True(t, strings.HasPrefix(accessTokenKey, "authuser:"))
		assert.True(t, strings.HasPrefix(refreshTokenKey, "authuser:"))

		// Each should have unique secondary prefix
		assert.Contains(t, sessionKey, ":session:")
		assert.Contains(t, tokenKey, ":token:")
		assert.Contains(t, accessTokenKey, ":access:")
		assert.Contains(t, refreshTokenKey, ":refresh:")

		// All keys should be different
		keys := []string{sessionKey, tokenKey, accessTokenKey, refreshTokenKey}
		for i, key1 := range keys {
			for j, key2 := range keys {
				if i != j {
					assert.NotEqual(t, key1, key2, "Keys should be unique: %s vs %s", key1, key2)
				}
			}
		}
	})

	t.Run("should maintain prefix constants", func(t *testing.T) {
		// These tests ensure that if prefixes change, we catch it in tests
		sessionKey := session.FormatSessionKey("test")
		assert.Equal(t, "authuser:session:test", sessionKey)

		tokenKey := session.FormatAccessTokenKey("test")
		assert.Equal(t, "authuser:token:test", tokenKey)

		accessTokenKey := session.FormatAccessTokenKey("test")
		assert.Equal(t, "authuser:access:test", accessTokenKey)

		refreshTokenKey := session.FormatRefreshTokenKey("test")
		assert.Equal(t, "authuser:refresh:test", refreshTokenKey)
	})
}
