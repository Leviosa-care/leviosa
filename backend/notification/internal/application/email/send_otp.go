package email

import (
	"context"

	"github.com/Leviosa-care/notification/internal/domain"
)

func (s *EmailService) SendOTP(ctx context.Context, email string, otp string) error {
	data := domain.NewOTPEmailData(
		otp,
		s.cache.GetCompanyLegalAddress(),
		s.cache.GetCompanyInstagram(),
	)

	request := &domain.EmailRequest{
		To:       email,
		Subject:  "Votre code de vérification Leviosa",
		Template: "otp",
		Data:     data,
	}

	return s.emailClient.SendEmail(ctx, request)
}
