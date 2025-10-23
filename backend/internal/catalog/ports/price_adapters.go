package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

type PricePaymentGateway interface {
	CreatePrice(ctx context.Context, input domain.CreateStripePriceRequest) (*domain.PaymentPrice, error)
	GetPrice(ctx context.Context, priceID string) (*domain.PaymentPrice, error)
	GetPricesByProductID(ctx context.Context, productID string, opts *domain.PriceListOptions) ([]*domain.PaymentPrice, error)
	UpdatePrice(ctx context.Context, stripePriceID string, req domain.UpdateStripePriceRequest) (*domain.PaymentPrice, error)
	DeactivatePrices(ctx context.Context, priceIDs []string) error
	ReactivatePrices(ctx context.Context, stripePriceIDs []string) error
}
