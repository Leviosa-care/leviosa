package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/core/contracts/services"
	"github.com/hashicorp/vault/api"
	"github.com/hengadev/encx"
	"github.com/hengadev/encx/providers/hashicorpvault"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	EncryptionKey = "leviosa-app-key"
)

type VaultContainer struct {
	testcontainers.Container
	HTTPSEndpoint string
	RootToken     string
}

func SetupVault(ctx context.Context, t *testing.T) (*VaultContainer, error) {
	rootToken := "test-root-token"

	req := testcontainers.ContainerRequest{
		Image:        "hashicorp/vault:1.19",
		ExposedPorts: []string{"8200/tcp"},
		WaitingFor: wait.ForHTTP("/v1/sys/health").WithPort("8200/tcp").WithStartupTimeout(60 * time.Second).WithStatusCodeMatcher(func(status int) bool {
			return status == 200 || status == 429 || status == 473 || status == 503
		}),
		Env: map[string]string{
			"VAULT_DEV_ROOT_TOKEN_ID":  rootToken,
			"VAULT_DEV_LISTEN_ADDRESS": "0.0.0.0:8200",
			"VAULT_ADDR":               "http://0.0.0.0:8200",
			"VAULT_API_ADDR":           "http://0.0.0.0:8200",
			"VAULT_DISABLE_MLOCK":      "true",
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start vault container: %w", err)
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault host IP: %w", err)
	}

	port, err := container.MappedPort(ctx, "8200")
	if err != nil {
		return nil, fmt.Errorf("failed to get vault mapped port: %w", err)
	}

	httpsEndpoint := fmt.Sprintf("http://%s:%s", hostIP, port.Port())

	vaultContainer := &VaultContainer{
		Container:     container,
		HTTPSEndpoint: httpsEndpoint,
		RootToken:     rootToken,
	}

	// Initialize required secrets and engines
	if err := initializeVaultSecrets(vaultContainer); err != nil {
		return nil, fmt.Errorf("failed to initialize vault secrets: %w", err)
	}

	// Verify the secret was created correctly
	if err := verifyVaultSecret(vaultContainer, "secret/data/pepper"); err != nil {
		return nil, fmt.Errorf("failed to verify vault secrets: %w", err)
	}

	return vaultContainer, nil
}

func TeardownVault(ctx context.Context, t *testing.T, container *VaultContainer) {
	if err := container.Terminate(ctx); err != nil {
		t.Fatalf("failed to terminate vault container: %v", err)
	}
}

// initializeVaultSecrets creates the required secrets and enables engines in Vault for testing
func initializeVaultSecrets(vaultContainer *VaultContainer) error {
	// 1. Enable Transit secrets engine (required for encx)
	if err := enableTransitEngine(vaultContainer); err != nil {
		return fmt.Errorf("failed to enable transit engine: %w", err)
	}

	// 2. Create the encryption key for encx
	if err := createTransitKey(vaultContainer, EncryptionKey); err != nil {
		return fmt.Errorf("failed to create transit key: %w", err)
	}

	// 3. Create the pepper secret required by encx (must be exactly 32 characters)
	pepperData := map[string]any{
		"data": map[string]any{
			"value": "testpepper123456testpepper123456", // Exactly 32 chars
		},
	}

	if err := createVaultSecret(vaultContainer, "secret/data/pepper", pepperData); err != nil {
		return fmt.Errorf("failed to create pepper secret: %w", err)
	}

	return nil
}

// enableTransitEngine enables the Transit secrets engine at the default path
func enableTransitEngine(vaultContainer *VaultContainer) error {
	url := fmt.Sprintf("%s/v1/sys/mounts/transit", vaultContainer.HTTPSEndpoint)

	mountData := map[string]any{
		"type":        "transit",
		"description": "Transit engine for encx testing",
	}

	jsonData, err := json.Marshal(mountData)
	if err != nil {
		return fmt.Errorf("failed to marshal mount data: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create mount request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Vault-Token", vaultContainer.RootToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to enable transit engine: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		// Check if engine already exists
		bodyBytes, _ := io.ReadAll(resp.Body)
		if strings.Contains(string(bodyBytes), "path is already in use") {
			fmt.Println("Transit engine already enabled, continuing...")
			return nil
		}
		return fmt.Errorf("vault returned status %d when enabling transit engine: %s", resp.StatusCode, string(bodyBytes))
	}
	
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("vault returned status %d when enabling transit engine", resp.StatusCode)
	}

	fmt.Println("Transit engine enabled successfully")
	return nil
}

// createTransitKey creates a new encryption key in the Transit engine
func createTransitKey(vaultContainer *VaultContainer, keyName string) error {
	url := fmt.Sprintf("%s/v1/transit/keys/%s", vaultContainer.HTTPSEndpoint, keyName)

	keyData := map[string]any{
		"type": "aes256-gcm96", // Compatible with encx requirements
	}

	jsonData, err := json.Marshal(keyData)
	if err != nil {
		return fmt.Errorf("failed to marshal key data: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create key request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Vault-Token", vaultContainer.RootToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create transit key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("vault returned status %d when creating transit key %s", resp.StatusCode, keyName)
	}

	fmt.Printf("Transit key '%s' created successfully\n", keyName)
	return nil
}

// createVaultSecret creates a secret in Vault using the HTTP API
func createVaultSecret(vaultContainer *VaultContainer, path string, data map[string]any) error {
	url := fmt.Sprintf("%s/v1/%s", vaultContainer.HTTPSEndpoint, path)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal secret data: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Vault-Token", vaultContainer.RootToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("vault returned status %d when creating secret at %s", resp.StatusCode, path)
	}

	return nil
}

// verifyVaultSecret verifies that a secret exists and can be read from Vault
func verifyVaultSecret(vaultContainer *VaultContainer, path string) error {
	url := fmt.Sprintf("%s/v1/%s", vaultContainer.HTTPSEndpoint, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create verification request: %w", err)
	}

	req.Header.Set("X-Vault-Token", vaultContainer.RootToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to verify secret: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("vault returned status %d when verifying secret at %s", resp.StatusCode, path)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode verification response: %w", err)
	}

	fmt.Printf("Vault secret verification for %s: %+v\n", path, result)
	return nil
}

// InitializeServiceKeys creates API keys for all services in the test Vault instance
func InitializeServiceKeys(vaultContainer *VaultContainer, cryptoService encx.CryptoService) (map[string]string, error) {
	// Create Vault API client
	config := &VaultClientConfig{
		Address: vaultContainer.HTTPSEndpoint,
		Token:   vaultContainer.RootToken,
	}

	vaultClient, err := NewVaultClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault client: %w", err)
	}

	// Create service key manager
	keyManager := services.NewServiceKeyManager(vaultClient, cryptoService)

	// Generate all service keys
	serviceKeys, err := keyManager.GenerateAllServiceKeys(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to generate service keys: %w", err)
	}

	fmt.Printf("Generated API keys for %d services\n", len(serviceKeys))
	return serviceKeys, nil
}

// VaultClientConfig holds configuration for creating a Vault client
type VaultClientConfig struct {
	Address string
	Token   string
}

// NewVaultClient creates a new Vault API client with the given configuration
func NewVaultClient(config *VaultClientConfig) (*api.Client, error) {
	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = config.Address

	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault client: %w", err)
	}

	client.SetToken(config.Token)
	return client, nil
}

// Helper functions for per-service encryption setup

// GetServiceEncryptionKeyName returns the transit key name for a service
func GetServiceEncryptionKeyName(serviceName string) string {
	return fmt.Sprintf("%s-encryption-key", serviceName)
}

// GetServicePepperPath returns the pepper secret path for a service
func GetServicePepperPath(serviceName string) string {
	return fmt.Sprintf("secret/data/peppers/%s", serviceName)
}

// createServiceEncryptionKey creates a service-specific encryption key in Vault transit engine
func createServiceEncryptionKey(vaultContainer *VaultContainer, serviceName string) error {
	keyName := GetServiceEncryptionKeyName(serviceName)
	url := fmt.Sprintf("%s/v1/transit/keys/%s", vaultContainer.HTTPSEndpoint, keyName)
	
	payload := map[string]any{
		"type": "aes256-gcm96",
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal key creation payload: %w", err)
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create key creation request: %w", err)
	}
	
	req.Header.Set("X-Vault-Token", vaultContainer.RootToken)
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create service encryption key: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		// Key already exists, which is fine
		fmt.Printf("✓ Service encryption key already exists: %s\n", keyName)
		return nil
	}
	
	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create service encryption key %s, status: %d, body: %s", keyName, resp.StatusCode, string(body))
	}
	
	fmt.Printf("✓ Created service encryption key: %s\n", keyName)
	return nil
}

// createServicePepper creates a service-specific pepper secret in Vault KV store
func createServicePepper(vaultContainer *VaultContainer, serviceName string) error {
	pepperPath := GetServicePepperPath(serviceName)
	
	// Generate a unique 32-character pepper for this service
	pepper := fmt.Sprintf("test%s%s", serviceName, strings.Repeat("x", 28-len(serviceName)))
	if len(pepper) > 32 {
		pepper = pepper[:32]
	}
	if len(pepper) < 32 {
		pepper = pepper + strings.Repeat("y", 32-len(pepper))
	}
	
	pepperData := map[string]any{
		"data": map[string]any{
			"value": pepper,
		},
	}
	
	return createVaultSecret(vaultContainer, pepperPath, pepperData)
}

// CreateServiceCryptoService creates a service-specific crypto service with isolated encryption
func CreateServiceCryptoService(ctx context.Context, vaultContainer *VaultContainer, serviceName string) (encx.CryptoService, error) {
	// Ensure the service encryption key exists
	if err := createServiceEncryptionKey(vaultContainer, serviceName); err != nil {
		return nil, fmt.Errorf("failed to create service encryption key: %w", err)
	}
	
	// Ensure the service pepper exists
	if err := createServicePepper(vaultContainer, serviceName); err != nil {
		return nil, fmt.Errorf("failed to create service pepper: %w", err)
	}
	
	// Create Vault client for this service
	config := &VaultClientConfig{
		Address: vaultContainer.HTTPSEndpoint,
		Token:   vaultContainer.RootToken,
	}
	
	// Create KMS provider using HashiCorp Vault
	// Set environment variables for the hashicorpvault.New() function
	originalAddr := os.Getenv("VAULT_ADDR")
	originalToken := os.Getenv("VAULT_TOKEN")
	
	os.Setenv("VAULT_ADDR", config.Address)
	os.Setenv("VAULT_TOKEN", config.Token)
	
	kms, err := hashicorpvault.New()
	
	// Restore original environment variables
	os.Setenv("VAULT_ADDR", originalAddr)
	os.Setenv("VAULT_TOKEN", originalToken)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create KMS provider: %w", err)
	}
	
	// Create crypto service with service-specific keys
	serviceKeyName := GetServiceEncryptionKeyName(serviceName)
	servicePepperPath := GetServicePepperPath(serviceName)
	
	crypto, err := encx.New(ctx, kms, serviceKeyName, servicePepperPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create crypto service for %s: %w", serviceName, err)
	}
	
	fmt.Printf("✓ Created crypto service for service: %s\n", serviceName)
	return crypto, nil
}

// ServiceVaultSetup contains the complete setup for service-specific Vault configuration
type ServiceVaultSetup struct {
	VaultContainer  *VaultContainer
	ServiceKeys     map[string]string               // service name -> API key
	CryptoServices  map[string]encx.CryptoService   // service name -> crypto service
	VaultClient     *api.Client                     // Vault client for auth middleware
}

// InitializeServiceVault sets up Vault with per-service encryption and service authentication
// This function creates the complete GDPR-compliant setup for multiple services
func InitializeServiceVault(ctx context.Context, vaultContainer *VaultContainer, serviceNames []string) (*ServiceVaultSetup, error) {
	fmt.Printf("=== Initializing Service Vault Setup ===\n")
	fmt.Printf("Services: %v\n", serviceNames)
	
	// Initialize the basic Vault setup (KV and Transit engines, shared encryption key)
	if err := initializeVaultSecrets(vaultContainer); err != nil {
		return nil, fmt.Errorf("failed to initialize basic Vault secrets: %w", err)
	}
	
	// Create Vault client for auth middleware and service key operations
	vaultClient, err := NewVaultClient(&VaultClientConfig{
		Address: vaultContainer.HTTPSEndpoint,
		Token:   vaultContainer.RootToken,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault client: %w", err)
	}
	
	// Create a shared crypto service for service key hashing
	// This uses the original shared key for API key operations
	// Set environment variables for the hashicorpvault.New() function
	originalAddr := os.Getenv("VAULT_ADDR")
	originalToken := os.Getenv("VAULT_TOKEN")
	
	os.Setenv("VAULT_ADDR", vaultContainer.HTTPSEndpoint)
	os.Setenv("VAULT_TOKEN", vaultContainer.RootToken)
	
	kms, err := hashicorpvault.New()
	
	// Restore original environment variables
	os.Setenv("VAULT_ADDR", originalAddr)
	os.Setenv("VAULT_TOKEN", originalToken)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create KMS provider: %w", err)
	}
	
	sharedCrypto, err := encx.New(ctx, kms, EncryptionKey, "secret/data/pepper")
	if err != nil {
		return nil, fmt.Errorf("failed to create shared crypto service: %w", err)
	}
	
	// Generate service API keys using the existing function
	serviceKeys, err := InitializeServiceKeys(vaultContainer, sharedCrypto)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize service keys: %w", err)
	}
	
	// Create per-service crypto services for data encryption
	cryptoServices := make(map[string]encx.CryptoService)
	
	for _, serviceName := range serviceNames {
		fmt.Printf("Creating service-specific crypto for: %s\n", serviceName)
		
		serviceCrypto, err := CreateServiceCryptoService(ctx, vaultContainer, serviceName)
		if err != nil {
			return nil, fmt.Errorf("failed to create crypto service for %s: %w", serviceName, err)
		}
		
		cryptoServices[serviceName] = serviceCrypto
	}
	
	setup := &ServiceVaultSetup{
		VaultContainer: vaultContainer,
		ServiceKeys:    serviceKeys,
		CryptoServices: cryptoServices,
		VaultClient:    vaultClient,
	}
	
	fmt.Printf("✅ Service Vault setup complete!\n")
	fmt.Printf("   - %d service API keys generated\n", len(serviceKeys))
	fmt.Printf("   - %d service-specific crypto services created\n", len(cryptoServices))
	fmt.Printf("   - GDPR-compliant data isolation enabled\n")
	
	return setup, nil
}

// SetupServiceVault is a convenience function that creates a Vault container and initializes service-specific setup
// This is the main function that tests should use for GDPR-compliant service testing
func SetupServiceVault(ctx context.Context, t *testing.T, serviceNames []string) (*ServiceVaultSetup, error) {
	// Create Vault container
	vaultContainer, err := SetupVault(ctx, t)
	if err != nil {
		return nil, fmt.Errorf("failed to setup Vault container: %w", err)
	}
	
	// Initialize service-specific configuration
	setup, err := InitializeServiceVault(ctx, vaultContainer, serviceNames)
	if err != nil {
		// Clean up container on failure
		if t != nil {
			TeardownVault(ctx, t, vaultContainer)
		}
		return nil, fmt.Errorf("failed to initialize service vault: %w", err)
	}
	
	return setup, nil
}

// GetServiceCrypto is a convenience method to get crypto service for a specific service
func (s *ServiceVaultSetup) GetServiceCrypto(serviceName string) (encx.CryptoService, bool) {
	crypto, exists := s.CryptoServices[serviceName]
	return crypto, exists
}

// GetServiceAPIKey is a convenience method to get API key for a specific service
func (s *ServiceVaultSetup) GetServiceAPIKey(serviceName string) (string, bool) {
	key, exists := s.ServiceKeys[serviceName]
	return key, exists
}
