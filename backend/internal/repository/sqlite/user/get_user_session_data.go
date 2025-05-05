package userRepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

// GetUserSessionData retrieves the user's session-related data (ID and role) by their email.
//
// Parameters:
//   - ctx: Context for managing the lifecycle of the database query.
//   - email: The user's email to look up.
//
// Returns:
//   - string: The user ID.
//   - models.Role: The user's role (e.g., ADMIN, USER).
//   - error: An error if the query fails or the user is not found.
//   - Returns a "not found" error if the user does not exist in the database.
//   - Returns a context error if the operation is canceled or the deadline is exceeded.
//   - Returns a database error for any other query-related issues.
func (u *Repository) GetUserSessionData(ctx context.Context, emailHash string) (string, models.Role, error) {
	var id, role string
	query := "SELECT id, role from users where email_hash = ?;"
	err := u.DB.QueryRowContext(ctx, query, emailHash).Scan(
		&id,
		&role,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", models.VISITOR, rp.NewNotFoundErr(err, "user session data")
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return "", models.VISITOR, rp.NewContextErr(err)
		default:
			return "", models.VISITOR, rp.NewDatabaseErr(err)
		}
	}
	return id, models.ConvertToRole(role), nil
}
