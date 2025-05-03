package userService

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

// CreateOAuthPendingUser adds a user to the pending user table for OAuth-based registration.
//
// Parameters:
//   - ctx: A context.Context instance to manage request lifecycle and cancellation.
//   - user: A pointer to a models.User instance containing the user's information.
//   - provider: A models.ProviderType value representing the OAuth provider used for registration.
//
// Returns:
//   - error: An error if the user encryption fails, the user cannot be added to the pending user table,
//     or an unexpected error occurs. Returns nil if the user is successfully added.
func (s *service) CreateOAuthPendingUser(ctx context.Context, user *models.User, provider models.ProviderType) error {
	// encrypt user
	if errs := s.EncryptUser(user); len(errs) > 0 {
		return domain.NewNotEncryptedErr("OAuth pending user", errs)
	}
	// add user to pending_user table
	if err := s.repo.AddPendingUser(ctx, user, provider); err != nil {
		switch {
		case errors.Is(err, rp.ErrContext):
			return err
		case errors.Is(err, rp.ErrDatabase):
			return domain.NewQueryFailedErr(err)
		case errors.Is(err, rp.ErrNotCreated):
			return domain.NewNotCreatedErr(err)
		default:
			return domain.NewUnexpectTypeErr(err)
		}
	}
	return nil
}
