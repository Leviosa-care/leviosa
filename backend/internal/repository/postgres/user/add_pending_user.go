package userRepository

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/user/models"
)

// AddPendingUser adds a user to the 'pending_users' table by invoking the addGenericUser function.
//
// Parameters:
//   - ctx: The context for managing the transaction lifecycle and cancelation.
//   - user: The user object containing details to be stored in the 'pending_users' table.
//   - provider: The provider type (e.g., 'apple', 'google', 'mail') used for authentication.
//
// Returns:
//   - error: An error if the user addition fails, as returned by the addGenericUser function. Returns nil if successful.
func (u *repository) AddPendingUser(ctx context.Context, user *models.User, provider models.ProviderType) error {
	if err := u.addGenericUser(ctx, user, provider, "pending_users"); err != nil {
		return err
	}
	return nil
}
