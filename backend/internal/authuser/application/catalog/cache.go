package catalog

import (
	"context"
	"sync"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/google/uuid"
)

// CatalogCache provides an in-memory cache for catalog data
// This follows the same pattern as the Settings cache in the AuthUser service
// Only items with "published" status are stored in the cache
type CatalogCache struct {
	mu         sync.RWMutex
	categories map[uuid.UUID]domain.CachedCategory
	products   map[uuid.UUID]domain.CachedProduct
}

// NewCatalogCache creates a new in-memory catalog cache
func NewCatalogCache() *CatalogCache {
	return &CatalogCache{
		categories: make(map[uuid.UUID]domain.CachedCategory),
		products:   make(map[uuid.UUID]domain.CachedProduct),
	}
}

// GetCategoryByID returns a category by ID and whether it exists
func (c *CatalogCache) GetCategoryByID(id uuid.UUID) (*domain.CachedCategory, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	category, exists := c.categories[id]
	if !exists {
		return nil, false
	}

	// Return a copy to prevent external mutation
	copy := category
	return &copy, true
}

// GetProductByID returns a product by ID and whether it exists
func (c *CatalogCache) GetProductByID(id uuid.UUID) (*domain.CachedProduct, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	product, exists := c.products[id]
	if !exists {
		return nil, false
	}

	// Return a copy to prevent external mutation
	copy := product
	return &copy, true
}

// ListCategories returns all cached categories
func (c *CatalogCache) ListCategories() []domain.CachedCategory {
	c.mu.RLock()
	defer c.mu.RUnlock()

	categories := make([]domain.CachedCategory, 0, len(c.categories))
	for _, category := range c.categories {
		categories = append(categories, category)
	}

	return categories
}

// ListProducts returns all cached products
func (c *CatalogCache) ListProducts() []domain.CachedProduct {
	c.mu.RLock()
	defer c.mu.RUnlock()

	products := make([]domain.CachedProduct, 0, len(c.products))
	for _, product := range c.products {
		products = append(products, product)
	}

	return products
}

// ListProductsByCategory returns products filtered by category ID
func (c *CatalogCache) ListProductsByCategory(categoryID uuid.UUID) []domain.CachedProduct {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var filtered []domain.CachedProduct
	for _, product := range c.products {
		if product.CategoryID == categoryID {
			filtered = append(filtered, product)
		}
	}

	return filtered
}

// IsValidCategory checks if a category exists in the cache
func (c *CatalogCache) IsValidCategory(categoryID uuid.UUID) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.categories[categoryID]
	return exists
}

// IsValidProduct checks if a product exists in the cache
func (c *CatalogCache) IsValidProduct(productID uuid.UUID) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.products[productID]
	return exists
}

// UpsertCategory adds or updates a category in the cache
// Only categories with "published" status will be stored
func (c *CatalogCache) UpsertCategory(ctx context.Context, category *domain.CachedCategory) error {
	if category == nil {
		return nil // Silently ignore nil categories
	}

	// Only cache published categories
	if !category.IsPublished() {
		// If category exists but is no longer published, remove it
		c.mu.Lock()
		delete(c.categories, category.ID)
		c.mu.Unlock()
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Store a copy to prevent external mutation
	c.categories[category.ID] = *category
	return nil
}

// UpsertProduct adds or updates a product in the cache
// Only products with "published" status will be stored
func (c *CatalogCache) UpsertProduct(ctx context.Context, product *domain.CachedProduct) error {
	if product == nil {
		return nil // Silently ignore nil products
	}

	// Only cache published products
	if !product.IsPublished() {
		// If product exists but is no longer published, remove it
		c.mu.Lock()
		delete(c.products, product.ID)
		c.mu.Unlock()
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Store a copy to prevent external mutation
	c.products[product.ID] = *product
	return nil
}

// DeleteCategory removes a category from the cache and cascades to delete its products
func (c *CatalogCache) DeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Remove the category
	delete(c.categories, categoryID)

	// Remove all products belonging to this category
	for productID, product := range c.products {
		if product.CategoryID == categoryID {
			delete(c.products, productID)
		}
	}

	return nil
}

// DeleteProduct removes a product from the cache
func (c *CatalogCache) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.products, productID)
	return nil
}

