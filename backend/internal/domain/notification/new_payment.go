package notification

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/user/models"

	"github.com/hengadev/errsx"
)

func (s *mailService) NewPayment(ctx context.Context, user *models.User, eventTime string) error {
	var errs errsx.Map
	return errs.AsError()
}
