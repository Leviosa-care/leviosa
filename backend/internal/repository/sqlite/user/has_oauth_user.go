package userRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

// HasOAuthUser checks whether a user with the specified email hash has an OAuth account (Google, Apple, etc.)
// associated with them in either the "users" or "pending_users" tables.
//
// Parameters:
//   - ctx: The context for managing the lifecycle of the database query.
//   - emailHash: The email hash of the user to check.
//   - p: The OAuth provider type (e.g., Google, Apple, etc.).
//
// Returns:
//   - error: Returns an error if any database issues occur or the user is not found.
func (u *repository) HasOAuthUser(ctx context.Context, emailHash string, p models.ProviderType) error {
	tx, err := u.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return rp.NewDatabaseErr(fmt.Errorf("failed to start transaction: %w", err))
	}
	defer tx.Rollback()
	// this is the string provider here
	provider := string(p)
	tables := []string{"users", "pending_users"}
	for _, table := range tables {
		query := fmt.Sprintf(`
        SELECT EXISTS (
            SELECT 1 
            FROM %s 
            WHERE email_hash = ? 
            AND %s_id_encrypted IS NOT NULL
            AND %s_id_encrypted != ''
        );`, table, provider, provider)
		var exists bool
		err = tx.QueryRowContext(ctx, query, emailHash).Scan(&exists)
		if err != nil {
			switch {
			case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
				return rp.NewContextErr(err)
			default:
				return rp.NewDatabaseErr(err)
			}
		}
		if !exists {
			return rp.NewNotFoundErr(err, "google account with specified email hash")
		}
	}

	if err := tx.Commit(); err != nil {
		return rp.NewDatabaseErr(fmt.Errorf("commit transaction: %w", err))
	}
	return nil
}
