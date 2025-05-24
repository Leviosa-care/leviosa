package voteRepository

import (
	"context"
	"errors"
	"fmt"

	rp "github.com/hengadev/leviosa/internal/repository"
)

func (v *repository) RemoveVote(ctx context.Context, userID string, month, year int) error {
	query := "DELETE FROM votes WHERE user_id=$1 AND month=$2 AND year=$3;"
	res, err := v.DB.ExecContext(ctx, query, userID, month, year)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return rp.NewContextErr(err)
		default:
			return rp.NewDatabaseErr(err)
		}
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return rp.NewDatabaseErr(err)
	}
	if rowsAffected == 0 {
		return rp.NewNotDeletedErr(fmt.Errorf("no row deleted"), fmt.Sprintf("vote for user with ID %s", userID))
	}
	return nil
}
