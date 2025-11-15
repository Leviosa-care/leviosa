package availability

import (
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	"github.com/hengadev/encx"
)

type AvailabilityService struct {
	availabilityRepo ports.AvailabilityRepository
	allocationRepo   ports.RoomAllocationRepository
	roomRepo         ports.RoomRepository
	crypto           encx.CryptoService
	authUserClient   ports.AuthUserClient
}

// New creates a new instance of the availability service
func New(
	availabilityRepo ports.AvailabilityRepository,
	allocationRepo ports.RoomAllocationRepository,
	roomRepo ports.RoomRepository,
	crypto encx.CryptoService,
	// authUserClient ports.AuthUserClient,
) ports.AvailabilityService {
	return &AvailabilityService{
		availabilityRepo: availabilityRepo,
		allocationRepo:   allocationRepo,
		roomRepo:         roomRepo,
		crypto:           crypto,
		// authUserClient:   authUserClient,
	}
}
