package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/hengadev/encx"
)

func (s *UserService) CheckEmailAvailability(ctx context.Context, request *domain.CheckEmailAvailabilityRequest) (bool, error) {
	if err := request.Valid(ctx); err != nil {
		return false, errs.NewInvalidValueErr(err.Error())
	}

	emailBytes, err := encx.SerializeValue(request.Email)
	if err != nil {
		return false, errs.NewInvalidValueErr(fmt.Sprintf("failed to serialize userID: %v", err))
	}
	emailHash := s.crypto.HashBasic(ctx, emailBytes)

	exists, err := s.repo.ExistsByEmailHash(ctx, emailHash)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			// User not found means email is available - this is success
			return true, nil
		}
		return false, fmt.Errorf("check email availability: %w", err)
	}

	// Email is available if user does NOT exist
	return !exists, nil
}
