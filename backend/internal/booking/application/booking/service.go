package booking

import (
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
)

type BookingService struct {
	bookingRepo      ports.BookingRepository
	availabilityRepo ports.AvailabilityRepository
	paymentService   ports.PaymentService
}

// New creates a new instance of the booking service
func New(
	bookingRepo ports.BookingRepository,
	availabilityRepo ports.AvailabilityRepository,
	paymentService ports.PaymentService,
) ports.BookingService {
	return &BookingService{
		bookingRepo:      bookingRepo,
		availabilityRepo: availabilityRepo,
		paymentService:   paymentService,
	}
}
