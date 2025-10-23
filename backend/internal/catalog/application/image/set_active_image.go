package image

import (
	"context"
	"errors"
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
		if errors.Is(existingErr, errs.ErrRepositoryNotFound) {
			return errs.NewNotFoundErr(existingErr, fmt.Sprintf("%s with ID %s", parentType, parentID))
		}
		return errs.NewUnexpectedError(fmt.Errorf("failed to retrieve %s with ID %s: %w", parentType, parentID, existingErr))
	}

	err := s.repo.SetActiveImage(ctx, imageID, parentID, parentType)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			// This means the imageID provided was not found for the given parent.
			return errs.NewNotFoundErr(err, fmt.Sprintf("image with ID %s for %s %s", imageID, parentType, parentID))
		case errors.Is(err, errs.ErrDBQuery), errors.Is(err, errs.ErrDatabase):
			// General database query or connection issue.
			return errs.NewQueryFailedErr(fmt.Errorf("repository query failed to set image active: %w", err))
		default:
			// Catch any unhandled repository errors.
			return errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during image activation: %w", err))
		}
	}

	return nil
}
