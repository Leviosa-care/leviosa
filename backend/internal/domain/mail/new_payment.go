package mailService

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	"github.com/hengadev/leviosa/pkg/errsx"
)

func (s *Service) NewPayment(ctx context.Context, user *models.User, eventTime string) errsx.Map {
	var errs errsx.Map
	return errs
}
