package voteRepository

import (
	"context"
	"errors"
	"fmt"

	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres"
)

// CreateVote inserts a new vote record into the 'votes' table for a specific user and date.
//
// Parameters:
//   - ctx: The context for managing the transaction lifecycle and cancellation.
//   - userID: The unique identifier of the user casting the vote.
//   - days: A string representing the days the vote applies to (format to be specified).
//   - month: An integer representing the month of the vote.
//   - year: An integer representing the year of the vote.
//
// Returns:
//   - error: An error if the vote creation fails, including database-related errors. Returns nil if successful.
//   - If there is a database execution error, a database-specific error is returned.
//   - If no rows are affected, a not-created error is returned, indicating the vote was not inserted.
func (v *repository) CreateVote(ctx context.Context, userID string, days string, month, year int) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, days, month, year) VALUES ($1, $2, $3, $4);", pg.QualifiedTable(v.schema, "votes"))
	result, err := v.DB.ExecContext(ctx, query, userID, days, month, year)
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
		return rp.NewNotCreatedErr(errors.New("no rows affected by insertion statement"), "vote")
	}
	return nil
}
