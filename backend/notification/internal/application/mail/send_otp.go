package mail

import (
	"context"
	"time"
)

func (s *MailService) SendOTP(ctx context.Context, email string, otp string) error {
	legalAddress := s.cache.GetCompanyLegalAddress()
	companyInstagram := s.cache.GetCompanyInstagram()
	data := struct {
		OTP           string
		Year          int
		Address       string
		InstagramPath string
	}{
		OTP:           otp,
		Year:          time.Now().Year(),
		Address:       legalAddress,
		InstagramPath: companyInstagram,
	}
	if err := s.sendMail(
		ctx,
		email,
		"Votre code de vérification Leviosa",
		"otp",
		data,
		nil,
		nil,
	); err != nil {
		return err
	}
	return nil
}
