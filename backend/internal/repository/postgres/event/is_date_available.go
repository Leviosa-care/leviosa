package eventRepository

import (
	"context"
	"errors"
	"fmt"

	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres"
)

func (e *repository) IsDateAvailable(ctx context.Context, day, month, year int) error {
	query := fmt.Sprintf(`
	SELECT EXISTS (
		SELECT 1
		FROM %s
		WHERE day = $1 
		AND month = $2
		AND year = $3
	);`, pg.QualifiedTable(e.schema, "events"))
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
