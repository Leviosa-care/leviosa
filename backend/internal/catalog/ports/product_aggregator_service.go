package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

type ProductAggregatorService interface {
	CreateProductWithPrice(ctx context.Context, request *domain.CreateProductWithPriceRequest) (string, string, error)
	GetAdminAllProducts(ctx context.Context) ([]*domain.ProductAggregator, error)
	GetAllPublishedProducts(ctx context.Context) ([]*domain.ProductAggregator, error)
	GetProductByID(ctx context.Context, productIDStr string) (*domain.ProductAggregator, error)
}
