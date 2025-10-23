package email

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/internal/domain/user/models"
	"github.com/Leviosa-care/notification/internal/domain"
)

func (s *EmailService) SendEventNotification(ctx context.Context, users []*models.User, eventTime string) error {
	if len(users) == 0 {
		return nil
	}

	data := domain.NewEventNotificationEmailData(
		eventTime,
		s.cache.GetCompanyLegalAddress(),
		s.cache.GetCompanyInstagram(),
	)

	var lastErr error
	for _, user := range users {
		request := &domain.EmailRequest{
			To:       user.Email,
			Subject:  "Nouvel événement Leviosa",
			Template: "event_notification",
			Data:     data,
		}

		if err := s.emailClient.SendEmail(ctx, request); err != nil {
			lastErr = fmt.Errorf("send event notification to %s: %w", user.Email, err)
		}
	}

	return lastErr
}
