package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetAccessTokenCookie(t *testing.T) {
	// Setup
	w := httptest.NewRecorder()
	accessToken := "new_access_token"
	expiry := time.Now().Add(time.Hour)

	// Execute
	SetAccessTokenCookie(w, accessToken, expiry)

	// Verify
	cookies := w.Result().Cookies()
	require.Len(t, cookies, 1, "Should set exactly 1 cookie")

	cookie := cookies[0]
	assert.Equal(t, AccessTokenCookieName, cookie.Name)
	assert.Equal(t, accessToken, cookie.Value)
	assert.Equal(t, "/", cookie.Path)
	assert.True(t, cookie.HttpOnly)
	assert.True(t, cookie.Secure)
	assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
	assert.WithinDuration(t, expiry, cookie.Expires, time.Second)
}