package media

import (
	"context"

	"github.com/Leviosa-care/settings/internal/ports"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type repository struct {
	Client     *s3.Client
	BucketName string
}

func New(ctx context.Context, client *s3.Client, bucketName string) ports.SettingsMedia {
	return &repository{
		Client:     client,
		BucketName: bucketName,
	}
}
