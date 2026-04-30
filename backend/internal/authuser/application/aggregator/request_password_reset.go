package aggregator

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *AuthAggregatorService) RequestPasswordReset(ctx context.Context, request *domain.RequestPasswordResetRequest) error {
	if err := request.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	available, err := s.user.CheckEmailAvailability(ctx, &domain.CheckEmailAvailabilityRequest{
		Email: request.Email,
	})
	if err != nil {
		return err
	}

	if available {
		// Email not registered — return success silently to prevent enumeration
		return nil
	}

	if err := s.otp.RequestOTP(ctx, request.Email); err != nil {
		return err
	}
	return nil
}
