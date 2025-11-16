package image

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

func (s *ImageService) AddImage(ctx context.Context, req *domain.CreateImageRequest, file io.Reader, fileSize int64, contentType string) (string, error) {
	if err := req.Valid(ctx); err != nil {
		return "", errs.NewInvalidValueErr(err.Error())
	}

	// file validation
	var fileErrs errsx.Map
	if file == nil {
		fileErrs.Set("file", "image file is required")
	}
	if fileSize <= 0 {
		fileErrs.Set("file size", "image file cannot be empty")
	}
	if !isValidImageContentType(contentType) { // Basic content type validation (you might want a more exhaustive list)
		fileErrs.Set("file content type", fmt.Sprintf("unsupported image content type: %s", contentType))
	}
	if fileErrs != nil {
		return "", errs.NewInvalidValueErr(fileErrs.Error())
	}

	parentID, _ := uuid.Parse(req.ParentID)

	// check for existing parent with parentID
	var existingErr error
	switch req.ParentType {
	case domain.CategoryType:
		_, existingErr = s.sharedRepo.GetCategoryByID(ctx, parentID)
	case domain.ProductType:
		_, existingErr = s.sharedRepo.GetProductByID(ctx, parentID)
	}
	if existingErr != nil {
		return "", fmt.Errorf("failed to retrieve %s with ID %s: %w", req.ParentType, parentID, existingErr)
	}

	imageID := uuid.New()

	// 2. Generate a Unique Key for the Image (S3 Object Key)
	// It's good practice to generate a unique key, often incoerrsorating the product ID
	// and a unique suffix, along with the file extension.
	// imageKey, err := createProductImagePrefix(productIDStr, contentType)
	imageKey, err := CreateParentImagePrefix(parentID, imageID, req.ParentType, contentType)
	if err != nil {
		return "", errs.NewInvalidValueErr(err.Error())
	}

	// 4. Upload Image to S3 (via Repository)
	// The repository function returns the key if successful, or an error.
	// _, err = s.mediaRepo.UploadFile(ctx, imageKey, file, size, contentType)
	_, err = s.mediaRepo.UploadFile(ctx, imageKey, file, fileSize, contentType)
	if err != nil {
		log.Printf("Service: Failed to create image to S3 bucket: %v", err)
		return "", errs.NewExternalServiceErr(err, "upload image file failed to S3 bucket")
	}

	isActiveRequested := false
	if req.IsActive != nil && *req.IsActive {
		isActiveRequested = *req.IsActive
	}

	image := &domain.Image{
		ID:          imageID,
		ParentID:    parentID,
		ParentType:  req.ParentType,
		Title:       req.Title,
		S3Key:       imageKey,
		Size:        fileSize,
		ContentType: contentType,
		IsActive:    false,
		CreatedAt:   time.Now(),
	}

	err = s.repo.CreateImage(ctx, image)
	if err != nil {
		if rollbackErr := s.mediaRepo.DeleteFile(ctx, imageKey); rollbackErr != nil {
			log.Printf("Service: Failed to rollback S3 image creation for %s with ID %s. Data inconsistency detected! Rollback error: %v", req.ParentType, parentID, rollbackErr)
		}
		return "", fmt.Errorf("create image in database: %w", err)
	}

	if isActiveRequested {
		if err := s.SetActiveImage(ctx, &domain.ImageModifierRequest{
			ImageID:    imageID.String(),
			ParentID:   req.ParentID,
			ParentType: string(req.ParentType),
		}); err != nil {
			log.Printf("Service: Failed to set new image %s as active: %v", image.ID, err)

			// Attempt to delete the new image from the database.
			if delErr := s.repo.DeleteImage(ctx, image.ID); delErr != nil {
				log.Printf("Service: FATAL ERROR - Failed to rollback DB record for image %s: %v", image.ID, delErr)
			}

			// Attempt to delete the S3 file.
			if delErr := s.mediaRepo.DeleteFile(ctx, imageKey); delErr != nil {
				log.Printf("Service: FATAL ERROR - Failed to rollback S3 file for image %s: %v", image.ID, delErr)
			}

			// Return a descriptive error to the user.
			return "", errs.NewUnexpectedError(fmt.Errorf("image uploaded, but failed to be set as active: %w", err))
		}
	}

	return image.ID.String(), nil
}
