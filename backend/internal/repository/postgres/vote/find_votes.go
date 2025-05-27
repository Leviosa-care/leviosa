package voteRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres"
)

func (v *repository) FindVotes(ctx context.Context, month, year int, userID string) (string, error) {
	var days string
	query := fmt.Sprintf("SELECT days FROM %s WHERE user_id=$1 and month=$2 and year=$3;", pg.QualifiedTable(v.schema, "votes"))
	err := v.DB.QueryRowContext(ctx, query, userID, month, year).Scan(&days)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", rp.NewNotFoundErr(err, "vote")
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return "", rp.NewContextErr(err)
		default:
			return "", rp.NewDatabaseErr(err)
		}
	}
	return days, nil
}
