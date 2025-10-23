package imageMedia_test

import (
	"context"
	"fmt"
	"testing"

	imageMedia "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/s3/image"
	td "github.com/Leviosa-care/catalog/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteFile(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully delete an existing file from S3", func(t *testing.T) {
		td.ClearBucket(t, ctx, s3Client)

		testKey := fmt.Sprintf("delete-test/%s.txt", uuid.New().String())
		testContent := "content to delete"

		// Setup: Upload a file to be deleted
		td.UploadFileHelper(t, ctx, s3Client, testKey, testContent, "text/plain")
		assert.True(t, td.CheckFileExistsInS3(t, ctx, s3Client, testKey), "Pre-condition: file should exist before deletion")

		// Act: Delete the file
		err := repo.DeleteFile(ctx, testKey)
		assert.NoError(t, err)

		// Assert: Verify the file no longer exists in S3
		assert.False(t, td.CheckFileExistsInS3(t, ctx, s3Client, testKey), "Post-condition: file should not exist after deletion")
	})

	t.Run("should succeed when attempting to delete a non-existent file from S3 (idempotent)", func(t *testing.T) {
		td.ClearBucket(t, ctx, s3Client)

		nonExistentKey := fmt.Sprintf("non-existent/%s.txt", uuid.New().String())

		// Pre-condition: Ensure the file does not exist
		assert.False(t, td.CheckFileExistsInS3(t, ctx, s3Client, nonExistentKey), "Pre-condition: file should not exist")

		// Act: Attempt to delete a non-existent file
		err := repo.DeleteFile(ctx, nonExistentKey)

		// Assert: No error should be returned, as S3 DeleteObject is idempotent
		assert.NoError(t, err)
		assert.False(t, td.CheckFileExistsInS3(t, ctx, s3Client, nonExistentKey), "Post-condition: file should still not exist")
	})

	t.Run("should return an error if S3 client operation fails (e.g., invalid bucket)", func(t *testing.T) {
		td.ClearBucket(t, ctx, s3Client)

		testKey := fmt.Sprintf("error-test/%s.txt", uuid.New().String())
		testContent := "content"
		td.UploadFileHelper(t, ctx, s3Client, testKey, testContent, "text/plain")

		// Create a mock repo with a client configured for a non-existent bucket
		// (This is a simplified way to induce an error for testing purposes.
		// In a real scenario, you might mock the S3 client interface directly.)
		mockS3Client := s3.NewFromConfig(aws.Config{}, func(o *s3.Options) {
			o.BaseEndpoint = aws.String("http://localhost:4566") // Localstack endpoint
			o.UsePathStyle = true
			o.Region = "us-east-1"
		})
		badBucketRepo := imageMedia.New(ctx, mockS3Client, "non-existent-bucket")

		// Act: Attempt to delete with the misconfigured repo
		err := badBucketRepo.DeleteFile(ctx, testKey)

		// Assert: Expect an ExternalStorageErr
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrExternalStorage)
		assert.Contains(t, err.Error(), "non-existent-bucket") // Specific error message check
	})
}
