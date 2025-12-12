package ports

import (
	"context"

	"github.com/Leviosa-care/notification/internal/domain"
)

type SMSService interface {
	SendSMS(ctx context.Context, request *domain.SMSRequest) error
}
