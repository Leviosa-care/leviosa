package domain

import (
	"context"
	"fmt"

	"github.com/hengadev/errsx"
)

func ValidateOTP(ctx context.Context, code string, expectedLength int) error {
	var errs errsx.Map
	// Check empty
	if code == "" {
		errs.Set("code missing", "code is required")
	}

	// Check length
	if len(code) != expectedLength {
		errs.Set("invalid length", fmt.Sprintf("code must be %d digits", expectedLength))
	}

	// Check numeric
	for _, r := range code {
		if r < '0' || r > '9' {
			errs.Set("invalid characters", "code must only contain digits")
			break
		}
	}
	return errs.AsError()
}
