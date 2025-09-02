package otpRepository_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Leviosa-care/authuser/internal/adapters/redis/otp"
	"github.com/Leviosa-care/authuser/internal/ports"

	tu "github.com/Leviosa-care/core/testutils"
	"github.com/redis/go-redis/v9"
)

var (
	redisContainer *tu.RedisContainer
	testClient     *redis.Client
	repo           ports.OTPRepository
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
	repo = otpRepository.New(testClient)
	log.Println("OTP Repository created.")

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
