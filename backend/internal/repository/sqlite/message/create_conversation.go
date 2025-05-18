package messageRepository

import (
	"context"
	"errors"

	messageService "github.com/hengadev/leviosa/internal/domain/message"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (m *repository) CreateConversation(ctx context.Context, conversation *messageService.Conversation) error {
	query := `
        INSERT INTO conversations (
            id,
            user_id,
            partner_id,
            created_at
        ) VALUES (?, ?, ?, ?);`
	result, err := m.DB.ExecContext(
		ctx,
		query,
		conversation.ID,
		conversation.UserID,
		conversation.PartnerID,
		conversation.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return rp.NewContextErr(err)
		default:
			return rp.NewDatabaseErr(err)
		}

	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return rp.NewDatabaseErr(err)
	}
	if rowsAffected == 0 {
		return rp.NewNotCreatedErr(errors.New("no rows affected by insertion statement"), "conversation")
	}
	return nil
}
