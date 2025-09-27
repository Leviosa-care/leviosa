package specialization

import (
	"github.com/Leviosa-care/authuser/internal/ports"
	"github.com/hengadev/encx"
)

type SpecializationService struct {
	repo   ports.SpecializationRepository
	crypto encx.CryptoService
}

// New creates a new instance of the specialization service.
func New(repo ports.SpecializationRepository, crypto encx.CryptoService) ports.SpecializationService {
	return &SpecializationService{
		repo:   repo,
		crypto: crypto,
	}
}