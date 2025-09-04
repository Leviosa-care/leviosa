package ports

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
)

// signup flow for the frontend pages
// 1. auth/email (email)
// 2. auth/otp (verify the OTP sent)
// 3. auth/general (firstname, lastname, age, gender etc...)
// 4. auth/address (self explanatory)
// 5. auth/password (self explanatory)
// 6. auth/pending (redirect user to a page that allows them to visit the website as a pending user because admin need to approve user)

type AuthAggregatorService interface {
	CheckEmailSendOTP(ctx context.Context, request *domain.CheckEmailAvailabilityRequest) error
	ValidateOTPCreatePendingUser(ctx context.Context, request *domain.ValidateOTPRequest) (*domain.CreateSessionResponse, error)
}
