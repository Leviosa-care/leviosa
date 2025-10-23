package imageRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/catalog/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSetActiveImage(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully set an inactive image as active when no other is active", func(t *testing.T) {
		td.ClearImagesTable(t, ctx, testPool)

		parentID := uuid.New()

		// Create an inactive image
		mockImage := td.NewValidImage(parentID)
		mockImage.IsActive = false
		imageID := mockImage.ID

		td.InsertImage(t, ctx, mockImage, testPool)
		assert.False(t, td.GetImageStatus(t, ctx, imageID, testPool))

		// Set the image as active
		err := repo.SetActiveImage(ctx, imageID, parentID, mockImage.ParentType)
		assert.NoError(t, err)

		// Verify it is now active
		assert.True(t, td.GetImageStatus(t, ctx, imageID, testPool))
	})

	t.Run("should successfully set a new image as active and deactivate the old one", func(t *testing.T) {
		td.ClearImagesTable(t, ctx, testPool)

		parentID := uuid.New()

		// Create the first active image
		mockImage1 := td.NewValidImage(parentID) // Created as active
		mockImage1.IsActive = true
		imageID1 := mockImage1.ID

		td.InsertImage(t, ctx, mockImage1, testPool)
		assert.True(t, td.GetImageStatus(t, ctx, imageID1, testPool)) // Verify it's active

		// Create a second inactive image
		mockImage2 := td.NewValidImage(parentID)
		mockImage2.IsActive = false
		imageID2 := mockImage2.ID
		td.InsertImage(t, ctx, mockImage2, testPool)
		assert.False(t, td.GetImageStatus(t, ctx, imageID2, testPool)) // Verify it's inactive

		// Set the second image as active
		err := repo.SetActiveImage(ctx, imageID2, parentID, mockImage2.ParentType)
		assert.NoError(t, err)

		// Verify image2 is active and image1 is inactive
		assert.True(t, td.GetImageStatus(t, ctx, imageID2, testPool))
		assert.False(t, td.GetImageStatus(t, ctx, imageID1, testPool))
	})

	t.Run("should return NotFound error if imageID to activate does not exist", func(t *testing.T) {
		td.ClearImagesTable(t, ctx, testPool)

		parentID := uuid.New()
		nonExistentImageID := uuid.New()

		// Attempt to set a non-existent image as active
		err := repo.SetActiveImage(ctx, nonExistentImageID, parentID, domain.ProductType) // ParentType doesn't matter much here
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("should return NotFound error if imageID exists but parentID/parentType mismatch", func(t *testing.T) {
		td.ClearImagesTable(t, ctx, testPool)

		parentID1 := uuid.New()
		parentID2 := uuid.New() // Different parent

		// Create an image for parentID1
		mockImage := td.NewValidImage(parentID1)
		mockImage.IsActive = false
		imageID := mockImage.ID

		td.InsertImage(t, ctx, mockImage, testPool)

		// Attempt to set it active with a mismatched parentID
		err := repo.SetActiveImage(ctx, imageID, parentID2, mockImage.ParentType)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)

		// Verify the image's status hasn't changed
		assert.False(t, td.GetImageStatus(t, ctx, imageID, testPool))
	})

}
