package otpRepository

import (
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"
	"github.com/redis/go-redis/v9"
)

const (
	OTPKeyPrefix = "authuser:otp:"
)

type OTPRepository struct {
	client *redis.Client
}

func New(client *redis.Client) ports.OTPRepository {
	return &OTPRepository{
		client: client,
	}
}

// FormatOTPKey formats an OTP key for Redis storage
// This function is public to allow consistent key formatting in tests
func FormatOTPKey(emailHash string) string {
	return fmt.Sprintf("%s%s", OTPKeyPrefix, emailHash)
}
