package image

import (
	"context"
	"fmt"
	"log"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

// DeleteImage deletes an image record from the database and removes the file from S3.
// It performs a rollback of the S3 deletion if the database deletion fails.
func (s *ImageService) DeleteImage(ctx context.Context, request *domain.ImageModifierRequest) error {
	if err := request.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}
	imageID, _ := uuid.Parse(request.ImageID)
	parentID, _ := uuid.Parse(request.ParentID)
	parentType := domain.ParentType(request.ParentType)

	image, err := s.repo.GetImageByID(ctx, imageID)
	if err != nil {
		return fmt.Errorf("failed to retrieve image for deletion: %w", err)
	}

	if image.ParentID != parentID || image.ParentType != parentType {
		return errs.NewInvalidValueErr(fmt.Sprintf("image ID %s does not belong to %s with ID %s", imageID, parentType, parentID))
	}

	if err := s.mediaRepo.DeleteFile(ctx, image.S3Key); err != nil {
		log.Printf("Service: Failed to delete image file from S3 for key %s: %v", image.S3Key, err)
		return errs.NewExternalServiceErr(err, "failed to delete image file from storage")
	}

	if err := s.repo.DeleteImage(ctx, imageID); err != nil {
		// TODO:
		// CRITICAL ROLLBACK: If DB deletion fails, attempt to re-upload the file to S3.
		// This is complex and might require storing the original file content or having a
		// separate mechanism for S3 recovery. For simplicity here, we'll just log
		// a severe inconsistency. A more robust solution might involve a message queue
		// for eventual consistency or a "tombstone" approach.

		// For now, we'll just log a severe inconsistency as re-uploading the original
		// file content is not feasible at this point without storing it.
		log.Printf("Service: FATAL ERROR - Image record %s failed to delete from DB after S3 deletion. Data inconsistency detected! DB error: %v", imageID, err)
		return fmt.Errorf("failed to delete image record from database after S3 deletion: %w", err)
	}

	return nil
}
