package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/services"
	"github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	"github.com/hashicorp/vault/api"
	"github.com/hengadev/encx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServiceAuthWithVault tests the complete service authentication flow with real Vault
func TestServiceAuthWithVault(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Vault integration test in short mode")
	}
	
	ctx := context.Background()
	
	// Setup Vault container
	vaultContainer, err := testutils.SetupVault(ctx, t)
	require.NoError(t, err, "Failed to setup Vault container")
	defer testutils.TeardownVault(ctx, t, vaultContainer)
	
	// Create Vault client
	config := api.DefaultConfig()
	config.Address = vaultContainer.HTTPSEndpoint
	vaultClient, err := api.NewClient(config)
	require.NoError(t, err, "Failed to create Vault client")
	vaultClient.SetToken(vaultContainer.RootToken)
	
	// Initialize encx crypto service for testing
	cryptoService, err := encx.NewTestCrypto(t)
	require.NoError(t, err, "Failed to create test crypto service")
	
	// Initialize service key manager
	keyManager := services.NewServiceKeyManager(vaultClient, cryptoService)
	
	// Generate and store a service API key
	testServiceName := services.Catalog
	apiKey, err := keyManager.GenerateAndStoreServiceKey(ctx, testServiceName)
	require.NoError(t, err, "Failed to generate service key")
	require.NotEmpty(t, apiKey, "Generated API key should not be empty")
	
	// Create auth middleware with Vault client
	middleware := NewSessionAuthMiddleware(nil, cryptoService, vaultClient)
	
	// Test successful authentication
	t.Run("Valid service authentication", func(t *testing.T) {
		// Create request with valid service headers
		req := httptest.NewRequest("GET", "/internal/test", nil)
		req.Header.Set(services.ServiceNameHeader, testServiceName)
		req.Header.Set(services.ServiceKeyHeader, apiKey)
		
		// Add basic context (no logger required for this test)
		ctx := req.Context()
		req = req.WithContext(ctx)
		
		// Create response recorder
		rr := httptest.NewRecorder()
		
		// Track if next handler was called
		nextCalled := false
		var capturedServiceInfo *ServiceInfo
		nextHandler := func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			
			// Extract service info from context
			serviceInfo, err := GetServiceInfoFromContext(r.Context())
			if err == nil {
				capturedServiceInfo = serviceInfo
			}
			
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "success"}`))
		}
		
		// Execute middleware
		handler := middleware.RequireServiceAuth(nextHandler)
		handler(rr, req)
		
		// Assertions
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, nextCalled, "Next handler should have been called")
		assert.NotNil(t, capturedServiceInfo, "Service info should be available in context")
		assert.Equal(t, testServiceName, capturedServiceInfo.Name)
	})
	
	// Test authentication with wrong key
	t.Run("Invalid service key", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/internal/test", nil)
		req.Header.Set(services.ServiceNameHeader, testServiceName)
		req.Header.Set(services.ServiceKeyHeader, "invalid-key-12345")
		
		// Add basic context (no logger required for this test)
		ctx := req.Context()
		req = req.WithContext(ctx)
		
		rr := httptest.NewRecorder()
		
		nextCalled := false
		nextHandler := func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			w.WriteHeader(http.StatusOK)
		}
		
		handler := middleware.RequireServiceAuth(nextHandler)
		handler(rr, req)
		
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.False(t, nextCalled, "Next handler should not have been called")
	})
	
	// Test authentication with non-existent service
	t.Run("Non-existent service", func(t *testing.T) {
		nonExistentService := "nonexistent"
		
		req := httptest.NewRequest("GET", "/internal/test", nil)
		req.Header.Set(services.ServiceNameHeader, nonExistentService)
		req.Header.Set(services.ServiceKeyHeader, apiKey)
		
		// Add basic context (no logger required for this test)
		ctx := req.Context()
		req = req.WithContext(ctx)
		
		rr := httptest.NewRecorder()
		
		nextCalled := false
		nextHandler := func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			w.WriteHeader(http.StatusOK)
		}
		
		handler := middleware.RequireServiceAuth(nextHandler)
		handler(rr, req)
		
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.False(t, nextCalled, "Next handler should not have been called")
	})
}

// TestServiceKeyManager tests the service key management functionality
func TestServiceKeyManager(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Vault integration test in short mode")
	}
	
	ctx := context.Background()
	
	// Setup Vault container
	vaultContainer, err := testutils.SetupVault(ctx, t)
	require.NoError(t, err, "Failed to setup Vault container")
	defer testutils.TeardownVault(ctx, t, vaultContainer)
	
	// Create Vault client
	config := api.DefaultConfig()
	config.Address = vaultContainer.HTTPSEndpoint
	vaultClient, err := api.NewClient(config)
	require.NoError(t, err, "Failed to create Vault client")
	vaultClient.SetToken(vaultContainer.RootToken)
	
	// Initialize crypto service
	cryptoService, err := encx.NewTestCrypto(t)
	require.NoError(t, err, "Failed to create crypto service")
	
	// Initialize service key manager
	keyManager := services.NewServiceKeyManager(vaultClient, cryptoService)
	
	t.Run("Generate and validate single service key", func(t *testing.T) {
		testService := services.Settings
		
		// Generate and store key
		apiKey, err := keyManager.GenerateAndStoreServiceKey(ctx, testService)
		require.NoError(t, err)
		require.NotEmpty(t, apiKey)
		
		// Validate key
		isValid, err := keyManager.ValidateServiceKey(ctx, testService, apiKey)
		require.NoError(t, err)
		assert.True(t, isValid, "Generated key should be valid")
		
		// Test with wrong key
		isValid, err = keyManager.ValidateServiceKey(ctx, testService, "wrong-key")
		require.NoError(t, err)
		assert.False(t, isValid, "Wrong key should not be valid")
	})
	
	t.Run("Generate all service keys", func(t *testing.T) {
		// Generate keys for all services
		serviceKeys, err := keyManager.GenerateAllServiceKeys(ctx)
		require.NoError(t, err)
		
		// Verify all services have keys
		expectedServices := services.AllServices()
		assert.Len(t, serviceKeys, len(expectedServices))
		
		for _, serviceName := range expectedServices {
			apiKey, exists := serviceKeys[serviceName]
			assert.True(t, exists, "Service %s should have an API key", serviceName)
			assert.NotEmpty(t, apiKey, "API key for %s should not be empty", serviceName)
			
			// Validate each key
			isValid, err := keyManager.ValidateServiceKey(ctx, serviceName, apiKey)
			require.NoError(t, err)
			assert.True(t, isValid, "Generated key for %s should be valid", serviceName)
		}
	})
}

// TestEndToEndServiceCommunication tests complete service-to-service communication
func TestEndToEndServiceCommunication(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping end-to-end integration test in short mode")
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Setup Vault container
	vaultContainer, err := testutils.SetupVault(ctx, t)
	require.NoError(t, err)
	defer testutils.TeardownVault(ctx, t, vaultContainer)
	
	// Create Vault client and crypto service
	config := api.DefaultConfig()
	config.Address = vaultContainer.HTTPSEndpoint
	vaultClient, err := api.NewClient(config)
	require.NoError(t, err)
	vaultClient.SetToken(vaultContainer.RootToken)
	
	cryptoService, err := encx.NewTestCrypto(t)
	require.NoError(t, err)
	
	// Generate service keys
	keyManager := services.NewServiceKeyManager(vaultClient, cryptoService)
	catalogKey, err := keyManager.GenerateAndStoreServiceKey(ctx, services.Catalog)
	require.NoError(t, err)
	
	// Create auth middleware
	middleware := NewSessionAuthMiddleware(nil, cryptoService, vaultClient)
	
	// Create a test server with protected endpoint
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		serviceInfo, err := GetServiceInfoFromContext(r.Context())
		if err != nil {
			http.Error(w, "Service info not found", http.StatusInternalServerError)
			return
		}
		
		response := map[string]string{
			"message":     "Hello from protected endpoint",
			"service":     serviceInfo.Name,
			"endpoint":    r.URL.Path,
			"method":      r.Method,
		}
		
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message":"%s","service":"%s","endpoint":"%s","method":"%s"}`,
			response["message"], response["service"], response["endpoint"], response["method"])
	}
	
	// Wrap handler with service auth middleware
	protectedHandler := func(w http.ResponseWriter, r *http.Request) {
		// Apply service auth middleware
		middleware.RequireServiceAuth(testHandler)(w, r)
	}
	
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(protectedHandler))
	defer server.Close()
	
	// Test successful service call
	t.Run("Successful service-to-service call", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, "GET", server.URL+"/internal/test", nil)
		require.NoError(t, err)
		
		// Add service authentication headers
		req.Header.Set(services.ServiceNameHeader, services.Catalog)
		req.Header.Set(services.ServiceKeyHeader, catalogKey)
		
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		
		// Note: In a real test, you'd decode the JSON response here
		// For this example, we'll just verify the status code
	})
	
	// Test failed authentication
	t.Run("Failed authentication with invalid key", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, "GET", server.URL+"/internal/test", nil)
		require.NoError(t, err)
		
		// Add invalid service authentication headers
		req.Header.Set(services.ServiceNameHeader, services.Catalog)
		req.Header.Set(services.ServiceKeyHeader, "invalid-key")
		
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}