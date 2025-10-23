package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

type PriceRepository interface {
	// reader
	GetPrice(ctx context.Context, priceID string) (*domain.Price, error)
	GetPriceByStripeID(ctx context.Context, stripePriceID string) (*domain.Price, error)
	GetPricesByProductID(ctx context.Context, productID string, activeOnly bool) ([]*domain.Price, error)
	GetProductIDByStripeProductID(ctx context.Context, stripeProductID string) (string, error)
	GetAllPrices(ctx context.Context) ([]*domain.Price, error)
	// writer
	CreatePrice(ctx context.Context, price *domain.Price) error
	UpdatePrice(ctx context.Context, priceID string, patch *domain.UpdatePriceRequest) error
	// DeactivatePrices(ctx context.Context, priceIDs []string) error
}
