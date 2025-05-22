package userRepository

import (
	"context"
	"database/sql"
	"errors"

	rp "github.com/hengadev/leviosa/internal/repository"
)

// GetHashedPasswordByEmail retrieves the hashed password of a user by their email address from the 'users' table.
//
// Parameters:
//   - ctx: The context for managing the transaction lifecycle and cancelation.
//   - email: The email address of the user for which the hashed password is to be retrieved.
//
// Returns:
//   - string: The hashed password of the user.
//   - error: An error if the query fails or the user is not found.
//   - Returns a "not found" error if no user exists with the given email.
//   - Returns a context error if the operation is canceled or the deadline is exceeded.
//   - Returns a database error for any other query-related issues.
func (u *repository) GetHashedPasswordByEmail(ctx context.Context, email string) (string, error) {
	var hashedPassword string
	query := "SELECT password_hash from users where email_hash = $1;"
	err := u.DB.QueryRowContext(ctx, query, email).Scan(&hashedPassword)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", rp.NewNotFoundErr(err, "user")
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return "", rp.NewContextErr(err)
		default:
			return "", rp.NewDatabaseErr(err)
		}
	}
	return hashedPassword, nil
}
