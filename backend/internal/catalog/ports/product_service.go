package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

type ProductService interface {
	PrivateProductService
	PublicProductService
}

type PublicProductService interface {
	GetProductByID(ctx context.Context, ID string) (*domain.ProductRes, error)
	GetAllPublishedProducts(ctx context.Context) ([]*domain.ProductRes, error)
	GetAllProducts(ctx context.Context) ([]*domain.ProductRes, error)
}

type PrivateProductService interface {
	CreateProduct(ctx context.Context, request *domain.CreateProductRequest) (string, error)
	UpdateProduct(ctx context.Context, productID string, request *domain.UpdateProductRequest) error
	RemoveProduct(ctx context.Context, ID string) error
}
