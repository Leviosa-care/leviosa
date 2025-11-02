package partner

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/application/catalog"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"

	// authRabbitMQ "github.com/Leviosa-care/leviosa/backend/internal/authuser/infrastructure/rabbitmq"
	"github.com/hengadev/encx"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PartnerService struct {
	partnerRepo    ports.PartnerRepository
	userRepo       ports.UserRepository
	catalogCache   *catalog.CatalogCache // Concrete type for both read and write access
	catalogService ports.CatalogService
	crypto         encx.CryptoService
	stripe         ports.StripeService
}

// New creates a new instance of the partner service.
//
// Catalog Integration Architecture (Modular Monolith):
// - catalogService: Direct in-process calls to catalog service for initial cache population
// - catalogCache: In-memory cache for fast validation (populated on startup + real-time updates)
// - RabbitMQ consumer: Receives real-time catalog updates (categories/products created/updated/deleted)
//
// Future Microservices Migration:
// - Replace catalogService with HTTP-based implementation (same interface)
// - Cache and RabbitMQ consumer remain unchanged
// - No changes to business logic required
func New(
	ctx context.Context,
	partnerRepo ports.PartnerRepository,
	userRepo ports.UserRepository,
	catalogService ports.CatalogService,
	conn *amqp.Connection,
	crypto encx.CryptoService,
	stripe ports.StripeService,
) (ports.PartnerService, error) {
	// Create catalog cache - using concrete type for both read and write access
	catalogCache := catalog.NewCatalogCache()

	partnerService := &PartnerService{
		partnerRepo:    partnerRepo,
		userRepo:       userRepo,
		catalogCache:   catalogCache,
		catalogService: catalogService,
		crypto:         crypto,
		stripe:         stripe,
	}

	// Load catalog data on startup using direct service call
	if err := partnerService.loadCatalogValues(ctx); err != nil {
		return nil, fmt.Errorf("failed to load catalog data: %w", err)
	}

	// TODO: Start RabbitMQ consumer for real-time catalog updates when catalog service is implemented
	// consumer := authRabbitMQ.NewCatalogConsumer(conn, catalogCache)
	// if err := consumer.Start(ctx); err != nil {
	// 	return nil, fmt.Errorf("failed to start catalog consumer: %w", err)
	// }

	return partnerService, nil
}

// loadCatalogValues populates the catalog cache with published categories and products.
//
// In modular monolith: Direct in-process call to catalog service
// In microservices: Would be HTTP GET to catalog API
//
// This method is called once on service initialization to populate the cache.
// Real-time updates are then handled by the RabbitMQ consumer.
func (s *PartnerService) loadCatalogValues(ctx context.Context) error {
	// Load published categories from catalog service
	categories, err := s.catalogService.ListPublishedCategories(ctx)
	if err != nil {
		return fmt.Errorf("load published categories: %w", err)
	}

	// Populate cache with categories
	for _, cat := range categories {
		if err := s.catalogCache.UpsertCategory(ctx, &cat); err != nil {
			return fmt.Errorf("cache category %s: %w", cat.ID, err)
		}
	}

	// Load published products from catalog service
	products, err := s.catalogService.ListPublishedProducts(ctx)
	if err != nil {
		return fmt.Errorf("load published products: %w", err)
	}

	// Populate cache with products
	for _, prod := range products {
		if err := s.catalogCache.UpsertProduct(ctx, &prod); err != nil {
			return fmt.Errorf("cache product %s: %w", prod.ID, err)
		}
	}

	return nil
}
