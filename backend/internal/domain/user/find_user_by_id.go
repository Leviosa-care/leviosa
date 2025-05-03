package userService

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

// FindUserByID retrieves a user account by its ID.
//
// Parameters:
//   - ctx: A context.Context instance to manage request lifecycle and cancellation.
//   - userID: A string representing the ID of the user whose account is being retrieved.
//
// Returns:
//   - *models.User: A pointer to the user account retrieved. If the user is not found, an empty user object is returned.
//   - error: An error if the user cannot be retrieved, the user data cannot be decrypted, or an unexpected error occurs.
//     Returns nil if the user is successfully retrieved and decrypted.
func (s *service) FindUserByID(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.repo.FindAccountByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrNotFound):
			return &models.User{}, domain.NewNotFoundErr(err)
		case errors.Is(err, rp.ErrContext):
			return nil, err
		case errors.Is(err, rp.ErrDatabase):
			return nil, domain.NewQueryFailedErr(err)
		default:
			return nil, domain.NewUnexpectTypeErr(err)
		}
	}
	if errs := s.DecryptUser(user); len(errs) > 0 {
		return nil, domain.NewInvalidValueErr(errs.Error())
	}
	return user, nil
}
