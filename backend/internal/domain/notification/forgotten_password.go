package notification

import (
	"context"
)

func (s *mailService) HandlePasswordForgotten(ctx context.Context, to string) error {
	// send an email to the user and when redirected to that link, give the user an opportunity to remake the password.
	return nil
}
