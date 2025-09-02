package aggregator

import (
	"github.com/Leviosa-care/authuser/internal/ports"
)

type AuthAggregatorService struct {
	otp     ports.OTPService
	user    ports.UserService
	session ports.SessionService
}

// NewAuthAggregatorService creates a new instance of the aggregator service.
func New(otp ports.OTPService, user ports.UserService, session ports.SessionService) ports.AuthAggregatorService {
	return &AuthAggregatorService{
		otp:     otp,
		user:    user,
		session: session,
	}
}
