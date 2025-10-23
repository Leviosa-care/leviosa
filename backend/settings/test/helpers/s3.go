package helpers

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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
