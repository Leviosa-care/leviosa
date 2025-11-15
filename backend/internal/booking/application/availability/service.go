package availability

import (
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
)

type AvailabilityService struct {
	availabilityRepo ports.AvailabilityRepository
	allocationRepo   ports.RoomAllocationRepository
	roomRepo         ports.RoomRepository
	authUserClient   ports.AuthUserClient
}

// New creates a new instance of the availability service
func New(
	availabilityRepo ports.AvailabilityRepository,
	allocationRepo ports.RoomAllocationRepository,
	roomRepo ports.RoomRepository,
	authUserClient ports.AuthUserClient,
) ports.AvailabilityService {
	return &AvailabilityService{
		availabilityRepo: availabilityRepo,
		allocationRepo:   allocationRepo,
		roomRepo:         roomRepo,
		authUserClient:   authUserClient,
	}
}
