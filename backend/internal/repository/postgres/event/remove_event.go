package eventRepository

import (
	"context"
	"errors"
	"fmt"

	rp "github.com/hengadev/leviosa/internal/repository"
)

func (e *repository) RemoveEvent(ctx context.Context, eventID string) error {
	query := "DELETE FROM events WHERE id=$1"
	result, err := e.DB.ExecContext(ctx, query, eventID)
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
		return rp.NewNotDeletedErr(err, fmt.Sprintf("event with ID %s", eventID))
	}

	return nil
}
