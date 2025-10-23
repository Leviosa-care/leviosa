package domain

import (
	"strings"

	"github.com/hengadev/errsx"
)

// City validation constants
const (
	CityMinLength = 2
	CityMaxLength = 100
)

// City validation error keys
const (
	CityRequiredKey     = "city_required"
	CityTooShortKey     = "city_too_short"
	CityTooLongKey      = "city_too_long"
	CityInvalidCharsKey = "city_invalid_chars"
)

// City validation error messages
const (
	CityRequiredMsg     = "city is required"
	CityTooShortMsg     = "city must be at least 2 characters long"
	CityTooLongMsg      = "city must be no more than 100 characters long"
	CityInvalidCharsMsg = "city contains invalid characters"
)

func validateCity(city string) error {
	var errs errsx.Map

	trimmed := strings.TrimSpace(city)

	if len(trimmed) == 0 {
		errs.Set(CityRequiredKey, CityRequiredMsg)
	}

	if len(trimmed) < CityMinLength {
		errs.Set(CityTooShortKey, CityTooShortMsg)
	}

	if len(trimmed) > CityMaxLength {
		errs.Set(CityTooLongKey, CityTooLongMsg)
	}

	// Reject dangerous patterns for GDPR compliance
	if strings.ContainsAny(trimmed, "<>;\"'&") {
		errs.Set(CityInvalidCharsKey, CityInvalidCharsMsg)
	}

	return errs.AsError()
}
