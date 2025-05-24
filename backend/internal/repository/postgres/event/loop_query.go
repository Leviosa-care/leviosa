package eventRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	rp "github.com/hengadev/leviosa/internal/repository"
)

func (e *repository) LoopQuery(
	ctx context.Context,
	tx *sql.Tx,
	name string,
	table string,
	values []string,
	eventID string,
) error {

	query := fmt.Sprintf(`INSERT INTO %s (
                event_id,
                %s
            ) VALUES ($1, $2);`, table, name)
	for _, value := range values {
		result, err := tx.ExecContext(ctx, query, table, name, value, eventID)
		// result, err := tx.ExecContext(ctx, query, eventID, value)
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
			return rp.NewNotCreatedErr(errors.New("no rows affected by insertion statement"), fmt.Sprintf("%s in %s table", name, table))
		}
	}
	return nil
}
