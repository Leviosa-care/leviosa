package userService

import (
	"context"
	"errors"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

// GetAllPendingUsers retrieves all pending users from the database.
//
// Parameters:
//   - ctx: A context.Context instance to manage request lifecycle and cancellation.
//
// Returns:
//   - []*models.UserPending: A slice of pointers to the pending users retrieved. If no pending users are found,
//     an empty slice is returned.
//   - error: An error if the users cannot be retrieved, the users cannot be decrypted, or an unexpected error occurs.
//     Returns nil if the users are successfully retrieved and decrypted.
func (s *service) GetAllPendingUsers(ctx context.Context) ([]*models.UserPending, error) {
	users, err := s.repo.GetPendingUsers(ctx)
	var pendingUsers []*models.UserPending
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrContext):
			return nil, err
		case errors.Is(err, rp.ErrDatabase):
			return nil, domain.NewQueryFailedErr(err)
		case errors.Is(err, rp.ErrNotFound):
			return pendingUsers, nil
		default:
			return nil, domain.NewUnexpectTypeErr(err)
		}
	}
	for _, user := range users {
		if err := s.crypto.DecryptStruct(ctx, user); err != nil {
			return nil, domain.NewInvalidValueErr(err.Error())
		}
		pendingUsers = append(pendingUsers, user.ToUserPending())
	}
	return pendingUsers, nil
}
