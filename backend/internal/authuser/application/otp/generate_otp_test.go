package otp

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateSecureCode(t *testing.T) {
	s := &OTPService{}

	t.Run("should generate code with correct length", func(t *testing.T) {
		lengths := []int{4, 6, 8, 10}

		for _, length := range lengths {
			code, err := s.generateSecureCode(length)
			require.NoError(t, err)
			assert.Len(t, code, length, "Code should have length %d", length)

			// Verify all characters are digits
			for _, char := range code {
				assert.True(t, char >= '0' && char <= '9', "All characters should be digits")
			}
		}
	})

	t.Run("should generate different codes on multiple calls", func(t *testing.T) {
		length := 6
		codes := make(map[string]bool)

		// Generate multiple codes
		for i := 0; i < 100; i++ {
			code, err := s.generateSecureCode(length)
			require.NoError(t, err)
			codes[code] = true
		}

		// Should have generated many different codes (at least 80% unique)
		assert.Greater(t, len(codes), 80, "Should generate diverse codes")
	})

	t.Run("should pad with leading zeros", func(t *testing.T) {
		length := 6

		for i := 0; i < 50; i++ {
			code, err := s.generateSecureCode(length)
			require.NoError(t, err)

			// Check that code is exactly the right length (padded if necessary)
			assert.Len(t, code, length)
			assert.Regexp(t, `^\d{6}$`, code, "Should be exactly 6 digits")
		}
	})

	t.Run("should handle single digit length", func(t *testing.T) {
		code, err := s.generateSecureCode(1)
		require.NoError(t, err)
		assert.Len(t, code, 1)
		assert.Regexp(t, `^\d$`, code, "Should be a single digit")
	})

	t.Run("should return error for invalid length", func(t *testing.T) {
		invalidLengths := []int{0, -1, -5}

		for _, length := range invalidLengths {
			code, err := s.generateSecureCode(length)
			assert.Error(t, err, "Should return error for length %d", length)
			assert.Empty(t, code)
			assert.Contains(t, err.Error(), "length must be positive")
		}
	})

	t.Run("should handle very large lengths", func(t *testing.T) {
		length := 15
		code, err := s.generateSecureCode(length)
		require.NoError(t, err)
		assert.Len(t, code, length)

		// Verify all digits
		for _, char := range code {
			assert.True(t, char >= '0' && char <= '9')
		}
	})

	t.Run("should generate codes with proper distribution", func(t *testing.T) {
		length := 2
		digitCount := make(map[rune]int)

		// Generate many codes to check distribution
		for i := 0; i < 1000; i++ {
			code, err := s.generateSecureCode(length)
			require.NoError(t, err)

			for _, digit := range code {
				digitCount[digit]++
			}
		}

		// Each digit (0-9) should appear at least once in 2000 total digit positions
		for digit := '0'; digit <= '9'; digit++ {
			assert.Greater(t, digitCount[digit], 0, "Digit %c should appear at least once", digit)
		}
	})
}

func TestGenerateOTP(t *testing.T) {
	// Mock OTPService with test values
	s := &OTPService{
		cache: NewMockOTPCache(6, 15, 3), // length, duration (15 minutes), maxAttempts
	}

	t.Run("should generate valid OTP with correct fields", func(t *testing.T) {
		email := "test@example.com"
		startTime := time.Now()

		otp, err := s.generateOTP(email)
		require.NoError(t, err)
		require.NotNil(t, otp)

		// Verify fields
		assert.Equal(t, email, otp.Email)
		assert.Len(t, otp.Code, 6, "Code should have default length of 6")
		assert.Equal(t, 0, otp.Attempts)

		// Verify timestamps
		assert.True(t, otp.CreatedAt.After(startTime.Add(-time.Second)), "CreatedAt should be around current time")
		assert.True(t, otp.CreatedAt.Before(time.Now().Add(time.Second)), "CreatedAt should be around current time")

		expectedExpiry := otp.CreatedAt.Add(15 * time.Minute)
		assert.Equal(t, expectedExpiry, otp.ExpiresAt, "ExpiresAt should be CreatedAt + duration")

		// Verify code is numeric
		assert.Regexp(t, `^\d{6}$`, otp.Code, "Code should be 6 digits")
	})

	t.Run("should generate different codes for same email", func(t *testing.T) {
		email := "test@example.com"
		codes := make(map[string]bool)

		for i := 0; i < 10; i++ {
			otp, err := s.generateOTP(email)
			require.NoError(t, err)
			codes[otp.Code] = true
		}

		assert.Greater(t, len(codes), 8, "Should generate different codes")
	})

	t.Run("should handle different email addresses", func(t *testing.T) {
		emails := []string{
			"user@example.com",
			"test.user@domain.co.uk",
			"user+tag@example.org",
		}

		for _, email := range emails {
			otp, err := s.generateOTP(email)
			require.NoError(t, err, "Should generate OTP for email: %s", email)
			assert.Equal(t, email, otp.Email)
			assert.NotEmpty(t, otp.Code)
		}
	})

	t.Run("should handle empty email", func(t *testing.T) {
		otp, err := s.generateOTP("")
		require.NoError(t, err)
		assert.Equal(t, "", otp.Email)
		assert.NotEmpty(t, otp.Code)
	})

	t.Run("should respect custom OTP length from cache", func(t *testing.T) {
		// Test with different length
		s.cache = NewMockOTPCache(8, 10, 3)

		otp, err := s.generateOTP("test@example.com")
		require.NoError(t, err)
		assert.Len(t, otp.Code, 8, "Should use custom length from cache")
	})

	t.Run("should respect custom duration from cache", func(t *testing.T) {
		customDuration := 30 // 30 minutes
		s.cache = NewMockOTPCache(6, customDuration, 3)

		startTime := time.Now()
		otp, err := s.generateOTP("test@example.com")
		require.NoError(t, err)

		expectedExpiry := otp.CreatedAt.Add(time.Duration(customDuration) * time.Minute)
		assert.Equal(t, expectedExpiry, otp.ExpiresAt, "Should use custom duration from cache")
		assert.True(t, otp.ExpiresAt.After(startTime.Add(25*time.Minute)), "Should be at least 25 minutes from start")
	})

	t.Run("should initialize attempts to zero", func(t *testing.T) {
		otp, err := s.generateOTP("test@example.com")
		require.NoError(t, err)
		assert.Equal(t, 0, otp.Attempts, "Attempts should be initialized to 0")
	})

	t.Run("should set CreatedAt before ExpiresAt", func(t *testing.T) {
		otp, err := s.generateOTP("test@example.com")
		require.NoError(t, err)
		assert.True(t, otp.CreatedAt.Before(otp.ExpiresAt), "CreatedAt should be before ExpiresAt")
	})
}


func TestGenerateOTP_EdgeCases(t *testing.T) {
	s := &OTPService{
		cache: NewMockOTPCache(6, 15, 3),
	}

	t.Run("should handle very long email addresses", func(t *testing.T) {
		longEmail := strings.Repeat("a", 100) + "@" + strings.Repeat("example", 20) + ".com"

		otp, err := s.generateOTP(longEmail)
		require.NoError(t, err)
		assert.Equal(t, longEmail, otp.Email)
		assert.NotEmpty(t, otp.Code)
	})

	t.Run("should handle special characters in email", func(t *testing.T) {
		specialEmail := "test+tag@example-domain.co.uk"

		otp, err := s.generateOTP(specialEmail)
		require.NoError(t, err)
		assert.Equal(t, specialEmail, otp.Email)
		assert.NotEmpty(t, otp.Code)
	})

	t.Run("should handle minimum duration", func(t *testing.T) {
		s.cache = NewMockOTPCache(6, 1, 3) // 1 minute duration

		otp, err := s.generateOTP("test@example.com")
		require.NoError(t, err)

		expectedExpiry := otp.CreatedAt.Add(1 * time.Minute)
		assert.Equal(t, expectedExpiry, otp.ExpiresAt)
	})

	t.Run("should handle large duration", func(t *testing.T) {
		s.cache = NewMockOTPCache(6, 1440, 3) // 24 hours duration

		otp, err := s.generateOTP("test@example.com")
		require.NoError(t, err)

		expectedExpiry := otp.CreatedAt.Add(1440 * time.Minute)
		assert.Equal(t, expectedExpiry, otp.ExpiresAt)
	})
}

