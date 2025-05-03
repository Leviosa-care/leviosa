package messageService

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/message/models"
	"github.com/hengadev/leviosa/internal/domain/message/security"
	"github.com/hengadev/leviosa/pkg/config"
)

type Service interface {
	CreateConversation(ctx context.Context, userID, adminID string) (string, error)
	GetMessages(ctx context.Context, conversationID string) ([]*models.Message, error)
	ListConversations(ctx context.Context, userID string) ([]*models.Conversation, error)
	SendMessage(ctx context.Context, conversationID, senderID, content string) error
}

type service struct {
	repo ReadWriter
	*security.SecureMessageData
}

func New(repo ReadWriter, conf *config.SecurityConfig) Service {
	return &service{
		repo,
		security.NewSecureMessageData(conf)}
}
