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

func TestDeleteImage(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully delete an existing image", func(t *testing.T) {
		td.ClearImagesTable(t, ctx, testPool)

		// Setup: Insert a mock image
		imageID := uuid.New()
		mockImage := &domain.Image{
			ID: imageID, ParentID: uuid.New(), ParentType: domain.ProductType,
			Title: "Test Image", S3Key: "test/key.jpg", Size: 100, ContentType: "image/jpeg",
			IsActive: false, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		td.InsertImage(t, ctx, mockImage, testPool)

		// Pre-condition: Verify image exists
		assert.Equal(t, 1, getCountByID(t, imageID))

		// Act: Delete the image
		err := repo.DeleteImage(ctx, imageID)
		assert.NoError(t, err)

		// Post-condition: Verify image no longer exists
		assert.Equal(t, 0, getCountByID(t, imageID))
	})

	t.Run("should return NotFound error when attempting to delete a non-existent image", func(t *testing.T) {
		td.ClearImagesTable(t, ctx, testPool)

		nonExistentImageID := uuid.New()

		// Act: Attempt to delete a non-existent image
		err := repo.DeleteImage(ctx, nonExistentImageID)

		// Assert: Expect a NotFound error
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)

		// Post-condition: Verify no images were affected
		assert.Equal(t, 0, getCountByID(t, nonExistentImageID))
	})
}

// getCountByID checks if an image with a given ID exists in the database.
func getCountByID(t *testing.T, id uuid.UUID) int {
	t.Helper()
	var count int
	err := testPool.QueryRow(context.Background(), "SELECT COUNT(*) FROM catalog.images WHERE id = $1", id).Scan(&count)
	assert.NoError(t, err, "Failed to query for image count")
	return count
}
