package messageService

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/message/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) SendMessage(ctx context.Context, conversationID, senderID, content string) error {
	// I need to encrypt the message when I get it
	message := &models.Message{}
	s.EncryptMessage(message)
	if err := s.repo.SendMessage(ctx, message); err != nil {
		switch {
		case errors.Is(err, rp.ErrNotCreated):
			return domain.NewNotCreatedErr(err)
		case errors.Is(err, rp.ErrContext):
			return err
		case errors.Is(err, rp.ErrDatabase):
			return domain.NewQueryFailedErr(err)
		}
	}
	return nil
}
