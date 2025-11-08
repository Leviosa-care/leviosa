package auth_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	th "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestOAuthCallback TEST_PATH=test/integration/authuser/auth/oauth_callback_test.go

func TestOAuthCallback(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("Error Handling", func(t *testing.T) {
		t.Run("should handle missing OAuth code", func(t *testing.T) {
			// Create OAuth callback request without code parameter
			req := th.NewOAuthCallbackRequest(t, testServerURL, "google", "", "valid_state")

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Without OAuth providers initialized, returns 500 for missing code
			assert.True(t, resp.StatusCode >= 400, "Should return error status for missing OAuth code")
		})

		t.Run("should handle missing OAuth state", func(t *testing.T) {
			// Create OAuth callback request without state parameter
			req := th.NewOAuthCallbackRequest(t, testServerURL, "google", "valid_code", "")

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Without OAuth providers initialized, returns 500 for missing state
			assert.True(t, resp.StatusCode >= 400, "Should return error status for missing OAuth state")
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
			t.Setenv("GOOGLE_CLIENT_ID", "")
			t.Setenv("GOOGLE_CLIENT_SECRET", "")

			// Create OAuth callback request
			req := th.NewOAuthCallbackRequest(t, testServerURL, "google", "code", "state")

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Should handle OAuth configuration errors
			assert.True(t, resp.StatusCode >= 400, "Should return error status for OAuth configuration issues")
		})

		t.Run("should handle malformed callback URLs", func(t *testing.T) {
			// Create OAuth callback request with special characters
			req := th.NewOAuthCallbackRequest(t, testServerURL, "google", "code<script>", "state")

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Should handle malformed input
			assert.True(t, resp.StatusCode >= 400, "Should return error for malformed input")
		})
	})

	t.Run("Route Accessibility", func(t *testing.T) {
		t.Run("should have accessible OAuth callback route for Google", func(t *testing.T) {
			th.ClearUsersTable(t, ctx, testPool)

			// Create a basic OAuth callback request
			req := th.NewOAuthCallbackRequest(t, testServerURL, "google", "test_code", "test_state")

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Route should be accessible (will return 400 for invalid code, but route exists)
			assert.NotEqual(t, http.StatusNotFound, resp.StatusCode, "OAuth callback route should exist")
		})

		t.Run("should have accessible OAuth callback route for Apple", func(t *testing.T) {
			th.ClearUsersTable(t, ctx, testPool)

			// Create a basic OAuth callback request
			req := th.NewOAuthCallbackRequest(t, testServerURL, "apple", "test_code", "test_state")

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Route should be accessible (will return 400 for invalid code, but route exists)
			assert.NotEqual(t, http.StatusNotFound, resp.StatusCode, "OAuth callback route should exist")
		})
	})
}
