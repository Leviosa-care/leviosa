package userService

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

// CreateUser creates a verified user in the system based on a pending user response.
//
// Parameters:
//   - ctx: A context.Context instance to manage request lifecycle and cancellation.
//   - userResponse: A pointer to a models.UserPendingResponse instance containing the pending user's response data.
//
// Returns:
//   - *models.User: A pointer to the newly created user object.
//   - error: An error if the pending user cannot be retrieved, the role is invalid, the user cannot be added to the
//     database, or an unexpected error occurs. Returns nil if the user is successfully created.
func (s *service) CreateUser(ctx context.Context, userResponse *models.UserPendingResponse) (*models.User, error) {
	// get encrypted user from hashed email
	user, err := s.repo.GetPendingUser(ctx, userResponse.Email, userResponse.Provider)
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrValidation):
			return nil, domain.NewInvalidValueErr(err.Error())
		case errors.Is(err, rp.ErrNotFound):
			return nil, domain.NewNotFoundErr(err)
		case errors.Is(err, rp.ErrContext):
			return nil, err
		case errors.Is(err, rp.ErrDatabase):
			return nil, domain.NewQueryFailedErr(err)
		default:
			return nil, domain.NewUnexpectTypeErr(err)
		}
	}
	// convert to Role to see if there is no problem with the role sent
	role := models.ConvertToRole(userResponse.Role)
	if role == models.VISITOR {
		return nil, domain.NewInvalidValueErr("invalid role")
	}
	// add remaining field to user
	user.Role = role.String()
	user.Create()
	user.Login()

	// add user to database
	err = s.repo.AddUser(ctx, user, userResponse.Provider)
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrNotCreated):
			return nil, domain.NewNotCreatedErr(err)
		case errors.Is(err, rp.ErrContext):
			return nil, err
		default:
			return nil, domain.NewUnexpectTypeErr(err)
		}
	}
	return user, nil
}
