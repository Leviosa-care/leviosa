package messageService

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) SendMessage(ctx context.Context, conversationID, senderID, content string) error {
	dek, err := s.crypto.GenerateDEK()
	if err != nil {
		return fmt.Errorf("generate DEK for message encryption: %w", err)
	}
	message := &Message{
		ID:             uuid.NewString(),
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		CreatedAt:      time.Now(),
		DEK:            dek,
	}
	if err := s.crypto.ProcessStruct(ctx, message); err != nil {
		return domain.NewNotEncryptedErr("message", err)
	}
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
