package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hengadev/encx"
	hashicorpkeys "github.com/hengadev/encx/providers/keys/hashicorp"
	hashicorpsecrets "github.com/hengadev/encx/providers/secrets/hashicorp"
)

func main() {
	ctx := context.Background()

	// Test ENCX v0.6.0 API creation
	fmt.Println("🔒 Testing ENCX v0.6.0 API...")

	// Create KMS provider (KeyManagementService) for cryptographic operations
	kms, err := hashicorpkeys.NewTransitService()
	if err != nil {
		log.Fatalf("Failed to create KMS provider: %v", err)
	}
	defer kms.Shutdown()

	// Create secrets provider (SecretManagementService) for pepper storage
	secrets, err := hashicorpsecrets.NewKVStore()
	if err != nil {
		log.Fatalf("Failed to create secrets provider: %v", err)
	}
	defer secrets.Shutdown()

	// Create Config struct
	cfg := encx.Config{
		KEKAlias:    "test-kek",     // Use test encryption key
		PepperAlias: "test-service", // Use service name
	}

	// Create crypto service with v0.6.0 API
	crypto, err := encx.NewCrypto(ctx, kms, secrets, cfg)
	if err != nil {
		log.Fatalf("Failed to create crypto service: %v", err)
	}

	fmt.Println("✅ ENCX v0.6.0 crypto service created successfully")

	// Test encryption/decryption
	testData := struct {
		Email    string `encx:"encrypt,hash_basic"`
		Name     string `encx:"encrypt"`
		Password string `encx:"hash_secure"`
	}{
		Email:    "test@example.com",
		Name:     "John Doe",
		Password: "super-secret-password",
	}

	// Test encryption
	encrypted, err := encx.ProcessStruct(ctx, crypto, &testData)
	if err != nil {
		log.Fatalf("Failed to encrypt data: %v", err)
	}

	fmt.Printf("✅ Data encrypted successfully: %+v\n", encrypted)

	// Test decryption
	var decrypted struct {
		Email    string `encx:"encrypt,hash_basic"`
		Name     string `encx:"encrypt"`
		Password string `encx:"hash_secure"`
	}

	err = encx.ProcessStructInverse(ctx, crypto, encrypted, &decrypted)
	if err != nil {
		log.Fatalf("Failed to decrypt data: %v", err)
	}

	fmt.Printf("✅ Data decrypted successfully: Email=%s, Name=%s, Password=%s\n",
		decrypted.Email, decrypted.Name, decrypted.Password)

	// Verify the data
	if decrypted.Email != testData.Email || decrypted.Name != testData.Name {
		log.Fatalf("❌ Decrypted data doesn't match original")
	}

	// Password should be hashed, not equal to original
	if decrypted.Password == testData.Password {
		log.Fatalf("❌ Password should be hashed, not equal to original")
	}

	fmt.Println("🎉 ENCX v0.6.0 migration verification completed successfully!")
	fmt.Printf("   - KEK Alias: %s\n", cfg.KEKAlias)
	fmt.Printf("   - Pepper Alias: %s\n", cfg.PepperAlias)
	fmt.Printf("   - Test completed at: %s\n", time.Now().Format(time.RFC3339))
}

