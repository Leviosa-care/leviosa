package price

import (
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
)

type PriceService struct {
	repo       ports.PriceRepository
	sharedRepo ports.SharedRepository
	stripe     ports.PricePaymentGateway
}

func New(repo ports.PriceRepository, productRepo ports.SharedRepository, stripe ports.PricePaymentGateway) ports.PriceService {
	return &PriceService{
		repo:       repo,
		sharedRepo: productRepo,
		stripe:     stripe,
	}
}
