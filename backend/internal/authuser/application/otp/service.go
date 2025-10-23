package otp

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"

	"github.com/hengadev/encx"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	defaultOTPDuration    = 10 // Default OTP duration in minutes
	defaultOTPLength      = 6  // Default OTP length
	defaultOTPMaxAttempts = 3  // Default max attempts
)

type OTPService struct {
	repo   ports.OTPRepository
	crypto encx.CryptoService
	cache  ports.OTPCache
	mq     *amqp.Connection
}

// New creates a new OTPService instance with hardcoded OTP configuration.
//
// OTP settings are currently hardcoded as constants (defaultOTPDuration, defaultOTPLength, defaultOTPMaxAttempts)
// for simplicity and production reliability. The cache layer and RabbitMQ consumer infrastructure
// (StartOTPSettingConsumer) are preserved but not used, allowing for future migration to a microservices
// architecture where dynamic settings may be needed.
//
// To re-enable dynamic settings in the future:
// 1. Call StartOTPSettingConsumer() after service initialization
// 2. Optionally implement HTTP-based settings loading for initial values
// 3. Update application code to use s.GetOTPDuration(), s.GetOTPLength(), s.GetOTPMaxAttempts() instead of constants
func New(ctx context.Context, repo ports.OTPRepository, crypto encx.CryptoService, rabbitConn *amqp.Connection) (ports.OTPService, error) {
	// Initialize cache with hardcoded default values
	cache := domain.NewOTPCache(
		defaultOTPDuration,
		defaultOTPLength,
		defaultOTPMaxAttempts,
	)

	otpService := &OTPService{
		repo:   repo,
		crypto: crypto,
		cache:  cache,
		mq:     rabbitConn,
	}

	return otpService, nil
}
