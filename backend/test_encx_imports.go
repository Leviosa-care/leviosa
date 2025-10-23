package main

import (
	"fmt"
	"log"

	"github.com/hengadev/encx"
	hashicorpkeys "github.com/hengadev/encx/providers/keys/hashicorp"
	hashicorpsecrets "github.com/hengadev/encx/providers/secrets/hashicorp"
)

func main() {
	fmt.Println("Testing ENCX v0.6.0 imports...")

	// Test that the new imports compile
	// Note: We can't actually run the code without proper Vault setup,
	// but if this compiles, our import changes are successful

	kms, err := hashicorpkeys.NewTransitService()
	if err != nil {
		log.Printf("Warning: KMS provider creation failed (expected without Vault): %v", err)
	} else {
		fmt.Println("✓ KMS provider import successful")
	}

	secrets, err := hashicorpsecrets.NewKVStore()
	if err != nil {
		log.Printf("Warning: Secrets provider creation failed (expected without Vault): %v", err)
	} else {
		fmt.Println("✓ Secrets provider import successful")
	}

	// Test Config struct creation
	cfg := encx.Config{
		KEKAlias:    "test-key",
		PepperAlias: "test-service",
	}
	fmt.Printf("✓ Config struct created: %+v\n", cfg)

	fmt.Println("✅ ENCX v0.6.0 imports working correctly!")
}