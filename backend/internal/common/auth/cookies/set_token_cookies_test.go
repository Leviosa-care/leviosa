package cookies_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetTokenCookies(t *testing.T) {
	// Setup
	w := httptest.NewRecorder()
	accessToken := "access_token_123"
	refreshToken := "refresh_token_456"
	accessExpiry := time.Now().Add(time.Hour)
	refreshExpiry := time.Now().Add(24 * time.Hour)

	// Execute
	cookies.SetTokenCookies(w, accessToken, refreshToken, accessExpiry, refreshExpiry)

	// Verify
	responseCookies := w.Result().Cookies()
	require.Len(t, responseCookies, 2, "Should set exactly 2 cookies")

	// Find and verify access token cookie
	var accessCookie, refreshCookie *http.Cookie
	for _, cookie := range responseCookies {
		if cookie.Name == cookies.AccessTokenCookieName {
			accessCookie = cookie
		} else if cookie.Name == cookies.RefreshTokenCookieName {
			refreshCookie = cookie
		}
	}

	// Verify access token cookie
	require.NotNil(t, accessCookie, "Access token cookie should be set")
	assert.Equal(t, accessToken, accessCookie.Value)
	assert.Equal(t, "/", accessCookie.Path)
	assert.True(t, accessCookie.HttpOnly)
	assert.True(t, accessCookie.Secure)
	assert.Equal(t, http.SameSiteStrictMode, accessCookie.SameSite)
	assert.WithinDuration(t, accessExpiry, accessCookie.Expires, time.Second)

	// Verify refresh token cookie
	require.NotNil(t, refreshCookie, "Refresh token cookie should be set")
	assert.Equal(t, refreshToken, refreshCookie.Value)
	assert.Equal(t, cookies.RefreshEndpoint, refreshCookie.Path)
	assert.True(t, refreshCookie.HttpOnly)
	assert.True(t, refreshCookie.Secure)
	assert.Equal(t, http.SameSiteStrictMode, refreshCookie.SameSite)
	assert.WithinDuration(t, refreshExpiry, refreshCookie.Expires, time.Second)
}

