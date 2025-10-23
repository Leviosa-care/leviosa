package product

import (
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
)

// compile time assertion check if *ServiceImpl implements Service interface
var _ ports.ProductService = (*ProductService)(nil)

type ProductService struct {
	repo        ports.ProductRepository
	sharedRepo  ports.SharedRepository
	stripe      ports.ProductPaymentGateway
	priceStripe ports.PricePaymentGateway
}

func New(repo ports.ProductRepository, sharedRepo ports.SharedRepository, stripe ports.ProductPaymentGateway, priceStripe ports.PricePaymentGateway) ports.ProductService {
	return &ProductService{
		repo:        repo,
		sharedRepo:  sharedRepo,
		stripe:      stripe,
		priceStripe: priceStripe,
	}
}
