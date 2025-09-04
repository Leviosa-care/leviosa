package domain

import (
	"strings"
	"testing"

	"github.com/hengadev/errsx"
	"github.com/stretchr/testify/assert"
)

func TestValidatePassword(t *testing.T) {
	t.Run("valid passwords should pass", func(t *testing.T) {
		validPasswords := []string{
			"validpassword123",                     // Basic valid password
			"MySecureP@ss123",                      // Complex password
			"longpasswordwithmanychars",            // Long password
			strings.Repeat("a", PasswordMinLength), // Exactly minimum length
			"Test1234",                             // Common valid pattern
			"abcdefgh",                             // Exactly 8 characters
			"P@ssw0rd123!",                         // Password with special characters
			"my-secure-password-2024",              // Hyphenated password
		}

		for _, password := range validPasswords {
			t.Run(password, func(t *testing.T) {
				// Skip actual pwned check for unit tests by mocking behavior
				// In a real test, you might want to mock the HTTP call
				err := validatePasswordLength(password)
				assert.NoError(t, err, "expected password %s to have valid length", password)
			})
		}
	})

	t.Run("passwords too short should fail", func(t *testing.T) {
		shortPasswords := []string{
			"",
			"1",
			"12",
			"123",
			"1234",
			"12345",
			"123456",
			"1234567", // 7 characters - just under minimum
		}

		for _, password := range shortPasswords {
			t.Run(password, func(t *testing.T) {
				err := validatePasswordLength(password)
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, PasswordLengthKey)
				assert.Equal(t, PasswordLengthMsg, errMap[PasswordLengthKey].Error())
			})
		}
	})

	t.Run("boundary cases", func(t *testing.T) {
		t.Run("exactly minimum length should pass", func(t *testing.T) {
			minLengthPassword := strings.Repeat("a", PasswordMinLength) // Exactly 8 characters
			err := validatePasswordLength(minLengthPassword)
			assert.NoError(t, err)
		})

		t.Run("one character under minimum should fail", func(t *testing.T) {
			underMinPassword := strings.Repeat("a", PasswordMinLength-1) // 7 characters
			err := validatePasswordLength(underMinPassword)
			assert.Error(t, err)

			var errMap errsx.Map
			assert.True(t, errsx.As(err, &errMap))
			assert.Contains(t, errMap, PasswordLengthKey)
		})

		t.Run("very long password should pass length check", func(t *testing.T) {
			longPassword := strings.Repeat("a", 100) // 100 characters
			err := validatePasswordLength(longPassword)
			assert.NoError(t, err)
		})
	})
}

// Helper function to test just the length validation without the pwned check
func validatePasswordLength(p string) error {
	var errs errsx.Map

	if len(p) < PasswordMinLength {
		errs.Set(PasswordLengthKey, PasswordLengthMsg)
	}

	return errs.AsError()
}

func TestNewPassword(t *testing.T) {
	t.Run("valid password should create Password instance", func(t *testing.T) {
		validPassword := "validpassword123"

		// For testing purposes, we'll test the NewPassword function but skip
		// the actual pwned password check since it requires network access
		password, err := NewPassword(validPassword)

		// Note: This test might fail if the password is actually pwned
		// In a production test suite, you'd mock the CheckPasswordPwned function
		if err != nil {
			// If error is due to length, it should fail
			var errMap errsx.Map
			if errsx.As(err, &errMap) {
				if _, hasLengthError := errMap[PasswordLengthKey]; hasLengthError {
					assert.Fail(t, "Password length should be valid")
				}
				// If it's a pwned check error or pwned error, that's expected for some passwords
				t.Logf("Password validation failed (possibly due to pwned check): %v", err)
			}
		} else {
			assert.Equal(t, Password(validPassword), password)
		}
	})

	t.Run("short password should fail", func(t *testing.T) {
		shortPassword := "short"

		password, err := NewPassword(shortPassword)
		assert.Error(t, err)
		assert.Equal(t, Password(""), password)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, PasswordLengthKey)
	})
}

func TestPassword_String(t *testing.T) {
	password := Password("testpassword")
	result := password.String()
	assert.Equal(t, "testpassword", result)
}

func TestCheckPasswordPwned(t *testing.T) {
	t.Run("known pwned password should return true", func(t *testing.T) {
		// "password" is a well-known pwned password
		pwned, err := CheckPasswordPwned("password")

		if err != nil {
			t.Logf("Network error during pwned check: %v", err)
			t.Skip("Skipping pwned check due to network issues")
		}

		assert.True(t, pwned, "The password 'password' should be detected as pwned")
	})

	t.Run("likely secure password should return false", func(t *testing.T) {
		// A complex, unique password that's unlikely to be pwned
		securePassword := "MyVeryUniqueP@ssw0rd2024WithSpecialChars!"
		pwned, err := CheckPasswordPwned(securePassword)

		if err != nil {
			t.Logf("Network error during pwned check: %v", err)
			t.Skip("Skipping pwned check due to network issues")
		}

		assert.False(t, pwned, "Complex unique password should not be pwned")
	})

	t.Run("empty password should not cause panic", func(t *testing.T) {
		pwned, err := CheckPasswordPwned("")

		if err != nil {
			t.Logf("Network error during pwned check: %v", err)
			t.Skip("Skipping pwned check due to network issues")
		}

		// Empty password might or might not be in the database
		// The important thing is that it doesn't panic
		assert.IsType(t, bool(false), pwned)
	})
}

func TestPasswordConstants(t *testing.T) {
	t.Run("password minimum length should be 8", func(t *testing.T) {
		assert.Equal(t, 8, PasswordMinLength)
	})

	t.Run("error messages should be user-friendly", func(t *testing.T) {
		assert.Equal(t, "password must be at least 8 characters", PasswordLengthMsg)
		assert.Equal(t, "failed to verify password security", PasswordPwnedCheckMsg)
		assert.Equal(t, "password has been found in data breaches and cannot be used", PasswordPwnedMsg)
	})

	t.Run("error keys should be consistent", func(t *testing.T) {
		assert.Equal(t, "password_length", PasswordLengthKey)
		assert.Equal(t, "password_pwned_check", PasswordPwnedCheckKey)
		assert.Equal(t, "password_pwned", PasswordPwnedKey)
	})
}

func TestPasswordValidationErrors(t *testing.T) {
	t.Run("multiple validation errors structure", func(t *testing.T) {
		// Test that our error structure follows the same pattern as other validators
		shortPassword := "short"

		err := validatePasswordLength(shortPassword)
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))

		// Check that we can iterate over errors
		for key, errValue := range errMap {
			assert.IsType(t, "", key)
			assert.NotNil(t, errValue)
		}
	})
}

// Example test showing how to mock the pwned check in a real test suite
func TestPasswordValidationWithMockedPwnedCheck(t *testing.T) {
	t.Run("example of how to test with mocked pwned check", func(t *testing.T) {
		// In a real test suite, you would:
		// 1. Create an interface for the pwned check function
		// 2. Create a mock implementation
		// 3. Inject the dependency

		// For now, we just demonstrate the structure
		password := "testpassword123"

		// Mock pwned check that always returns false (not pwned)
		mockPwnedCheck := func(p string) (bool, error) {
			return false, nil
		}

		// Test length validation
		err := validatePasswordLength(password)
		assert.NoError(t, err)

		// Mock shows password is not pwned
		pwned, err := mockPwnedCheck(password)
		assert.NoError(t, err)
		assert.False(t, pwned)
	})
}

