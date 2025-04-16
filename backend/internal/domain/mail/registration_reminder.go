package mailService

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	"github.com/hengadev/leviosa/pkg/errsx"
)

// Function that send an email to user to remind them of an event incoming.
func (s *Service) SendRegistrationReminderEmail(ctx context.Context, user *models.User, registrationName string, daysLeft int) errsx.Map {
	// TODO: Add the call to a certain function to handle using that value
	var errs errsx.Map
	switch daysLeft {
	case 2:
	case 7:
	}
	return errs
}
