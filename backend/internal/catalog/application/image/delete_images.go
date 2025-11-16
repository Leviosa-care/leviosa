package image

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

// DeleteImages deletes all image records for a given parent and removes their files from S3.
// It handles potential inconsistencies by logging severe errors if S3 deletion succeeds but DB fails.
func (s *ImageService) DeleteImages(ctx context.Context, parentIDStr string, parentTypeStr string) error {
	parentID, err := uuid.Parse(parentIDStr)
	if err != nil {
		return errs.NewInvalidValueErr("parent ID must be a valid UUID")
	}
	parentType := domain.ParentType(parentTypeStr)
	if !parentType.IsValid() {
		return errs.NewInvalidValueErr("invalid parent type")
	}

	var existingErr error
	switch parentType {
	case domain.CategoryType:
		_, existingErr = s.sharedRepo.GetCategoryByID(ctx, parentID)
	case domain.ProductType:
		_, existingErr = s.sharedRepo.GetProductByID(ctx, parentID)
	default:
		// This case should ideally be caught by parentType.IsValid(), but acts as a safeguard.
		return errs.NewInvalidValueErr("unsupported parent type")
	}

	if existingErr != nil {
		if errors.Is(existingErr, errs.ErrRepositoryNotFound) {
			// If the parent itself is not found, there are no images to delete for it.
			// This can be treated as a successful no-op, as the goal is to ensure cleanup.
			log.Printf("Service: Attempted to delete images for non-existent %s with ID %s. No-op.", parentType, parentID)
			return nil
		}
		return errs.NewUnexpectedError(fmt.Errorf("failed to retrieve %s with ID %s: %w", parentType, parentID, existingErr))
	}

	images, err := s.repo.GetImagesByParentID(ctx, parentID, parentType)
	if err != nil {
		return fmt.Errorf("failed to retrieve images for deletion for parent %s %s: %w", parentType, parentID, err)
	}

	if len(images) == 0 {
		log.Printf("Service: No images found for %s with ID %s. No deletion performed.", parentType, parentID)
		return nil
	}

	for _, img := range images {
		if err := s.mediaRepo.DeleteFile(ctx, img.S3Key); err != nil {
			log.Printf("Service: Failed to delete S3 file %s for image %s: %v", img.S3Key, img.ID, err)
			// Return an external service error as S3 is the primary storage.
			return fmt.Errorf("failed to delete image file %s from S3", img.S3Key)
		}
	}

	_, err = s.repo.DeleteImagesByParentID(ctx, parentID, parentType)
	if err != nil {
		// CRITICAL INCONSISTENCY: S3 files are gone, but DB records remain.
		// Re-uploading files is not feasible here. Log a severe error.
		log.Printf("Service: FATAL ERROR - Failed to delete DB image records for %s with ID %s after S3 deletion. Data inconsistency detected! DB error: %v", parentType, parentID, err)
		return fmt.Errorf("failed to delete image records from database for %s %s: %w", parentType, parentID, err)
	}

	return nil
}
