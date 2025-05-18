package eventRepository

import (
	"context"
	"errors"

	rp "github.com/hengadev/leviosa/internal/repository"
)

func (e *repository) IsDateAvailable(ctx context.Context, day, month, year int) error {
	query := `
	SELECT EXISTS (
		SELECT 1
		FROM events
		WHERE day = ? 
		AND month = ?
		AND year = ?
	);`
	var exists bool
	err := e.DB.QueryRowContext(ctx, query, day, month, year).Scan(&exists)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return rp.NewContextErr(err)
		default:
			return rp.NewDatabaseErr(err)
		}
	}
	if exists {
		return rp.NewValidationErr(nil, "event already exists in database this date")
	}
	return nil
}
