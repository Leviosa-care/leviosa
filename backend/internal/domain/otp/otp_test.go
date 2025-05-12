package otpService

import (
	"testing"
	"time"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestIncreaseAttempt(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name             string
		initialAttempts  int
		expectedAttempts int
		expectedError    bool
	}{
		{
			name:             "Initial attempts less than max",
			initialAttempts:  0,
			expectedAttempts: 1,
			expectedError:    false,
		},
		{
			name:             "One attempt away from max",
			initialAttempts:  MaxOTPAttempts - 1,
			expectedAttempts: MaxOTPAttempts,
			expectedError:    false,
		},
		{
			name:             "At max attempts",
			initialAttempts:  MaxOTPAttempts,
			expectedAttempts: MaxOTPAttempts,
			expectedError:    true,
		},
		{
			name:             "Already exceeded max attempts",
			initialAttempts:  MaxOTPAttempts + 1,
			expectedAttempts: MaxOTPAttempts + 1,
			expectedError:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			otp := &OTP{
				Email:     "test@example.com",
				Code:      "123456",
				Attempts:  tt.initialAttempts,
				ExpiresAt: now.Add(OTPDURATION),
				CreatedAt: now,
			}

			err := otp.increaseAttempt()

			if tt.expectedError {
				assert.Error(t, err, "Expected an error")
				assert.ErrorIs(t, err, domain.ErrInvalidValue)
				assert.Equal(t, tt.initialAttempts, otp.Attempts, "Attempts should not have increased")
			} else {
				assert.NoError(t, err, "Did not expect an error")
				assert.Equal(t, tt.expectedAttempts, otp.Attempts, "Attempts should have increased")
			}
		})
	}
}
