package domain

import (
	"strings"
	"testing"

	"github.com/hengadev/errsx"
	"github.com/stretchr/testify/assert"
)

func TestValidatePostalCode(t *testing.T) {
	t.Run("valid postal codes should pass", func(t *testing.T) {
		validPostalCodes := []string{
			// US ZIP codes
			"12345",
			"12345-6789",

			// Canadian postal codes
			"K1A 0A6",
			"M5V 3L9",

			// UK postal codes
			"SW1A 1AA",
			"M1 1AA",
			"B33 8TH",

			// French postal codes
			"75001",
			"69002",

			// German postal codes
			"10115",
			"80331",

			// Australian postal codes
			"2000",
			"3000",

			// Other international formats
			"1010",                                   // 4 digits
			"12345678",                               // 8 digits
			"ABC 123",                                // Letters and numbers with space
			"A1B-2C3",                                // Letters and numbers with dash
			strings.Repeat("A", PostalCodeMaxLength), // Exactly 10 characters
		}

		for _, postalCode := range validPostalCodes {
			t.Run(postalCode, func(t *testing.T) {
				err := validatePostalCode(postalCode)
				assert.NoError(t, err, "expected %s to be valid", postalCode)
			})
		}
	})

	t.Run("postal codes with whitespace should pass after trimming", func(t *testing.T) {
		postalCodesWithSpaces := []string{
			"  12345  ",
			"\t12345-6789\t",
			" K1A 0A6 ",
		}

		for _, postalCode := range postalCodesWithSpaces {
			t.Run(postalCode, func(t *testing.T) {
				err := validatePostalCode(postalCode)
				assert.NoError(t, err, "expected %s to be valid after trimming", postalCode)
			})
		}
	})

	t.Run("empty postal code should fail", func(t *testing.T) {
		err := validatePostalCode("")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, PostalCodeRequiredKey)
		assert.Equal(t, PostalCodeRequiredMsg, errMap[PostalCodeRequiredKey].Error())
	})

	t.Run("whitespace-only postal code should fail", func(t *testing.T) {
		err := validatePostalCode("   ")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, PostalCodeRequiredKey)
		assert.Equal(t, PostalCodeRequiredMsg, errMap[PostalCodeRequiredKey].Error())
	})

	t.Run("postal code too short should fail", func(t *testing.T) {
		shortPostalCodes := []string{
			"12",
			"A",
			"AB",
		}

		for _, postalCode := range shortPostalCodes {
			t.Run(postalCode, func(t *testing.T) {
				err := validatePostalCode(postalCode)
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, PostalCodeTooShortKey)
				assert.Equal(t, PostalCodeTooShortMsg, errMap[PostalCodeTooShortKey].Error())
			})
		}
	})

	t.Run("postal code too long should fail", func(t *testing.T) {
		longPostalCode := strings.Repeat("A", PostalCodeMaxLength+1) // 11 characters
		err := validatePostalCode(longPostalCode)
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))
		assert.Contains(t, errMap, PostalCodeTooLongKey)
		assert.Equal(t, PostalCodeTooLongMsg, errMap[PostalCodeTooLongKey].Error())
	})

	t.Run("postal code with invalid characters should fail", func(t *testing.T) {
		invalidPostalCodes := []string{
			"12345@",   // @ symbol
			"ABC!123",  // ! symbol
			"12#45",    // # symbol
			"AB$CD",    // $ symbol
			"123%456",  // % symbol
			"A1B*2C3",  // * symbol
			"12(34)5",  // Parentheses
			"[12345]",  // Square brackets
			"123/456",  // Forward slash
			"123\\456", // Backslash
			"123+456",  // Plus sign
			"123=456",  // Equals sign
			"123?456",  // Question mark
			"123&456",  // Ampersand
			"123<456>", // Angle brackets
			"123|456",  // Pipe symbol
			"123~456",  // Tilde
			"123`456",  // Backtick
			"123^456",  // Caret
			"123{456}", // Curly braces
			"123_456",  // Underscore
			"123.456",  // Period
			"123,456",  // Comma
			"123:456",  // Colon
			"123;456",  // Semicolon
			"123\"456", // Double quote
			"123'456",  // Single quote
		}

		for _, postalCode := range invalidPostalCodes {
			t.Run(postalCode, func(t *testing.T) {
				err := validatePostalCode(postalCode)
				assert.Error(t, err)

				var errMap errsx.Map
				assert.True(t, errsx.As(err, &errMap))
				assert.Contains(t, errMap, PostalCodeInvalidFmtKey)
				assert.Equal(t, PostalCodeInvalidFmtMsg, errMap[PostalCodeInvalidFmtKey].Error())
			})
		}
	})

	t.Run("boundary cases", func(t *testing.T) {
		t.Run("minimum valid length", func(t *testing.T) {
			minLengthCode := strings.Repeat("A", PostalCodeMinLength) // Exactly 3 characters
			err := validatePostalCode(minLengthCode)
			assert.NoError(t, err)
		})

		t.Run("maximum valid length", func(t *testing.T) {
			maxLengthCode := strings.Repeat("A", PostalCodeMaxLength) // Exactly 10 characters
			err := validatePostalCode(maxLengthCode)
			assert.NoError(t, err)
		})

		t.Run("one character too short", func(t *testing.T) {
			tooShortCode := strings.Repeat("A", PostalCodeMinLength-1) // 2 characters
			err := validatePostalCode(tooShortCode)
			assert.Error(t, err)

			var errMap errsx.Map
			assert.True(t, errsx.As(err, &errMap))
			assert.Contains(t, errMap, PostalCodeTooShortKey)
		})

		t.Run("one character too long", func(t *testing.T) {
			tooLongCode := strings.Repeat("A", PostalCodeMaxLength+1) // 11 characters
			err := validatePostalCode(tooLongCode)
			assert.Error(t, err)

			var errMap errsx.Map
			assert.True(t, errsx.As(err, &errMap))
			assert.Contains(t, errMap, PostalCodeTooLongKey)
		})
	})

	t.Run("multiple validation errors should be reported", func(t *testing.T) {
		// Too short AND contains invalid characters
		err := validatePostalCode("A@")
		assert.Error(t, err)

		var errMap errsx.Map
		assert.True(t, errsx.As(err, &errMap))

		// Should have both length and format errors
		assert.Contains(t, errMap, PostalCodeTooShortKey)
		assert.Contains(t, errMap, PostalCodeInvalidFmtKey)
		assert.Equal(t, PostalCodeTooShortMsg, errMap[PostalCodeTooShortKey].Error())
		assert.Equal(t, PostalCodeInvalidFmtMsg, errMap[PostalCodeInvalidFmtKey].Error())
	})

	t.Run("international postal code formats", func(t *testing.T) {
		t.Run("US formats", func(t *testing.T) {
			usFormats := []string{
				"12345",      // Standard 5-digit
				"12345-6789", // ZIP+4
			}

			for _, postalCode := range usFormats {
				err := validatePostalCode(postalCode)
				assert.NoError(t, err, "expected US format %s to be valid", postalCode)
			}
		})

		t.Run("Canadian formats", func(t *testing.T) {
			canadianFormats := []string{
				"K1A 0A6", // Standard Canadian
				"M5V 3L9", // Standard Canadian
				"K1A0A6",  // Without space
			}

			for _, postalCode := range canadianFormats {
				err := validatePostalCode(postalCode)
				assert.NoError(t, err, "expected Canadian format %s to be valid", postalCode)
			}
		})

		t.Run("UK formats", func(t *testing.T) {
			ukFormats := []string{
				"SW1A 1AA", // Standard UK
				"M1 1AA",   // Standard UK
				"B33 8TH",  // Standard UK
				"SW1A1AA",  // Without space
			}

			for _, postalCode := range ukFormats {
				err := validatePostalCode(postalCode)
				assert.NoError(t, err, "expected UK format %s to be valid", postalCode)
			}
		})

		t.Run("European formats", func(t *testing.T) {
			europeanFormats := []string{
				"75001", // France
				"69002", // France
				"10115", // Germany
				"80331", // Germany
				"1010",  // Austria/Netherlands
			}

			for _, postalCode := range europeanFormats {
				err := validatePostalCode(postalCode)
				assert.NoError(t, err, "expected European format %s to be valid", postalCode)
			}
		})

		t.Run("mixed alphanumeric formats", func(t *testing.T) {
			mixedFormats := []string{
				"A1B 2C3", // Canadian-style with space
				"A1B-2C3", // With dash
				"ABC123",  // Letters then numbers
				"123ABC",  // Numbers then letters
			}

			for _, postalCode := range mixedFormats {
				err := validatePostalCode(postalCode)
				assert.NoError(t, err, "expected mixed format %s to be valid", postalCode)
			}
		})
	})

	t.Run("edge case whitespace and dash handling", func(t *testing.T) {
		t.Run("internal spaces should be preserved", func(t *testing.T) {
			spacedCodes := []string{
				"K1A 0A6",  // Canadian format with internal space
				"SW1A 1AA", // UK format with internal space
				"A B C",    // Multiple spaces
			}

			for _, postalCode := range spacedCodes {
				err := validatePostalCode(postalCode)
				assert.NoError(t, err, "expected spaced format %s to be valid", postalCode)
			}
		})

		t.Run("internal dashes should be preserved", func(t *testing.T) {
			dashedCodes := []string{
				"12345-6789", // US ZIP+4
				"A1B-2C3",    // Custom format with dash
				"123-456",    // Generic dash format
			}

			for _, postalCode := range dashedCodes {
				err := validatePostalCode(postalCode)
				assert.NoError(t, err, "expected dashed format %s to be valid", postalCode)
			}
		})
	})
}

