package validation

import (
	"strings"
	"testing"

	"github.com/hengadev/errsx"
	"github.com/stretchr/testify/assert"
)

func TestValidatePhone(t *testing.T) {
	t.Run("valid French phone numbers should pass", func(t *testing.T) {
		validPhones := []string{
			"0123456789", // Fixed line (01-05)
			"0223456789", // Fixed line (01-05)
			"0323456789", // Fixed line (01-05)
			"0423456789", // Fixed line (01-05)
			"0523456789", // Fixed line (01-05)
			"0612345678", // Mobile (06)
			"0723456789", // Mobile (07)
		}

		for _, phone := range validPhones {
			t.Run(phone, func(t *testing.T) {
				err := ValidatePhone(phone)
				assert.NoError(t, err, "expected %s to be valid", phone)
			})
		}
	})

	t.Run("valid phone numbers with whitespace should pass", func(t *testing.T) {
		phoneWithSpaces := []string{
			"  0123456789  ",
			"\t0612345678\t",
			" 0723456789 ",
		}

		for _, phone := range phoneWithSpaces {
			t.Run(phone, func(t *testing.T) {
				err := ValidatePhone(phone)
				assert.NoError(t, err, "expected %s to be valid after trimming", phone)
			})
		}
	})

	t.Run("phone numbers too short should fail", func(t *testing.T) {
		shortPhones := []string{
			"",
			"0",
			"01",
			"012345678", // 9 digits
		}

		for _, phone := range shortPhones {
			t.Run(phone, func(t *testing.T) {
				err := ValidatePhone(phone)
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, PhoneTooShortKey)
				assert.Equal(t, PhoneTooShortMsg, errMap[PhoneTooShortKey].Error())
			})
		}
	})

	t.Run("phone numbers too long should fail", func(t *testing.T) {
		longPhone := strings.Repeat("1", 21)
		err := ValidatePhone(longPhone)
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, PhoneTooLongKey)
		assert.Equal(t, PhoneTooLongMsg, errMap[PhoneTooLongKey].Error())
	})

	t.Run("invalid French phone number formats should fail", func(t *testing.T) {
		invalidPhones := []string{
			"1234567890",   // Not starting with 0
			"0000000000",   // Invalid prefix 00
			"06123456789",  // Too long for mobile
			"0812345678",   // Invalid prefix 08
			"0912345678",   // Invalid prefix 09
			"0012345678",   // Invalid prefix 00
			"0a12345678",   // Contains letter
			"06 12345678",  // Contains space (after trim check)
			"06-12345678",  // Contains dash
			"06.12345678",  // Contains dot
			"(06)12345678", // Contains parentheses
			"+33612345678", // International format not supported
		}

		for _, phone := range invalidPhones {
			t.Run(phone, func(t *testing.T) {
				err := ValidatePhone(phone)
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, PhoneFormatKey)
				assert.Equal(t, PhoneFormatMsg, errMap[PhoneFormatKey].Error())
			})
		}
	})

	t.Run("edge cases for French prefixes", func(t *testing.T) {
		t.Run("boundary valid prefixes", func(t *testing.T) {
			validBoundaryPhones := []string{
				"0123456789", // 01 - minimum fixed line
				"0523456789", // 05 - maximum fixed line
				"0612345678", // 06 - minimum mobile
				"0723456789", // 07 - maximum mobile
			}

			for _, phone := range validBoundaryPhones {
				err := ValidatePhone(phone)
				assert.NoError(t, err, "expected %s to be valid", phone)
			}
		})

		t.Run("boundary invalid prefixes", func(t *testing.T) {
			invalidBoundaryPhones := []string{
				"0012345678", // 00 - below minimum
				"0812345678", // 08 - above maximum mobile
			}

			for _, phone := range invalidBoundaryPhones {
				err := ValidatePhone(phone)
				assert.Error(t, err, "expected %s to be invalid", phone)
			}
		})
	})

	t.Run("whitespace handling", func(t *testing.T) {
		t.Run("only whitespace should fail", func(t *testing.T) {
			whitespaceOnlyPhones := []string{
				"   ",
				"\t\t",
				"\n\n",
				"",
			}

			for _, phone := range whitespaceOnlyPhones {
				err := ValidatePhone(phone)
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, PhoneTooShortKey)
				assert.Equal(t, PhoneTooShortMsg, errMap[PhoneTooShortKey].Error())
			}
		})

		t.Run("valid phone with leading/trailing whitespace should pass", func(t *testing.T) {
			err := ValidatePhone("   0123456789   ")
			assert.NoError(t, err)
		})
	})
}
