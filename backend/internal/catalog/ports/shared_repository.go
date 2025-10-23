package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/google/uuid"
)

type SharedRepository interface {
	GetProductByID(ctx context.Context, productID uuid.UUID) (*domain.Product, error)
	GetStripeProductAndPriceIDs(ctx context.Context, productID uuid.UUID) (string, []string, error)
	GetCategoryByID(ctx context.Context, categoryID uuid.UUID) (*domain.Category, error)
	DeactivatePrices(ctx context.Context, priceIDs []string) error // Specific deactivation if needed, or handled by UpdatePrice
	ReactivatePrices(ctx context.Context, stripePriceIDs []string) error
}
