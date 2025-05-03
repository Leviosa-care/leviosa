package mediaService

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// TODO:
// https://dev.to/aws-builders/get-objects-from-aws-s3-bucket-with-golang-2mne

// Function that gets all the photos from a S3 bucket associated with an event of ID eventID.
func (s *service) GetAllObjects(ctx context.Context, eventID string) ([]types.Object, error) {
	objects, err := s.Repo.FindAllObjects(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("get photos: %w", err)
	}
	return objects, nil
}
