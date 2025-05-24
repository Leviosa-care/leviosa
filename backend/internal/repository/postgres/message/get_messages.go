package messageRepository

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain/message"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (m *repository) GetMessages(ctx context.Context, conversationID string) ([]*messageService.Message, error) {
	query := `
        SELECT 
            id,
            conversation_id,
            sender_id,
            content_encrypted,
            created_at
        FROM messages;`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return nil, rp.NewContextErr(err)
		default:
			return nil, rp.NewDatabaseErr(err)
		}
	}
	defer rows.Close()

	var messages []*messageService.Message
	for rows.Next() {
		var message messageService.Message
		err := rows.Scan(
			&message.ID,
			&message.ConversationID,
			&message.SenderID,
			&message.Content,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, rp.NewDatabaseErr(err)
		}
		messages = append(messages, &message)
	}
	if err := rows.Err(); err != nil {
		return nil, rp.NewDatabaseErr(err)
	}
	if len(messages) == 0 {
		return []*messageService.Message{}, nil
	}
	return messages, nil
}
