package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	SetTokenCookies(w, accessToken, refreshToken, accessExpiry, refreshExpiry)

	// Verify
	cookies := w.Result().Cookies()
	require.Len(t, cookies, 2, "Should set exactly 2 cookies")

	// Find and verify access token cookie
	var accessCookie, refreshCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == AccessTokenCookieName {
			accessCookie = cookie
		} else if cookie.Name == RefreshTokenCookieName {
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
	assert.Equal(t, RefreshEndpoint, refreshCookie.Path)
	assert.True(t, refreshCookie.HttpOnly)
	assert.True(t, refreshCookie.Secure)
	assert.Equal(t, http.SameSiteStrictMode, refreshCookie.SameSite)
	assert.WithinDuration(t, refreshExpiry, refreshCookie.Expires, time.Second)
}

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

func TestGetAccessTokenFromCookies(t *testing.T) {
	tests := []struct {
		name        string
		setupCookie func(*http.Request)
		expectedVal string
		expectError bool
	}{
		{
			name: "valid access token cookie",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  AccessTokenCookieName,
					Value: "test_access_token",
				})
			},
			expectedVal: "test_access_token",
			expectError: false,
		},
		{
			name:        "missing access token cookie",
			setupCookie: func(req *http.Request) {},
			expectedVal: "",
			expectError: true,
		},
		{
			name: "empty access token cookie",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  AccessTokenCookieName,
					Value: "",
				})
			},
			expectedVal: "",
			expectError: false,
		},
		{
			name: "access token with special characters",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  AccessTokenCookieName,
					Value: "token-123_456.789",
				})
			},
			expectedVal: "token-123_456.789",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			req := httptest.NewRequest("GET", "/test", nil)
			tt.setupCookie(req)

			// Execute
			token, err := GetAccessTokenFromCookies(req)

			// Verify
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedVal, token)
			}
		})
	}
}

func TestGetRefreshTokenFromCookies(t *testing.T) {
	tests := []struct {
		name        string
		setupCookie func(*http.Request)
		expectedVal string
		expectError bool
	}{
		{
			name: "valid refresh token cookie",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  RefreshTokenCookieName,
					Value: "test_refresh_token",
				})
			},
			expectedVal: "test_refresh_token",
			expectError: false,
		},
		{
			name:        "missing refresh token cookie",
			setupCookie: func(req *http.Request) {},
			expectedVal: "",
			expectError: true,
		},
		{
			name: "empty refresh token cookie",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  RefreshTokenCookieName,
					Value: "",
				})
			},
			expectedVal: "",
			expectError: false,
		},
		{
			name: "refresh token with special characters",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  RefreshTokenCookieName,
					Value: "refresh-987_654.321",
				})
			},
			expectedVal: "refresh-987_654.321",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			req := httptest.NewRequest("GET", "/test", nil)
			tt.setupCookie(req)

			// Execute
			token, err := GetRefreshTokenFromCookies(req)

			// Verify
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedVal, token)
			}
		})
	}
}

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

