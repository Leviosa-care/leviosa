package image

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

func (s *ImageService) SetActiveImage(ctx context.Context, request *domain.ImageModifierRequest) error {
	if err := request.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}
	imageID, _ := uuid.Parse(request.ImageID)
	parentID, _ := uuid.Parse(request.ParentID)
	parentType := domain.ParentType(request.ParentType)

	var existingErr error
	switch parentType {
	case domain.CategoryType:
		_, existingErr = s.sharedRepo.GetCategoryByID(ctx, parentID)
	case domain.ProductType:
		_, existingErr = s.sharedRepo.GetProductByID(ctx, parentID)
	}
	if existingErr != nil {
		return fmt.Errorf("failed to retrieve %s with ID %s: %w", parentType, parentID, existingErr)
	}

	err := s.repo.SetActiveImage(ctx, imageID, parentID, parentType)
	if err != nil {
		return fmt.Errorf("set image active: %w", err)
	}

	return nil
}
