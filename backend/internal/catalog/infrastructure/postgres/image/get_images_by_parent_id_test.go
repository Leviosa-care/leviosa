package imageRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetImagesByParentID(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve multiple images for a parent", func(t *testing.T) {
		td.ClearImagesTable(t, ctx, testPool)

		parentID := uuid.New()

		// Setup: Insert multiple images for the same parent
		image1 := &domain.Image{
			ID: uuid.New(), ParentID: parentID, ParentType: domain.ProductType,
			Title: "Image 1", S3Key: "parent/img1.jpg", Size: 100, ContentType: "image/jpeg",
			IsActive: true, CreatedAt: time.Now().Add(-2 * time.Hour), UpdatedAt: time.Now().Add(-2 * time.Hour),
		}
		image2 := &domain.Image{
			ID: uuid.New(), ParentID: parentID, ParentType: domain.ProductType,
			Title: "Image 2", S3Key: "parent/img2.png", Size: 150, ContentType: "image/png",
			IsActive: false, CreatedAt: time.Now().Add(-1 * time.Hour), UpdatedAt: time.Now().Add(-1 * time.Hour),
		}
		image3 := &domain.Image{
			ID: uuid.New(), ParentID: parentID, ParentType: domain.ProductType,
			Title: "Image 3", S3Key: "parent/img3.gif", Size: 200, ContentType: "image/gif",
			IsActive: false, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		td.InsertImage(t, ctx, image1, testPool)
		td.InsertImage(t, ctx, image2, testPool)
		td.InsertImage(t, ctx, image3, testPool)

		// Insert an image for a different parent to ensure filtering works
		otherParentID := uuid.New()
		otherImage := &domain.Image{
			ID: uuid.New(), ParentID: otherParentID, ParentType: domain.CategoryType,
			Title: "Other Image", S3Key: "other/img.jpg", Size: 50, ContentType: "image/jpeg",
			IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		td.InsertImage(t, ctx, otherImage, testPool)

		// Act: Retrieve images for the target parent
		retrievedImages, err := repo.GetImagesByParentID(ctx, parentID, domain.ProductType)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedImages)
		assert.Len(t, retrievedImages, 3)

		// Assert: Verify retrieved images match and are in correct order
		assert.Equal(t, image1.ID, retrievedImages[0].ID)
		assert.Equal(t, image2.ID, retrievedImages[1].ID)
		assert.Equal(t, image3.ID, retrievedImages[2].ID)

		// Verify content of one image in detail
		assert.Equal(t, image1.Title, retrievedImages[0].Title)
		assert.Equal(t, image1.S3Key, retrievedImages[0].S3Key)
		assert.Equal(t, image1.IsActive, retrievedImages[0].IsActive)
	})

	t.Run("should return an empty slice if parent exists but has no images", func(t *testing.T) {
		td.ClearImagesTable(t, ctx, testPool)

		parentID := uuid.New()
		// No images are inserted for this parent

		// Act: Retrieve images for the parent
		retrievedImages, err := repo.GetImagesByParentID(ctx, parentID, domain.ProductType)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedImages) // Should be an empty slice, not nil
		assert.Len(t, retrievedImages, 0)
	})

	t.Run("should return an empty slice if parent does not exist", func(t *testing.T) {
		td.ClearImagesTable(t, ctx, testPool)

		nonExistentParentID := uuid.New()

		// Act: Retrieve images for a non-existent parent
		retrievedImages, err := repo.GetImagesByParentID(ctx, nonExistentParentID, domain.ProductType)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedImages) // Should be an empty slice, not nil
		assert.Len(t, retrievedImages, 0)
	})
}
