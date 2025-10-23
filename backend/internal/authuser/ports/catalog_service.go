package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
)

// CatalogService provides access to catalog data for in-process communication.
//
// This interface supports the modular monolith architecture by allowing direct
// service-to-service calls without HTTP overhead. When migrating to microservices,
// implement an HTTPCatalogService adapter that makes HTTP calls to the catalog API.
//
// The interface is read-only from the perspective of the authuser service - catalog
// data is managed exclusively by the catalog service and synchronized via RabbitMQ events.
type CatalogService interface {
	// ListPublishedCategories returns all categories with "published" status.
	// Used to populate the catalog cache on service initialization.
	//
	// In modular monolith: Direct database query within same process
	// In microservices: HTTP GET /api/v1/categories?status=published
	ListPublishedCategories(ctx context.Context) ([]domain.CachedCategory, error)

	// ListPublishedProducts returns all products with "published" status.
	// Used to populate the catalog cache on service initialization.
	//
	// In modular monolith: Direct database query within same process
	// In microservices: HTTP GET /api/v1/products?status=published
	ListPublishedProducts(ctx context.Context) ([]domain.CachedProduct, error)
}
