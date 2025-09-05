package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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