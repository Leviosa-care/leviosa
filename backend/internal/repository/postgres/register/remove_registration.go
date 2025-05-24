package registerRepository

import (
	"context"
	"errors"
	"fmt"

	rp "github.com/hengadev/leviosa/internal/repository"
)

func (r *repository) RemoveRegistration(ctx context.Context, day, year int, month string) error {
	tablename := getTablename(day, year, month)
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=?", tablename)
	result, err := r.DB.ExecContext(ctx, query)
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
		// TODO: add the user ID once it is in the function signature
		return rp.NewNotDeletedErr(err, fmt.Sprintf("registration for %s for user", tablename))
	}
	return nil
}
