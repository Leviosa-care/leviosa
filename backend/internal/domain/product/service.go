package productService

import "context"

type Service interface {
	CreateOffer(ctx context.Context, offer *Offer) error
	CreateProduct(ctx context.Context, product *Product) error
	GetProduct(ctx context.Context, productID string) (*Product, error)
	RemoveOffer(ctx context.Context, offerID int) error
	RemoveProduct(ctx context.Context, productID string) error
	UpdateProductType(ctx context.Context, product *Offer) error
	UpdateProduct(ctx context.Context, product *Product) error
}

type service struct {
	repo ReadWriter
}

func New(repo ReadWriter) Service {
	return &service{repo}
}
