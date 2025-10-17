package auth_test

// import (
// 	"context"
// 	"net/http"
// 	"strings"
// 	"testing"
//
// 	"github.com/Leviosa-care/authuser/test/helpers"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )
//
// // TEST=TestOAuth_DatabaseErrors make test-integration-auth-test
//
// func TestOAuth_DatabaseErrors(t *testing.T) {
// 	ctx := context.Background()
// 	client := &http.Client{Timeout: testTimeout}
//
// 	t.Run("should handle database connection failures during OAuth", func(t *testing.T) {
// 		// Note: This test would require database connection mocking or temporary connection issues
// 		// For integration tests, we'll test the happy path and simulate errors where possible
//
// 		// Clear users table
// 		helpers.ClearUsersTable(t, ctx, testPool)
//
// 		// Setup minimal OAuth configuration
// 		t.Setenv("GOOGLE_CLIENT_ID", "test_client")
// 		t.Setenv("GOOGLE_CLIENT_SECRET", "test_secret")
// 		t.Setenv("BASE_URL", testServerURL)
//
// 		// Test OAuth start with valid configuration
// 		req := helpers.NewOAuthStartRequest(t, testServerURL, "google")
// 		resp, err := client.Do(req)
// 		require.NoError(t, err)
// 		defer resp.Body.Close()
//
// 		// Should handle request even if subsequent database operations might fail
// 		assert.True(t, resp.StatusCode == http.StatusFound || resp.StatusCode >= 400,
// 			"Should handle OAuth start request")
// 	})
// }
//
// // TEST=TestOAuth_ConfigurationErrors make test-integration-auth-test
//
// func TestOAuth_ConfigurationErrors(t *testing.T) {
// 	client := &http.Client{
// 		Timeout: testTimeout,
// 		CheckRedirect: func(req *http.Request, via []*http.Request) error {
// 			return http.ErrUseLastResponse
// 		},
// 	}
//
// 	testCases := []struct {
// 		name           string
// 		provider       string
// 		envVars        map[string]string
// 		expectedStatus int
// 		description    string
// 	}{
// 		{
// 			name:     "missing Google client ID",
// 			provider: "google",
// 			envVars: map[string]string{
// 				"GOOGLE_CLIENT_ID":     "",
// 				"GOOGLE_CLIENT_SECRET": "test_secret",
// 			},
// 			expectedStatus: http.StatusInternalServerError,
// 			description:    "Should fail when Google client ID is missing",
// 		},
// 		{
// 			name:     "missing Google client secret",
// 			provider: "google",
// 			envVars: map[string]string{
// 				"GOOGLE_CLIENT_ID":     "test_client",
// 				"GOOGLE_CLIENT_SECRET": "",
// 			},
// 			expectedStatus: http.StatusInternalServerError,
// 			description:    "Should fail when Google client secret is missing",
// 		},
// 		{
// 			name:     "missing Apple client ID",
// 			provider: "apple",
// 			envVars: map[string]string{
// 				"APPLE_CLIENT_ID":     "",
// 				"APPLE_CLIENT_SECRET": "test_secret",
// 				"APPLE_TEAM_ID":       "TEAMID123",
// 				"APPLE_KEY_ID":        "KEYID123",
// 				"APPLE_PRIVATE_KEY":   "private_key",
// 			},
// 			expectedStatus: http.StatusInternalServerError,
// 			description:    "Should fail when Apple client ID is missing",
// 		},
// 		{
// 			name:     "missing Apple team ID",
// 			provider: "apple",
// 			envVars: map[string]string{
// 				"APPLE_CLIENT_ID":     "test.client.id",
// 				"APPLE_CLIENT_SECRET": "test_secret",
// 				"APPLE_TEAM_ID":       "",
// 				"APPLE_KEY_ID":        "KEYID123",
// 				"APPLE_PRIVATE_KEY":   "private_key",
// 			},
// 			expectedStatus: http.StatusInternalServerError,
// 			description:    "Should fail when Apple team ID is missing",
// 		},
// 		{
// 			name:     "missing Google client ID with Nextcloud testing enabled",
// 			provider: "google",
// 			envVars: map[string]string{
// 				"GOOGLE_CLIENT_ID":          "",
// 				"GOOGLE_CLIENT_SECRET":      "test_secret",
// 				"USE_NEXTCLOUD_FOR_TESTING": "true",
// 				"NEXTCLOUD_TEST_URL":        "https://nextcloud.example.com",
// 			},
// 			expectedStatus: http.StatusInternalServerError,
// 			description:    "Should fail when Google client ID is missing even with Nextcloud testing",
// 		},
// 		{
// 			name:     "missing Nextcloud test URL when testing enabled",
// 			provider: "google",
// 			envVars: map[string]string{
// 				"GOOGLE_CLIENT_ID":          "test_client",
// 				"GOOGLE_CLIENT_SECRET":      "test_secret",
// 				"USE_NEXTCLOUD_FOR_TESTING": "true",
// 				"NEXTCLOUD_TEST_URL":        "",
// 			},
// 			expectedStatus: http.StatusInternalServerError,
// 			description:    "Should fail when Nextcloud test URL is missing but testing is enabled",
// 		},
// 	}
//
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Set environment variables for this test case
// 			for key, value := range tc.envVars {
// 				t.Setenv(key, value)
// 			}
// 			t.Setenv("BASE_URL", testServerURL)
//
// 			// Create OAuth start request
// 			req := helpers.NewOAuthStartRequest(t, testServerURL, tc.provider)
//
// 			// Make request
// 			resp, err := client.Do(req)
// 			require.NoError(t, err)
// 			defer resp.Body.Close()
//
// 			// Verify expected error response
// 			assert.Equal(t, tc.expectedStatus, resp.StatusCode, tc.description)
// 		})
// 	}
// }
//
// // TEST=TestOAuth_InvalidProviders make test-integration-auth-test
//
// func TestOAuth_InvalidProviders(t *testing.T) {
// 	client := &http.Client{Timeout: testTimeout}
//
// 	testCases := []struct {
// 		name         string
// 		provider     string
// 		expectedCode int
// 	}{
// 		{
// 			name:         "unsupported provider",
// 			provider:     "unsupported_provider",
// 			expectedCode: http.StatusBadRequest,
// 		},
// 		{
// 			name:         "empty provider",
// 			provider:     "",
// 			expectedCode: http.StatusNotFound, // Route not found
// 		},
// 		{
// 			name:         "special characters in provider",
// 			provider:     "provider@#$%",
// 			expectedCode: http.StatusBadRequest,
// 		},
// 		{
// 			name:         "case sensitivity",
// 			provider:     "GOOGLE", // Should be lowercase "google"
// 			expectedCode: http.StatusBadRequest,
// 		},
// 	}
//
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Create OAuth start request with invalid provider
// 			var req *http.Request
// 			if tc.provider == "" {
// 				// For empty provider, construct URL manually
// 				req, _ = http.NewRequest("GET", testServerURL+"/auth/oauth//start", nil)
// 			} else {
// 				req = helpers.NewOAuthStartRequest(t, testServerURL, tc.provider)
// 			}
//
// 			resp, err := client.Do(req)
// 			require.NoError(t, err)
// 			defer resp.Body.Close()
//
// 			assert.Equal(t, tc.expectedCode, resp.StatusCode,
// 				"Provider '%s' should return %d", tc.provider, tc.expectedCode)
// 		})
// 	}
// }
//
// // TEST=TestOAuth_CallbackErrorScenarios make test-integration-auth-test
//
// func TestOAuth_CallbackErrorScenarios(t *testing.T) {
// 	ctx := context.Background()
// 	client := &http.Client{Timeout: testTimeout}
//
// 	// Setup valid OAuth configuration
// 	t.Setenv("GOOGLE_CLIENT_ID", "test_client")
// 	t.Setenv("GOOGLE_CLIENT_SECRET", "test_secret")
// 	t.Setenv("BASE_URL", testServerURL)
//
// 	testCases := []struct {
// 		name         string
// 		provider     string
// 		code         string
// 		state        string
// 		extraParams  map[string]string
// 		expectedCode int
// 		description  string
// 	}{
// 		{
// 			name:         "missing authorization code",
// 			provider:     "google",
// 			code:         "",
// 			state:        "valid_state",
// 			expectedCode: http.StatusBadRequest,
// 			description:  "Should reject callback without authorization code",
// 		},
// 		{
// 			name:         "missing state parameter",
// 			provider:     "google",
// 			code:         "valid_code",
// 			state:        "",
// 			expectedCode: http.StatusBadRequest,
// 			description:  "Should reject callback without state (CSRF protection)",
// 		},
// 		{
// 			name:         "OAuth error from provider",
// 			provider:     "google",
// 			code:         "",
// 			state:        "valid_state",
// 			extraParams:  map[string]string{"error": "access_denied"},
// 			expectedCode: http.StatusUnauthorized,
// 			description:  "Should handle OAuth provider errors",
// 		},
// 		{
// 			name:     "OAuth error with description",
// 			provider: "google",
// 			code:     "",
// 			state:    "valid_state",
// 			extraParams: map[string]string{
// 				"error":             "invalid_request",
// 				"error_description": "The request is missing a required parameter",
// 			},
// 			expectedCode: http.StatusBadRequest,
// 			description:  "Should handle detailed OAuth errors",
// 		},
// 		{
// 			name:         "expired state parameter",
// 			provider:     "google",
// 			code:         "valid_code",
// 			state:        "expired_state_from_another_session",
// 			expectedCode: http.StatusBadRequest,
// 			description:  "Should reject expired or invalid state",
// 		},
// 	}
//
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Clear users table
// 			helpers.ClearUsersTable(t, ctx, testPool)
//
// 			// Create OAuth callback request
// 			req := helpers.NewOAuthCallbackRequest(t, testServerURL, tc.provider, tc.code, tc.state)
//
// 			// Add extra parameters if provided
// 			if len(tc.extraParams) > 0 {
// 				query := req.URL.Query()
// 				for key, value := range tc.extraParams {
// 					query.Add(key, value)
// 				}
// 				req.URL.RawQuery = query.Encode()
// 			}
//
// 			resp, err := client.Do(req)
// 			require.NoError(t, err)
// 			defer resp.Body.Close()
//
// 			assert.Equal(t, tc.expectedCode, resp.StatusCode, tc.description)
// 		})
// 	}
// }
//
// // TEST=TestOAuth_SessionErrors make test-integration-auth-test
//
// func TestOAuth_SessionErrors(t *testing.T) {
// 	ctx := context.Background()
// 	client := &http.Client{Timeout: testTimeout}
//
// 	t.Run("should handle session creation failures", func(t *testing.T) {
// 		// Clear users table
// 		helpers.ClearUsersTable(t, ctx, testPool)
//
// 		// Setup OAuth configuration
// 		t.Setenv("GOOGLE_CLIENT_ID", "test_client")
// 		t.Setenv("GOOGLE_CLIENT_SECRET", "test_secret")
// 		t.Setenv("SESSION_SECRET", "") // Missing session secret should cause errors
// 		t.Setenv("BASE_URL", testServerURL)
//
// 		// Create OAuth start request
// 		req := helpers.NewOAuthStartRequest(t, testServerURL, "google")
//
// 		resp, err := client.Do(req)
// 		require.NoError(t, err)
// 		defer resp.Body.Close()
//
// 		// Should handle session configuration errors
// 		assert.True(t, resp.StatusCode >= 400, "Should return error for session configuration issues")
// 	})
//
// 	t.Run("should handle invalid session data", func(t *testing.T) {
// 		// Setup OAuth configuration
// 		t.Setenv("GOOGLE_CLIENT_ID", "test_client")
// 		t.Setenv("GOOGLE_CLIENT_SECRET", "test_secret")
// 		t.Setenv("SESSION_SECRET", "valid_session_secret_32_bytes!!")
// 		t.Setenv("BASE_URL", testServerURL)
//
// 		// Create OAuth callback request with potentially corrupted session
// 		req := helpers.NewOAuthCallbackRequest(t, testServerURL, "google", "code", "invalid_state")
//
// 		// Add malformed session cookie
// 		req.AddCookie(&http.Cookie{
// 			Name:  "session",
// 			Value: "corrupted_session_data",
// 		})
//
// 		resp, err := client.Do(req)
// 		require.NoError(t, err)
// 		defer resp.Body.Close()
//
// 		// Should handle corrupted session data gracefully
// 		assert.True(t, resp.StatusCode >= 400, "Should handle corrupted session data")
// 	})
// }
//
// // TEST=TestOAuth_NetworkErrors make test-integration-auth-test
//
// func TestOAuth_NetworkErrors(t *testing.T) {
// 	client := &http.Client{Timeout: testTimeout}
//
// 	t.Run("should handle unreachable OAuth provider", func(t *testing.T) {
// 		// Setup Google OAuth with Nextcloud testing using unreachable URL
// 		t.Setenv("GOOGLE_CLIENT_ID", "test_client")
// 		t.Setenv("GOOGLE_CLIENT_SECRET", "test_secret")
// 		t.Setenv("USE_NEXTCLOUD_FOR_TESTING", "true")
// 		t.Setenv("NEXTCLOUD_TEST_URL", "https://unreachable.invalid.domain")
// 		t.Setenv("BASE_URL", testServerURL)
//
// 		// Create OAuth callback request (this would fail when trying to exchange code)
// 		req := helpers.NewOAuthCallbackRequest(t, testServerURL, "google", "valid_code", "valid_state")
//
// 		resp, err := client.Do(req)
// 		require.NoError(t, err)
// 		defer resp.Body.Close()
//
// 		// Should handle network errors gracefully
// 		assert.True(t, resp.StatusCode >= 400, "Should handle unreachable OAuth provider")
// 	})
// }
//
// // TEST=TestOAuth_RateLimiting make test-integration-auth-test
//
// func TestOAuth_RateLimiting(t *testing.T) {
// 	client := &http.Client{Timeout: testTimeout}
//
// 	// Setup OAuth configuration
// 	t.Setenv("GOOGLE_CLIENT_ID", "test_client")
// 	t.Setenv("GOOGLE_CLIENT_SECRET", "test_secret")
// 	t.Setenv("BASE_URL", testServerURL)
//
// 	t.Run("should handle multiple rapid OAuth requests", func(t *testing.T) {
// 		// Make multiple rapid OAuth start requests to test rate limiting
// 		const numRequests = 10
// 		responses := make([]*http.Response, numRequests)
//
// 		for i := 0; i < numRequests; i++ {
// 			req := helpers.NewOAuthStartRequest(t, testServerURL, "google")
// 			resp, err := client.Do(req)
// 			require.NoError(t, err)
// 			responses[i] = resp
// 		}
//
// 		// Check that all requests were handled
// 		successCount := 0
// 		for _, resp := range responses {
// 			if resp.StatusCode == http.StatusFound {
// 				successCount++
// 			}
// 			resp.Body.Close()
// 		}
//
// 		// Should handle multiple requests gracefully (may implement rate limiting)
// 		assert.True(t, successCount > 0, "Should handle at least some OAuth requests")
//
// 		// Note: If rate limiting is implemented, some requests might return 429 Too Many Requests
// 		// This test verifies the system doesn't crash under rapid requests
// 	})
// }
//
// // TEST=TestOAuth_SecurityScenarios make test-integration-auth-test
//
// func TestOAuth_SecurityScenarios(t *testing.T) {
// 	client := &http.Client{Timeout: testTimeout}
//
// 	// Setup OAuth configuration
// 	t.Setenv("GOOGLE_CLIENT_ID", "test_client")
// 	t.Setenv("GOOGLE_CLIENT_SECRET", "test_secret")
// 	t.Setenv("BASE_URL", testServerURL)
//
// 	t.Run("should reject malicious redirect attempts", func(t *testing.T) {
// 		// Create request with potentially malicious redirect
// 		req, err := http.NewRequest("GET", testServerURL+"/auth/oauth/google/start", nil)
// 		require.NoError(t, err)
//
// 		// Add malicious parameters
// 		query := req.URL.Query()
// 		query.Add("redirect_uri", "https://malicious.com/steal-tokens")
// 		req.URL.RawQuery = query.Encode()
//
// 		resp, err := client.Do(req)
// 		require.NoError(t, err)
// 		defer resp.Body.Close()
//
// 		// Should not allow arbitrary redirect URLs
// 		if resp.StatusCode == http.StatusFound {
// 			location := resp.Header.Get("Location")
// 			assert.False(t, strings.Contains(location, "malicious.com"),
// 				"Should not allow malicious redirect URLs")
// 		}
// 	})
//
// 	t.Run("should validate CSRF state parameter", func(t *testing.T) {
// 		// Create OAuth callback with mismatched state
// 		req := helpers.NewOAuthCallbackRequest(t, testServerURL, "google", "code", "tampered_state")
//
// 		resp, err := client.Do(req)
// 		require.NoError(t, err)
// 		defer resp.Body.Close()
//
// 		// Should reject requests with invalid state
// 		assert.True(t, resp.StatusCode >= 400, "Should validate CSRF state parameter")
// 	})
// }
