package sessionRepository_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/infrastructure/redis/session"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"

	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	"github.com/redis/go-redis/v9"
)

var (
	redisContainer *tu.RedisContainer
	vaultContainer *tu.VaultContainer
	testClient     *redis.Client
	repo           ports.SessionRepository
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error

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
