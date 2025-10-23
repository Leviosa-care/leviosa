package domain

import (
	"fmt"
	"strings"

	"github.com/hengadev/errsx"
)

type Gender string

const (
	GenderMan            Gender = "man"
	GenderWoman          Gender = "woman"
	GenderNonBinary      Gender = "non_binary"
	GenderPreferNotToSay Gender = "prefer_not_to_say"
	GenderCustom         Gender = "custom"
)

// Gender validation constants
const (
	CustomGenderMaxLength = 50
)

// Gender validation error keys
const (
	GenderInvalidKey        = "gender_invalid"
	CustomGenderRequiredKey = "custom_gender_required"
	CustomGenderTooLongKey  = "custom_gender_too_long"
	CustomGenderCharsKey    = "custom_gender_invalid_chars"
)

// Gender validation error messages
const (
	GenderInvalidMsg        = "invalid gender value"
	CustomGenderRequiredMsg = "custom gender is required when gender is 'custom'"
	CustomGenderTooLongMsg  = "custom gender must be no more than 50 characters"
	CustomGenderCharsMsg    = "custom gender contains invalid characters"
)

type GenderInput struct {
	Gender       Gender `json:"gender" validate:"required"`
	CustomGender string `json:"customGender,omitempty"`
}

func (g *GenderInput) ValidateGender() error {
	var errs errsx.Map

	switch g.Gender {
	case GenderMan, GenderWoman, GenderNonBinary, GenderPreferNotToSay:
		// Valid predefined gender values
		break
	case GenderCustom:
		if customErr := validateCustomGender(g); customErr != nil {
			// Merge custom gender validation errors
			var customErrs errsx.Map
			if errsx.As(customErr, &customErrs) {
				for key, err := range customErrs {
					errs.Set(key, err)
				}
			} else {
				errs.Set(CustomGenderRequiredKey, customErr)
			}
		}
	default:
		errs.Set(GenderInvalidKey, fmt.Sprintf("%s: %s", GenderInvalidMsg, g.Gender))
	}

	return errs.AsError()
}

func (g Gender) String() string {
	return string(g)
}

func validateCustomGender(gender *GenderInput) error {
	var errs errsx.Map

	trimmed := strings.TrimSpace(gender.CustomGender)

	if len(trimmed) == 0 {
		errs.Set(CustomGenderRequiredKey, CustomGenderRequiredMsg)
	}

	if len(trimmed) > CustomGenderMaxLength {
		errs.Set(CustomGenderTooLongKey, CustomGenderTooLongMsg)
	}

	// Reject dangerous patterns for GDPR compliance and security
	if strings.ContainsAny(trimmed, "<>;\"'&") {
		errs.Set(CustomGenderCharsKey, CustomGenderCharsMsg)
	}

	return errs.AsError()
}
