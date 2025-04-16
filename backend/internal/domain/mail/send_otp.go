package mailService

import (
	"context"

	otpService "github.com/hengadev/leviosa/internal/domain/otp"
	"github.com/hengadev/leviosa/pkg/errsx"
)

// TODO: make the right email template for that mail domain service
func (s *Service) SendOTP(ctx context.Context, email, firstname string, otp *otpService.OTP) errsx.Map {
	var errs errsx.Map

	// data used in the email
	templData := struct {
		Firstname string
		Value     string
	}{
		Firstname: firstname,
		Value:     otp.Code,
	}
	if err := s.sendMail(
		email,
		"[Leviosa] Confirmation d'addresse email",
		"/internal/domain/mail/templates/otp.html",
		templData,
		nil,
		nil,
	); err != nil {
		errs.Set("send email:", err)
	}
	return errs
}
