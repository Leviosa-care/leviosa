package messageRepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain/message"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres"
)

func (m *repository) SendMessage(ctx context.Context, message *messageService.Message) error {
	query := fmt.Sprintf(`
        INSERT INTO %s (
            id,
            conversation_id,
            sender_id,
            content_encrypted,
            created_at
        ) VALUES ($1, $2, $3, $4, $5);`, pg.QualifiedTable(m.schema, "messages"))
	result, err := m.DB.ExecContext(
		ctx,
		query,
		message.ID,
		message.ConversationID,
		message.SenderID,
		message.Content,
		message.CreatedAt,
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
		return rp.NewNotCreatedErr(errors.New("no rows affected by insertion statement"), "message")
	}
	return nil
}
