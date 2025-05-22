package settingsMedia

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type repository struct {
	Client     *s3.Client
	BucketName string
}

func New(ctx context.Context, client *s3.Client, bucketName string) (*repository, error) {
	return &repository{
		Client:     client,
		BucketName: bucketName,
	}, nil
}
