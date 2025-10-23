package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

type ProductService interface {
	CreateProduct(ctx context.Context, p *domain.CreateProductRequest) (string, error)
	GetProductByID(ctx context.Context, productID string) (*domain.ProductRes, error)
	GetAllPublishedProducts(ctx context.Context) ([]*domain.ProductRes, error)
	GetAllProducts(ctx context.Context) ([]*domain.ProductRes, error)
	UpdateProduct(ctx context.Context, productID string, p *domain.UpdateProductRequest) error
	RemoveProduct(ctx context.Context, ID string) error
}
