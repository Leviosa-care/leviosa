package userRepository

import (
	"context"
	"errors"
	"fmt"

	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/postgres"
)

// DeleteUser removes a user from the 'users' table in the database based on the provided user ID.
// It executes a DELETE SQL statement and checks if the operation was successful by evaluating
// the number of rows affected. If no rows were affected, it returns an error indicating that the
// user was not deleted.
//
// Parameters:
//   - ctx: The context for managing the transaction lifecycle and cancelation.
//   - userID: The unique identifier of the user to be deleted from the database.
//
// Returns:
//   - error: An error if the deletion fails or if no rows were affected.
//   - If the delete operation encounters context-related errors (e.g., deadline exceeded, canceled),
//     a context error is returned.
//   - If no rows were affected, a "not deleted" error is returned.
//   - If the deletion fails for any other reason, a database error is returned.
func (u *repository) DeleteUser(ctx context.Context, userID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1;", pg.QualifiedTable(u.schema, "users"))
	result, err := u.DB.ExecContext(ctx, query, userID)
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
		return rp.NewNotDeletedErr(err, "user")
	}
	return nil
}
