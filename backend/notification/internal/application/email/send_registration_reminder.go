package email

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/internal/domain/user/models"
	"github.com/Leviosa-care/notification/internal/domain"
)

func (s *EmailService) SendRegistrationReminder(ctx context.Context, user *models.User, registrationName string, daysLeft int) error {
	data := domain.NewRegistrationReminderEmailData(
		registrationName,
		daysLeft,
		s.cache.GetCompanyLegalAddress(),
		s.cache.GetCompanyInstagram(),
	)

	request := &domain.EmailRequest{
		To:       user.Email,
		Subject:  fmt.Sprintf("Rappel d'inscription - %s", registrationName),
		Template: "registration_reminder",
		Data:     data,
	}

	return s.emailClient.SendEmail(ctx, request)
}
