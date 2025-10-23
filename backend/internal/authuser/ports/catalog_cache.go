package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/google/uuid"
)

// CatalogCache provides read-only access to catalog data for services
// This interface is used by services like Partner to validate against catalog items
type CatalogCache interface {
	// GetCategoryByID returns a category by ID and whether it exists
	GetCategoryByID(id uuid.UUID) (*domain.CachedCategory, bool)

	// GetProductByID returns a product by ID and whether it exists
	GetProductByID(id uuid.UUID) (*domain.CachedProduct, bool)

	// ListCategories returns all cached categories
	ListCategories() []domain.CachedCategory

	// ListProducts returns all cached products
	ListProducts() []domain.CachedProduct

	// ListProductsByCategory returns products filtered by category ID
	ListProductsByCategory(categoryID uuid.UUID) []domain.CachedProduct

	// IsValidCategory checks if a category exists in the cache
	IsValidCategory(categoryID uuid.UUID) bool

	// IsValidProduct checks if a product exists in the cache
	IsValidProduct(productID uuid.UUID) bool
}

// CatalogCacheUpdater provides write access to catalog data for consumers
// This interface is used by the RabbitMQ consumer to update the cache
type CatalogCacheUpdater interface {
	// UpsertCategory adds or updates a category in the cache
	// Only categories with "published" status will be stored
	UpsertCategory(ctx context.Context, category *domain.CachedCategory) error

	// UpsertProduct adds or updates a product in the cache
	// Only products with "published" status will be stored
	UpsertProduct(ctx context.Context, product *domain.CachedProduct) error

	// DeleteCategory removes a category from the cache and cascades to delete its products
	DeleteCategory(ctx context.Context, categoryID uuid.UUID) error

	// DeleteProduct removes a product from the cache
	DeleteProduct(ctx context.Context, productID uuid.UUID) error
}
