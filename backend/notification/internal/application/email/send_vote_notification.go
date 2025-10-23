package email

import (
	"context"

	"github.com/Leviosa-care/internal/domain/user/models"
	"github.com/Leviosa-care/notification/internal/domain"
)

func (s *EmailService) SendVoteNotification(ctx context.Context, user *models.User, eventTime string) error {
	data := domain.NewVoteNotificationEmailData(
		eventTime,
		s.cache.GetCompanyLegalAddress(),
		s.cache.GetCompanyInstagram(),
	)

	request := &domain.EmailRequest{
		To:       user.Email,
		Subject:  "Nouveau vote Leviosa",
		Template: "new_vote",
		Data:     data,
	}

	return s.emailClient.SendEmail(ctx, request)
}
