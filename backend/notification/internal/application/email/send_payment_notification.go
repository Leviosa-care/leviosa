package email

import (
	"context"

	"github.com/Leviosa-care/internal/domain/user/models"
	"github.com/Leviosa-care/notification/internal/domain"
)

func (s *EmailService) SendPaymentNotification(ctx context.Context, user *models.User, eventTime string) error {
	data := domain.NewPaymentNotificationEmailData(
		eventTime,
		s.cache.GetCompanyLegalAddress(),
		s.cache.GetCompanyInstagram(),
	)

	request := &domain.EmailRequest{
		To:       user.Email,
		Subject:  "Confirmation de paiement Leviosa",
		Template: "payment",
		Data:     data,
	}

	return s.emailClient.SendEmail(ctx, request)
}

