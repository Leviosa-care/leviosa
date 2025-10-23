package catalog

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
)

func TestCatalogCache_UpsertCategory(t *testing.T) {
	cache := NewCatalogCache()
	ctx := context.Background()

	t.Run("should insert published category", func(t *testing.T) {
		category := &domain.CachedCategory{
			ID:          uuid.New(),
			Name:        "Test Category",
			Description: "Test Description",
			Status:      "published",
			Metadata:    map[string]any{"key": "value"},
		}

		err := cache.UpsertCategory(ctx, category)
		require.NoError(t, err)

		// Verify category was stored
		stored, exists := cache.GetCategoryByID(category.ID)
		require.True(t, exists)
		assert.Equal(t, category.Name, stored.Name)
		assert.Equal(t, category.Description, stored.Description)
		assert.Equal(t, category.Status, stored.Status)
	})

	t.Run("should not store unpublished category", func(t *testing.T) {
		category := &domain.CachedCategory{
			ID:          uuid.New(),
			Name:        "Unpublished Category",
			Description: "Test Description",
			Status:      "draft",
			Metadata:    map[string]any{"key": "value"},
		}

		err := cache.UpsertCategory(ctx, category)
		require.NoError(t, err)

		// Verify category was not stored
		_, exists := cache.GetCategoryByID(category.ID)
		assert.False(t, exists)
	})

	t.Run("should update existing category", func(t *testing.T) {
		categoryID := uuid.New()
		original := &domain.CachedCategory{
			ID:          categoryID,
			Name:        "Original Name",
			Description: "Original Description",
			Status:      "published",
			Metadata:    map[string]any{"key": "original"},
		}

		// Insert original
		err := cache.UpsertCategory(ctx, original)
		require.NoError(t, err)

		// Update with new data
		updated := &domain.CachedCategory{
			ID:          categoryID,
			Name:        "Updated Name",
			Description: "Updated Description",
			Status:      "published",
			Metadata:    map[string]any{"key": "updated"},
		}

		err = cache.UpsertCategory(ctx, updated)
		require.NoError(t, err)

		// Verify update
		stored, exists := cache.GetCategoryByID(categoryID)
		require.True(t, exists)
		assert.Equal(t, "Updated Name", stored.Name)
		assert.Equal(t, "Updated Description", stored.Description)
	})

	t.Run("should handle nil category gracefully", func(t *testing.T) {
		err := cache.UpsertCategory(ctx, nil)
		require.NoError(t, err)
	})
}

func TestCatalogCache_UpsertProduct(t *testing.T) {
	cache := NewCatalogCache()
	ctx := context.Background()

	t.Run("should insert published product", func(t *testing.T) {
		product := &domain.CachedProduct{
			ID:                 uuid.New(),
			Name:               "Test Product",
			Description:        "Test Description",
			CategoryID:         uuid.New(),
			Duration:           60,
			Status:             "published",
			Availability:       "available",
			BufferTime:         15,
			CancellationHours:  24,
			StripeProductID:    "prod_test123",
			Metadata:           map[string]any{"key": "value"},
		}

		err := cache.UpsertProduct(ctx, product)
		require.NoError(t, err)

		// Verify product was stored
		stored, exists := cache.GetProductByID(product.ID)
		require.True(t, exists)
		assert.Equal(t, product.Name, stored.Name)
		assert.Equal(t, product.Description, stored.Description)
		assert.Equal(t, product.Status, stored.Status)
	})

	t.Run("should not store unpublished product", func(t *testing.T) {
		product := &domain.CachedProduct{
			ID:                 uuid.New(),
			Name:               "Unpublished Product",
			Description:        "Test Description",
			CategoryID:         uuid.New(),
			Duration:           60,
			Status:             "draft",
			Availability:       "available",
			BufferTime:         15,
			CancellationHours:  24,
			StripeProductID:    "prod_test123",
			Metadata:           map[string]any{"key": "value"},
		}

		err := cache.UpsertProduct(ctx, product)
		require.NoError(t, err)

		// Verify product was not stored
		_, exists := cache.GetProductByID(product.ID)
		assert.False(t, exists)
	})

	t.Run("should update existing product", func(t *testing.T) {
		productID := uuid.New()
		original := &domain.CachedProduct{
			ID:                 productID,
			Name:               "Original Product",
			Description:        "Original Description",
			CategoryID:         uuid.New(),
			Duration:           60,
			Status:             "published",
			Availability:       "available",
			BufferTime:         15,
			CancellationHours:  24,
			StripeProductID:    "prod_original",
			Metadata:           map[string]any{"key": "original"},
		}

		// Insert original
		err := cache.UpsertProduct(ctx, original)
		require.NoError(t, err)

		// Update with new data
		updated := &domain.CachedProduct{
			ID:                 productID,
			Name:               "Updated Product",
			Description:        "Updated Description",
			CategoryID:         uuid.New(),
			Duration:           90,
			Status:             "published",
			Availability:       "limited",
			BufferTime:         30,
			CancellationHours:  48,
			StripeProductID:    "prod_updated",
			Metadata:           map[string]any{"key": "updated"},
		}

		err = cache.UpsertProduct(ctx, updated)
		require.NoError(t, err)

		// Verify update
		stored, exists := cache.GetProductByID(productID)
		require.True(t, exists)
		assert.Equal(t, "Updated Product", stored.Name)
		assert.Equal(t, 90, stored.Duration)
		assert.Equal(t, "limited", stored.Availability)
	})

	t.Run("should handle nil product gracefully", func(t *testing.T) {
		err := cache.UpsertProduct(ctx, nil)
		require.NoError(t, err)
	})
}

func TestCatalogCache_DeleteCategory(t *testing.T) {
	cache := NewCatalogCache()
	ctx := context.Background()

	t.Run("should delete category and its products", func(t *testing.T) {
		categoryID := uuid.New()

		// Insert category
		category := &domain.CachedCategory{
			ID:          categoryID,
			Name:        "Test Category",
			Description: "Test Description",
			Status:      "published",
			Metadata:    map[string]any{"key": "value"},
		}
		err := cache.UpsertCategory(ctx, category)
		require.NoError(t, err)

		// Insert products in this category
		product1 := &domain.CachedProduct{
			ID:         uuid.New(),
			Name:       "Product 1",
			CategoryID: categoryID,
			Status:     "published",
		}
		product2 := &domain.CachedProduct{
			ID:         uuid.New(),
			Name:       "Product 2",
			CategoryID: categoryID,
			Status:     "published",
		}
		err = cache.UpsertProduct(ctx, product1)
		require.NoError(t, err)
		err = cache.UpsertProduct(ctx, product2)
		require.NoError(t, err)

		// Verify everything exists
		_, catExists := cache.GetCategoryByID(categoryID)
		assert.True(t, catExists)
		_, prod1Exists := cache.GetProductByID(product1.ID)
		assert.True(t, prod1Exists)
		_, prod2Exists := cache.GetProductByID(product2.ID)
		assert.True(t, prod2Exists)

		// Delete category
		err = cache.DeleteCategory(ctx, categoryID)
		require.NoError(t, err)

		// Verify category and its products are deleted
		_, catExists = cache.GetCategoryByID(categoryID)
		assert.False(t, catExists)
		_, prod1Exists = cache.GetProductByID(product1.ID)
		assert.False(t, prod1Exists)
		_, prod2Exists = cache.GetProductByID(product2.ID)
		assert.False(t, prod2Exists)
	})

	t.Run("should handle deleting non-existent category", func(t *testing.T) {
		err := cache.DeleteCategory(ctx, uuid.New())
		require.NoError(t, err)
	})
}

func TestCatalogCache_DeleteProduct(t *testing.T) {
	cache := NewCatalogCache()
	ctx := context.Background()

	t.Run("should delete product", func(t *testing.T) {
		product := &domain.CachedProduct{
			ID:         uuid.New(),
			Name:       "Test Product",
			CategoryID: uuid.New(),
			Status:     "published",
		}
		err := cache.UpsertProduct(ctx, product)
		require.NoError(t, err)

		// Verify product exists
		_, exists := cache.GetProductByID(product.ID)
		assert.True(t, exists)

		// Delete product
		err = cache.DeleteProduct(ctx, product.ID)
		require.NoError(t, err)

		// Verify product is deleted
		_, exists = cache.GetProductByID(product.ID)
		assert.False(t, exists)
	})

	t.Run("should handle deleting non-existent product", func(t *testing.T) {
		err := cache.DeleteProduct(ctx, uuid.New())
		require.NoError(t, err)
	})
}

func TestCatalogCache_ListOperations(t *testing.T) {
	cache := NewCatalogCache()
	ctx := context.Background()

	// Setup test data
	category1ID := uuid.New()
	category2ID := uuid.New()
	product1ID := uuid.New()
	product2ID := uuid.New()
	product3ID := uuid.New()

	category1 := &domain.CachedCategory{
		ID:          category1ID,
		Name:        "Category 1",
		Description: "Description 1",
		Status:      "published",
	}
	category2 := &domain.CachedCategory{
		ID:          category2ID,
		Name:        "Category 2",
		Description: "Description 2",
		Status:      "published",
	}

	product1 := &domain.CachedProduct{
		ID:         product1ID,
		Name:       "Product 1",
		CategoryID: category1ID,
		Status:     "published",
	}
	product2 := &domain.CachedProduct{
		ID:         product2ID,
		Name:       "Product 2",
		CategoryID: category1ID,
		Status:     "published",
	}
	product3 := &domain.CachedProduct{
		ID:         product3ID,
		Name:       "Product 3",
		CategoryID: category2ID,
		Status:     "published",
	}

	// Insert test data
	err := cache.UpsertCategory(ctx, category1)
	require.NoError(t, err)
	err = cache.UpsertCategory(ctx, category2)
	require.NoError(t, err)
	err = cache.UpsertProduct(ctx, product1)
	require.NoError(t, err)
	err = cache.UpsertProduct(ctx, product2)
	require.NoError(t, err)
	err = cache.UpsertProduct(ctx, product3)
	require.NoError(t, err)

	t.Run("should list all categories", func(t *testing.T) {
		categories := cache.ListCategories()
		assert.Len(t, categories, 2)

		// Verify both categories are present
		categoryNames := make(map[string]bool)
		for _, cat := range categories {
			categoryNames[cat.Name] = true
		}
		assert.True(t, categoryNames["Category 1"])
		assert.True(t, categoryNames["Category 2"])
	})

	t.Run("should list all products", func(t *testing.T) {
		products := cache.ListProducts()
		assert.Len(t, products, 3)

		// Verify all products are present
		productNames := make(map[string]bool)
		for _, prod := range products {
			productNames[prod.Name] = true
		}
		assert.True(t, productNames["Product 1"])
		assert.True(t, productNames["Product 2"])
		assert.True(t, productNames["Product 3"])
	})

	t.Run("should list products by category", func(t *testing.T) {
		category1Products := cache.ListProductsByCategory(category1ID)
		assert.Len(t, category1Products, 2)

		// Verify correct products are returned
		productNames := make(map[string]bool)
		for _, prod := range category1Products {
			productNames[prod.Name] = true
		}
		assert.True(t, productNames["Product 1"])
		assert.True(t, productNames["Product 2"])
		assert.False(t, productNames["Product 3"])

		category2Products := cache.ListProductsByCategory(category2ID)
		assert.Len(t, category2Products, 1)

		productNames = make(map[string]bool)
		for _, prod := range category2Products {
			productNames[prod.Name] = true
		}
		assert.True(t, productNames["Product 3"])
		assert.False(t, productNames["Product 1"])
		assert.False(t, productNames["Product 2"])
	})
}

func TestCatalogCache_ValidationMethods(t *testing.T) {
	cache := NewCatalogCache()
	ctx := context.Background()

	// Setup test data
	categoryID := uuid.New()
	productID := uuid.New()

	category := &domain.CachedCategory{
		ID:     categoryID,
		Name:   "Test Category",
		Status: "published",
	}
	product := &domain.CachedProduct{
		ID:         productID,
		Name:       "Test Product",
		CategoryID: categoryID,
		Status:     "published",
	}

	err := cache.UpsertCategory(ctx, category)
	require.NoError(t, err)
	err = cache.UpsertProduct(ctx, product)
	require.NoError(t, err)

	t.Run("should validate existing category", func(t *testing.T) {
		assert.True(t, cache.IsValidCategory(categoryID))
	})

	t.Run("should not validate non-existent category", func(t *testing.T) {
		assert.False(t, cache.IsValidCategory(uuid.New()))
	})

	t.Run("should validate existing product", func(t *testing.T) {
		assert.True(t, cache.IsValidProduct(productID))
	})

	t.Run("should not validate non-existent product", func(t *testing.T) {
		assert.False(t, cache.IsValidProduct(uuid.New()))
	})
}