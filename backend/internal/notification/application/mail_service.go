package application

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/notification/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/notification/ports"
)

// MailService handles email sending logic
type MailService struct {
	emailClient      ports.EmailService
	settingsProvider ports.SettingsProvider
}

func NewMailService(emailClient ports.EmailService, settingsProvider ports.SettingsProvider) *MailService {
	return &MailService{
		emailClient:      emailClient,
		settingsProvider: settingsProvider,
	}
}

func (s *MailService) SendOTPEmail(ctx context.Context, email, otp string) error {
	companyEmail, err := s.settingsProvider.GetCompanyEmail(ctx)
	if err != nil {
		return fmt.Errorf("failed to get company email for OTP: %w", err)
	}

	logoURL, err := s.settingsProvider.GetCompanyLogo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get company logo for OTP: %w", err)
	}

	request := domain.OTPEmailRequest{
		ToEmail:   email,
		OTP:       otp,
		FromEmail: companyEmail,
		LogoURL:   logoURL,
	}

	if err := s.emailClient.SendOTPEmail(ctx, request); err != nil {
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	return nil
}

func (s *MailService) SendWelcomeEmail(ctx context.Context, email, firstName, lastName string) error {
	companyEmail, err := s.settingsProvider.GetCompanyEmail(ctx)
	if err != nil {
		return fmt.Errorf("failed to get company email for welcome: %w", err)
	}

	companyName, err := s.settingsProvider.GetCompanyName(ctx)
	if err != nil {
		return fmt.Errorf("failed to get company name for welcome: %w", err)
	}

	logoURL, err := s.settingsProvider.GetCompanyLogo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get company logo for welcome: %w", err)
	}

	request := domain.WelcomeEmailRequest{
		ToEmail:     email,
		ToFirstName: firstName,
		ToLastName:  lastName,
		FromEmail:   companyEmail,
		CompanyName: companyName,
		LogoURL:     logoURL,
	}

	if err := s.emailClient.SendWelcomeEmail(ctx, request); err != nil {
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	return nil
}

func (s *MailService) SendVerifyEmailEmail(ctx context.Context, email, firstName, lastName string) error {
	companyEmail, err := s.settingsProvider.GetCompanyEmail(ctx)
	if err != nil {
		return fmt.Errorf("failed to get company email for verification: %w", err)
	}

	companyName, err := s.settingsProvider.GetCompanyName(ctx)
	if err != nil {
		return fmt.Errorf("failed to get company name for verification: %w", err)
	}

	logoURL, err := s.settingsProvider.GetCompanyLogo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get company logo for verification: %w", err)
	}

	request := domain.VerifyEmailRequest{
		ToEmail:     email,
		ToFirstName: firstName,
		ToLastName:  lastName,
		FromEmail:   companyEmail,
		CompanyName: companyName,
		LogoURL:     logoURL,
	}

	if err := s.emailClient.SendVerifyEmailEmail(ctx, request); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

func (s *MailService) SendEventNotificationEmail(ctx context.Context, email, firstName, lastName, eventName, eventDetails string) error {
	companyEmail, err := s.settingsProvider.GetCompanyEmail(ctx)
	if err != nil {
		return fmt.Errorf("failed to get company email for event notification: %w", err)
	}

	companyName, err := s.settingsProvider.GetCompanyName(ctx)
	if err != nil {
		return fmt.Errorf("failed to get company name for event notification: %w", err)
	}

	logoURL, err := s.settingsProvider.GetCompanyLogo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get company logo for event notification: %w", err)
	}

	request := domain.EventNotificationRequest{
		ToEmail:     email,
		ToFirstName: firstName,
		ToLastName:  lastName,
		Event:       eventName,
		Details:     eventDetails,
		FromEmail:   companyEmail,
		CompanyName: companyName,
		LogoURL:     logoURL,
	}

	if err := s.emailClient.SendEventNotificationEmail(ctx, request); err != nil {
		return fmt.Errorf("failed to send event notification email: %w", err)
	}

	return nil
}

func (s *MailService) SendPaymentNotificationEmail(ctx context.Context, email, firstName, lastName, amount, product, paymentDate string) error {
	companyEmail, err := s.settingsProvider.GetCompanyEmail(ctx)
	if err != nil {
		return fmt.Errorf("failed to get company email for payment notification: %w", err)
	}

	companyName, err := s.settingsProvider.GetCompanyName(ctx)
	if err != nil {
		return fmt.Errorf("failed to get company name for payment notification: %w", err)
	}

	logoURL, err := s.settingsProvider.GetCompanyLogo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get company logo for payment notification: %w", err)
	}

	request := domain.PaymentNotificationRequest{
		ToEmail:     email,
		ToFirstName: firstName,
		ToLastName:  lastName,
		Amount:      amount,
		Product:     product,
		PaymentDate: paymentDate,
		FromEmail:   companyEmail,
		CompanyName: companyName,
		LogoURL:     logoURL,
	}

	if err := s.emailClient.SendPaymentNotificationEmail(ctx, request); err != nil {
		return fmt.Errorf("failed to send payment notification email: %w", err)
	}

	return nil
}
