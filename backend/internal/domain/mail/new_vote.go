package mailService

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/user/models"

	"github.com/hengadev/errsx"
)

// Function that send an email to user after receiving payment.
func (s *service) NewVote(ctx context.Context, user *models.User, eventTime string) error {
	var errs errsx.Map

	templData := struct {
		Username string
		Heure    string
	}{Username: user.FirstName, Heure: eventTime}

	if err := s.sendMail(
		user.Email,
		"[Leviosa] Nouveau votes disponibles",
		"/internal/domain/mail/newRegistry.html",
		templData,
		nil,
		nil,
	); err != nil {
		errs.Set("send mail", err)
	}
	return errs.AsError()
}
