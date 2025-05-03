package userService

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"
)

// DeleteUser deletes a user from the database based on their user ID.
//
// Parameters:
//   - ctx: A context.Context instance to manage request lifecycle and cancellation.
//   - userID: A string representing the ID of the user to be deleted.
//
// Returns:
//   - error: An error if the user ID is invalid, the user cannot be deleted, or an unexpected error occurs.
//     Returns nil if the user is successfully deleted.
func (s *service) DeleteUser(ctx context.Context, userID string) error {
	if err := uuid.Validate(userID); err != nil {
		return domain.NewInvalidValueErr(fmt.Sprintf("invalid user ID: %s", err))
	}
	if err := s.repo.DeleteUser(ctx, userID); err != nil {
		switch {
		case errors.Is(err, rp.ErrContext):
			return err
		case errors.Is(err, rp.ErrDatabase):
			return domain.NewQueryFailedErr(err)
		case errors.Is(err, rp.ErrNotDeleted):
			return domain.NewNotDeletedErr(err)
		default:
			return domain.NewUnexpectTypeErr(err)
		}
	}

	return nil
}
