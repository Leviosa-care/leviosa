package userRepository

import (
	"context"
	"errors"

	rp "github.com/hengadev/leviosa/internal/repository"
)

// HasUser checks whether a user exists in the "users" table based on their email hash.
//
// Parameters:
//   - ctx: The context for managing the lifecycle of the database query.
//   - emailHash: The email hash of the user to check.
//
// Returns:
//   - error: Returns an error if any database issues occur or if the user is not found.
func (u *repository) HasUser(ctx context.Context, emailHash string) error {
	query := `
        SELECT EXISTS (
            SELECT 1 
            FROM users 
            WHERE email_hash = $1
        );`
	var exists bool
	err := u.DB.QueryRowContext(ctx, query, emailHash).Scan(&exists)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return rp.NewContextErr(err)
		default:
			return rp.NewDatabaseErr(err)
		}
	}
	if !exists {
		return rp.NewNotFoundErr(err, "user")
	}
	return nil
}
