package image

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// GetAllImages gets all images for a parent entity, sorted by creation date.
// It returns an empty slice if no images are found for the given parent.
func (s *ImageService) GetAllActiveImages(ctx context.Context, parentTypeStr string) ([]*domain.Image, error) {
	parentType := domain.ParentType(parentTypeStr)
	if !parentType.IsValid() {
		return nil, errs.NewInvalidValueErr("")
	}
	images, err := s.repo.GetAllActiveImages(ctx, parentType)
	if err != nil {
		return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed to get all images for category '%s'", parentType))

	}
	return images, nil
}
