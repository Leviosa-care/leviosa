package aggregator

import (
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"
)

type AuthAggregatorService struct {
	otp     ports.OTPService
	user    ports.UserService
	session ports.SessionService
	partner ports.PartnerService
}

// NewAuthAggregatorService creates a new instance of the aggregator service.
func New(otp ports.OTPService, user ports.UserService, session ports.SessionService, partner ports.PartnerService) ports.AuthAggregatorService {
	return &AuthAggregatorService{
		otp:     otp,
		user:    user,
		session: session,
		partner: partner,
	}
}
