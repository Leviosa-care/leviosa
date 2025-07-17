package notification

import (
	"context"
	"time"
)

func (s *mailService) SendOTP(ctx context.Context, email string, otp string) error {
	legalAddress := s.cache.getCompanyLegalAddress()
	companyInstagram := s.cache.getCompanyInstagram()
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
	// TODO: just make a get request to the public api to get the logo right ?
	// Move that part to the send mail function, this is where it belongs
	// TODO: use the helper function for that in the sendMail thing
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
