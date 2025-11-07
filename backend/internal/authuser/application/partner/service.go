package partner

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"
	catalogPorts "github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"

	"github.com/hengadev/encx"
)

type PartnerService struct {
	partnerRepo     ports.PartnerRepository
	userRepo        ports.UserRepository
	productService  catalogPorts.PublicProductService
	categoryService catalogPorts.PublicCategoryService
	crypto          encx.CryptoService
	stripe          ports.StripeService
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
	productService catalogPorts.PublicProductService,
	categoryService catalogPorts.PublicCategoryService,
	crypto encx.CryptoService,
	stripe ports.StripeService,
) (ports.PartnerService, error) {
	partnerService := &PartnerService{
		partnerRepo:     partnerRepo,
		userRepo:        userRepo,
		productService:  productService,
		categoryService: categoryService,
		crypto:          crypto,
		stripe:          stripe,
	}

	return partnerService, nil
}
