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

func TestGetAllActiveImages(t *testing.T) {
	ctx := context.Background()
	td.ClearImagesTable(t, ctx, testPool)

	// Create a single active image for a category
	categoryID1 := uuid.New()
	activeCategoryImg := &domain.Image{
		ID: uuid.New(), ParentID: categoryID1, ParentType: domain.CategoryType,
		Title: "Active Category Image", S3Key: "cat/active.jpg", Size: 100, ContentType: "image/jpeg",
		IsActive: true, CreatedAt: time.Now().Add(-time.Hour), UpdatedAt: time.Now().Add(-time.Hour),
	}
	td.InsertImage(t, ctx, activeCategoryImg, testPool)

	// Create an inactive image for the same category
	inactiveCategoryImg := &domain.Image{
		ID: uuid.New(), ParentID: categoryID1, ParentType: domain.CategoryType,
		Title: "Inactive Category Image", S3Key: "cat/inactive.jpg", Size: 100, ContentType: "image/jpeg",
		IsActive: false, CreatedAt: time.Now().Add(-2 * time.Hour), UpdatedAt: time.Now().Add(-2 * time.Hour),
	}
	td.InsertImage(t, ctx, inactiveCategoryImg, testPool)

	// Create a single active image for a product
	productID1 := uuid.New()
	activeProductImg := &domain.Image{
		ID: uuid.New(), ParentID: productID1, ParentType: domain.ProductType,
		Title: "Active Product Image", S3Key: "prod/active.jpg", Size: 200, ContentType: "image/png",
		IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	td.InsertImage(t, ctx, activeProductImg, testPool)

	// Create another inactive image for a different product
	productID2 := uuid.New()
	inactiveProductImg := &domain.Image{
		ID: uuid.New(), ParentID: productID2, ParentType: domain.ProductType,
		Title: "Inactive Product Image", S3Key: "prod/inactive.png", Size: 200, ContentType: "image/png",
		IsActive: false, CreatedAt: time.Now().Add(-3 * time.Hour), UpdatedAt: time.Now().Add(-3 * time.Hour),
	}
	td.InsertImage(t, ctx, inactiveProductImg, testPool)

	t.Run("should retrieve all active images for a category", func(t *testing.T) {
		images, err := repo.GetAllActiveImages(ctx, domain.CategoryType)
		assert.NoError(t, err)
		assert.Len(t, images, 1)
		assert.Equal(t, activeCategoryImg.ID, images[0].ID)
		assert.True(t, images[0].IsActive)
	})

	t.Run("should retrieve all active images for a product", func(t *testing.T) {
		images, err := repo.GetAllActiveImages(ctx, domain.ProductType)
		assert.NoError(t, err)
		assert.Len(t, images, 1)
		assert.Equal(t, activeProductImg.ID, images[0].ID)
		assert.True(t, images[0].IsActive)
	})

	t.Run("should return an empty slice if no active images exist for a type", func(t *testing.T) {
		// Use a parent type that has no active images (or no images at all)
		images, err := repo.GetAllActiveImages(ctx, "inventory")
		assert.NoError(t, err)
		assert.NotNil(t, images)
		assert.Len(t, images, 0)
	})

	t.Run("should return an empty slice if no images exist in the table", func(t *testing.T) {
		td.ClearImagesTable(t, ctx, testPool)
		images, err := repo.GetAllActiveImages(ctx, domain.CategoryType)
		assert.NoError(t, err)
		assert.NotNil(t, images)
		assert.Len(t, images, 0)
	})
}
