package availability

import (
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	catalogPorts "github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"

	"github.com/hengadev/encx"
)

type AvailabilityService struct {
	availabilityRepo ports.AvailabilityRepository
	allocationRepo   ports.RoomAllocationRepository
	roomRepo         ports.RoomRepository
	productService   catalogPorts.PublicProductService
	crypto           encx.CryptoService
	// authUserClient   ports.AuthUserClient
}

// New creates a new instance of the availability service
func New(
	availabilityRepo ports.AvailabilityRepository,
	allocationRepo ports.RoomAllocationRepository,
	roomRepo ports.RoomRepository,
	productService catalogPorts.PublicProductService,
	crypto encx.CryptoService,
	// authUserClient ports.AuthUserClient,
) ports.AvailabilityService {
	return &AvailabilityService{
		availabilityRepo: availabilityRepo,
		allocationRepo:   allocationRepo,
		roomRepo:         roomRepo,
		productService:   productService,
		crypto:           crypto,
		// authUserClient:   authUserClient,
	}
}
