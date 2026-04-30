package cookies_test

import (
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
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
	cookies.SetRefreshTokenCookie(w, refreshToken, expiry)

	// Verify
	responseCookies := w.Result().Cookies()
	require.Len(t, responseCookies, 1, "Should set exactly 1 cookie")

	cookie := responseCookies[0]
	assert.Equal(t, cookies.RefreshTokenCookieName, cookie.Name)
	assert.Equal(t, refreshToken, cookie.Value)
	assert.Equal(t, "/", cookie.Path)
	assert.True(t, cookie.HttpOnly)
	assert.True(t, cookie.Secure)
	assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
	assert.WithinDuration(t, expiry, cookie.Expires, time.Second)
}

