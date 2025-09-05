package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClearTokenCookies(t *testing.T) {
	// Setup
	w := httptest.NewRecorder()

	// Execute
	ClearTokenCookies(w)

	// Verify
	cookies := w.Result().Cookies()
	require.Len(t, cookies, 2, "Should clear exactly 2 cookies")

	// Find and verify cookies
	var accessCookie, refreshCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == AccessTokenCookieName {
			accessCookie = cookie
		} else if cookie.Name == RefreshTokenCookieName {
			refreshCookie = cookie
		}
	}

	// Verify access token cookie is cleared
	require.NotNil(t, accessCookie, "Access token cookie should be set for clearing")
	assert.Equal(t, "", accessCookie.Value)
	assert.Equal(t, "/", accessCookie.Path)
	assert.True(t, accessCookie.HttpOnly)
	assert.True(t, accessCookie.Secure)
	assert.Equal(t, http.SameSiteStrictMode, accessCookie.SameSite)
	assert.Equal(t, -1, accessCookie.MaxAge)

	// Verify refresh token cookie is cleared
	require.NotNil(t, refreshCookie, "Refresh token cookie should be set for clearing")
	assert.Equal(t, "", refreshCookie.Value)
	assert.Equal(t, RefreshEndpoint, refreshCookie.Path)
	assert.True(t, refreshCookie.HttpOnly)
	assert.True(t, refreshCookie.Secure)
	assert.Equal(t, http.SameSiteStrictMode, refreshCookie.SameSite)
	assert.Equal(t, -1, refreshCookie.MaxAge)
}