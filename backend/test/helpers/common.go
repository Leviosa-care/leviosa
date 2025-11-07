package helpers

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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

// ClearAuthTestData cleans all test data from both PostgreSQL and Redis
func ClearAuthTestData(t *testing.T, ctx context.Context, pool *pgxpool.Pool, redisClient *redis.Client) {
	t.Helper()

	// Clear database tables
	ClearUsersTable(t, ctx, pool)

	// Clear Redis OTP keys
	ClearOTPKeys(t, ctx, redisClient)

	// Clear Redis session keys
	ClearSessionsRedis(t, ctx, redisClient)
}

// ClearSettingTestData cleans all test data from both PostgreSQL and Redis
func ClearSettingTestData(t *testing.T, ctx context.Context, pool *pgxpool.Pool, s3Client *s3.Client) {
	t.Helper()

	ClearSettingsTable(t, ctx, pool)
	ClearS3Bucket(t, ctx, s3Client)
}

// ClearS3Bucket removes all objects from the test S3 bucket
func ClearS3Bucket(t *testing.T, ctx context.Context, s3Client *s3.Client) {
	t.Helper()

	// List all objects in the bucket
	listObjectsInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(BUCKETNAME),
	}
	objects, err := s3Client.ListObjectsV2(ctx, listObjectsInput)
	if err != nil {
		t.Fatalf("Failed to list objects in bucket %s: %v", BUCKETNAME, err)
	}

	if *objects.KeyCount > 0 {
		var objectIds []types.ObjectIdentifier
		for _, object := range objects.Contents {
			objectIds = append(objectIds, types.ObjectIdentifier{
				Key: object.Key,
			})
		}

		// Delete the objects
		deleteInput := &s3.DeleteObjectsInput{
			Bucket: aws.String(BUCKETNAME),
			Delete: &types.Delete{
				Objects: objectIds,
			},
		}
		_, err = s3Client.DeleteObjects(ctx, deleteInput)
		if err != nil {
			t.Fatalf("Failed to delete objects from bucket %s: %v", BUCKETNAME, err)
		}
	}
}

// CreateContextWithTimeout creates a context with a specific timeout duration
// Useful for testing timeout scenarios
func CreateContextWithTimeout(parent context.Context, duration time.Duration) context.Context {
	ctx, _ := context.WithTimeout(parent, duration)
	return ctx
}

// CreateCancelledContext creates a pre-cancelled context
// Useful for testing request cancellation scenarios
func CreateCancelledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Immediately cancel
	return ctx
}
