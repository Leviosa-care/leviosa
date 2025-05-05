package userService

import (
	"context"
	"errors"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

// GetUserSessionData retrieves the session data for a user based on their email.
//
// Parameters:
//   - ctx: A context.Context instance to manage request lifecycle and cancellation.
//   - email: A string representing the email address of the user whose session data is being retrieved.
//
// Returns:
//   - string: The user's ID if found. If an error occurs or the user is not found, an empty string is returned.
//   - models.Role: The user's role. If an error occurs or the user is not found, models.UNKNOWN is returned.
//   - error: An error if the session data cannot be retrieved, the email is invalid, or an unexpected error occurs.
//     Returns nil if the session data is successfully retrieved.
func (s *service) GetUserSessionData(ctx context.Context, email string) (string, models.Role, error) {
	if _, err := models.NewEmail(email); err != nil {
		return "", models.VISITOR, domain.NewInvalidValueErr(fmt.Sprintf("invalid email: %q", err))
	}
	emailHash := s.crypto.HashBasic(ctx, []byte(email))
	ID, role, err := s.repo.GetUserSessionData(ctx, emailHash)
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrNotFound):
			return "", models.VISITOR, domain.NewNotFoundErr(err)
		case errors.Is(err, rp.ErrContext):
			return "", models.VISITOR, err
		case errors.Is(err, rp.ErrDatabase):
			return "", models.VISITOR, domain.NewQueryFailedErr(err)
		default:
			return "", models.VISITOR, domain.NewUnexpectTypeErr(err)
		}
	}
	return ID, role, nil
}
