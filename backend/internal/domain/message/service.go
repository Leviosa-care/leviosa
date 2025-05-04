package messageService

import (
	"context"

	"github.com/hengadev/encx"
)

type Service interface {
	CreateConversation(ctx context.Context, userID, adminID string) (string, error)
	GetMessages(ctx context.Context, conversationID string) ([]*Message, error)
	ListConversations(ctx context.Context, userID string) ([]*Conversation, error)
	SendMessage(ctx context.Context, conversationID, senderID, content string) error
}

type service struct {
	repo   ReadWriter
	crypto *encx.Crypto
}

func New(repo ReadWriter, crypto *encx.Crypto) Service {
	return &service{
		repo,
		crypto,
	}
}
