package aggregator

import (
	"context"
	"errors"

	"github.com/Leviosa-care/authuser/internal/domain"

	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/middleware/auth"
)

func (s *AuthAggregatorService) ValidateOTPCreatePendingUser(ctx context.Context, request *domain.ValidateOTPRequest) (*domain.CreateSessionResponse, error) {
	// First, validate the OTP
	if err := s.otp.ValidateOTP(ctx, request); err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			return nil, err // Pass through validation errors (invalid OTP format, etc.)
		case errors.Is(err, errs.ErrDomainNotFound):
			return nil, err // Pass through not found errors (OTP doesn't exist)
		case errors.Is(err, errs.ErrExpiredToken):
			return nil, err // Pass through expired token errors
		case errors.Is(err, errs.ErrRateLimit):
			return nil, err // Pass through rate limit errors (max attempts exceeded)
		case errors.Is(err, errs.ErrValueMismatch):
			return nil, err // Pass through value mismatch errors (wrong code)
		case errors.Is(err, errs.ErrAlreadyConsumed):
			return nil, err // Pass through already consumed errors (concurrent request consumed OTP)
		case errors.Is(err, errs.ErrExternalService):
			return nil, err // Pass through external service errors (database issues)
		default:
			return nil, errs.NewInternalErr(err) // Wrap unexpected errors
		}
	}

	// OTP is valid, now create the pending user
	userID, err := s.user.CreatePendingUser(ctx, request.Email)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrConflict):
			// User already exists - this could be a race condition
			// Since OTP was valid and we got the existing user's ID, we can proceed
			// The user has verified their email and we can create a session for them
			// Note: userID contains the existing user's ID even with the conflict error
		case errors.Is(err, errs.ErrInvalidValue):
			return nil, err // Pass through validation errors
		case errors.Is(err, errs.ErrExternalService):
			return nil, err // Pass through external service errors (database issues)
		default:
			return nil, errs.NewInternalErr(err) // Wrap unexpected errors
		}
	}

	response, err := s.session.CreateSession(ctx, &domain.CreateSessionRequest{
		UserID: userID.String(),
		Role:   identity.Visitor,
		State:  auth.SessionPending,
	})
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			return nil, err // Pass through validation errors
		case errors.Is(err, errs.ErrDomainNotFound):
			return nil, err // Pass through not found errors
		case errors.Is(err, errs.ErrQueryFailed):
			return nil, errs.NewExternalServiceErr(err, "database error during session creation") // Database query issues
		case errors.Is(err, errs.ErrMarshalJSON):
			return nil, errs.NewInternalErr(err) // JSON marshaling issues
		case errors.Is(err, errs.ErrUnexpectedError):
			return nil, errs.NewInternalErr(err) // Unexpected errors
		default:
			return nil, errs.NewInternalErr(err) // Wrap unexpected errors
		}
	}

	return response, nil
}
