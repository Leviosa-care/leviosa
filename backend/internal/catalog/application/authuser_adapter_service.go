package application

import (
	"context"

	catalogDomain "github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
	"github.com/google/uuid"
)

// AuthUserCatalogAdapter provides catalog data for authuser service.
//
// This adapter implements the CatalogService interface expected by the authuser service,
// converting catalog domain types to the cached types used by authuser.
//
// In modular monolith: Direct in-process calls (this adapter)
// In microservices: HTTP-based adapter making API calls
type AuthUserCatalogAdapter struct {
	categoryService ports.CategoryService
	productService  ports.ProductService
}

// NewAuthUserCatalogAdapter creates a new adapter for authuser service.
func NewAuthUserCatalogAdapter(
	categoryService ports.CategoryService,
	productService ports.ProductService,
) *AuthUserCatalogAdapter {
	return &AuthUserCatalogAdapter{
		categoryService: categoryService,
		productService:  productService,
	}
}

// CachedCategory represents a simplified category for in-memory caching in authuser service.
// This type mirrors the authuser domain.CachedCategory to avoid import cycles.
type CachedCategory struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	Metadata    map[string]any `json:"metadata"`
}

// CachedProduct represents a simplified product for in-memory caching in authuser service.
// This type mirrors the authuser domain.CachedProduct to avoid import cycles.
type CachedProduct struct {
	ID                uuid.UUID      `json:"id"`
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	CategoryID        uuid.UUID      `json:"category_id"`
	Duration          int            `json:"duration"`
	Status            string         `json:"status"`
	Availability      string         `json:"availability"`
	BufferTime        int            `json:"buffer_time"`
	CancellationHours int            `json:"cancellation_hours"`
	StripeProductID   string         `json:"stripe_product_id"`
	Metadata          map[string]any `json:"metadata"`
}

// ListPublishedCategories returns all published categories converted to cached format.
func (a *AuthUserCatalogAdapter) ListPublishedCategories(ctx context.Context) ([]CachedCategory, error) {
	categories, err := a.categoryService.GetAllPublishedCategories(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]CachedCategory, 0, len(categories))
	for _, cat := range categories {
		result = append(result, convertCategoryToCached(cat))
	}

	return result, nil
}

// ListPublishedProducts returns all published products converted to cached format.
func (a *AuthUserCatalogAdapter) ListPublishedProducts(ctx context.Context) ([]CachedProduct, error) {
	products, err := a.productService.GetAllPublishedProducts(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]CachedProduct, 0, len(products))
	for _, prod := range products {
		result = append(result, convertProductToCached(prod))
	}

	return result, nil
}

// convertCategoryToCached converts catalog Category to CachedCategory.
func convertCategoryToCached(cat *catalogDomain.Category) CachedCategory {
	return CachedCategory{
		ID:          cat.ID,
		Name:        cat.Name,
		Description: cat.Description,
		Status:      string(cat.Status), // Convert PublishedStatus to string
		Metadata:    cat.Metadata,
	}
}

// convertProductToCached converts catalog ProductRes to CachedProduct.
func convertProductToCached(prod *catalogDomain.ProductRes) CachedProduct {
	return CachedProduct{
		ID:          prod.ID,
		Name:        prod.Name,
		Description: prod.Description,
		// CategoryID:        prod.CategoryID,
		Duration:          prod.Duration,
		Status:            string(prod.Status), // Convert PublishedStatus to string
		Availability:      string(prod.Availability),
		BufferTime:        prod.BufferTime,
		CancellationHours: prod.CancellationHours,
		StripeProductID:   prod.StripeProductID,
		Metadata:          prod.Metadata,
	}
}
