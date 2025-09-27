package partner

import (
	"github.com/Leviosa-care/authuser/internal/ports"
	"github.com/hengadev/encx"
)

type PartnerService struct {
	partnerRepo        ports.PartnerRepository
	userRepo           ports.UserRepository
	specializationRepo ports.SpecializationRepository
	crypto             encx.CryptoService
	stripe             ports.StripeService
}

// New creates a new instance of the partner service.
func New(
	partnerRepo ports.PartnerRepository,
	userRepo ports.UserRepository,
	specializationRepo ports.SpecializationRepository,
	crypto encx.CryptoService,
	stripe ports.StripeService,
) ports.PartnerService {
	return &PartnerService{
		partnerRepo:        partnerRepo,
		userRepo:           userRepo,
		specializationRepo: specializationRepo,
		crypto:             crypto,
		stripe:             stripe,
	}
}