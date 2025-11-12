package room

import (
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/hengadev/encx"
)

type RoomService struct {
	roomRepo     ports.RoomRepository
	buildingRepo ports.BuildingRepository
	crypto       encx.CryptoService
}

// New creates a new instance of the room service
func New(roomRepo ports.RoomRepository, buildingRepo ports.BuildingRepository, crypto encx.CryptoService) ports.RoomService {
	return &RoomService{
		roomRepo:     roomRepo,
		buildingRepo: buildingRepo,
		crypto:       crypto,
	}
}

