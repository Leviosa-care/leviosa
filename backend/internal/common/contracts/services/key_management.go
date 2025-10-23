package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/hengadev/encx"
)

// ServiceKeyManager provides utilities for managing service API keys in Vault
type ServiceKeyManager struct {
	vaultClient *api.Client
	crypto      encx.CryptoService
}

// NewServiceKeyManager creates a new service key manager
func NewServiceKeyManager(vaultClient *api.Client, crypto encx.CryptoService) *ServiceKeyManager {
	return &ServiceKeyManager{
		vaultClient: vaultClient,
		crypto:      crypto,
	}
}

// GenerateServiceKey generates a cryptographically secure API key
func (skm *ServiceKeyManager) GenerateServiceKey() (string, error) {
	// Generate 32 bytes of random data (256 bits)
	keyBytes := make([]byte, 32)
	_, err := rand.Read(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}

	// Encode as base64 for safe transmission
	key := base64.URLEncoding.EncodeToString(keyBytes)
	return key, nil
}

// StoreServiceKey stores a service API key in Vault with proper hashing
func (skm *ServiceKeyManager) StoreServiceKey(ctx context.Context, serviceName, apiKey string) error {
	if !IsValidService(serviceName) {
		return fmt.Errorf("invalid service name: %s", serviceName)
	}

	// Hash the API key for storage
	apiKeyBytes, err := encx.SerializeValue(apiKey)
	if err != nil {
		return fmt.Errorf("failed to serialize API key for hashing: %w", err)
	}
	keyHash := skm.crypto.HashBasic(ctx, apiKeyBytes)

	// Prepare Vault data structure
	vaultData := map[string]interface{}{
		"data": map[string]interface{}{
			"key_hash":     keyHash,
			"service_name": serviceName,
			"created_at":   fmt.Sprintf("%d", ctx.Value("timestamp")),
		},
	}

	// Store in Vault at the correct path
	vaultPath := ServiceAPIKeyPath(serviceName)
	_, err = skm.vaultClient.Logical().Write(vaultPath, vaultData)
	if err != nil {
		return fmt.Errorf("failed to store service key in Vault at %s: %w", vaultPath, err)
	}

	return nil
}

// GenerateAndStoreServiceKey generates a new API key and stores it in Vault
func (skm *ServiceKeyManager) GenerateAndStoreServiceKey(ctx context.Context, serviceName string) (string, error) {
	// Generate new API key
	apiKey, err := skm.GenerateServiceKey()
	if err != nil {
		return "", fmt.Errorf("failed to generate service key for %s: %w", serviceName, err)
	}

	// Store in Vault
	err = skm.StoreServiceKey(ctx, serviceName, apiKey)
	if err != nil {
		return "", fmt.Errorf("failed to store service key for %s: %w", serviceName, err)
	}

	return apiKey, nil
}

// GenerateAllServiceKeys generates and stores API keys for all defined services
func (skm *ServiceKeyManager) GenerateAllServiceKeys(ctx context.Context) (map[string]string, error) {
	serviceKeys := make(map[string]string)

	for _, serviceName := range AllServices() {
		apiKey, err := skm.GenerateAndStoreServiceKey(ctx, serviceName)
		if err != nil {
			return nil, fmt.Errorf("failed to generate key for service %s: %w", serviceName, err)
		}
		serviceKeys[serviceName] = apiKey
	}

	return serviceKeys, nil
}

// ValidateServiceKey validates a service key against Vault storage (for testing/verification)
func (skm *ServiceKeyManager) ValidateServiceKey(ctx context.Context, serviceName, providedKey string) (bool, error) {
	if !IsValidService(serviceName) {
		return false, fmt.Errorf("invalid service name: %s", serviceName)
	}

	// Read from Vault
	vaultPath := ServiceAPIKeyPath(serviceName)
	secret, err := skm.vaultClient.Logical().Read(vaultPath)
	if err != nil {
		return false, fmt.Errorf("failed to read service key from Vault: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return false, fmt.Errorf("service key not found in Vault for service: %s", serviceName)
	}

	// Extract stored key hash
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("invalid Vault response format")
	}

	storedKeyHash, ok := data["key_hash"].(string)
	if !ok || storedKeyHash == "" {
		return false, fmt.Errorf("missing or invalid key_hash in Vault")
	}

	// Compare hashes
	providedKeyBytes, err := encx.SerializeValue(providedKey)
	if err != nil {
		return false, fmt.Errorf("failed to serialize provided key for hashing: %w", err)
	}
	providedKeyHash := skm.crypto.HashBasic(ctx, providedKeyBytes)
	return storedKeyHash == providedKeyHash, nil
}

// DeleteServiceKey removes a service API key from Vault
func (skm *ServiceKeyManager) DeleteServiceKey(serviceName string) error {
	if !IsValidService(serviceName) {
		return fmt.Errorf("invalid service name: %s", serviceName)
	}

	vaultPath := ServiceAPIKeyPath(serviceName)
	_, err := skm.vaultClient.Logical().Delete(vaultPath)
	if err != nil {
		return fmt.Errorf("failed to delete service key from Vault at %s: %w", vaultPath, err)
	}

	return nil
}

// ListServiceKeys returns a list of services that have API keys stored in Vault
func (skm *ServiceKeyManager) ListServiceKeys() ([]string, error) {
	// List all service API key paths
	listPath := "secret/metadata/services"
	secret, err := skm.vaultClient.Logical().List(listPath)
	if err != nil {
		return nil, fmt.Errorf("failed to list service keys: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return []string{}, nil
	}

	// Extract service names from the keys
	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return []string{}, nil
	}

	serviceNames := make([]string, 0, len(keys))
	for _, key := range keys {
		if serviceName, ok := key.(string); ok {
			serviceNames = append(serviceNames, serviceName)
		}
	}

	return serviceNames, nil
}
