package mailService

import (
	"context"

	"github.com/hengadev/errsx"
	"github.com/hengadev/leviosa/internal/domain/user/models"
)

func (s *service) WelcomeUser(ctx context.Context, user *models.User) error {
	var errs errsx.Map
	return errs.AsError()
}
