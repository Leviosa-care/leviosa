package messageRepository

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain/message"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (m *repository) ListConversations(ctx context.Context, userID string) ([]*messageService.Conversation, error) {
	query := `
        SELECT 
            id,
            user_id,
            partner_id,
            created_at
        FROM conversations;`
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
	var conversations []*messageService.Conversation
	for rows.Next() {
		var conversation messageService.Conversation
		err := rows.Scan(
			&conversation.ID,
			&conversation.UserID,
			&conversation.PartnerID,
			&conversation.CreatedAt,
		)

		if err != nil {
			return nil, rp.NewDatabaseErr(err)
		}
		conversations = append(conversations, &conversation)
	}
	if err := rows.Err(); err != nil {
		return nil, rp.NewDatabaseErr(err)
	}
	if len(conversations) == 0 {
		return []*messageService.Conversation{}, nil
	}
	return conversations, nil
}
