package settingsMedia

import (
	"bytes"
	"context"
	"io"

	rp "github.com/hengadev/leviosa/internal/repository"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const logoKey = "public/assets/logo.jpg"

func (r *repository) GetLogo(ctx context.Context) ([]byte, error) {
	result, err := r.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(logoKey),
	})
	if err != nil {
		return nil, rp.NewNotFoundErr(err, "logo from S3")
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

func (r *repository) SetLogo(ctx context.Context, logo []byte) error {
	reader := bytes.NewReader(logo)

	_, err := r.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.BucketName),
		Key:         aws.String(logoKey),
		Body:        reader,
		ContentType: aws.String("image/jpeg"),
		ACL:         "public-read",
	})

	if err != nil {
		return rp.NewNotCreatedErr(err, "logo to S3")
	}
	return nil
}
