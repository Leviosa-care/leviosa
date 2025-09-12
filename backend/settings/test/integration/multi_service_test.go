package helpers

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/core/contracts/services"
	"github.com/Leviosa-care/core/httpx"
	"github.com/Leviosa-care/settings/internal/domain"
	th "github.com/Leviosa-care/settings/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMultiServiceIntegration tests service-to-service communication scenarios
func TestMultiServiceIntegration(t *testing.T) {
	ctx := context.Background()
	
	// Setup test data
	setupMultiServiceTestData(t, ctx)

	t.Run("Service-to-Service Communication", func(t *testing.T) {
		t.Run("catalog service fetching company settings", func(t *testing.T) {
			// Create service client for catalog service
			catalogAPIKey, exists := vaultSetup.GetServiceAPIKey(services.Catalog)
			require.True(t, exists, "Catalog service API key must exist")

			serviceClient, err := httpx.NewServiceClient(httpx.ServiceClientConfig{
				ServiceName: services.Catalog,
				APIKey:      catalogAPIKey,
				BaseURL:     testServerURL,
				Timeout:     10 * time.Second,
			})
			require.NoError(t, err, "Should create service client successfully")

			// Test fetching company name
			resp, err := serviceClient.Get(ctx, "/internal/settings/name")
			require.NoError(t, err, "Should make request successfully")
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK")

			var nameResp domain.GetCompanyNameResponse
			err = json.NewDecoder(resp.Body).Decode(&nameResp)
			require.NoError(t, err, "Should decode response successfully")
			assert.Equal(t, "Multi-Service Test Company", nameResp.Name)
			resp.Body.Close()

			// Test fetching company email
			resp, err = serviceClient.Get(ctx, "/internal/settings/email")
			require.NoError(t, err, "Should make request successfully")
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK")

			var emailResp domain.GetCompanyEmailResponse
			err = json.NewDecoder(resp.Body).Decode(&emailResp)
			require.NoError(t, err, "Should decode response successfully")
			assert.Equal(t, "multiservice@testcompany.com", emailResp.Email)
			resp.Body.Close()
		})

		t.Run("notification service fetching email settings", func(t *testing.T) {
			// Create service client for notification service
			notificationAPIKey, exists := vaultSetup.GetServiceAPIKey(services.Notification)
			require.True(t, exists, "Notification service API key must exist")

			serviceClient, err := httpx.NewServiceClient(httpx.ServiceClientConfig{
				ServiceName: services.Notification,
				APIKey:      notificationAPIKey,
				BaseURL:     testServerURL,
				Timeout:     10 * time.Second,
			})
			require.NoError(t, err, "Should create service client successfully")

			// Notification service would typically fetch email settings
			resp, err := serviceClient.Get(ctx, "/internal/settings/email")
			require.NoError(t, err, "Should make request successfully")
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK")

			var emailResp domain.GetCompanyEmailResponse
			err = json.NewDecoder(resp.Body).Decode(&emailResp)
			require.NoError(t, err, "Should decode response successfully")
			assert.Equal(t, "multiservice@testcompany.com", emailResp.Email)
			resp.Body.Close()
		})

		t.Run("authuser service fetching OTP settings", func(t *testing.T) {
			// Create service client for authuser service
			authuserAPIKey, exists := vaultSetup.GetServiceAPIKey(services.AuthUser)
			require.True(t, exists, "AuthUser service API key must exist")

			serviceClient, err := httpx.NewServiceClient(httpx.ServiceClientConfig{
				ServiceName: services.AuthUser,
				APIKey:      authuserAPIKey,
				BaseURL:     testServerURL,
				Timeout:     10 * time.Second,
			})
			require.NoError(t, err, "Should create service client successfully")

			// AuthUser service would fetch OTP configuration
			resp, err := serviceClient.Get(ctx, "/internal/settings/otp/duration")
			require.NoError(t, err, "Should make request successfully")
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK")

			var otpResp domain.GetOTPDurationResponse
			err = json.NewDecoder(resp.Body).Decode(&otpResp)
			require.NoError(t, err, "Should decode response successfully")
			assert.Equal(t, 600, otpResp.Duration) // 10 minutes
			resp.Body.Close()
		})

		t.Run("settings service self-access", func(t *testing.T) {
			// Settings service should be able to access its own internal endpoints
			settingsAPIKey, exists := vaultSetup.GetServiceAPIKey(services.Settings)
			require.True(t, exists, "Settings service API key must exist")

			serviceClient, err := httpx.NewServiceClient(httpx.ServiceClientConfig{
				ServiceName: services.Settings,
				APIKey:      settingsAPIKey,
				BaseURL:     testServerURL,
				Timeout:     10 * time.Second,
			})
			require.NoError(t, err, "Should create service client successfully")

			// Settings service accessing its own bulk endpoint
			resp, err := serviceClient.Get(ctx, "/internal/settings/bulk")
			require.NoError(t, err, "Should make request successfully")
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK")
			resp.Body.Close()
		})
	})

	t.Run("Cross-Service Security Validation", func(t *testing.T) {
		t.Run("service cannot use another service's API key", func(t *testing.T) {
			// Try to use notification service key with catalog service name
			notificationAPIKey, _ := vaultSetup.GetServiceAPIKey(services.Notification)
			
			client := &http.Client{Timeout: 10 * time.Second}
			req := th.NewInternalGetCompanyNameRequest(t, ctx, testServerURL)
			req.Header.Set(services.ServiceNameHeader, services.Catalog) // Wrong service name
			req.Header.Set(services.ServiceKeyHeader, notificationAPIKey) // Wrong key for catalog

			resp, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Should reject cross-service key usage")
		})

		t.Run("invalid service name should be rejected", func(t *testing.T) {
			catalogAPIKey, _ := vaultSetup.GetServiceAPIKey(services.Catalog)
			
			client := &http.Client{Timeout: 10 * time.Second}
			req := th.NewInternalGetCompanyNameRequest(t, ctx, testServerURL)
			req.Header.Set(services.ServiceNameHeader, "invalid-service") // Invalid service name
			req.Header.Set(services.ServiceKeyHeader, catalogAPIKey)

			resp, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Should reject invalid service names")
		})
	})

	t.Run("Service Client Configuration Validation", func(t *testing.T) {
		t.Run("should reject invalid service name in client", func(t *testing.T) {
			_, err := httpx.NewServiceClient(httpx.ServiceClientConfig{
				ServiceName: "invalid-service",
				APIKey:      "some-key",
				BaseURL:     testServerURL,
			})
			assert.Error(t, err, "Should reject invalid service name")
			assert.Contains(t, err.Error(), "invalid service name", "Error should mention invalid service name")
		})

		t.Run("should reject missing API key", func(t *testing.T) {
			_, err := httpx.NewServiceClient(httpx.ServiceClientConfig{
				ServiceName: services.Catalog,
				APIKey:      "", // Empty API key
				BaseURL:     testServerURL,
			})
			assert.Error(t, err, "Should reject empty API key")
			assert.Contains(t, err.Error(), "API key is required", "Error should mention missing API key")
		})

		t.Run("should reject missing base URL", func(t *testing.T) {
			catalogAPIKey, _ := vaultSetup.GetServiceAPIKey(services.Catalog)
			
			_, err := httpx.NewServiceClient(httpx.ServiceClientConfig{
				ServiceName: services.Catalog,
				APIKey:      catalogAPIKey,
				BaseURL:     "", // Empty base URL
			})
			assert.Error(t, err, "Should reject empty base URL")
			assert.Contains(t, err.Error(), "base URL is required", "Error should mention missing base URL")
		})
	})
}

// TestServiceDiscoveryScenarios tests realistic service discovery patterns
func TestServiceDiscoveryScenarios(t *testing.T) {
	ctx := context.Background()
	setupMultiServiceTestData(t, ctx)

	t.Run("Realistic Service Communication Patterns", func(t *testing.T) {
		t.Run("catalog service workflow", func(t *testing.T) {
			// Realistic scenario: Catalog service needs company info for product displays
			catalogAPIKey, _ := vaultSetup.GetServiceAPIKey(services.Catalog)
			serviceClient, err := httpx.NewServiceClient(httpx.ServiceClientConfig{
				ServiceName: services.Catalog,
				APIKey:      catalogAPIKey,
				BaseURL:     testServerURL,
				Timeout:     10 * time.Second,
			})
			require.NoError(t, err)

			// 1. Fetch company name for product page headers
			resp, err := serviceClient.Get(ctx, "/internal/settings/name")
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			resp.Body.Close()

			// 2. Fetch company address for shipping calculations
			resp, err = serviceClient.Get(ctx, "/internal/settings/address")
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			resp.Body.Close()
		})

		t.Run("notification service workflow", func(t *testing.T) {
			// Realistic scenario: Notification service sending emails
			notificationAPIKey, _ := vaultSetup.GetServiceAPIKey(services.Notification)
			serviceClient, err := httpx.NewServiceClient(httpx.ServiceClientConfig{
				ServiceName: services.Notification,
				APIKey:      notificationAPIKey,
				BaseURL:     testServerURL,
				Timeout:     10 * time.Second,
			})
			require.NoError(t, err)

			// 1. Fetch company name for email headers
			resp, err := serviceClient.Get(ctx, "/internal/settings/name")
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			resp.Body.Close()

			// 2. Fetch company email for from address
			resp, err = serviceClient.Get(ctx, "/internal/settings/email")
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			resp.Body.Close()
		})

		t.Run("authuser service workflow", func(t *testing.T) {
			// Realistic scenario: Auth service configuring OTP
			authuserAPIKey, _ := vaultSetup.GetServiceAPIKey(services.AuthUser)
			serviceClient, err := httpx.NewServiceClient(httpx.ServiceClientConfig{
				ServiceName: services.AuthUser,
				APIKey:      authuserAPIKey,
				BaseURL:     testServerURL,
				Timeout:     10 * time.Second,
			})
			require.NoError(t, err)

			// Fetch OTP configuration
			resp, err := serviceClient.Get(ctx, "/internal/settings/otp/duration")
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			resp.Body.Close()
		})
	})
}

// setupMultiServiceTestData creates test data for multi-service integration tests
func setupMultiServiceTestData(t *testing.T, ctx context.Context) {
	// Clear existing data
	th.ClearSettingsTable(t, ctx, testPool)
	
	// Insert comprehensive test data
	th.InsertTestCompanyName(t, ctx, "Multi-Service Test Company", testPool)
	th.InsertTestCompanyEmail(t, ctx, "multiservice@testcompany.com", testPool)
	th.InsertTestCompanyPhone(t, ctx, "+1-555-MULTI", testPool)
	th.InsertTestCompanyAddress(t, ctx, "456 Service Integration Blvd", testPool)
	th.InsertTestOTPDuration(t, ctx, 600, testPool) // 10 minutes
}