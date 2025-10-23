package domain

import (
	"time"

	"github.com/hengadev/errsx"
)

// Birth date validation constants
const (
	MinAgeYears = 13  // GDPR compliance - minimum age
	MaxAgeYears = 120 // Reasonable maximum age
)

// Birth date validation error keys
const (
	BirthDateFutureKey   = "birthdate_future"
	BirthDateTooYoungKey = "birthdate_too_young"
	BirthDateTooOldKey   = "birthdate_too_old"
)

// Birth date validation error messages
const (
	BirthDateFutureMsg   = "birth date cannot be in the future"
	BirthDateTooYoungMsg = "user must be at least 13 years old"
	BirthDateTooOldMsg   = "birth date is not valid"
)

func validateBirthDate(birthdate time.Time) error {
	var errs errsx.Map

	now := time.Now()

	// Check if birth date is in the future
	if birthdate.After(now) {
		errs.Set(BirthDateFutureKey, BirthDateFutureMsg)
	}

	// Check minimum age (13 years for GDPR compliance)
	minBirthDate := now.AddDate(-MinAgeYears, 0, 0)
	if birthdate.After(minBirthDate) {
		errs.Set(BirthDateTooYoungKey, BirthDateTooYoungMsg)
	}

	// Check maximum age (reasonable limit of 120 years)
	maxBirthDate := now.AddDate(-MaxAgeYears, 0, 0)
	if birthdate.Before(maxBirthDate) {
		errs.Set(BirthDateTooOldKey, BirthDateTooOldMsg)
	}

	return errs.AsError()
}
