package media_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	td "github.com/Leviosa-care/catalog/test/helpers"

	"github.com/Leviosa-care/core/errs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUploadLogo(t *testing.T) {
	ctx := context.Background()
	t.Run("should successfully upload a file to S3", func(t *testing.T) {

		td.ClearBucket(t, ctx, s3Client)

		// Create mock file data
		testData := "this is a test file content"
		fileReader := bytes.NewReader([]byte(testData))
		fileSize := int64(len(testData))
		contentType := "text/plain"
		key := fmt.Sprintf("test-uploads/%s.txt", uuid.New().String())

		// Call the function to test
		returnedKey, err := repo.UploadLogo(ctx, key, fileReader, fileSize, contentType)
		assert.NoError(t, err)
		assert.Equal(t, key, returnedKey)

		// Verify the file exists in the S3 bucket
		getObjectInput := &s3.GetObjectInput{
			Bucket: aws.String(td.BUCKETNAME),
			Key:    aws.String(key),
		}
		result, err := s3Client.GetObject(ctx, getObjectInput)
		assert.NoError(t, err)

		defer result.Body.Close()

		// Read the content and verify it matches the original data
		uploadedData, err := io.ReadAll(result.Body)
		assert.NoError(t, err)
		assert.Equal(t, testData, string(uploadedData))

		// Verify metadata
		assert.Equal(t, aws.String(contentType), result.ContentType)
		assert.Equal(t, &fileSize, result.ContentLength)
	})

	t.Run("should return an error when uploading a file with an empty key", func(t *testing.T) {
		td.ClearBucket(t, ctx, s3Client)

		// Create mock file data
		testData := "this is a test file content"
		fileReader := bytes.NewReader([]byte(testData))
		fileSize := int64(len(testData))
		contentType := "text/plain"
		key := "" // Empty key, which is an invalid S3 key

		// Call the function to test
		returnedKey, err := repo.UploadLogo(ctx, key, fileReader, fileSize, contentType)
		assert.Error(t, err)
		assert.Equal(t, "", returnedKey)

		// Verify the returned error is the custom external storage error
		assert.ErrorIs(t, err, errs.ErrExternalStorage)
	})
}
