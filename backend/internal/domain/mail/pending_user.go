package mailService

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	"github.com/hengadev/leviosa/pkg/errsx"
)

func (s *service) PendingUser(ctx context.Context, user *models.User) errsx.Map {
	var errs errsx.Map
	return errs
}
