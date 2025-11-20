package allocation

import (
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/hengadev/encx"
)

type RoomAllocationService struct {
	allocationRepo ports.RoomAllocationRepository
	roomRepo       ports.RoomRepository
	authUserClient ports.AuthUserClient
	crypto         encx.CryptoService
}

// New creates a new instance of the room allocation service
func New(allocationRepo ports.RoomAllocationRepository, roomRepo ports.RoomRepository, authUserClient ports.AuthUserClient, crypto encx.CryptoService) ports.RoomAllocationService {
	return &RoomAllocationService{
		allocationRepo: allocationRepo,
		roomRepo:       roomRepo,
		authUserClient: authUserClient,
		crypto:         crypto,
	}
}
