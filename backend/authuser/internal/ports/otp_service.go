package ports

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
)

type OTPService interface {
	// Settings cache management (already implemented)
	SetOTPDuration(duration int)
	SetOTPLength(length int)
	SetOTPMaxAttempts(maxAttempts int)
	GetOTPDuration() int
	GetOTPLength() int
	GetOTPMaxAttempts() int
	// application
	// HACK: loading
	ValidateOTP(ctx context.Context, request *domain.ValidateOTPRequest) error
	CreateOTP(ctx context.Context, email string) error
	CancelOTP(ctx context.Context, email string) error
}
