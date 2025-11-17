package allocation

import (
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
)

type RoomAllocationService struct {
	allocationRepo ports.RoomAllocationRepository
	roomRepo       ports.RoomRepository
	authUserClient ports.AuthUserClient
}

// New creates a new instance of the room allocation service
func New(allocationRepo ports.RoomAllocationRepository, roomRepo ports.RoomRepository, authUserClient ports.AuthUserClient) ports.RoomAllocationService {
	return &RoomAllocationService{
		allocationRepo: allocationRepo,
		roomRepo:       roomRepo,
		authUserClient: authUserClient,
	}
}
