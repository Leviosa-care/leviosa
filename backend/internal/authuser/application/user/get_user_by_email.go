package user

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/validation"

	"github.com/hengadev/encx"
)

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.UserResponse, error) {
	if err := validation.ValidateEmail(email); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	emailBytes, err := encx.SerializeValue(email)
	if err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}
	emailHash := s.crypto.HashBasic(ctx, emailBytes)

	userEncx, err := s.repo.GetUserByEmailHash(ctx, emailHash)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	// Decrypt user data
	user, err := domain.DecryptUserEncx(ctx, s.crypto, userEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("user retrieved by email", err)
	}

	return user.ToResponse(), nil
}
