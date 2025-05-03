package userService

import (
	"context"
	"errors"
	"fmt"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

// ValidateCredentials verifies the user's credentials by checking if the provided email and password
// match the stored data in the repository.
//
// Parameters:
//   - ctx: The context used for the operation.
//   - user: The user object containing the email and password to be validated.
//
// Returns:
//   - error: An error indicating any issues during the validation process. It can return:
//   - NewNotFoundErr if the email is not found.
//   - NewQueryFailedErr if there is a query failure.
//   - NewInvalidValueErr if the password verification fails.
func (s *service) ValidateCredentials(ctx context.Context, user *models.UserSignIn) error {
	emailHash := s.crypto.HashBasic(ctx, []byte(user.Email))
	hashedPassword, err := s.repo.GetHashedPasswordByEmail(ctx, emailHash)
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrNotFound):
			return domain.NewNotFoundErr(err)
		case errors.Is(err, rp.ErrContext):
			return err
		case errors.Is(err, rp.ErrDatabase):
			return domain.NewQueryFailedErr(err)
		default:
			return domain.NewUnexpectTypeErr(err)
		}
	}
	// ok, err := s.crypto.VerifyPassword(ctx, user.Password, hashedPassword)
	ok, err := s.crypto.CompareSecureHashAndValue(ctx, user.Password, hashedPassword)
	if err != nil {
		return domain.NewInvalidValueErr(fmt.Sprintf("invalid password verification: %s", err.Error()))
	}
	if !ok {
		return domain.NewInvalidValueErr("password does not match")
	}
	return nil
}
