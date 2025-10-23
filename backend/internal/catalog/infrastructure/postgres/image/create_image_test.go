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

func TestCreateImage(t *testing.T) {
	t.Run("should successfully create an image record with IsActive=false", func(t *testing.T) {
		ctx := context.Background()
		td.ClearImagesTable(t, ctx, testPool)

		parentID := uuid.New()
		mockImage := td.NewValidImage(parentID)

		err := repo.CreateImage(ctx, mockImage)
		assert.NoError(t, err)

		// Verify the image was created and has the correct data
		var createdImage domain.Image
		row := testPool.QueryRow(ctx, "SELECT id, parent_id, parent_type, title, s3_key, size, content_type, is_active, created_at FROM catalog.images WHERE id = $1", mockImage.ID)
		err = row.Scan(&createdImage.ID, &createdImage.ParentID, &createdImage.ParentType, &createdImage.Title, &createdImage.S3Key, &createdImage.Size, &createdImage.ContentType, &createdImage.IsActive, &createdImage.CreatedAt)
		assert.NoError(t, err)

		// Assert that the retrieved data matches the mock data
		assert.Equal(t, mockImage.ID, createdImage.ID)
		assert.Equal(t, mockImage.ParentID, createdImage.ParentID)
		assert.Equal(t, mockImage.ParentType, createdImage.ParentType)
		assert.Equal(t, mockImage.Title, createdImage.Title)
		assert.Equal(t, mockImage.S3Key, createdImage.S3Key)
		assert.Equal(t, mockImage.Size, createdImage.Size)
		assert.Equal(t, mockImage.ContentType, createdImage.ContentType)
		assert.Equal(t, false, createdImage.IsActive) // Ensure it was created as inactive
		assert.True(t, time.Since(createdImage.CreatedAt) < time.Minute)
	})

	t.Run("should fail with a unique constraint error if S3Key already exists", func(t *testing.T) {
		ctx := context.Background()
		td.ClearImagesTable(t, ctx, testPool)

		parentID := uuid.New()
		mockImage1 := td.NewValidImage(parentID)

		err := repo.CreateImage(ctx, mockImage1)
		assert.NoError(t, err)

		// Attempt to create a second image with the same S3Key
		mockImage2 := td.NewValidImage(parentID)
		mockImage2.S3Key = mockImage1.S3Key // Duplicate S3Key

		err = repo.CreateImage(ctx, mockImage2)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrUniqueViolation)
	})

	t.Run("should fail with a CHECK constraint error for an invalid ParentType", func(t *testing.T) {
		ctx := context.Background()
		td.ClearImagesTable(t, ctx, testPool)

		mockImage := td.NewValidImage(uuid.New())
		mockImage.ParentType = "invalid_type"

		err := repo.CreateImage(ctx, mockImage)
		assert.Error(t, err)
		// Assuming a helper exists to check for this specific error type
		assert.ErrorIs(t, err, errs.ErrCheckViolation, "expected check constraint error, got: %v", err)
	})

	t.Run("should fail with unique partial index violation if two active images are created", func(t *testing.T) {
		ctx := context.Background()
		td.ClearImagesTable(t, ctx, testPool)

		parentID := uuid.New()

		// Create the first active image
		mockImage1 := td.NewValidImage(parentID)
		mockImage1.IsActive = true
		err := repo.CreateImage(ctx, mockImage1)
		assert.NoError(t, err)

		// Attempt to create a second active image for the same parent
		mockImage2 := td.NewValidImage(parentID)
		mockImage2.IsActive = true
		mockImage2.S3Key = "another/unique/key.jpg" // Ensure S3Key is not the unique key violation cause

		err = repo.CreateImage(ctx, mockImage2)
		assert.Error(t, err)
		// The error should be a unique constraint violation from the partial index
		assert.ErrorIs(t, err, errs.ErrUniqueViolation)
	})
}
