package settingsMedia

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type repository struct {
	Client     *s3.Client
	BucketName string
}

func New(ctx context.Context, bucketName string) (*repository, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load default configuration for S3 repository: %w", err)
	}
	client := s3.NewFromConfig(cfg)
	return &repository{
		Client:     client,
		BucketName: bucketName,
	}, nil
}
