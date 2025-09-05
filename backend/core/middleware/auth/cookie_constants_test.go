package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCookieConstants(t *testing.T) {
	// Verify cookie name constants are set correctly
	assert.Equal(t, "leviosa_access_token", AccessTokenCookieName)
	assert.Equal(t, "leviosa_refresh_token", RefreshTokenCookieName)

	// Verify they're not empty
	assert.NotEmpty(t, AccessTokenCookieName)
	assert.NotEmpty(t, RefreshTokenCookieName)

	// Verify they're different
	assert.NotEqual(t, AccessTokenCookieName, RefreshTokenCookieName)
}