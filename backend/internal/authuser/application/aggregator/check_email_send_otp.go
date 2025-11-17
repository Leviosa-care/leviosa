package aggregator

import (
	"context"
	"errors"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *AuthAggregatorService) CheckEmailSendOTP(ctx context.Context, request *domain.CheckEmailAvailabilityRequest) error {
	if err := request.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	available, err := s.user.CheckEmailAvailability(ctx, request)
	if err != nil {
		return err
	}

	if !available {
		return errs.NewConflictErr(errors.New("email is already registered"))
	}

	if err := s.otp.RequestOTP(ctx, request.Email); err != nil {
		return err
	}

	return nil
}
