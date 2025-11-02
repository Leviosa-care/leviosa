package services

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateServiceKey(t *testing.T) {
	// We can test key generation without Vault
	skm := &ServiceKeyManager{} // No Vault client needed for key generation

	key, err := skm.GenerateServiceKey()
	require.NoError(t, err)

	// Test key properties
	assert.NotEmpty(t, key, "Generated key should not be empty")
	assert.Greater(t, len(key), 40, "Generated key should be sufficiently long")

	// Test that key is valid base64
	_, err = base64.URLEncoding.DecodeString(key)
	assert.NoError(t, err, "Generated key should be valid base64")

	// Test uniqueness - generate multiple keys
	keys := make(map[string]bool)
	for i := 0; i < 100; i++ {
		key, err := skm.GenerateServiceKey()
		require.NoError(t, err)
		assert.False(t, keys[key], "Generated keys should be unique")
		keys[key] = true
	}
}

func TestServiceKeyValidation(t *testing.T) {
	// Test service name validation logic without requiring Vault
	tests := []struct {
		name        string
		serviceName string
		isValid     bool
	}{
		{"Valid service - authuser", AuthUser, true},
		{"Valid service - catalog", Catalog, true},
		{"Valid service - settings", Settings, true},
		{"Valid service - notification", Notification, true},
		{"Invalid service", "invalid-service", false},
		{"Empty service name", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the IsValidService function directly
			result := IsValidService(tt.serviceName)
			assert.Equal(t, tt.isValid, result)
		})
	}
}

func TestKeyManagementVaultPaths(t *testing.T) {
	// Test that our key management uses the correct Vault paths
	tests := []struct {
		service      string
		expectedPath string
	}{
		{AuthUser, "secret/data/services/authuser/api-key"},
		{Catalog, "secret/data/services/catalog/api-key"},
		{Settings, "secret/data/services/settings/api-key"},
		{Notification, "secret/data/services/notification/api-key"},
	}

	for _, tt := range tests {
		t.Run(tt.service, func(t *testing.T) {
			path := ServiceAPIKeyPath(tt.service)
			assert.Equal(t, tt.expectedPath, path)
		})
	}
}

func TestAllServicesConstant(t *testing.T) {
	// Ensure our key management covers all defined services
	allServices := AllServices()

	// Test that we have the expected services
	expectedServices := []string{"authuser", "catalog", "settings", "notification"}
	assert.ElementsMatch(t, expectedServices, allServices)

	// Test that each service is valid
	for _, service := range allServices {
		assert.True(t, IsValidService(service), "All services from AllServices() should be valid")
	}
}

