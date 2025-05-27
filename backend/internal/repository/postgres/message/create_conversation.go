package messageRepository

import (
	"context"
	"errors"

	"fmt"

	"github.com/hengadev/leviosa/internal/domain/message"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres"
)

func (m *repository) CreateConversation(ctx context.Context, conversation *messageService.Conversation) error {
	query := fmt.Sprintf(`
        INSERT INTO %s (
            id,
            user_id,
            partner_id,
            created_at
        ) VALUES ($1, $2, $3, $4);`, pg.QualifiedTable(m.schema, "conversations"))
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
