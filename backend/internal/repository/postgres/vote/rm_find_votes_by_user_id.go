package voteRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	rp "github.com/hengadev/leviosa/internal/repository"
)

func (v *repository) FindVotesByUserID(ctx context.Context, month string, year int, userID string) (string, error) {
	var votes string
	tableName := fmt.Sprintf("votes_%s_%d", month, year)
	query := fmt.Sprintf("SELECT * FROM %s WHERE userid=$1;", tableName)
	if err := v.DB.QueryRowContext(ctx, query).Scan(&votes); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", rp.NewNotFoundErr(err, "votes")
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return "", rp.NewContextErr(err)
		default:
			return "", rp.NewDatabaseErr(err)
		}
	}
	return votes, nil
}
