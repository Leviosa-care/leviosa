package building

import (
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
)

type BuildingService struct {
	buildingRepo ports.BuildingRepository
}

// New creates a new instance of the building service
func New(buildingRepo ports.BuildingRepository) ports.BuildingService {
	return &BuildingService{
		buildingRepo: buildingRepo,
	}
}
