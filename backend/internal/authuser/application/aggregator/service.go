package aggregator

import (
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"
)

type AuthAggregatorService struct {
	otp           ports.OTPService
	user          ports.UserService
	session       ports.SessionService
	partner       ports.PartnerService
	bookingClient ports.BookingClient
}

// New creates a new AuthAggregatorService. The returned concrete pointer satisfies
// ports.AuthAggregatorService and also exposes SetBookingClient for late injection.
func New(otp ports.OTPService, user ports.UserService, session ports.SessionService, partner ports.PartnerService, bookingClient ports.BookingClient) *AuthAggregatorService {
	return &AuthAggregatorService{
		otp:           otp,
		user:          user,
		session:       session,
		partner:       partner,
		bookingClient: bookingClient,
	}
}

// SetBookingClient allows late injection of the BookingClient dependency.
// This is needed because BookingService is created after AuthAggregatorService
// in the dependency graph.
func (s *AuthAggregatorService) SetBookingClient(client ports.BookingClient) {
	s.bookingClient = client
}
