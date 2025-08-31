package validation

import (
	"strings"
	"testing"

	"github.com/hengadev/errsx"
	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	t.Run("valid emails should pass", func(t *testing.T) {
		validEmails := []string{
			"test@example.com",
			"user@domain.org",
			"test.email@example.co.uk",
			"user+tag@example.com",
			"user_name@domain.net",
			"user-name@domain.com",
			"a@b.com",
			"test123@example123.com",
			"user~test@example.com",
		}

		for _, email := range validEmails {
			t.Run(email, func(t *testing.T) {
				err := ValidateEmail(email)
				assert.NoError(t, err, "expected %s to be valid", email)
			})
		}
	})

	t.Run("empty email should fail", func(t *testing.T) {
		err := ValidateEmail("")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, EmailRequiredKey)
		assert.Equal(t, EmailRequiredMsg, errMap[EmailRequiredKey].Error())
	})

	t.Run("whitespace-only email should fail", func(t *testing.T) {
		err := ValidateEmail("   ")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, EmailRequiredKey)
	})

	t.Run("email exceeding max length should fail", func(t *testing.T) {
		longEmail := strings.Repeat("a", EmailMaxLength) + "@example.com"
		err := ValidateEmail(longEmail)
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, EmailLengthKey)
		assert.Equal(t, "email cannot exceed 255 characters", errMap[EmailLengthKey].Error())
	})

	t.Run("email with whitespace should fail", func(t *testing.T) {
		whitespaceEmails := []string{
			"test @example.com",
			"test@ example.com",
			"test@example .com",
			"test\t@example.com",
			"test@example.com\n",
			"test@example.com\r",
		}

		for _, email := range whitespaceEmails {
			t.Run(email, func(t *testing.T) {
				err := ValidateEmail(email)
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, EmailWhitespaceKey)
				assert.Equal(t, EmailWhitespaceMsg, errMap[EmailWhitespaceKey].Error())
			})
		}
	})

	t.Run("email with quotes should fail", func(t *testing.T) {
		quoteEmails := []string{
			`"test"@example.com`,
			`test'@example.com`,
			`test@"example".com`,
			`test@example'com`,
		}

		for _, email := range quoteEmails {
			t.Run(email, func(t *testing.T) {
				err := ValidateEmail(email)
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, EmailQuotesKey)
				assert.Equal(t, EmailQuotesMsg, errMap[EmailQuotesKey].Error())
			})
		}
	})

	t.Run("email missing @ sign should fail", func(t *testing.T) {
		err := ValidateEmail("testexample.com")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, EmailFormatKey)
		assert.Equal(t, EmailMissingAtMsg, errMap[EmailFormatKey].Error())
	})

	t.Run("email starting with @ should fail", func(t *testing.T) {
		err := ValidateEmail("@example.com")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, EmailFormatKey)
		assert.Equal(t, EmailMissingBeforeMsg, errMap[EmailFormatKey].Error())
	})

	t.Run("email ending with @ should fail", func(t *testing.T) {
		err := ValidateEmail("test@")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, EmailFormatKey)
		assert.Equal(t, EmailMissingAfterMsg, errMap[EmailFormatKey].Error())
	})

	t.Run("email with name should fail", func(t *testing.T) {
		err := ValidateEmail("Test User <test@example.com>")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, EmailNameKey)
		assert.Equal(t, EmailNameNotAllowedMsg, errMap[EmailNameKey].Error())
	})

	t.Run("email with invalid characters should fail", func(t *testing.T) {
		invalidCharEmails := []string{
			"test#@example.com",
			"test$@example.com",
			"test%@example.com",
			"test&@example.com",
			"test*@example.com",
			"test!@example.com",
		}

		for _, email := range invalidCharEmails {
			t.Run(email, func(t *testing.T) {
				err := ValidateEmail(email)
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, EmailCharsKey)
				assert.Contains(t, errMap[EmailCharsKey].Error(), "email contains invalid characters")
			})
		}
	})

	t.Run("email missing TLD should fail", func(t *testing.T) {
		err := ValidateEmail("test@example")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, EmailFormatKey)
		assert.Equal(t, EmailMissingTLDMsg, errMap[EmailFormatKey].Error())
	})

	t.Run("multiple validation errors should be reported", func(t *testing.T) {
		err := ValidateEmail("test @example")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))

		// Should have both whitespace and format errors
		assert.Contains(t, errMap, EmailWhitespaceKey)
		assert.Contains(t, errMap, EmailFormatKey)
	})
}
