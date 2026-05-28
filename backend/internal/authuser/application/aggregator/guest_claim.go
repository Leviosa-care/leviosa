package aggregator

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// GuestClaim validates the guest claim request, checks email availability,
// and sends an OTP to the provided email. Returns 202 on success.
func (s *AuthAggregatorService) GuestClaim(ctx context.Context, req *domain.GuestClaimRequest) error {
	if err := req.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	// Check if email is already registered
	if err := s.user.ValidateGuestClaimEmail(ctx, req.Email); err != nil {
		return err
	}

	// Send OTP to the provided email
	if err := s.otp.RequestOTP(ctx, req.Email); err != nil {
		return err
	}

	return nil
}
