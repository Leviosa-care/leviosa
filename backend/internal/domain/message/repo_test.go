package messageService_test

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/message"
)

type MockRepo struct {
	GetMessagesFunc        func(ctx context.Context, conversationID string) ([]*messageService.Message, error)
	ListConversationsFunc  func(ctx context.Context, userID string) ([]*messageService.Conversation, error)
	CreateConversationFunc func(ctx context.Context, conversation *messageService.Conversation) error
	SendMessageFunc        func(ctx context.Context, message *messageService.Message) error
}

func (m *MockRepo) GetMessages(ctx context.Context, conversationID string) ([]*messageService.Message, error) {
	if m.GetMessagesFunc != nil {
		return m.GetMessagesFunc(ctx, conversationID)
	}
	return nil, nil
}

func (m *MockRepo) ListConversations(ctx context.Context, userID string) ([]*messageService.Conversation, error) {
	if m.ListConversationsFunc != nil {
		return m.ListConversationsFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockRepo) CreateConversation(ctx context.Context, conversation *messageService.Conversation) error {
	if m.CreateConversationFunc != nil {
		return m.CreateConversationFunc(ctx, conversation)
	}
	return nil
}

func (m *MockRepo) SendMessage(ctx context.Context, message *messageService.Message) error {
	if m.SendMessageFunc != nil {
		return m.SendMessageFunc(ctx, message)
	}
	return nil
}
