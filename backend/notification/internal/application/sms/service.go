package sms

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/notification/internal/domain"
	"github.com/Leviosa-care/notification/internal/ports"
)

type SMSService struct {
	smsClient ports.SMSService
}

func New(smsClient ports.SMSService) *SMSService {
	return &SMSService{
		smsClient: smsClient,
	}
}

func (s *SMSService) SendSMS(ctx context.Context, phone, message string) error {
	request, err := domain.NewSMSRequest(phone, message)
	if err != nil {
		return errs.NewInvalidValueErr(fmt.Errorf("create SMS request: %w", err))
	}

	return s.smsClient.SendSMS(ctx, request)
}

