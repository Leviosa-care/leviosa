package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/google/uuid"
)

type ProductRepository interface {
	// reader
	GetAllProducts(ctx context.Context) ([]*domain.ProductRes, error)
	GetAllPublishedProducts(ctx context.Context) ([]*domain.ProductRes, error)
	// writer
	AddProduct(ctx context.Context, p *domain.Product) (string, error)
	UpdateProduct(ctx context.Context, productID uuid.UUID, p *domain.UpdateProductRequest) error
	DeleteProduct(ctx context.Context, productID uuid.UUID) error
}
