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

// GetActiveImage retrieves the single active image record for a parent entity.
// It returns errs.NewRepositoryNotFoundErr if no active image is found for the given parent.
func TestGetActiveImage(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve the active image for a parent", func(t *testing.T) {
		td.ClearImagesTable(t, ctx, testPool)

		parentID := uuid.New()

		// Setup: Insert one active image and some inactive ones for the same parent
		activeImage := &domain.Image{
			ID: uuid.New(), ParentID: parentID, ParentType: domain.ProductType,
			Title: "Active Image", S3Key: "active/img.jpg", Size: 100, ContentType: "image/jpeg",
			IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		inactiveImage1 := &domain.Image{
			ID: uuid.New(), ParentID: parentID, ParentType: domain.ProductType,
			Title: "Inactive Image 1", S3Key: "inactive/img1.png", Size: 150, ContentType: "image/png",
			IsActive: false, CreatedAt: time.Now().Add(-1 * time.Hour), UpdatedAt: time.Now().Add(-1 * time.Hour),
		}
		inactiveImage2 := &domain.Image{
			ID: uuid.New(), ParentID: parentID, ParentType: domain.ProductType,
			Title: "Inactive Image 2", S3Key: "inactive/img2.gif", Size: 200, ContentType: "image/gif",
			IsActive: false, CreatedAt: time.Now().Add(-2 * time.Hour), UpdatedAt: time.Now().Add(-2 * time.Hour),
		}
		td.InsertImage(t, ctx, activeImage, testPool)
		td.InsertImage(t, ctx, inactiveImage1, testPool)
		td.InsertImage(t, ctx, inactiveImage2, testPool)

		// Act: Retrieve the active image
		retrievedImage, err := repo.GetActiveImage(ctx, parentID, domain.ProductType)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedImage)

		// Assert: Verify it's the correct active image
		assert.Equal(t, activeImage.ID, retrievedImage.ID)
		assert.True(t, retrievedImage.IsActive)
		assert.Equal(t, activeImage.Title, retrievedImage.Title)
	})

	t.Run("should return NotFound error if no active image exists for the parent", func(t *testing.T) {
		td.ClearImagesTable(t, ctx, testPool)

		parentID := uuid.New()

		// Setup: Insert only inactive images for the parent
		inactiveImage1 := &domain.Image{
			ID: uuid.New(), ParentID: parentID, ParentType: domain.ProductType,
			Title: "Inactive Image 1", S3Key: "inactive/img1.png", Size: 150, ContentType: "image/png",
			IsActive: false, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		inactiveImage2 := &domain.Image{
			ID: uuid.New(), ParentID: parentID, ParentType: domain.ProductType,
			Title: "Inactive Image 2", S3Key: "inactive/img2.gif", Size: 200, ContentType: "image/gif",
			IsActive: false, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		td.InsertImage(t, ctx, inactiveImage1, testPool)
		td.InsertImage(t, ctx, inactiveImage2, testPool)

		// Act: Attempt to retrieve the active image
		retrievedImage, err := repo.GetActiveImage(ctx, parentID, domain.ProductType)

		// Assert: Expect a NotFound error and nil image
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
		assert.Nil(t, retrievedImage)
	})

	t.Run("should return NotFound error if parent does not exist", func(t *testing.T) {
		td.ClearImagesTable(t, ctx, testPool)

		nonExistentParentID := uuid.New()

		// Act: Attempt to retrieve an active image for a non-existent parent
		retrievedImage, err := repo.GetActiveImage(ctx, nonExistentParentID, domain.ProductType)

		// Assert: Expect a NotFound error and nil image
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
		assert.Nil(t, retrievedImage)
	})
}
