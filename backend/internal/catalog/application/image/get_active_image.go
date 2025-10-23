package image

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetActiveImage gets the single active image for a parent entity.
// It returns errs.NewNotFoundErr if no active image is found or the parent does not exist.
func (s *ImageService) GetActiveImage(ctx context.Context, parentIDStr string, parentTypeStr string) (*domain.Image, error) {
	parentID, err := uuid.Parse(parentIDStr)
	if err != nil {
		return nil, errs.NewInvalidValueErr("parent ID must be a valid UUID")
	}
	parentType := domain.ParentType(parentTypeStr)
	if !parentType.IsValid() {
		return nil, errs.NewInvalidValueErr("invalid parent type")
	}

	var existingErr error
	switch parentType {
	case domain.CategoryType:
		// Assuming GetCategoryByID takes uuid.UUID for consistency
		_, existingErr = s.sharedRepo.GetCategoryByID(ctx, parentID)
	case domain.ProductType:
		// Assuming GetProductByID takes uuid.UUID for consistency
		_, existingErr = s.sharedRepo.GetProductByID(ctx, parentID)
	default:
		// This case should ideally be caught by parentType.IsValid(), but acts as a safeguard.
		return nil, errs.NewInvalidValueErr("unsupported parent type")
	}

	if existingErr != nil {
		if errors.Is(existingErr, errs.ErrRepositoryNotFound) {
			return nil, errs.NewNotFoundErr(existingErr, fmt.Sprintf("%s with ID %s", parentType, parentID))
		}
		return nil, errs.NewUnexpectedError(fmt.Errorf("failed to retrieve %s with ID %s: %w", parentType, parentID, existingErr))
	}

	image, err := s.repo.GetActiveImage(ctx, parentID, parentType)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			// This means no active image was found for the existing parent.
			return nil, errs.NewNotFoundErr(err, fmt.Sprintf("active image for %s with ID %s", parentType, parentID))
		case errors.Is(err, errs.ErrDBQuery), errors.Is(err, errs.ErrDatabase):
			// General database query or connection issue.
			return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed to get active image: %w", err))
		default:
			// Catch any unhandled repository errors.
			return nil, errs.NewUnexpectedError(fmt.Errorf("unhandled repository error during active image retrieval: %w", err))
		}
	}

	return image, nil
}
