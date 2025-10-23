package aggregator

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
)

func (s *AuthAggregatorService) ValidatePasswordResetOTP(ctx context.Context, request *domain.ValidatePasswordResetOTPRequest) (*domain.ValidatePasswordResetOTPResponse, error) {
	// First, validate the OTP
	if err := s.otp.ValidateOTP(ctx, &domain.ValidateOTPRequest{
		Email: request.Email,
		Code:  request.Code,
	}); err != nil {
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
			return nil, fmt.Errorf("validate password reset OTP: %w", err) // Wrap unexpected errors with context
		}
	}

	// OTP is valid, now generate reset token
	resetToken, err := session.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("generate password reset token: %w", err)
	}

	// Store reset session with 15-minute expiration
	resetTokenTTL := 15 * time.Minute
	// if err := s.session.CreateResetSession(ctx, tokenHash, userEmailHash, resetTokenTTL); err != nil {
	if err := s.session.CreateResetSession(ctx, resetToken, request.Email, resetTokenTTL); err != nil {
		switch {
		case errors.Is(err, errs.ErrExternalService):
			return nil, err // Pass through external service errors
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections), errors.Is(err, errs.ErrResourceExhausted), errors.Is(err, errs.ErrTransactionFailure):
			return nil, err // Pass through infrastructure errors
		case errors.Is(err, context.Canceled):
			return nil, fmt.Errorf("create password reset session cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return nil, fmt.Errorf("create password reset session timeout: %w", err)
		default:
			return nil, fmt.Errorf("create password reset session: %w", err) // Wrap unexpected errors with context
		}
	}

	// Return the unhashed token to the client
	expiresAt := time.Now().Add(resetTokenTTL)
	return &domain.ValidatePasswordResetOTPResponse{
		Token:     resetToken,
		ExpiresAt: expiresAt,
	}, nil
}
