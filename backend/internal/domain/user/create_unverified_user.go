package userService

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

// CreateUnverifiedUser creates an unverified user entry in the system from the provided signup data.
//
// Parameters:
//   - ctx: A context.Context instance to manage request lifecycle and cancellation.
//   - userSignUp: A pointer to a models.UserSignUp instance containing the user's signup information.
//
// Returns:
//   - string: The hashed email of the newly created unverified user.
//   - error: An error if the user encryption fails, the user cannot be added to the unverified user table,
//     or an unexpected error occurs. Returns nil if the user is successfully added.
func (s *service) CreateUnverifiedUser(ctx context.Context, userSignUp *models.UserSignUp) (string, error) {
	user := userSignUp.ToUser()
	if errs := s.EncryptUser(user); len(errs) > 0 {
		return "", domain.NewNotEncryptedErr("unverified user", errs)
	}
	err := s.repo.AddUnverifiedUser(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrDatabase):
			return "", domain.NewQueryFailedErr(err)
		case errors.Is(err, rp.ErrNotCreated):
			return "", domain.NewNotCreatedErr(err)
		case errors.Is(err, rp.ErrContext):
			return "", err
		default:
			return "", domain.NewUnexpectTypeErr(err)
		}
	}
	return user.EmailHash, nil
}
