package email

import (
	"context"

	"github.com/Leviosa-care/notification/internal/domain"
)

func (s *EmailService) SendWelcomeEmail(ctx context.Context, email string) error {
	data := domain.NewWelcomeEmailData(
		email,
		s.cache.GetCompanyLegalAddress(),
		s.cache.GetCompanyInstagram(),
	)

	request := &domain.EmailRequest{
		To:       email,
		Subject:  "Bienvenue chez Leviosa",
		Template: "welcome",
		Data:     data,
	}

	return s.emailClient.SendEmail(ctx, request)
}
