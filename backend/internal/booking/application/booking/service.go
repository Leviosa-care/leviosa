package booking

import (
	"github.com/Leviosa-care/leviosa/backend/internal/booking/ports"
	catalogPorts "github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
	"github.com/hengadev/encx"
)

type BookingService struct {
	bookingRepo         ports.BookingRepository
	availabilityRepo    ports.AvailabilityRepository
	paymentService      ports.PaymentService
	productService      catalogPorts.PublicProductService
	priceService        catalogPorts.PublicPriceService
	notificationService ports.BookingNotificationService
	roomService         ports.RoomService
	authUserClient      ports.AuthUserClient
	crypto              encx.CryptoService
	tokenSecret         []byte
}

// New creates a new instance of the booking service
func New(
	bookingRepo ports.BookingRepository,
	availabilityRepo ports.AvailabilityRepository,
	paymentService ports.PaymentService,
	productService catalogPorts.PublicProductService,
	priceService catalogPorts.PublicPriceService,
	notificationService ports.BookingNotificationService,
	crypto encx.CryptoService,
	opts ...Option,
) ports.BookingService {
	svc := &BookingService{
		bookingRepo:         bookingRepo,
		availabilityRepo:    availabilityRepo,
		paymentService:      paymentService,
		productService:      productService,
		priceService:        priceService,
		notificationService: notificationService,
		crypto:              crypto,
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

// Option configures optional BookingService dependencies.
type Option func(*BookingService)

// WithRoomService injects the room service dependency.
func WithRoomService(rs ports.RoomService) Option {
	return func(s *BookingService) { s.roomService = rs }
}

// WithAuthUserClient injects the auth-user client dependency.
func WithAuthUserClient(c ports.AuthUserClient) Option {
	return func(s *BookingService) { s.authUserClient = c }
}

// WithTokenSecret sets the HMAC secret used to sign booking tokens.
func WithTokenSecret(secret []byte) Option {
	return func(s *BookingService) { s.tokenSecret = secret }
}
