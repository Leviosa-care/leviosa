package models

import (
	"errors"
	"time"
)

// ValidateBirthDate checks if the user is at least 18 years old.
func ValidateBirthDate(birthdate time.Time) error {
	now := time.Now()
	eighteenYearsAgo := now.AddDate(-18, 0, 0)

	if birthdate.After(eighteenYearsAgo) {
		return errors.New("user must be at least 18 years old")
	}
	return nil
}
