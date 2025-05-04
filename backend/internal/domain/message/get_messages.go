package messageService

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) GetMessages(ctx context.Context, conversationID string) ([]*Message, error) {
	messages, err := s.repo.GetMessages(ctx, conversationID)
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrNotFound):
			return nil, domain.NewNotFoundErr(err)
		case errors.Is(err, rp.ErrContext):
			return nil, err
		case errors.Is(err, rp.ErrDatabase):
			return nil, domain.NewQueryFailedErr(err)
		}
	}
	for _, message := range messages {
		if err := s.crypto.DecryptStruct(ctx, message); err != nil {
			return nil, domain.NewNotEncryptedErr("message content decryption", err)
		}
	}
	return messages, nil
}
