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

// Email validation error keys
const (
	EmailLengthKey     = "email_length"
	EmailRequiredKey   = "email_required"
	EmailWhitespaceKey = "email_whitespace"
	EmailQuotesKey     = "email_quotes"
	EmailFormatKey     = "email_format"
	EmailNameKey       = "email_name"
	EmailCharsKey      = "email_chars"
)

// Email validation error messages
const (
	EmailRequiredMsg       = "email cannot be empty"
	EmailWhitespaceMsg     = "email cannot contain whitespace"
	EmailQuotesMsg         = "email cannot contain quotes"
	EmailMissingBeforeMsg  = "email missing part before @ sign"
	EmailMissingAfterMsg   = "email missing part after @ sign"
	EmailMissingAtMsg      = "email missing @ sign"
	EmailInvalidFormatMsg  = "invalid email format"
	EmailNameNotAllowedMsg = "email cannot include a name"
	EmailMissingTLDMsg     = "email missing top-level domain (e.g. .com, .org)"
)

var (
	invalidEmailChars = regexp.MustCompile(`[^a-zA-Z0-9+.@_~\-]`)
	validEmailSeq     = regexp.MustCompile(`^[a-zA-Z0-9+._~\-]+@[a-zA-Z0-9+._~\-]+(\.[a-zA-Z0-9+._~\-]+)+$`)
)

func ValidateEmail(email string) error {
	var errs errsx.Map

	// Length check first (before expensive parsing)
	if rc := utf8.RuneCountInString(email); rc > EmailMaxLength {
		errs.Set(EmailLengthKey, fmt.Sprintf("email cannot exceed %v characters", EmailMaxLength))
	}

	if strings.TrimSpace(email) == "" {
		errs.Set(EmailRequiredKey, EmailRequiredMsg)
	}

	if strings.ContainsAny(email, " \t\n\r") {
		errs.Set(EmailWhitespaceKey, EmailWhitespaceMsg)
	}

	if strings.ContainsAny(email, `"'`) {
		errs.Set(EmailQuotesKey, EmailQuotesMsg)
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		email = strings.TrimSpace(email)
		msg := strings.TrimPrefix(strings.ToLower(err.Error()), "mail: ")

		switch {
		case strings.HasPrefix(email, "@"):
			errs.Set(EmailFormatKey, EmailMissingBeforeMsg)
		case strings.HasSuffix(email, "@"):
			errs.Set(EmailFormatKey, EmailMissingAfterMsg)
		case strings.Contains(msg, "missing '@'"):
			errs.Set(EmailFormatKey, EmailMissingAtMsg)
		default:
			errs.Set(EmailFormatKey, EmailInvalidFormatMsg)
		}
	}

	if addr != nil {
		if addr.Name != "" {
			errs.Set(EmailNameKey, EmailNameNotAllowedMsg)
		}

		if matches := invalidEmailChars.FindAllString(addr.Address, -1); len(matches) != 0 {
			errs.Set(EmailCharsKey, fmt.Sprintf("email contains invalid characters: %v", matches))
		}

		// Stricter validation - require TLD
		if !validEmailSeq.MatchString(addr.Address) {
			_, end, _ := strings.Cut(addr.Address, "@")
			if !strings.Contains(end, ".") {
				errs.Set(EmailFormatKey, EmailMissingTLDMsg)
			} else {
				errs.Set(EmailFormatKey, EmailInvalidFormatMsg)
			}
		}
	}

	return errs.AsError()
}
