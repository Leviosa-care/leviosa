package application

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/notification/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/notification/ports"
)

// SMSService handles SMS sending logic
type SMSService struct {
	smsClient ports.SMSService
}

func NewSMSService(smsClient ports.SMSService) *SMSService {
	return &SMSService{
		smsClient: smsClient,
	}
}

func (s *SMSService) SendOTPBySMS(ctx context.Context, phoneNumber, otp string) error {
	request := domain.OTPSMSRequest{
		PhoneNumber: phoneNumber,
		OTP:         otp,
	}

	if err := s.smsClient.SendOTP(ctx, request); err != nil {
		return fmt.Errorf("failed to send OTP SMS: %w", err)
	}

	return nil
}

func (s *SMSService) SendGenericSMS(ctx context.Context, phoneNumber, message string) error {
	request := domain.GenericSMSRequest{
		PhoneNumber: phoneNumber,
		Message:     message,
	}

	if err := s.smsClient.SendSMS(ctx, request); err != nil {
		return fmt.Errorf("failed to send generic SMS: %w", err)
	}

	return nil
}
