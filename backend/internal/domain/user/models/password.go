package models

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

const passwordMinLength = 8

func ValidatePassword(p string) error {
	var errs errsx.Map
	if len(p) < passwordMinLength {
		errs.Set("password length", fmt.Sprintf("expect at least %d caracter", passwordMinLength))
	}
	pwned, err := CheckPasswordPwned(p)
	if err != nil {
		errs.Set("password powned verification attempt", err)
	}
	if pwned {
		errs.Set("password powned", err)
	}

	return errs.AsError()

}
func NewPassword(p string) (Password, error) {
	var errs errsx.Map
	if err := ValidatePassword(p); err != nil {
		errs.Set("validate password", err)
	}
	return Password(p), errs.AsError()
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
