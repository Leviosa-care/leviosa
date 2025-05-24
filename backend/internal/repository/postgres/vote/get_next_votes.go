package voteRepository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hengadev/leviosa/internal/domain/vote"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (v *repository) GetNextVotes(ctx context.Context, month, year int) ([]*vote.AvailableVote, error) {
	var votes []*vote.AvailableVote
	query := `SELECT days, month, year 
              FROM available_votes 
              WHERE (year = $1 AND month > $2) OR year = $3 
              ORDER BY year, month 
              LIMIT 3;`
	rows, err := v.DB.QueryContext(ctx, query, year, month, year+1)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return nil, rp.NewContextErr(err)
		default:
			return nil, rp.NewDatabaseErr(err)
		}
	}
	defer rows.Close()
	for rows.Next() {
		var days string
		var monthDB int
		var yearDB int
		err := rows.Scan(
			&days,
			&monthDB,
			&yearDB,
		)
		if err != nil {
			return nil, rp.NewDatabaseErr(err)
		}
		availableDays, err := parseDays(days)
		if err != nil {
			return nil, rp.NewValidationErr(err, "days for available votes")
		}
		for _, day := range availableDays {
			votes = append(votes, &vote.AvailableVote{
				Day:   day,
				Month: monthDB,
				Year:  yearDB,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, rp.NewDatabaseErr(err)
	}
	return votes, nil
}

func parseDays(days string) ([]int, error) {
	var res []int
	daysStr := strings.Split(days, vote.VoteSeparator)
	for _, dayStr := range daysStr {
		day, err := strconv.Atoi(dayStr)
		if err != nil {
			return nil, fmt.Errorf("conversion days string to int")
		}
		res = append(res, day)
	}
	return res, nil
}
