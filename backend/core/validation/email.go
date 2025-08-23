package validation

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/hengadev/errsx"
)

const EmailMaxLength = 255

var (
	invalidEmailChars = regexp.MustCompile(`[^a-zA-Z0-9+.@_~\-]`)
	validEmailSeq     = regexp.MustCompile(`^[a-zA-Z0-9+._~\-]+@[a-zA-Z0-9+._~\-]+(\.[a-zA-Z0-9+._~\-]+)+$`)
)

func ValidateEmail(email string) error {
	var errs errsx.Map

	// Length check first (before expensive parsing)
	if rc := utf8.RuneCountInString(email); rc > EmailMaxLength {
		errs.Set("email_length", fmt.Sprintf("email cannot exceed %v characters", EmailMaxLength))
	}

	if strings.TrimSpace(email) == "" {
		errs.Set("email_required", "email cannot be empty")
	}

	if strings.ContainsAny(email, " \t\n\r") {
		errs.Set("email_whitespace", "email cannot contain whitespace")
	}

	if strings.ContainsAny(email, `"'`) {
		errs.Set("email_quotes", "email cannot contain quotes")
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		email = strings.TrimSpace(email)
		msg := strings.TrimPrefix(strings.ToLower(err.Error()), "mail: ")

		switch {
		case strings.Contains(msg, "missing '@'"):
			errs.Set("email_format", "email missing @ sign")
		case strings.HasPrefix(email, "@"):
			errs.Set("email_format", "email missing part before @ sign")
		case strings.HasSuffix(email, "@"):
			errs.Set("email_format", "email missing part after @ sign")
		default:
			errs.Set("email_format", "invalid email format")
		}
	}

	if addr != nil {
		if addr.Name != "" {
			errs.Set("email_name", "email cannot include a name")
		}

		if matches := invalidEmailChars.FindAllString(addr.Address, -1); len(matches) != 0 {
			errs.Set("email_chars", fmt.Sprintf("email contains invalid characters: %v", matches))
		}

		// Stricter validation - require TLD
		if !validEmailSeq.MatchString(addr.Address) {
			_, end, _ := strings.Cut(addr.Address, "@")
			if !strings.Contains(end, ".") {
				errs.Set("email_format", "email missing top-level domain (e.g. .com, .org)")
			} else {
				errs.Set("email_format", "invalid email format")
			}
		}
	}

	return errs.AsError()
}

