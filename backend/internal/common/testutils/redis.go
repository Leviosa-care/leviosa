package testutils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RedisContainer struct {
	testcontainers.Container
	ConnectionString string
	Host             string
	Port             string
}

func SetupRedis(ctx context.Context, t *testing.T) (*RedisContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("6379/tcp"),
			wait.ForLog("Ready to accept connections"),
		).WithDeadline(30 * time.Second),
		Cmd: []string{"redis-server", "--appendonly", "yes"}, // Enable persistence for testing
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host IP: %w", err)
	}

	port, err := container.MappedPort(ctx, "6379")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	connStr := fmt.Sprintf("redis://%s:%s/0", hostIP, port.Port())

	// Test connection to ensure Redis is ready
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", hostIP, port.Port()),
		DB:   0,
	})

	// Ping Redis to ensure it's ready
	if err := client.Ping(ctx).Err(); err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}
	client.Close() // Close initial test connection

	return &RedisContainer{
		Container:        container,
		ConnectionString: connStr,
		Host:             hostIP,
		Port:             port.Port(),
	}, nil
}

func TeardownRedis(ctx context.Context, t *testing.T, container *RedisContainer) {
	if err := container.Terminate(ctx); err != nil {
		if t != nil {
			t.Fatalf("failed to terminate container: %v", err)
		}
	}
}

// NewClient creates a new Redis client for tests
func (r *RedisContainer) NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", r.Host, r.Port),
		DB:   0, // Use default database
	})
}

// FlushDB clears all data from the Redis database
func (r *RedisContainer) FlushDB(ctx context.Context) error {
	client := r.NewClient()
	defer client.Close()

	return client.FlushDB(ctx).Err()
}

// GetConnectionOptions returns Redis connection options for testing
func (r *RedisContainer) GetConnectionOptions() *redis.Options {
	return &redis.Options{
		Addr: fmt.Sprintf("%s:%s", r.Host, r.Port),
		DB:   0,
	}
}
