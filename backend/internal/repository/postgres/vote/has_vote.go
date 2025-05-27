package voteRepository

import (
	"context"
	"errors"
	"fmt"

	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres"
)

func (v *repository) HasVote(ctx context.Context, month, year int, userID string) error {
	var exists bool
	query := fmt.Sprintf(`
        SELECT EXISTS (
            SELECT 1
            FROM %s
            WHERE user_id = $1 AND month = $2 AND year = $3
        )`, pg.QualifiedTable(v.schema, "votes"))
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
