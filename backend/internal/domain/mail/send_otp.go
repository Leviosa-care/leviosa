package mailService

import (
	"context"
	"time"

	otpService "github.com/hengadev/leviosa/internal/domain/otp"
)

// these things are in all mails
type CompanyInfo struct {
	LegalAddress  string
	InstagramPath string
}

func (s *service) SendOTP(ctx context.Context, email string, otp *otpService.OTP) error {
	legalAddress, err := s.repo.GetCompanyLegalAddress(ctx)
	if err != nil {

	}
	companyInstagram, err := s.repo.GetCompanyInstagram(ctx)
	if err != nil {

	}
	data := struct {
		OTP           string
		Year          int
		Address       string
		InstagramPath string
	}{
		OTP:           otp.Code,
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
