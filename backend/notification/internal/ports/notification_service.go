package ports

import (
	"context"

	"github.com/Leviosa-care/internal/domain/user/models"
)

type NotificationService interface {
	// Email notifications
	SendOTP(ctx context.Context, email string, otp string) error
	SendWelcomeEmail(ctx context.Context, email string) error
	SendPasswordResetEmail(ctx context.Context, email string) error
	SendEventNotification(ctx context.Context, users []*models.User, eventTime string) error
	SendPaymentNotification(ctx context.Context, user *models.User, eventTime string) error
	SendVoteNotification(ctx context.Context, user *models.User, eventTime string) error
	SendRegistrationReminder(ctx context.Context, user *models.User, registrationName string, daysLeft int) error

	// SMS notifications
	SendSMS(ctx context.Context, phone, message string) error

	// Settings cache management
	UpdateCompanyEmail(email string)
	UpdateCompanyLogo(logo []byte)
	UpdateCompanyLegalAddress(address string)
	UpdateCompanyInstagram(instagram string)
}
