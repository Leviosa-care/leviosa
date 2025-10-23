package imageRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetImageByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve an image by ID", func(t *testing.T) {
		td.ClearImagesTable(t, ctx, testPool)

		// Setup: Insert a mock image
		parentID := uuid.New()
		imageID := uuid.New()
		mockImage := &domain.Image{
			ID: imageID, ParentID: parentID, ParentType: domain.ProductType,
			Title: "Test Image for Get", S3Key: "get/test/key.jpg", Size: 200, ContentType: "image/png",
			IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		td.InsertImage(t, ctx, mockImage, testPool)

		// Act: Retrieve the image
		retrievedImage, err := repo.GetImageByID(ctx, imageID)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedImage)

		// Assert: Verify retrieved data matches mock data
		assert.Equal(t, mockImage.ID, retrievedImage.ID)
		assert.Equal(t, mockImage.ParentID, retrievedImage.ParentID)
		assert.Equal(t, mockImage.ParentType, retrievedImage.ParentType)
		assert.Equal(t, mockImage.Title, retrievedImage.Title)
		assert.Equal(t, mockImage.S3Key, retrievedImage.S3Key)
		assert.Equal(t, mockImage.Size, retrievedImage.Size)
		assert.Equal(t, mockImage.ContentType, retrievedImage.ContentType)
		assert.Equal(t, mockImage.IsActive, retrievedImage.IsActive)
		// For timestamps, check within a reasonable delta
		assert.WithinDuration(t, mockImage.CreatedAt, retrievedImage.CreatedAt, time.Second)
		assert.WithinDuration(t, mockImage.UpdatedAt, retrievedImage.UpdatedAt, time.Second)
	})

	t.Run("should return NotFound error if image does not exist", func(t *testing.T) {
		td.ClearImagesTable(t, ctx, testPool)

		nonExistentImageID := uuid.New()

		// Act: Attempt to retrieve a non-existent image
		retrievedImage, err := repo.GetImageByID(ctx, nonExistentImageID)

		// Assert: Expect a NotFound error and nil image
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
		assert.Nil(t, retrievedImage)
	})
}
