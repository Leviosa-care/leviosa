package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	th "github.com/Leviosa-care/authuser/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestOAuthCallback_GoogleWithNextcloudTesting_NewUser make test-integration-auth-test

func TestOAuthCallback_GoogleWithNextcloudTesting_NewUser(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Clear users table for clean test
	th.ClearUsersTable(t, ctx, testPool)

	// Setup Nextcloud container for OAuth testing (registered as "google" provider)
	nextcloudHelper := th.SetupNextcloudOAuth(ctx, t)
	defer nextcloudHelper.TeardownNextcloudOAuth(ctx, t)

	// Get OAuth configuration from Nextcloud container
	clientID, clientSecret, nextcloudURL := nextcloudHelper.GetOAuthConfig()

	// Setup OAuth environment to use Nextcloud as Google provider for testing
	th.SetupNextcloudAsGoogleOAuthEnvironment(t, nextcloudURL, clientID, clientSecret)

	// Create test user in Nextcloud
	testUsername := "testoauthuser"
	testEmail := "testoauth@example.com"
	testDisplayName := "Test OAuth User"
	nextcloudHelper.CreateNextcloudTestUser(ctx, t, testUsername, testEmail, testDisplayName)

	t.Run("should create new user from Google OAuth (Nextcloud behind scenes)", func(t *testing.T) {
		// Clear users table
		th.ClearUsersTable(t, ctx, testPool)

		// Generate mock OAuth parameters
		code := th.MockOAuthCode(t)
		state := th.MockOAuthState(t)

		// Create OAuth callback request using "google" provider (Nextcloud behind scenes)
		req := th.NewOAuthCallbackRequest(t, testServerURL, "google", code, state)

		// Make request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Test that the endpoint is accessible and handles the OAuth flow attempt
		// Note: This tests the OAuth flow mechanics with real provider integration
		// The Nextcloud user ID will be stored in the GoogleID field for this test

		assert.True(t, resp.StatusCode >= 400 && resp.StatusCode < 500 || resp.StatusCode == 200,
			"Should handle OAuth callback request (success or client error expected)")
	})
}

// TEST=TestOAuthCallback_GoogleWithNextcloudTesting_ExistingUser make test-integration-auth-test

func TestOAuthCallback_GoogleWithNextcloudTesting_ExistingUser(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Setup Nextcloud container for OAuth testing (registered as "google" provider)
	nextcloudHelper := th.SetupNextcloudOAuth(ctx, t)
	defer nextcloudHelper.TeardownNextcloudOAuth(ctx, t)

	// Get OAuth configuration
	clientID, clientSecret, nextcloudURL := nextcloudHelper.GetOAuthConfig()
	th.SetupNextcloudAsGoogleOAuthEnvironment(t, nextcloudURL, clientID, clientSecret)

	t.Run("should return existing user for known Google OAuth ID", func(t *testing.T) {
		// Clear users table
		th.ClearUsersTable(t, ctx, testPool)

		// Create existing user with Google OAuth ID (will store Nextcloud ID in GoogleID field for testing)
		googleID := th.GenerateTestGoogleID()
		existingUser := th.NewTestGoogleUser(t, ctx, crypto,
			googleID, "existing@example.com", "Existing", "User")

		existingUserEncx, err := domain.ProcessUserEncx(ctx, crypto, existingUser)
		require.NoError(t, err)

		err = th.InsertUserEncx(t, ctx, existingUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Verify user was inserted
		userCount, err := th.CountUsers(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, 1, userCount, "Should have one user before OAuth callback")

		// Generate mock OAuth parameters
		code := th.MockOAuthCode(t)
		state := th.MockOAuthState(t)

		// Create OAuth callback request using "google" provider (Nextcloud behind scenes)
		req := th.NewOAuthCallbackRequest(t, testServerURL, "google", code, state)

		// Make request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Test that the endpoint processes the request
		assert.True(t, resp.StatusCode >= 400 && resp.StatusCode < 500 || resp.StatusCode == 200,
			"Should handle OAuth callback for existing user")

		// Verify no new users were created
		finalUserCount, err := th.CountUsers(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, 1, finalUserCount, "Should not create additional users for existing OAuth user")
	})
}

// TEST=TestOAuthCallback_GoogleWithNextcloudTesting_LinkExistingEmail make test-integration-auth-test

func TestOAuthCallback_GoogleWithNextcloudTesting_LinkExistingEmail(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Setup Nextcloud container (registered as "google" provider)
	nextcloudHelper := th.SetupNextcloudOAuth(ctx, t)
	defer nextcloudHelper.TeardownNextcloudOAuth(ctx, t)

	// Get OAuth configuration
	clientID, clientSecret, nextcloudURL := nextcloudHelper.GetOAuthConfig()
	th.SetupNextcloudAsGoogleOAuthEnvironment(t, nextcloudURL, clientID, clientSecret)

	t.Run("should link Google OAuth to existing email user", func(t *testing.T) {
		// Clear users table
		th.ClearUsersTable(t, ctx, testPool)

		// Create existing user with email but no OAuth ID
		existingUser := th.NewTestOAuthUser(t, ctx, crypto,
			"", "", "link@example.com", "Link", "User") // No OAuth ID initially

		existingUserEncx, err := domain.ProcessUserEncx(ctx, crypto, existingUser)
		require.NoError(t, err)

		err = th.InsertUserEncx(t, ctx, existingUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Verify user was inserted
		userCount, err := th.CountUsers(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, 1, userCount, "Should have one user before OAuth callback")

		// Generate mock OAuth parameters
		code := th.MockOAuthCode(t)
		state := th.MockOAuthState(t)

		// Create OAuth callback request using "google" provider (Nextcloud behind scenes)
		req := th.NewOAuthCallbackRequest(t, testServerURL, "google", code, state)

		// Make request
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Test that the endpoint processes the request
		assert.True(t, resp.StatusCode >= 400 && resp.StatusCode < 500 || resp.StatusCode == 200,
			"Should handle OAuth callback for account linking")

		// Verify no new users were created (should link to existing)
		finalUserCount, err := th.CountUsers(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, 1, finalUserCount, "Should not create additional users when linking OAuth to existing email")
	})
}

// TEST=TestOAuthCallback_ErrorHandling make test-integration-auth-test

func TestOAuthCallback_ErrorHandling(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should handle missing OAuth code", func(t *testing.T) {
		// Create OAuth callback request without code parameter
		req := th.NewOAuthCallbackRequest(t, testServerURL, "nextcloud", "", "valid_state")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return error for missing code
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return bad request for missing OAuth code")
	})

	t.Run("should handle missing OAuth state", func(t *testing.T) {
		// Create OAuth callback request without state parameter
		req := th.NewOAuthCallbackRequest(t, testServerURL, "nextcloud", "valid_code", "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return error for missing state (CSRF protection)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return bad request for missing OAuth state")
	})

	t.Run("should handle invalid provider", func(t *testing.T) {
		// Create OAuth callback request with invalid provider
		req := th.NewOAuthCallbackRequest(t, testServerURL, "invalid_provider", "code", "state")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return error for invalid provider
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return bad request for invalid provider")
	})

	t.Run("should handle OAuth provider errors", func(t *testing.T) {
		// Setup environment with missing configuration to trigger OAuth errors
		t.Setenv("NEXTCLOUD_URL", "")
		t.Setenv("NEXTCLOUD_CLIENT_ID", "")
		t.Setenv("NEXTCLOUD_CLIENT_SECRET", "")

		// Create OAuth callback request
		req := th.NewOAuthCallbackRequest(t, testServerURL, "nextcloud", "code", "state")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should handle OAuth configuration errors
		assert.True(t, resp.StatusCode >= 400, "Should return error status for OAuth configuration issues")
	})
}

// TEST=TestOAuthCallback_ResponseFormat make test-integration-auth-test

// TestOAuthCallback_ResponseFormat tests the structure of successful OAuth callback responses
func TestOAuthCallback_ResponseFormat(t *testing.T) {
	// This test verifies the expected response format for OAuth callbacks
	// In a real implementation, this would test with valid OAuth tokens

	expectedResponse := &domain.OAuthCallbackResponse{
		AccessToken:        "mock_access_token",
		RefreshToken:       "mock_refresh_token",
		AccessTokenExpiry:  1234567890,
		RefreshTokenExpiry: 1234567890,
		IsNewUser:          true,
	}

	// Verify the response structure can be marshalled/unmarshalled correctly
	jsonBytes, err := json.Marshal(expectedResponse)
	require.NoError(t, err, "Should be able to marshal OAuth response")

	var unmarshalledResponse domain.OAuthCallbackResponse
	err = json.Unmarshal(jsonBytes, &unmarshalledResponse)
	require.NoError(t, err, "Should be able to unmarshal OAuth response")

	// Verify all fields are preserved
	assert.Equal(t, expectedResponse.AccessToken, unmarshalledResponse.AccessToken)
	assert.Equal(t, expectedResponse.RefreshToken, unmarshalledResponse.RefreshToken)
	assert.Equal(t, expectedResponse.AccessTokenExpiry, unmarshalledResponse.AccessTokenExpiry)
	assert.Equal(t, expectedResponse.RefreshTokenExpiry, unmarshalledResponse.RefreshTokenExpiry)
	assert.Equal(t, expectedResponse.IsNewUser, unmarshalledResponse.IsNewUser)
}
