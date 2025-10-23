package cookies_test

import (
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCookieConstants(t *testing.T) {
	// Verify cookie name constants are set correctly
	assert.Equal(t, "leviosa_access_token", cookies.AccessTokenCookieName)
	assert.Equal(t, "leviosa_refresh_token", cookies.RefreshTokenCookieName)

	// Verify they're not empty
	assert.NotEmpty(t, cookies.AccessTokenCookieName)
	assert.NotEmpty(t, cookies.RefreshTokenCookieName)

	// Verify they're different
	assert.NotEqual(t, cookies.AccessTokenCookieName, cookies.RefreshTokenCookieName)
}

