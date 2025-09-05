package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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