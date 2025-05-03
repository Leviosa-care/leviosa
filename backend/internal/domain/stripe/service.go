package stripeService

import "context"

type Service interface {
	CreatePrice(ctx context.Context, productID string, priceValue int64) (string, error)
	CreateProduct(ctx context.Context, object Payment) (string, error)
	CreateCheckoutSession(ctx context.Context, priceID string, quantity int64) (string, error)
	RemovePrice(ctx context.Context, priceID string) error
	RemoveProduct(ctx context.Context, productID string) (string, error)
}

type service struct{}

func New() Service {
	return &service{}
}
