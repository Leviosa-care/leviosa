package models

import (
	"context"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/hengadev/errsx"
)

const EmailMaxLength = 100

var (
	invalidEmailChars = regexp.MustCompile(`[^a-zA-Z0-9+.@_~\-]`)
	validEmailSeq     = regexp.MustCompile(`^[a-zA-Z0-9+._~\-]+@[a-zA-Z0-9+._~\-]+(\.[a-zA-Z0-9+._~\-]+)+$`)
)

type Email string

func (e Email) Valid(ctx context.Context) error {
	return ValidateEmail(e.String())
}

func ValidateEmail(email string) error {
	var errs errsx.Map
	if strings.TrimSpace(email) == "" {
		errs.Set("emptiness", "cannot be empty")
	}
	if strings.ContainsAny(email, " \t\n\r") {
		errs.Set("whitespace", "cannot contain whitespace")
	}
	if strings.ContainsAny(email, `"'`) {
		errs.Set("quotes", "cannot contain quotes")
	}
	if rc := utf8.RuneCountInString(email); rc > EmailMaxLength {
		errs.Set("max length", fmt.Sprintf("cannot be a over %v characters in length", EmailMaxLength))
	}
	addr, err := mail.ParseAddress(email)
	if err != nil {
		email = strings.TrimSpace(email)
		msg := strings.TrimPrefix(strings.ToLower(err.Error()), "mail: ")

		switch {
		case strings.Contains(msg, "missing '@'"):
			errs.Set("@ sign", "missing the @ sign")

		case strings.HasPrefix(email, "@"):
			errs.Set("@ sign", "missing part before the @ sign")

		case strings.HasSuffix(email, "@"):
			errs.Set("@ sign", "missing part after the @ sign")
		}
	}
	if addr != nil {
		if addr.Name != "" {
			errs.Set("include name", "cannot not include a name")
		}
		if matches := invalidEmailChars.FindAllString(addr.Address, -1); len(matches) != 0 {
			errs.Set("invalid characters", fmt.Sprintf("cannot contain: %v", matches))
		}
		if !validEmailSeq.MatchString(addr.Address) {
			_, end, _ := strings.Cut(addr.Address, "@")
			if !strings.Contains(end, ".") {
				errs.Set("top level domain", "missing top-level domain, e.g. .com, .co.uk, etc.")
			}

			errs.Set("not email address", "must be an email address, e.g. email@example.com")
		}
	}
	return errs.AsError()
}

func NewEmail(email string) (Email, error) {
	if err := ValidateEmail(email); err != nil {
		return "", err
	}
	return Email(email), nil
}

func (e Email) String() string {
	return string(e)
}
