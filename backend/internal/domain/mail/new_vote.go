package mailService

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	"github.com/hengadev/leviosa/pkg/errsx"
)

// Function that send an email to user after receiving payment.
func (s *service) NewVote(ctx context.Context, user *models.User, eventTime string) errsx.Map {
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
	return errs
}
