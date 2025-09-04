package domain

import (
	"fmt"
	"strings"

	"github.com/hengadev/errsx"
)

// Address validation constants
const (
	AddressMinLength = 5
	AddressMaxLength = 200
)

// Address validation error keys
const (
	AddressRequiredKey     = "address_required"
	AddressTooShortKey     = "address_too_short"
	AddressTooLongKey      = "address_too_long"
	AddressInvalidCharsKey = "address_invalid_chars"
)

// Address validation error messages
const (
	AddressRequiredMsg     = "is required"
	AddressTooShortMsg     = "must be at least 5 characters long"
	AddressTooLongMsg      = "must be no more than 200 characters long"
	AddressInvalidCharsMsg = "contains invalid characters"
)

func validateAddress1(address string) error {
	return validateAddress(address, "address1")
}

func validateAddress2(address string) error {
	if address != "" {
		return validateAddress(address, "address2")
	}
	return nil
}

func validateAddress(address, fieldName string) error {
	var errs errsx.Map

	trimmed := strings.TrimSpace(address)

	if len(trimmed) == 0 {
		errs.Set(AddressRequiredKey, fmt.Sprintf("%s %s", fieldName, AddressRequiredMsg))
	}

	if len(trimmed) < AddressMinLength {
		errs.Set(AddressTooShortKey, fmt.Sprintf("%s %s", fieldName, AddressTooShortMsg))
	}

	if len(trimmed) > AddressMaxLength {
		errs.Set(AddressTooLongKey, fmt.Sprintf("%s %s", fieldName, AddressTooLongMsg))
	}

	// Reject dangerous patterns for GDPR compliance
	if strings.ContainsAny(trimmed, "<>;\"'&") {
		errs.Set(AddressInvalidCharsKey, fmt.Sprintf("%s %s", fieldName, AddressInvalidCharsMsg))
	}

	return errs.AsError()
}
