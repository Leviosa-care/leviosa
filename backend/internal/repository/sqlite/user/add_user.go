package userRepository

import (
	"context"

	"github.com/hengadev/leviosa/internal/domain/user/models"
)

// AddUser inserts a new user into the 'users' table, either creating a new user or linking an authentication method
// to an existing user depending on the provider and user existence.
//
// Parameters:
//   - ctx: The context for managing the transaction lifecycle and cancelation.
//   - user: The user object containing details to be stored in the 'users' table.
//     This includes email hash, password hash, and personal details.
//   - provider: The provider type (e.g., 'apple', 'google', or 'mail') used for user authentication.
//
// Returns:
//   - error: An error if the user creation or linking fails, including database-related errors. Returns nil if successful.
//   - If the provider is invalid, a validation error is returned.
//   - If the user already exists, the authentication method is linked to the existing user.
func (u *repository) AddUser(ctx context.Context, user *models.User, provider models.ProviderType) error {
	if err := u.addGenericUser(ctx, user, provider, "users"); err != nil {
		return err
	}
	return nil
}
