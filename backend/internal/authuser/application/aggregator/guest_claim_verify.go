package aggregator

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

// GuestClaimVerify validates the OTP, creates a fully active user account
// with profile_incomplete=true, and returns an active standard-role session.
func (s *AuthAggregatorService) GuestClaimVerify(ctx context.Context, req *domain.GuestClaimVerifyRequest, claimData *domain.GuestClaimRequest) (*domain.CreateSessionResponse, error) {
	// Validate claimData before consuming the OTP so a bad payload (weak password,
	// invalid phone, etc.) doesn't burn the OTP and force the user to restart.
	if err := claimData.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Validate the OTP (this also consumes it on success).
	if err := s.otp.ValidateOTP(ctx, &domain.ValidateOTPRequest{
		Email: req.Email,
		Code:  req.Code,
	}); err != nil {
		return nil, err
	}

	// OTP is valid — create the guest user account (active, standard, profile_incomplete=true)
	userID, err := s.user.CreateGuestUser(ctx, claimData)
	if err != nil {
		return nil, err
	}

	// Create an active session with standard role
	response, err := s.session.CreateSession(ctx, &domain.CreateSessionRequest{
		UserID: userID.String(),
		Role:   identity.Standard,
		State:  session.SessionActive,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
