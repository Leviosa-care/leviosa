package imageMedia

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// DeleteFile removes a file from S3 using its key.
// It returns errs.NewExternalStorageErr if the deletion fails.
func (r *repository) DeleteFile(ctx context.Context, key string) error {
	_, err := r.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		// Classify S3-specific errors if needed, otherwise use generic external storage error.
		return errs.NewExternalStorageErr(err, fmt.Sprintf("S3 DeleteObject for key %s", key), key)
	}
	// S3 DeleteObject is idempotent: it returns success even if the object doesn't exist.
	// So, we don't need to check RowsAffected equivalent here.
	return nil
}
