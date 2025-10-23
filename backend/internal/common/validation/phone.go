package validation

import (
	"regexp"
	"strings"

	"github.com/hengadev/errsx"
)

// Phone validation constants
const (
	PhoneMinLength = 10
	PhoneMaxLength = 20
)

// Phone validation error keys
const (
	PhoneTooShortKey = "phone_too_short"
	PhoneTooLongKey  = "phone_too_long"
	PhoneFormatKey   = "phone_format"
)

// Phone validation error messages
const (
	PhoneTooShortMsg = "phone must be at least 10 characters"
	PhoneTooLongMsg  = "phone cannot exceed 20 characters"
	PhoneFormatMsg   = "invalid phone number format"
)

// French phone number validation (can be extended for international later)
var frenchPhoneStrict = regexp.MustCompile(`^(0[1-5]|06|07)\d{8}$`)

func ValidatePhone(phone string) error {
	var errs errsx.Map

	// Trim whitespace for validation
	trimmed := strings.TrimSpace(phone)

	// Length checks
	if len(trimmed) < PhoneMinLength {
		errs.Set(PhoneTooShortKey, PhoneTooShortMsg)
	}
	if len(trimmed) > PhoneMaxLength {
		errs.Set(PhoneTooLongKey, PhoneTooLongMsg)
	}

	// French phone format validation
	if !frenchPhoneStrict.MatchString(trimmed) {
		errs.Set(PhoneFormatKey, PhoneFormatMsg)
	}

	return errs.AsError()
}
