package messageService

import (
	"context"
)

type Reader interface {
	GetMessages(ctx context.Context, conversationID string) ([]*Message, error)
	ListConversations(ctx context.Context, userID string) ([]*Conversation, error)
}

type Writer interface {
	CreateConversation(ctx context.Context, conversation *Conversation) error
	SendMessage(ctx context.Context, message *Message) error
}

type ReadWriter interface {
	Reader
	Writer
}
