package helpers

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/services"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	httpEndpoints "github.com/Leviosa-care/settings/internal/adapters/http"
	"github.com/Leviosa-care/settings/internal/domain"
	th "github.com/Leviosa-care/settings/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestServiceAuthentication make test-integration-test

// TestServiceAuthentication tests the /internal/* endpoints that require service authentication
func TestServiceAuthentication(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Setup test data
	setupServiceAuthTestData(t, ctx)

	t.Run("Valid Service Authentication", func(t *testing.T) {
		t.Run("should allow catalog service to access company name", func(t *testing.T) {
			// Get catalog service API key
			catalogAPIKey, exists := vaultSetup.GetServiceAPIKey(services.Catalog)
			require.True(t, exists, "Catalog service API key should exist")
			require.NotEmpty(t, catalogAPIKey, "Catalog API key should not be empty")

			// Create authenticated request
			req := th.NewInternalGetCompanyNameRequest(t, ctx, testServerURL)
			req.Header.Set(services.ServiceNameHeader, services.Catalog)
			req.Header.Set(services.ServiceKeyHeader, catalogAPIKey)

			// Make request
			resp, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Verify response
			var response domain.GetCompanyNameResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			assert.Equal(t, "Test Company Inc", response.Name)
		})

		t.Run("should allow notification service to access company email", func(t *testing.T) {
			// Get notification service API key
			notificationAPIKey, exists := vaultSetup.GetServiceAPIKey(services.Notification)
			require.True(t, exists, "Notification service API key should exist")

			// Create authenticated request
			req := th.NewInternalGetCompanyEmailRequest(t, ctx, testServerURL)
			req.Header.Set(services.ServiceNameHeader, services.Notification)
			req.Header.Set(services.ServiceKeyHeader, notificationAPIKey)

			// Make request
			resp, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Verify response
			var response domain.GetCompanyEmailResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			assert.Equal(t, "test@company.com", response.Email)
		})

		t.Run("should allow settings service to access its own endpoints", func(t *testing.T) {
			// Get settings service API key
			settingsAPIKey, exists := vaultSetup.GetServiceAPIKey(services.Settings)
			require.True(t, exists, "Settings service API key should exist")

			// Create authenticated request for bulk settings
			req := th.NewInternalBulkSettingsRequest(t, ctx, testServerURL)
			req.Header.Set(services.ServiceNameHeader, services.Settings)
			req.Header.Set(services.ServiceKeyHeader, settingsAPIKey)

			// Make request
			resp, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})

	t.Run("Invalid Service Authentication", func(t *testing.T) {
		t.Run("should reject request with missing service name header", func(t *testing.T) {
			catalogAPIKey, _ := vaultSetup.GetServiceAPIKey(services.Catalog)

			req := th.NewInternalGetCompanyNameRequest(t, ctx, testServerURL)
			// Missing X-Service-Name header
			req.Header.Set(services.ServiceKeyHeader, catalogAPIKey)

			resp, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should reject request with missing API key header", func(t *testing.T) {
			req := th.NewInternalGetCompanyNameRequest(t, ctx, testServerURL)
			req.Header.Set(services.ServiceNameHeader, services.Catalog)
			// Missing X-Service-Key header

			resp, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should reject request with invalid service name", func(t *testing.T) {
			catalogAPIKey, _ := vaultSetup.GetServiceAPIKey(services.Catalog)

			req := th.NewInternalGetCompanyNameRequest(t, ctx, testServerURL)
			req.Header.Set(services.ServiceNameHeader, "invalid-service")
			req.Header.Set(services.ServiceKeyHeader, catalogAPIKey)

			resp, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("should reject request with invalid API key", func(t *testing.T) {
			req := th.NewInternalGetCompanyNameRequest(t, ctx, testServerURL)
			req.Header.Set(services.ServiceNameHeader, services.Catalog)
			req.Header.Set(services.ServiceKeyHeader, "invalid-api-key")

			resp, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("should reject request with wrong service API key", func(t *testing.T) {
			// Try to use notification service key with catalog service name
			notificationAPIKey, _ := vaultSetup.GetServiceAPIKey(services.Notification)

			req := th.NewInternalGetCompanyNameRequest(t, ctx, testServerURL)
			req.Header.Set(services.ServiceNameHeader, services.Catalog)
			req.Header.Set(services.ServiceKeyHeader, notificationAPIKey) // Wrong key for service

			resp, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})

	t.Run("Service Client Integration", func(t *testing.T) {
		t.Run("should work with httpx.ServiceClient", func(t *testing.T) {
			// Create service client for catalog service
			catalogAPIKey, _ := vaultSetup.GetServiceAPIKey(services.Catalog)
			serviceClient, err := httpx.NewServiceClient(httpx.ServiceClientConfig{
				ServiceName: services.Catalog,
				APIKey:      catalogAPIKey,
				BaseURL:     testServerURL,
				Timeout:     10 * time.Second,
			})
			require.NoError(t, err)

			// Make authenticated request using service client
			resp, err := serviceClient.Get(ctx, httpEndpoints.InternalGetCompanyNameEndpoint)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Verify response
			var response domain.GetCompanyNameResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			assert.Equal(t, "Test Company Inc", response.Name)
		})
	})
}

// TEST=TestInternalEndpointsComprehensive make test-integration-test

// TestInternalEndpointsComprehensive tests all internal endpoints with service auth

func TestInternalEndpointsComprehensive(t *testing.T) {
	ctx := context.Background()

	// Setup test data
	setupServiceAuthTestData(t, ctx)

	// Get API key for testing
	catalogAPIKey, exists := vaultSetup.GetServiceAPIKey(services.Catalog)
	require.True(t, exists)

	// Test all internal endpoints
	testCases := []struct {
		name         string
		endpoint     string
		requestFunc  func(t *testing.T, ctx context.Context, serverURL string) *http.Request
		expectedCode int
	}{
		{
			name:         "internal company name",
			endpoint:     httpEndpoints.InternalGetCompanyNameEndpoint,
			requestFunc:  th.NewInternalGetCompanyNameRequest,
			expectedCode: http.StatusOK,
		},
		{
			name:         "internal company email",
			endpoint:     httpEndpoints.InternalGetCompanyEmailEndpoint,
			requestFunc:  th.NewInternalGetCompanyEmailRequest,
			expectedCode: http.StatusOK,
		},
		{
			name:         "internal company phone",
			endpoint:     httpEndpoints.InternalGetCompanyPhoneEndpoint,
			requestFunc:  th.NewInternalGetCompanyPhoneRequest,
			expectedCode: http.StatusOK,
		},
		{
			name:         "internal company address",
			endpoint:     httpEndpoints.InternalGetCompanyAddressEndpoint,
			requestFunc:  th.NewInternalGetCompanyAddressRequest,
			expectedCode: http.StatusOK,
		},
		{
			name:         "internal OTP duration",
			endpoint:     httpEndpoints.InternalGetOTPDurationEndpoint,
			requestFunc:  th.NewInternalGetOTPDurationRequest,
			expectedCode: http.StatusOK,
		},
		{
			name:         "internal bulk settings",
			endpoint:     httpEndpoints.InternalBulkEndpoint,
			requestFunc:  th.NewInternalBulkSettingsRequest,
			expectedCode: http.StatusOK,
		},
	}

	client := &http.Client{Timeout: 10 * time.Second}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := tc.requestFunc(t, ctx, testServerURL)
			req.Header.Set(services.ServiceNameHeader, services.Catalog)
			req.Header.Set(services.ServiceKeyHeader, catalogAPIKey)

			resp, err := client.Do(req)
			require.NoError(t, err, "Request should not fail")
			assert.Equal(t, tc.expectedCode, resp.StatusCode, "Should return expected status code")

			// All internal endpoints should return JSON
			contentType := resp.Header.Get("Content-Type")
			if tc.expectedCode == http.StatusOK {
				assert.Contains(t, contentType, "application/json", "Should return JSON content type")
			}
		})
	}
}

// setupServiceAuthTestData creates test data needed for service authentication tests
func setupServiceAuthTestData(t *testing.T, ctx context.Context) {
	t.Helper()

	// Clear existing data
	th.ClearSettingsTable(t, ctx, testPool)

	// Insert test company data
	th.InsertTestCompanyName(t, ctx, "Test Company Inc", testPool)
	th.InsertTestCompanyEmail(t, ctx, "test@company.com", testPool)
	th.InsertTestCompanyAddress(t, ctx, "123 Test Street", testPool)
	th.InsertTestOTPDuration(t, ctx, 300, testPool) // 5 minutes
	now := time.Now()
	phoneSettings := domain.SettingEncrypted{
		ID:        "858374573",
		Key:       settings.CompanyPhone,
		Value:     "0612345679",
		CreatedAt: now,
		UpdatedAt: now,
	}
	phoneEncrypted, err := domain.ProcessSettingEncryptedEncx(ctx, crypto, &phoneSettings)
	require.NoError(t, err)
	th.InsertTestCompanyPhone(t, ctx, phoneEncrypted, testPool)
}
