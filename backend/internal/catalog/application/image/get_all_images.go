package image

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// GetAllImages gets all images for a parent entity, sorted by creation date.
// It returns an empty slice if no images are found for the given parent.
func (s *ImageService) GetAllImages(ctx context.Context, parentIDStr string, parentTypeStr string) ([]*domain.Image, error) {
	// 1. Input Validation
	parentID, err := uuid.Parse(parentIDStr)
	if err != nil {
		return nil, errs.NewInvalidValueErr("parent ID must be a valid UUID")
	}
	parentType := domain.ParentType(parentTypeStr)
	if !parentType.IsValid() {
		return nil, errs.NewInvalidValueErr("invalid parent type")
	}

	// 2. Check if the parent entity exists (business rule)
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
			// If the parent itself is not found, it's a NotFound error for the parent,
			// not just that it has no images.
			return nil, errs.NewNotFoundErr(existingErr, fmt.Sprintf("%s with ID %s", parentType, parentID))
		}
		return nil, errs.NewUnexpectedError(fmt.Errorf("failed to retrieve %s with ID %s: %w", parentType, parentID, existingErr))
	}

	// 3. Call the repository to get all images for the parent
	images, err := s.repo.GetImagesByParentID(ctx, parentID, parentType)
	if err != nil {
		// Classify any database errors from the repository.
		// For a "get all" operation, ErrRepositoryNotFound is not expected from the repo,
		// as it should return an empty slice. So, any error here is a DB query failure.
		return nil, errs.NewQueryFailedErr(fmt.Errorf("repository query failed to get all images for parent: %w", err))
	}

	// The repository should return an empty slice if no images are found, so no need to check for nil.
	return images, nil
}
