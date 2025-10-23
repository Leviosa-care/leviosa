package ports

import (
	"context"

	"github.com/Leviosa-care/notification/internal/domain"
)

type EmailService interface {
	SendEmail(ctx context.Context, request *domain.EmailRequest) error
}
