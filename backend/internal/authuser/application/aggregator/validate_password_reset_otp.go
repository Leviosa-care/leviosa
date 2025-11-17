package aggregator

import (
	"context"
	"fmt"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
)

func (s *AuthAggregatorService) ValidatePasswordResetOTP(ctx context.Context, request *domain.ValidatePasswordResetOTPRequest) (*domain.ValidatePasswordResetOTPResponse, error) {
	// First, validate the OTP
	if err := s.otp.ValidateOTP(ctx, &domain.ValidateOTPRequest{
		Email: request.Email,
		Code:  request.Code,
	}); err != nil {
		return nil, err
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
		return nil, err
	}

	// Return the unhashed token to the client
	expiresAt := time.Now().Add(resetTokenTTL)
	return &domain.ValidatePasswordResetOTPResponse{
		Token:     resetToken,
		ExpiresAt: expiresAt,
	}, nil
}
