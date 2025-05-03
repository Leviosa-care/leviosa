package mailService

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/user/models"

	"github.com/hengadev/errsx"
)

func (s *service) NewPayment(ctx context.Context, user *models.User, eventTime string) error {
	var errs errsx.Map
	return errs.AsError()
}
