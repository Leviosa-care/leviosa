package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

type ProductPaymentGateway interface {
	CreateProduct(ctx context.Context, req domain.CreateStripeProductRequest) (*domain.PaymentProduct, error)
	GetProduct(ctx context.Context, productID string) (*domain.PaymentProduct, error)
	UpdateProduct(ctx context.Context, productID string, req *domain.UpdateStripeProductRequest) (*domain.PaymentProduct, error)
	DeactivateProduct(ctx context.Context, productID string) error
	ReactivateProduct(ctx context.Context, productID string) error
}
