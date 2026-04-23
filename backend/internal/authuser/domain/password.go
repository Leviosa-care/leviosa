package domain

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hengadev/errsx"
)

type Password string

// Password validation constants
const (
	PasswordMinLength = 8
)

// Password validation error keys
const (
	PasswordLengthKey     = "password_length"
	PasswordPwnedCheckKey = "password_pwned_check"
	PasswordPwnedKey      = "password_pwned"
)

// Password validation error messages
const (
	PasswordLengthMsg     = "password must be at least 8 characters"
	PasswordPwnedCheckMsg = "failed to verify password security"
	PasswordPwnedMsg      = "password has been found in data breaches and cannot be used"
)

// ValidatePasswordFormat validates only the password format (length)
// Used for login where we only need to verify the password is well-formed
func ValidatePasswordFormat(p string) error {
	var errs errsx.Map

	if len(p) < PasswordMinLength {
		errs.Set(PasswordLengthKey, PasswordLengthMsg)
	}

	return errs.AsError()
}

// ValidatePassword validates password format AND checks if it has been pwned
// Used for registration and password resets where we want to enforce strong passwords
func ValidatePassword(p string) error {
	var errs errsx.Map

	if len(p) < PasswordMinLength {
		errs.Set(PasswordLengthKey, PasswordLengthMsg)
	}

	pwned, err := CheckPasswordPwned(p)
	if err != nil {
		errs.Set(PasswordPwnedCheckKey, PasswordPwnedCheckMsg)
	}
	if pwned {
		errs.Set(PasswordPwnedKey, PasswordPwnedMsg)
	}

	return errs.AsError()
}
func NewPassword(p string) (Password, error) {
	if err := ValidatePassword(p); err != nil {
		return Password(""), err
	}
	return Password(p), nil
}

func (p Password) String() string {
	return string(p)
}

// CheckPasswordPwned uses the k-Anonymity model from Have I Been Pwned
func CheckPasswordPwned(password string) (bool, error) {
	// SHA1 hash of the password
	hash := sha1.Sum([]byte(password))
	hashStr := strings.ToUpper(hex.EncodeToString(hash[:]))

	// First 5 characters go to the API
	prefix := hashStr[:5]
	suffix := hashStr[5:]

	// Request to the API
	resp, err := http.Get("https://api.pwnedpasswords.com/range/" + prefix)
	if err != nil {
		return false, fmt.Errorf("failed to query Pwned Passwords API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read API response: %w", err)
	}

	// Check if the suffix is in the list
	for line := range strings.SplitSeq(string(body), "\n") {
		parts := strings.Split(line, ":")
		if len(parts) > 1 && parts[0] == suffix {
			return true, nil // password is compromised
		}
	}

	return false, nil // password not found }
}
