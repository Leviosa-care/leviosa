package imageRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/catalog/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// getImageCountByParent retrieves the number of images for a given parent.
func getImageCountByParent(t *testing.T, parentID uuid.UUID, parentType domain.ParentType) int {
	t.Helper()
	var count int
	err := testPool.QueryRow(context.Background(), "SELECT COUNT(*) FROM catalog.images WHERE parent_id = $1 AND parent_type = $2", parentID, parentType).Scan(&count)
	assert.NoError(t, err, "Failed to get image count by parent")
	return count
}

func TestDeleteImagesByParentID(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully delete all images for a given parent", func(t *testing.T) {
		td.ClearImagesTable(t, ctx, testPool)

		parentID := uuid.New()

		// Setup: Insert multiple images for the same parent
		td.InsertImage(t, ctx, &domain.Image{
			ID: uuid.New(), ParentID: parentID, ParentType: domain.ProductType,
			Title: "Image 1", S3Key: "p1/img1.jpg", Size: 100, ContentType: "image/jpeg", IsActive: false, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}, testPool)
		td.InsertImage(t, ctx, &domain.Image{
			ID: uuid.New(), ParentID: parentID, ParentType: domain.ProductType,
			Title: "Image 2", S3Key: "p1/img2.jpg", Size: 100, ContentType: "image/jpeg", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}, testPool)
		td.InsertImage(t, ctx, &domain.Image{
			ID: uuid.New(), ParentID: parentID, ParentType: domain.ProductType,
			Title: "Image 3", S3Key: "p1/img3.jpg", Size: 100, ContentType: "image/jpeg", IsActive: false, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}, testPool)

		// Insert images for a different parent to ensure filtering works
		otherParentID := uuid.New()
		td.InsertImage(t, ctx, &domain.Image{
			ID: uuid.New(), ParentID: otherParentID, ParentType: domain.CategoryType,
			Title: "Other Image", S3Key: "c1/img1.jpg", Size: 100, ContentType: "image/jpeg", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}, testPool)

		// Pre-condition: Verify initial counts
		assert.Equal(t, 3, getImageCountByParent(t, parentID, domain.ProductType))
		assert.Equal(t, 1, getImageCountByParent(t, otherParentID, domain.CategoryType))

		// Act: Delete images for the target parent
		rowsAffected, err := repo.DeleteImagesByParentID(ctx, parentID, domain.ProductType)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), rowsAffected) // Ensure 3 rows were affected

		// Post-condition: Verify images for target parent are gone, others remain
		assert.Equal(t, 0, getImageCountByParent(t, parentID, domain.ProductType))
		assert.Equal(t, 1, getImageCountByParent(t, otherParentID, domain.CategoryType))
	})

	t.Run("should return 0 rows affected if no images exist for the parent", func(t *testing.T) {
		td.ClearImagesTable(t, ctx, testPool)

		parentID := uuid.New()

		// Pre-condition: No images for this parent

		// Act: Attempt to delete images for a parent with no images
		rowsAffected, err := repo.DeleteImagesByParentID(ctx, parentID, domain.ProductType)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), rowsAffected) // Expect 0 rows affected
	})

	t.Run("should return 0 rows affected if parent does not exist (no matching images)", func(t *testing.T) {
		td.ClearImagesTable(t, ctx, testPool)

		nonExistentParentID := uuid.New()

		// Pre-condition: No images for this parent

		// Act: Attempt to delete images for a non-existent parent
		rowsAffected, err := repo.DeleteImagesByParentID(ctx, nonExistentParentID, domain.ProductType)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), rowsAffected) // Expect 0 rows affected
	})
}
