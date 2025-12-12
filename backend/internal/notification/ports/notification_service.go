package ports

import "context"

// NotificationService defines the public interface for the notification module
type NotificationService interface {
	// Email methods with explicit parameters (no User struct dependency)
	SendOTPEmail(ctx context.Context, email, otp string) error
	SendWelcomeEmail(ctx context.Context, email, firstName, lastName string) error
	SendVerifyEmailEmail(ctx context.Context, email, firstName, lastName string) error
	SendEventNotificationEmail(ctx context.Context, email, firstName, lastName, eventName, eventDetails string) error
	SendPaymentNotificationEmail(ctx context.Context, email, firstName, lastName, amount, product, paymentDate string) error

	// SMS methods
	SendOTPBySMS(ctx context.Context, phoneNumber, otp string) error
	SendGenericSMS(ctx context.Context, phoneNumber, message string) error

	// Cache management
	InvalidateSettingsCache()
}
