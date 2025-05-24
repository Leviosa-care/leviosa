package voteRepository

import (
	"context"
	"errors"

	rp "github.com/hengadev/leviosa/internal/repository"
)

func (v *repository) HasVote(ctx context.Context, month, year int, userID string) error {
	var exists bool
	query := `
        SELECT EXISTS (
            SELECT 1 
            FROM votes 
            WHERE user_id = $1 AND month = $2 AND year = $3
        )`
	err := v.DB.QueryRowContext(ctx, query, userID, month, year).Scan(&exists)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return rp.NewContextErr(err)
		default:
			return rp.NewDatabaseErr(err)
		}
	}
	if !exists {
		return rp.NewNotFoundErr(errors.New("no vote found for the user with specified ID"), "votes")
	}
	return nil
}
