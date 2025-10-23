package imageMedia

import (
	"context"
	"io"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// UploadFile uploads an image file to S3 and returns the generated S3 key.
// Parameters:
//
//	ctx: The context for the operation.
//	file: The content of the file (an io.Reader).
//	key: The original name of the file (e.g., "myimage.jpg").
//	size: The size of the file in bytes.
//	contentType: The MIME type of the file (e.g., "image/jpeg"). This needs to be passed down.
func (r *repository) UploadFile(ctx context.Context, key string, file io.Reader, size int64, contentType string) (string, error) {
	_, err := r.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(r.BucketName),
		Key:           aws.String(key),
		Body:          file, // multipart.File is an io.Reader
		ContentLength: &size,
		ContentType:   aws.String(contentType),
	})
	if err != nil {
		return "", errs.NewExternalStorageErr(err, "S3 PutObject for parent image", key)
	}
	return key, nil
}
