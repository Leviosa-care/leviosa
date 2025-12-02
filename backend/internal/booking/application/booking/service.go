package booking

import (
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	catalogPorts "github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
	"github.com/hengadev/encx"
)

type BookingService struct {
	bookingRepo      ports.BookingRepository
	availabilityRepo ports.AvailabilityRepository
	paymentService   ports.PaymentService
	productService   catalogPorts.PublicProductService
	crypto           encx.CryptoService
}

// New creates a new instance of the booking service
func New(
	bookingRepo ports.BookingRepository,
	availabilityRepo ports.AvailabilityRepository,
	paymentService ports.PaymentService,
	productService catalogPorts.PublicProductService,
	crypto encx.CryptoService,
) ports.BookingService {
	return &BookingService{
		bookingRepo:      bookingRepo,
		availabilityRepo: availabilityRepo,
		paymentService:   paymentService,
		productService:   productService,
		crypto:           crypto,
	}
}
