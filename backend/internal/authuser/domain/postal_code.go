package domain

import (
	"strings"

	"github.com/hengadev/errsx"
)

// Postal code validation constants
const (
	PostalCodeMinLength = 3
	PostalCodeMaxLength = 10
)

// Postal code validation error keys
const (
	PostalCodeRequiredKey   = "postal_code_required"
	PostalCodeTooShortKey   = "postal_code_too_short"
	PostalCodeTooLongKey    = "postal_code_too_long"
	PostalCodeInvalidFmtKey = "postal_code_invalid_format"
)

// Postal code validation error messages
const (
	PostalCodeRequiredMsg   = "postal code is required"
	PostalCodeTooShortMsg   = "postal code must be at least 3 characters"
	PostalCodeTooLongMsg    = "postal code must be no more than 10 characters"
	PostalCodeInvalidFmtMsg = "postal code contains invalid characters"
)

func validatePostalCode(postalCode string) error {
	var errs errsx.Map

	trimmed := strings.TrimSpace(postalCode)

	if len(trimmed) == 0 {
		errs.Set(PostalCodeRequiredKey, PostalCodeRequiredMsg)
	}

	if len(trimmed) < PostalCodeMinLength {
		errs.Set(PostalCodeTooShortKey, PostalCodeTooShortMsg)
	}

	if len(trimmed) > PostalCodeMaxLength {
		errs.Set(PostalCodeTooLongKey, PostalCodeTooLongMsg)
	}

	// Basic alphanumeric validation for international postal codes
	for _, r := range trimmed {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == ' ' || r == '-') {
			errs.Set(PostalCodeInvalidFmtKey, PostalCodeInvalidFmtMsg)
			break
		}
	}

	return errs.AsError()
}
