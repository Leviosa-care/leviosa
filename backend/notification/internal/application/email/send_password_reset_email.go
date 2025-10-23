package email

import (
	"context"

	"github.com/Leviosa-care/notification/internal/domain"
)

func (s *EmailService) SendPasswordResetEmail(ctx context.Context, email string) error {
	data := domain.NewPasswordResetEmailData(
		email,
		s.cache.GetCompanyLegalAddress(),
		s.cache.GetCompanyInstagram(),
	)

	request := &domain.EmailRequest{
		To:       email,
		Subject:  "Réinitialisation de votre mot de passe Leviosa",
		Template: "verify_email",
		Data:     data,
	}

	return s.emailClient.SendEmail(ctx, request)
}
