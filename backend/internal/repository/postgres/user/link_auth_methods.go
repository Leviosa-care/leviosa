package userRepository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

// linkAuthMethod links an authentication method (Google, Apple, or email) to a user's account.
// This function updates the user record in the specified table with the corresponding authentication information.
//
// Parameters:
//   - ctx: The context for managing the lifecycle of the database query.
//   - tx: The transaction to execute the query within.
//   - userID: The ID of the user to link the authentication method to.
//   - user: The user struct containing the authentication data (GoogleID, AppleID, or PasswordHash).
//   - provider: The provider type (Google, Apple, or Mail) indicating which authentication method to link.
//   - table: The name of the table (either `users` or `pending_users`) to update.
//
// Returns:
//   - error: Returns an error if any database issues occur or if the provider type is unsupported.
func (u *repository) linkAuthMethod(ctx context.Context, tx *sql.Tx, userID string, user *models.User, provider models.ProviderType, table string) error {
	switch provider {
	case models.Google:
		// Add Apple authentication
		_, err := tx.ExecContext(ctx,
			fmt.Sprintf(`UPDATE %s SET 
                google_id_encrypted = $1
            WHERE id = $2`, table),
			user.GoogleID,
			userID,
		)
		return err
	case models.Apple:
		// Add Apple authentication
		_, err := tx.ExecContext(ctx,
			fmt.Sprintf(`UPDATE %s SET 
                apple_id_encrypted = $1
            WHERE id = $2`, table),
			user.AppleID,
			userID,
		)
		return err
	case models.Mail:
		// Add email authentication
		_, err := tx.ExecContext(ctx,
			fmt.Sprintf(`UPDATE %s SET 
                password_hash = $1
            WHERE id = $2`, table),
			user.PasswordHash,
			userID,
		)
		return err
	default:
		return rp.NewValidationErr(fmt.Errorf("unsupported provider type: %v", provider), "provider")
	}
}
