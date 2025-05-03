package mediaService

import (
	"context"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// put all service function in here to make it better brother
type Service interface {
	// FindAllObjects(ctx context.Context, eventID string) ([]types.Object, error)
	GetAllObjects(ctx context.Context, eventID string) ([]types.Object, error)
	PostFile(ctx context.Context, file multipart.File, filename, folder string) (string, error)
}

type service struct {
	Repo ReadWriter
}

func New(repo ReadWriter) Service {
	return &service{repo}
}

// TODO: the thing that I am going to handle with S3
// - crud user picture (for the user that do not use oauth)
// - crud event banner
// - crud events photos
// - crud videos for the exercices that a client can do
// - crud offers photos - crud offers videos ?
