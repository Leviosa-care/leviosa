package models

import (
	"fmt"
	"strings"
)

type Gender string

const (
	GenderMan            Gender = "man"
	GenderWoman          Gender = "woman"
	GenderNonBinary      Gender = "non_binary"
	GenderPreferNotToSay Gender = "prefer_not_to_say"
	GenderCustom         Gender = "custom"
)

type GenderInput struct {
	Gender       Gender `json:"gender" validate:"required"`
	CustomGender string `json:"customGender,omitempty"`
}

func (g *GenderInput) ValidateGender() error {
	switch g.Gender {
	case GenderMan, GenderWoman, GenderNonBinary, GenderPreferNotToSay:
		break
	case GenderCustom:
		return validateCustomGender(g)
	default:
		return fmt.Errorf("invalid gender value: %s", g.Gender)
	}
	return nil
}

func (g Gender) String() string {
	return string(g)
}

func validateCustomGender(gender *GenderInput) error {
	if strings.TrimSpace(gender.CustomGender) == "" {
		return fmt.Errorf("customGender is required when gender is 'custom'")
	}
	trimmed := strings.TrimSpace(string(gender.Gender))
	if len(trimmed) == 0 || len(trimmed) > 50 {
		return fmt.Errorf("invalid gender value")
	}
	// Optional: reject dangerous patterns (e.g., scripts, SQL-y strings)
	if strings.ContainsAny(trimmed, "<>;") {
		return fmt.Errorf("Invalid characters in gender value")
	}
	return nil
}
