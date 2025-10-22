package otp

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
)

func (s *OTPService) generateOTP(email string) (*domain.OTP, error) {
	length := s.effectiveOTPLength()

	code, err := s.generateSecureCode(length)
	if err != nil {
		return nil, fmt.Errorf("failed to generate secure code: %w", err)
	}

	now := time.Now()
	duration := time.Duration(defaultOTPDuration) * time.Minute

	return &domain.OTP{Email: email,
		Code:      code,
		Attempts:  0,
		CreatedAt: now,
		ExpiresAt: now.Add(duration),
		// DEK:       dek,
	}, nil
}

func (s *OTPService) generateSecureCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be positive")
	}

	// Calculate maximum value for the given length
	max := int64(1)
	for range length {
		max *= 10
	}

	// Generate secure random bytes
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate secure random bytes: %w", err)
	}

	// Convert to number within range
	num := int64(binary.BigEndian.Uint64(bytes)) % max
	if num < 0 {
		num = -num
	}

	// Format with leading zeros
	format := fmt.Sprintf("%%0%dd", length)
	return fmt.Sprintf(format, num), nil
}
