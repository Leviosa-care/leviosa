package building

import (
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"

	"github.com/hengadev/encx"
)

type BuildingService struct {
	buildingRepo ports.BuildingRepository
	crypto       encx.CryptoService
}

// New creates a new instance of the building service
func New(buildingRepo ports.BuildingRepository, crypto encx.CryptoService) ports.BuildingService {
	return &BuildingService{
		buildingRepo: buildingRepo,
		crypto:       crypto,
	}
}
