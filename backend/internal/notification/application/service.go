package application

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/notification/ports"
)

// NotificationService is the main service coordinator
type NotificationService struct {
	mailService      *MailService
	smsService       *SMSService
	settingsProvider ports.SettingsProvider
}

func NewNotificationService(
	mailService *MailService,
	smsService *SMSService,
	settingsProvider ports.SettingsProvider,
) ports.NotificationService {
	return &NotificationService{
		mailService:      mailService,
		smsService:       smsService,
		settingsProvider: settingsProvider,
	}
}

// Email methods - delegate to MailService
func (n *NotificationService) SendOTPEmail(ctx context.Context, email, otp string) error {
	return n.mailService.SendOTPEmail(ctx, email, otp)
}

func (n *NotificationService) SendWelcomeEmail(ctx context.Context, email, firstName, lastName string) error {
	return n.mailService.SendWelcomeEmail(ctx, email, firstName, lastName)
}

func (n *NotificationService) SendVerifyEmailEmail(ctx context.Context, email, firstName, lastName string) error {
	return n.mailService.SendVerifyEmailEmail(ctx, email, firstName, lastName)
}

func (n *NotificationService) SendEventNotificationEmail(ctx context.Context, email, firstName, lastName, eventName, eventDetails string) error {
	return n.mailService.SendEventNotificationEmail(ctx, email, firstName, lastName, eventName, eventDetails)
}

func (n *NotificationService) SendPaymentNotificationEmail(ctx context.Context, email, firstName, lastName, amount, product, paymentDate string) error {
	return n.mailService.SendPaymentNotificationEmail(ctx, email, firstName, lastName, amount, product, paymentDate)
}

// SMS methods - delegate to SMSService
func (n *NotificationService) SendOTPBySMS(ctx context.Context, phoneNumber, otp string) error {
	return n.smsService.SendOTPBySMS(ctx, phoneNumber, otp)
}

func (n *NotificationService) SendGenericSMS(ctx context.Context, phoneNumber, message string) error {
	return n.smsService.SendGenericSMS(ctx, phoneNumber, message)
}

// Cache management
func (n *NotificationService) InvalidateSettingsCache() {
	n.settingsProvider.InvalidateAllCache()
}
