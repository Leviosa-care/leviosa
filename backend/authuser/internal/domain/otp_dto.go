package domain

import (
	"time"
)

// OTPSentEvent represents the data returned to client after resending OTP
type OTPSentEvent struct {
	Code      string    `json:"code"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expires_at"`
}
