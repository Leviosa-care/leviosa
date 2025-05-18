package userRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

// addGenericUser checks if a user exists in the specified database table and either creates a new user or links an authentication method to an existing user, all within a transaction.
//
// Parameters:
//   - ctx: The context for managing the transaction lifecycle and cancelation.
//   - user: The user object containing details to be stored or linked.
//   - provider: The provider type (e.g., 'apple', 'google', 'mail') used for authentication.
//   - table: The name of the database table to check and store the user data.
//
// Returns:
//   - error: An error if the transaction fails, the user creation or linking process fails, or any database-related errors occur. Returns nil if successful.
func (u *repository) addGenericUser(ctx context.Context, user *models.User, provider models.ProviderType, table string) error {
	tx, err := u.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return rp.NewDatabaseErr(fmt.Errorf("failed to start transaction: %w", err))
	}
	defer tx.Rollback()

	if !provider.IsValid() {
		return rp.NewValidationErr(fmt.Errorf("provider type value can only be 'apple', 'google' or 'mail', got : %q", provider), "provider type")
	}

	// Check if user exists
	var userID string
	err = tx.QueryRowContext(ctx,
		fmt.Sprintf("SELECT id FROM %s WHERE email_hash = ?", table),
		user.EmailHash,
	).Scan(&userID)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		// Create new user
		err = u.createNewUser(ctx, tx, user, provider, table)
		if err != nil {
			return rp.NewNotCreatedErr(err, "new user")
		}
	} else {
		// Link authentication method to existing user
		if err := u.linkAuthMethod(ctx, tx, userID, user, provider, table); err != nil {
			return rp.NewNotCreatedErr(err, "linked user account")
		}
	}

	if err := tx.Commit(); err != nil {
		return rp.NewDatabaseErr(fmt.Errorf("commit transaction: %w", err))
	}
	return nil
}
