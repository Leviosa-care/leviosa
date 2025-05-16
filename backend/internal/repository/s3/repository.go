package mediaRepository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Repository struct {
	Uploader   *manager.Uploader
	Downloader *manager.Downloader
	Client     *s3.Client
	BucketName string
}

func New(ctx context.Context, bucketName string) (*Repository, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load default configuration for S3 repository: %w", err)
	}
	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)
	downloader := manager.NewDownloader(client)
	return &Repository{
		Uploader:   uploader,
		Downloader: downloader,
		Client:     client,
		BucketName: bucketName,
	}, nil
}
