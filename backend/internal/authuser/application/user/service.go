package user

import (
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"
	"github.com/hengadev/encx"
)

type UserService struct {
	repo   ports.UserRepository
	crypto encx.CryptoService
	stripe ports.StripeService
}

// New creates a new instance of the aggregator service.
func New(user ports.UserRepository, crypto encx.CryptoService, stripe ports.StripeService) ports.UserService {
	return &UserService{
		repo:   user,
		crypto: crypto,
		stripe: stripe,
	}
}
