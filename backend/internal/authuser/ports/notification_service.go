package ports

import "context"

// NotificationService defines the interface for sending OTP notifications.
// This is an anti-corruption port — authuser does not depend on the notification
// module's concrete types.
type NotificationService interface {
	SendOTPEmail(ctx context.Context, email, otp string) error
}
