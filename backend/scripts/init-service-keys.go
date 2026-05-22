package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/services"
	"github.com/hengadev/encx"
	"github.com/hashicorp/vault/api"
)

// initServiceKeys initializes API keys for all services in Vault
func main() {
	// Get Vault configuration from environment
	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		vaultAddr = "http://localhost:8200"
	}

	vaultToken := os.Getenv("VAULT_TOKEN")
	if vaultToken == "" {
		log.Fatal("VAULT_TOKEN environment variable is required")
	}

	// Initialize Vault client
	config := api.DefaultConfig()
	config.Address = vaultAddr

	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Vault client: %v", err)
	}

	client.SetToken(vaultToken)

	// Initialize crypto service (mock for key generation)
	crypto := &MockCryptoService{}

	// Initialize service key manager
	skm := services.NewServiceKeyManager(client, crypto)

	// Generate keys for all services
	fmt.Println("Generating service API keys...")
	serviceKeys, err := skm.GenerateAllServiceKeys(context.Background())
	if err != nil {
		log.Fatalf("Failed to generate service keys: %v", err)
	}

	// Print generated keys (in production, these should be securely distributed)
	fmt.Println("\n=== GENERATED SERVICE API KEYS ===")
	fmt.Println("IMPORTANT: Store these keys securely and distribute to respective services")
	fmt.Println("These keys will not be displayed again!")

	for serviceName, apiKey := range serviceKeys {
		fmt.Printf("%s_SERVICE_API_KEY=%s\n",
			fmt.Sprintf("%s", serviceName),
			apiKey)
	}

	fmt.Println("\n=== VAULT STORAGE VERIFICATION ===")
	// Verify that keys were stored correctly
	for serviceName := range serviceKeys {
		storedKeys, err := skm.ListServiceKeys()
		if err != nil {
			log.Printf("Warning: Could not verify stored keys: %v", err)
			continue
		}

		found := false
		for _, stored := range storedKeys {
			if stored == serviceName {
				found = true
				break
			}
		}

		if found {
			fmt.Printf("✓ %s API key stored successfully in Vault\n", serviceName)
		} else {
			fmt.Printf("✗ %s API key NOT found in Vault\n", serviceName)
		}
	}

	fmt.Println("\n=== NEXT STEPS ===")
	fmt.Println("1. Add the API keys to your service configurations/environment variables")
	fmt.Println("2. Update your service initializations to use the new middleware constructor")
	fmt.Println("3. Test service-to-service authentication")
	fmt.Println("4. Consider setting up automatic key rotation")
}

// MockCryptoService provides basic hashing for key generation
type MockCryptoService struct{}

func (m *MockCryptoService) HashBasic(ctx context.Context, data []byte) string {
	// Simple implementation for demonstration
	// In production, use the actual encx implementation
	hash := fmt.Sprintf("%x", data)
	return hash
}

// Implement other required encx.CryptoService methods as no-ops
func (m *MockCryptoService) GetPepper() []byte {
	return []byte("mock-pepper")
}

func (m *MockCryptoService) GetArgon2Params() *encx.Argon2Params {
	params, _ := encx.NewArgon2Params(64*1024, 2, 4, 16, 32)
	return params
}

func (m *MockCryptoService) GetAlias() string {
	return "mock-alias"
}

func (m *MockCryptoService) GenerateDEK() ([]byte, error) {
	return []byte("mock-dek"), nil
}

func (m *MockCryptoService) EncryptData(ctx context.Context, plaintext []byte, dek []byte) ([]byte, error) {
	return plaintext, nil
}

func (m *MockCryptoService) DecryptData(ctx context.Context, ciphertext []byte, dek []byte) ([]byte, error) {
	return ciphertext, nil
}

func (m *MockCryptoService) EncryptDEK(ctx context.Context, plaintextDEK []byte) ([]byte, error) {
	return plaintextDEK, nil
}

func (m *MockCryptoService) DecryptDEKWithVersion(ctx context.Context, ciphertextDEK []byte, kekVersion int) ([]byte, error) {
	return ciphertextDEK, nil
}

func (m *MockCryptoService) RotateKEK(ctx context.Context) error {
	return nil
}

func (m *MockCryptoService) HashSecure(ctx context.Context, value []byte) (string, error) {
	return fmt.Sprintf("%x", value), nil
}

func (m *MockCryptoService) CompareSecureHashAndValue(ctx context.Context, value any, hashValue string) (bool, error) {
	return true, nil
}

func (m *MockCryptoService) CompareBasicHashAndValue(ctx context.Context, value any, hashValue string) (bool, error) {
	return true, nil
}

func (m *MockCryptoService) EncryptStream(ctx context.Context, reader io.Reader, writer io.Writer, dek []byte) error {
	return nil
}

func (m *MockCryptoService) DecryptStream(ctx context.Context, reader io.Reader, writer io.Writer, dek []byte) error {
	return nil
}

func (m *MockCryptoService) GetCurrentKEKVersion(ctx context.Context, alias string) (int, error) {
	return 1, nil
}

func (m *MockCryptoService) GetKMSKeyIDForVersion(ctx context.Context, alias string, version int) (string, error) {
	return "mock-key-id", nil
}
