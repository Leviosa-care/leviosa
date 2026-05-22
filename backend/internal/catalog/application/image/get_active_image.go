package image

import (
	"context"
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
		return nil, fmt.Errorf("failed to retrieve %s with ID %s: %w", parentType, parentID, existingErr)
	}

	image, err := s.repo.GetActiveImage(ctx, parentID, parentType)
	if err != nil {
		return nil, fmt.Errorf("get active image: %w", err)
	}

	return image, nil
}
