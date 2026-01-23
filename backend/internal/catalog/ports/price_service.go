package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

type PriceService interface {
	PrivatePriceService
	PublicPriceService
}

// PublicPriceService contains read-only methods for price lookups.
// Used by other services (e.g., booking) for cross-module queries.
type PublicPriceService interface {
	GetPrice(ctx context.Context, priceID string) (*domain.Price, error)
	GetPriceByStripeID(ctx context.Context, stripePriceID string) (*domain.Price, error)
	GetAllPrices(ctx context.Context) ([]*domain.Price, error)
	GetPricesByProductID(ctx context.Context, productID string) ([]*domain.Price, error)
	GetActiveOneTimePriceByProductID(ctx context.Context, productID string, currency string) (*domain.Price, error)
}

// PrivatePriceService contains write methods for price management.
// Used by catalog handlers for admin operations.
type PrivatePriceService interface {
	CreatePrice(ctx context.Context, productID string, request *domain.CreatePriceRequest) (string, error)
	UpdatePrice(ctx context.Context, priceID string, input domain.UpdatePriceRequest) (*domain.Price, error)
	DeactivatePrice(ctx context.Context, priceID string) error
}
