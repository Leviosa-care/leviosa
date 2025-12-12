package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/notification/domain"
)

// EmailService defines the interface for email sending implementations (SMTP)
type EmailService interface {
	SendOTPEmail(ctx context.Context, req domain.OTPEmailRequest) error
	SendWelcomeEmail(ctx context.Context, req domain.WelcomeEmailRequest) error
	SendVerifyEmailEmail(ctx context.Context, req domain.VerifyEmailRequest) error
	SendEventNotificationEmail(ctx context.Context, req domain.EventNotificationRequest) error
	SendPaymentNotificationEmail(ctx context.Context, req domain.PaymentNotificationRequest) error
}
