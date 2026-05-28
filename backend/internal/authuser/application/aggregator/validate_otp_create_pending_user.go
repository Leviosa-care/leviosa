package aggregator

import (
	"context"
	"errors"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *AuthAggregatorService) ValidateOTPCreatePendingUser(ctx context.Context, request *domain.ValidateOTPRequest) (*domain.CreateSessionResponse, error) {
	// First, validate the OTP
	if err := s.otp.ValidateOTP(ctx, request); err != nil {
		return nil, err
	}

	// OTP is valid, now create the pending user.
	// CreatePendingUser returns the existing user's ID on ErrConflict (race condition
	// where the same email was verified twice), so we can proceed to create a session.
	userID, err := s.user.CreatePendingUser(ctx, request.Email)
	if err != nil {
		if !errors.Is(err, errs.ErrConflict) {
			return nil, err
		}
	}

	response, err := s.session.CreateSession(ctx, &domain.CreateSessionRequest{
		UserID: userID.String(),
		Role:   identity.Visitor,
		State:  session.SessionPending,
	})
	if err != nil {
		return nil, err
	}

	// Fire-and-forget: claim guest bookings linked to this email.
	// A failure must never block or roll back account creation.
	if s.bookingClient != nil {
		s.bookingClient.ClaimBookings(ctx, userID.String(), request.Email)
	}

	return response, nil
}
