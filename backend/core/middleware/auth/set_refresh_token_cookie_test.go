package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetRefreshTokenCookie(t *testing.T) {
	// Setup
	w := httptest.NewRecorder()
	refreshToken := "new_refresh_token"
	expiry := time.Now().Add(24 * time.Hour)

	// Execute
	SetRefreshTokenCookie(w, refreshToken, expiry)

	// Verify
	cookies := w.Result().Cookies()
	require.Len(t, cookies, 1, "Should set exactly 1 cookie")

	cookie := cookies[0]
	assert.Equal(t, RefreshTokenCookieName, cookie.Name)
	assert.Equal(t, refreshToken, cookie.Value)
	assert.Equal(t, RefreshEndpoint, cookie.Path)
	assert.True(t, cookie.HttpOnly)
	assert.True(t, cookie.Secure)
	assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
	assert.WithinDuration(t, expiry, cookie.Expires, time.Second)
}