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
		WHERE day = $1 
		AND month = $2
		AND year = $3
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
