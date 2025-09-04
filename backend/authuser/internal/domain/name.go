package domain

import (
	"fmt"
	"strings"

	"github.com/hengadev/errsx"
)

// Name validation constants
const (
	NameMinLength = 2
	NameMaxLength = 50
)

// Name validation error keys
const (
	NameRequiredKey     = "name_required"
	NameTooShortKey     = "name_too_short"
	NameTooLongKey      = "name_too_long"
	NameInvalidCharsKey = "name_invalid_chars"
)

// Name validation error messages
const (
	NameRequiredMsg     = "is required"
	NameTooShortMsg     = "must be at least 2 characters long"
	NameTooLongMsg      = "must be no more than 50 characters long"
	NameInvalidCharsMsg = "contains invalid characters"
)

func validateName(name, fieldName string) error {
	var errs errsx.Map

	trimmed := strings.TrimSpace(name)

	if len(trimmed) == 0 {
		errs.Set(NameRequiredKey, fmt.Sprintf("%s %s", fieldName, NameRequiredMsg))
	}

	if len(trimmed) < NameMinLength {
		errs.Set(NameTooShortKey, fmt.Sprintf("%s %s", fieldName, NameTooShortMsg))
	}

	if len(trimmed) > NameMaxLength {
		errs.Set(NameTooLongKey, fmt.Sprintf("%s %s", fieldName, NameTooLongMsg))
	}

	// Reject dangerous patterns for GDPR compliance
	if strings.ContainsAny(trimmed, "<>;\"'&") {
		errs.Set(NameInvalidCharsKey, fmt.Sprintf("%s %s", fieldName, NameInvalidCharsMsg))
	}

	return errs.AsError()
}
