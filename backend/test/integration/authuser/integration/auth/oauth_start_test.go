package auth_test

// import (
// 	"context"
// 	"net/http"
// 	"net/url"
// 	"strings"
// 	"testing"
// 	"time"
//
// 	"github.com/Leviosa-care/leviosa/backend/test/helpers"
//
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )
//
// const testTimeout = 30 * time.Second
//
// // TEST=TestOAuthStart_GoogleWithNextcloudTesting make test-integration-auth-test
//
// func TestOAuthStart_GoogleWithNextcloudTesting(t *testing.T) {
// 	ctx := context.Background()
// 	client := &http.Client{
// 		Timeout: testTimeout,
// 		CheckRedirect: func(req *http.Request, via []*http.Request) error {
// 			return http.ErrUseLastResponse // Don't follow redirects automatically
// 		},
// 	}
//
// 	// Setup Nextcloud container for OAuth testing (behind the scenes as "google" provider)
// 	nextcloudHelper := helpers.SetupNextcloudOAuth(ctx, t)
// 	defer nextcloudHelper.TeardownNextcloudOAuth(ctx, t)
//
// 	// Get OAuth configuration from Nextcloud container
// 	clientID, clientSecret, nextcloudURL := nextcloudHelper.GetOAuthConfig()
//
// 	// Setup OAuth environment to use Nextcloud as Google provider for testing
// 	helpers.SetupNextcloudAsGoogleOAuthEnvironment(t, nextcloudURL, clientID, clientSecret)
//
// 	t.Run("should successfully redirect to Nextcloud OAuth (registered as google)", func(t *testing.T) {
// 		// Create OAuth start request using "google" provider (but Nextcloud behind the scenes)
// 		req := helpers.NewOAuthStartRequest(t, testServerURL, "google")
//
// 		// Make request
// 		resp, err := client.Do(req)
// 		require.NoError(t, err)
// 		defer resp.Body.Close()
//
// 		// Should redirect successfully
// 		require.Equal(t, http.StatusFound, resp.StatusCode, "Should redirect to OAuth provider")
//
// 		// Extract and validate redirect URL parameters
// 		redirectURL := helpers.ExtractOAuthRedirectURL(t, resp)
// 		query := redirectURL.Query()
//
// 		// Validate OAuth parameters
// 		assert.Equal(t, clientID, query.Get("client_id"), "Client ID should match configuration")
// 		assert.Equal(t, "code", query.Get("response_type"), "Response type should be 'code'")
// 		assert.NotEmpty(t, query.Get("state"), "State parameter should be present for CSRF protection")
// 		assert.Contains(t, query.Get("redirect_uri"), "/auth/oauth/google/callback", "Redirect URI should point to google callback endpoint")
//
// 		// The redirect URL should point to Nextcloud (because we're using Nextcloud for testing)
// 		assert.True(t, strings.HasPrefix(redirectURL.String(), nextcloudURL),
// 			"Should redirect to Nextcloud OAuth endpoint for testing")
// 	})
//
// 	t.Run("should handle invalid provider", func(t *testing.T) {
// 		// Create OAuth start request with invalid provider
// 		req := helpers.NewOAuthStartRequest(t, testServerURL, "invalid_provider")
//
// 		// Make request
// 		resp, err := client.Do(req)
// 		require.NoError(t, err)
// 		defer resp.Body.Close()
//
// 		// Should return error response
// 		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return bad request for invalid provider")
// 	})
//
// 	t.Run("should handle missing OAuth configuration", func(t *testing.T) {
// 		// Temporarily unset OAuth environment variables
// 		t.Setenv("GOOGLE_CLIENT_ID", "")
// 		t.Setenv("GOOGLE_CLIENT_SECRET", "")
// 		t.Setenv("USE_NEXTCLOUD_FOR_TESTING", "false")
//
// 		// Create OAuth start request
// 		req := helpers.NewOAuthStartRequest(t, testServerURL, "google")
//
// 		// Make request
// 		resp, err := client.Do(req)
// 		require.NoError(t, err)
// 		defer resp.Body.Close()
//
// 		// Should return error response due to missing configuration
// 		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return internal server error for missing configuration")
// 	})
// }
//
// // TEST=TestOAuthStart_Google make test-integration-auth-test
//
// func TestOAuthStart_Google(t *testing.T) {
// 	client := &http.Client{
// 		Timeout: testTimeout,
// 		CheckRedirect: func(req *http.Request, via []*http.Request) error {
// 			return http.ErrUseLastResponse // Don't follow redirects automatically
// 		},
// 	}
//
// 	// Setup mock Google OAuth configuration
// 	t.Setenv("GOOGLE_CLIENT_ID", "test_google_client_id")
// 	t.Setenv("GOOGLE_CLIENT_SECRET", "test_google_client_secret")
// 	t.Setenv("BASE_URL", testServerURL)
//
// 	t.Run("should successfully redirect to Google OAuth", func(t *testing.T) {
// 		// Create OAuth start request
// 		req := helpers.NewOAuthStartRequest(t, testServerURL, "google")
//
// 		// Make request
// 		resp, err := client.Do(req)
// 		require.NoError(t, err)
// 		defer resp.Body.Close()
//
// 		// Validate redirect response
// 		require.Equal(t, http.StatusFound, resp.StatusCode, "Should redirect to Google OAuth")
//
// 		// Extract and validate redirect URL
// 		location := resp.Header.Get("Location")
// 		require.NotEmpty(t, location, "Location header should be present")
//
// 		redirectURL, err := url.Parse(location)
// 		require.NoError(t, err)
//
// 		// Should redirect to Google OAuth endpoint
// 		assert.True(t, strings.HasPrefix(redirectURL.String(), "https://accounts.google.com/oauth/authorize"),
// 			"Should redirect to Google OAuth endpoint")
//
// 		// Validate OAuth parameters
// 		query := redirectURL.Query()
// 		assert.Equal(t, "test_google_client_id", query.Get("client_id"), "Client ID should match configuration")
// 		assert.Equal(t, "code", query.Get("response_type"), "Response type should be 'code'")
// 		assert.NotEmpty(t, query.Get("state"), "State parameter should be present")
// 		assert.Contains(t, query.Get("redirect_uri"), "/auth/oauth/google/callback", "Redirect URI should point to callback endpoint")
// 		assert.Contains(t, query.Get("scope"), "email", "Should request email scope")
// 		assert.Contains(t, query.Get("scope"), "profile", "Should request profile scope")
// 	})
// }
//
// // TEST=TestOAuthStart_Apple make test-integration-auth-test
//
// func TestOAuthStart_Apple(t *testing.T) {
// 	client := &http.Client{
// 		Timeout: testTimeout,
// 		CheckRedirect: func(req *http.Request, via []*http.Request) error {
// 			return http.ErrUseLastResponse // Don't follow redirects automatically
// 		},
// 	}
//
// 	// Setup mock Apple OAuth configuration
// 	t.Setenv("APPLE_CLIENT_ID", "test.apple.client.id")
// 	t.Setenv("APPLE_CLIENT_SECRET", "test_apple_client_secret")
// 	t.Setenv("APPLE_TEAM_ID", "TESTTEAMID")
// 	t.Setenv("APPLE_KEY_ID", "TESTKEYID")
// 	t.Setenv("APPLE_PRIVATE_KEY", "test_private_key")
// 	t.Setenv("BASE_URL", testServerURL)
//
// 	t.Run("should successfully redirect to Apple OAuth", func(t *testing.T) {
// 		// Create OAuth start request
// 		req := helpers.NewOAuthStartRequest(t, testServerURL, "apple")
//
// 		// Make request
// 		resp, err := client.Do(req)
// 		require.NoError(t, err)
// 		defer resp.Body.Close()
//
// 		// Validate redirect response
// 		require.Equal(t, http.StatusFound, resp.StatusCode, "Should redirect to Apple OAuth")
//
// 		// Extract and validate redirect URL
// 		location := resp.Header.Get("Location")
// 		require.NotEmpty(t, location, "Location header should be present")
//
// 		redirectURL, err := url.Parse(location)
// 		require.NoError(t, err)
//
// 		// Should redirect to Apple OAuth endpoint
// 		assert.True(t, strings.HasPrefix(redirectURL.String(), "https://appleid.apple.com/auth/authorize"),
// 			"Should redirect to Apple OAuth endpoint")
//
// 		// Validate OAuth parameters
// 		query := redirectURL.Query()
// 		assert.Equal(t, "test.apple.client.id", query.Get("client_id"), "Client ID should match configuration")
// 		assert.Equal(t, "code", query.Get("response_type"), "Response type should be 'code'")
// 		assert.NotEmpty(t, query.Get("state"), "State parameter should be present")
// 		assert.Contains(t, query.Get("redirect_uri"), "/auth/oauth/apple/callback", "Redirect URI should point to callback endpoint")
// 	})
// }
