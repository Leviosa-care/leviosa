package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

type PriceService interface {
	CreatePrice(ctx context.Context, productID string, request *domain.CreatePriceRequest) (string, error)
	GetPrice(ctx context.Context, priceID string) (*domain.Price, error)
	GetPriceByStripeID(ctx context.Context, stripePriceID string) (*domain.Price, error) // Use this for stripe API lookups
	GetAllPrices(ctx context.Context) ([]*domain.Price, error)
	GetPricesByProductID(ctx context.Context, productID string) ([]*domain.Price, error)
	UpdatePrice(ctx context.Context, priceID string, input domain.UpdatePriceRequest) (*domain.Price, error)
	DeactivatePrice(ctx context.Context, priceID string) error // Specific deactivation if needed, or handled by UpdatePrice
}
