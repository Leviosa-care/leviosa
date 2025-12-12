package application

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/internal/domain/user/models"
	"github.com/Leviosa-care/notification/internal/application/email"
	"github.com/Leviosa-care/notification/internal/application/sms"
	"github.com/Leviosa-care/notification/internal/ports"
)

type NotificationService struct {
	emailService *email.EmailService
	smsService   *sms.SMSService
}

func New(emailService *email.EmailService, smsService *sms.SMSService) ports.NotificationService {
	return &NotificationService{
		emailService: emailService,
		smsService:   smsService,
	}
}

func (s *NotificationService) SendOTP(ctx context.Context, email string, otp string) error {
	return s.emailService.SendOTP(ctx, email, otp)
}

func (s *NotificationService) SendWelcomeEmail(ctx context.Context, email string) error {
	return s.emailService.SendWelcomeEmail(ctx, email)
}

func (s *NotificationService) SendPasswordResetEmail(ctx context.Context, email string) error {
	return s.emailService.SendPasswordResetEmail(ctx, email)
}

func (s *NotificationService) SendEventNotification(ctx context.Context, users []*models.User, eventTime string) error {
	return s.emailService.SendEventNotification(ctx, users, eventTime)
}

func (s *NotificationService) SendPaymentNotification(ctx context.Context, user *models.User, eventTime string) error {
	return s.emailService.SendPaymentNotification(ctx, user, eventTime)
}

func (s *NotificationService) SendVoteNotification(ctx context.Context, user *models.User, eventTime string) error {
	return s.emailService.SendVoteNotification(ctx, user, eventTime)
}

func (s *NotificationService) SendRegistrationReminder(ctx context.Context, user *models.User, registrationName string, daysLeft int) error {
	return s.emailService.SendRegistrationReminder(ctx, user, registrationName, daysLeft)
}

func (s *NotificationService) SendSMS(ctx context.Context, phone, message string) error {
	return s.smsService.SendSMS(ctx, phone, message)
}

func (s *NotificationService) UpdateCompanyEmail(email string) {
	// This will be handled by the RabbitMQ consumer
}

func (s *NotificationService) UpdateCompanyLogo(logo []byte) {
	// This will be handled by the RabbitMQ consumer
}

func (s *NotificationService) UpdateCompanyLegalAddress(address string) {
	// This will be handled by the RabbitMQ consumer
}

func (s *NotificationService) UpdateCompanyInstagram(instagram string) {
	// This will be handled by the RabbitMQ consumer
}
