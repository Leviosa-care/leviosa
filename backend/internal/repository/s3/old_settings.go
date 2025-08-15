package mediaRepository

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const logoKey = "public/assets/logo.jpg"

func (r *Repository) GetLogo(ctx context.Context) ([]byte, error) {
	// Create a buffer to write the logo data to
	buf := manager.NewWriteAtBuffer([]byte{})

	// Create the download request
	_, err := r.Downloader.Download(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(logoKey),
	})

	if err != nil {
		return nil, fmt.Errorf("download logo from S3: %w", err)
	}

	return buf.Bytes(), nil
}

func (r *Repository) SetLogo(ctx context.Context, logo []byte) error {
	// Create a reader from the logo bytes
	reader := bytes.NewReader(logo)

	// Upload the logo to S3
	_, err := r.Uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.BucketName),
		Key:         aws.String(logoKey),
		Body:        reader,
		ContentType: aws.String("image/jpeg"), // Assuming the logo is a PNG file
		ACL:         "public-read",
	})

	if err != nil {
		return fmt.Errorf("upload logo to S3: %w", err)
	}
	return nil
}
