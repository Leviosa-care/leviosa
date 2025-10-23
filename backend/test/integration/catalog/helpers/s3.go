package helpers

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
)

// clearBucket deletes all objects from the test S3 bucket.
func ClearBucket(t *testing.T, ctx context.Context, s3Client *s3.Client) {
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

// UploadFileHelper is a helper to put an object into S3 for test setup.
func UploadFileHelper(t *testing.T, ctx context.Context, s3Client *s3.Client, key, content, contentType string) {
	t.Helper()
	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(BUCKETNAME),
		Key:           aws.String(key),
		Body:          bytes.NewReader([]byte(content)),
		ContentLength: aws.Int64(int64(len(content))),
		ContentType:   aws.String(contentType),
	})
	assert.NoError(t, err, "Failed to upload file for test setup: %s", key)
}

// CheckFileExistsInS3 checks if an object exists in S3.
func CheckFileExistsInS3(t *testing.T, ctx context.Context, s3Client *s3.Client, key string) bool {
	t.Helper()
	_, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(BUCKETNAME),
		Key:    aws.String(key),
	})
	if err != nil {
		var noSuchKey *types.NotFound
		if assert.ErrorAs(t, err, &noSuchKey) {
			return false // Object not found
		}
		assert.Fail(t, fmt.Sprintf("Failed to check object existence for key %s: %v", key, err))
	}
	return true // Object found
}
