package sessionRepository_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Leviosa-care/authuser/internal/adapters/redis/session"
	"github.com/Leviosa-care/authuser/internal/ports"

	tu "github.com/Leviosa-care/core/testutils"
	"github.com/hengadev/encx"
	"github.com/hengadev/encx/providers/hashicorpvault"
	"github.com/redis/go-redis/v9"
)

var (
	redisContainer *tu.RedisContainer
	vaultContainer *tu.VaultContainer
	testClient     *redis.Client
	crypto         encx.CryptoService
	repo           ports.SessionRepository
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error

	// Setup Vault testcontainer for crypto
	log.Println("Setting up Vault testcontainer...")
	vaultContainer, err = tu.SetupVault(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to setup Vault container: %v", err))
	}
	defer tu.TeardownVault(ctx, nil, vaultContainer)

	// Set environment variables for Vault
	os.Setenv("VAULT_ADDR", vaultContainer.HTTPSEndpoint)
	os.Setenv("VAULT_TOKEN", vaultContainer.RootToken)

	// Crypto service
	log.Println("Creating crypto service...")
	kms, err := hashicorpvault.New()
	if err != nil {
		panic(fmt.Sprintf("Failed to create vault provider: %v", err))
	}

	crypto, err = encx.New(
		ctx,
		kms,
		tu.EncryptionKey,
		"secret/data/pepper",
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create crypto service: %v", err))
	}
	if crypto == nil {
		panic("Crypto service is nil after creation")
	}
	log.Println("Crypto service created successfully")

	// Redis container
	redisContainer, err = tu.SetupRedis(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to setup redis container: %v", err))
	}
	defer tu.TeardownRedis(ctx, nil, redisContainer)

	// Redis client
	log.Println("Creating Redis client...")
	testClient = redisContainer.NewClient()

	// Test Redis connection
	if err = testClient.Ping(ctx).Err(); err != nil {
		tu.TeardownRedis(ctx, nil, redisContainer)
		panic(fmt.Sprintf("Failed to ping Redis: %v", err))
	}
	log.Println("Redis client connected successfully.")

	// Create repository instance
	repo = sessionRepository.New(testClient)
	log.Println("Session Repository created.")

	// Run tests
	code := m.Run()

	// Cleanup
	if testClient != nil {
		testClient.Close()
	}

	log.Println("Test(s) executed")

	// Exit with the test result code
	os.Exit(code)
}

// reconnectRedis is a helper function to reconnect Redis when tests close the connection
func reconnectRedis() {
	if testClient != nil {
		testClient.Close()
	}
	testClient = redisContainer.NewClient()
	if err := testClient.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("Failed to reconnect to Redis: %v", err))
	}
	// Update repository with new client
	repo = sessionRepository.New(testClient)
}
