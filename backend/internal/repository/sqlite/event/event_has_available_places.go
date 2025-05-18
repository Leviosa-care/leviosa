package eventRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	rp "github.com/hengadev/leviosa/internal/repository"
)

// EventHasAvailablePlaces checks if an event with the given eventID exists in the database
// and has a place count greater than 0.
func (e *repository) EventHasAvailablePlaces(ctx context.Context, eventID string) (bool, error) {
	var placecount int
	query := `
        SELECT
            placecount
        FROM events 
        WHERE id = ?;`

	if err := e.DB.QueryRowContext(ctx, query, eventID).Scan(&placecount); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return false, rp.NewNotFoundErr(err, fmt.Sprintf("event with ID %s", eventID))
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return false, rp.NewContextErr(err)
		default:
			return false, rp.NewDatabaseErr(err)
		}
	}
	return placecount > 0, nil
}
