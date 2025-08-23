package validation

import (
	"errors"
	"regexp"
	"strings"
)

// French phone number validation (can be extended for international later)
var frenchPhoneStrict = regexp.MustCompile(`^(0[1-5]|06|07)\d{8}$`)

func ValidatePhone(phone string) error {
	// Trim whitespace for validation
	trimmed := strings.TrimSpace(phone)

	// Length checks
	if len(trimmed) < 10 {
		return errors.New("phone must be at least 10 characters")
	}
	if len(trimmed) > 20 {
		return errors.New("phone cannot exceed 20 characters")
	}

	// French phone format validation
	if !frenchPhoneStrict.MatchString(trimmed) {
		return errors.New("invalid phone number format")
	}

	return nil
}

