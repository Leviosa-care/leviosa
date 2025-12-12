package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/notification/domain"
)

// SMSService defines the interface for SMS sending implementations (Twilio)
type SMSService interface {
	SendOTP(ctx context.Context, req domain.OTPSMSRequest) error
	SendSMS(ctx context.Context, req domain.GenericSMSRequest) error
}
