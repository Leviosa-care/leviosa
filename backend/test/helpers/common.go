package helpers

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Helper to create a pointer to a string
func StrPtr(s string) *string { return &s }

// Helper to create a pointer to a string
func IntPtr(i int) *int { return &i }

// Helper to create a pointer to a string
func BoolPtr(b bool) *bool { return &b }

// Helper to create a pointer to a PublishedStatus
func StatusStrPtr(s string) *string { return &s }

// ClearAllTestData cleans all test data from both PostgreSQL and Redis
func ClearAllTestData(t *testing.T, ctx context.Context, pool *pgxpool.Pool, redisClient *redis.Client) {
	t.Helper()

	// Clear database tables
	ClearUsersTable(t, ctx, pool)

	// Clear Redis OTP keys
	ClearOTPKeys(t, ctx, redisClient)

	// Clear Redis session keys
	ClearSessionsRedis(t, ctx, redisClient)
}

// ClearAuthTestData clears only auth-related test data
func ClearAuthTestData(t *testing.T, ctx context.Context, pool *pgxpool.Pool, redisClient *redis.Client) {
	t.Helper()

	// Clear users table
	ClearUsersTable(t, ctx, pool)

	// Clear OTP keys
	ClearOTPKeys(t, ctx, redisClient)
}
